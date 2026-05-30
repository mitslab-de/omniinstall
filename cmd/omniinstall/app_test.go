package main

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	adapter "github.com/mitslab-de/omniinstall/internal/adapters"
	"github.com/mitslab-de/omniinstall/internal/engine"
	"github.com/mitslab-de/omniinstall/internal/install"
	"github.com/mitslab-de/omniinstall/internal/source"
	"github.com/mitslab-de/omniinstall/internal/state"
)

// mockStore is an in-memory state.Store for testing.
type mockStore struct {
	records map[string]state.LocalInstallation
}

func newMockStore() *mockStore {
	return &mockStore{records: make(map[string]state.LocalInstallation)}
}

func (m *mockStore) Record(r state.LocalInstallation) error {
	m.records[r.ApplicationID] = r
	return nil
}

func (m *mockStore) Get(appID string) (state.LocalInstallation, bool, error) {
	r, ok := m.records[appID]
	return r, ok, nil
}

func (m *mockStore) List() ([]state.LocalInstallation, error) {
	list := make([]state.LocalInstallation, 0, len(m.records))
	for _, v := range m.records {
		list = append(list, v)
	}
	return list, nil
}

func (m *mockStore) MarkRemoved(appID string) error {
	r, ok := m.records[appID]
	if !ok {
		return state.ErrNotFound
	}
	r.InstallStatus = state.StatusRemoved
	m.records[appID] = r
	return nil
}

// mockAdapter is a configurable test double for adapter.Adapter.
type mockAdapter struct {
	name          string
	available     bool
	handledTypes  []source.Type
	installResult *adapter.Result
	removeResult  *adapter.Result
	removedIDs    []string
}

func (m *mockAdapter) Name() string      { return m.name }
func (m *mockAdapter) IsAvailable() bool { return m.available }
func (m *mockAdapter) CanHandle(t source.Type) bool {
	for _, ht := range m.handledTypes {
		if ht == t {
			return true
		}
	}
	return false
}
func (m *mockAdapter) CheckInstalled(id string) (adapter.InstalledState, error) {
	return adapter.StateNotInstalled, nil
}
func (m *mockAdapter) Install(plan *install.Plan) (*adapter.Result, error) {
	if m.installResult != nil {
		return m.installResult, nil
	}
	return &adapter.Result{Success: true, Message: "installed", ChangedSystem: true}, nil
}
func (m *mockAdapter) Remove(id string) (*adapter.Result, error) {
	m.removedIDs = append(m.removedIDs, id)
	if m.removeResult != nil {
		return m.removeResult, nil
	}
	return &adapter.Result{Success: true, Message: "removed", ChangedSystem: true}, nil
}
func (m *mockAdapter) Verify(plan *install.Plan) (*adapter.VerificationResult, error) {
	return &adapter.VerificationResult{Verified: true, Details: "ok"}, nil
}

// successEngine creates an engine with a mock adapter that succeeds.
func successEngine(t source.Type) *engine.DefaultEngine {
	a := &mockAdapter{
		name:         "mock",
		available:    true,
		handledTypes: []source.Type{t},
	}
	return engine.NewDefaultEngine([]adapter.Adapter{a}, nil)
}

func aptManagers() []source.Type {
	return []source.Type{source.TypeAPT}
}

func flatpakManagers() []source.Type {
	return []source.Type{source.TypeFlatpak}
}

// ----- Tests -----

func TestAppSearch_Found(t *testing.T) {
	buf := &strings.Builder{}
	app := newAppWithStore(buf, newMockStore(), aptManagers())
	if err := app.search("obs"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "obs-studio") {
		t.Errorf("expected output to contain 'obs-studio', got:\n%s", buf.String())
	}
}

func TestAppSearch_NotFound(t *testing.T) {
	buf := &strings.Builder{}
	app := newAppWithStore(buf, newMockStore(), aptManagers())
	if err := app.search("xyznonexistent"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No applications found") {
		t.Errorf("expected 'No applications found' in output, got:\n%s", buf.String())
	}
}

func TestAppSearch_EmptyQuery(t *testing.T) {
	buf := &strings.Builder{}
	app := newAppWithStore(buf, newMockStore(), aptManagers())
	if err := app.search(""); err == nil {
		t.Error("expected error for empty search query")
	}
}

func TestAppSearch_JSONFound(t *testing.T) {
	buf := &strings.Builder{}
	app := newAppWithStore(buf, newMockStore(), aptManagers())
	if err := app.searchJSON("obs"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out struct {
		Query   string `json:"query"`
		Results []struct {
			ApplicationID string `json:"application_id"`
		} `json:"results"`
	}
	if err := json.Unmarshal([]byte(buf.String()), &out); err != nil {
		t.Fatalf("expected valid json output, got error: %v", err)
	}
	if out.Query != "obs" {
		t.Fatalf("expected query obs, got %q", out.Query)
	}
	if len(out.Results) == 0 || out.Results[0].ApplicationID == "" {
		t.Fatalf("expected at least one json result, got %+v", out.Results)
	}
}

func TestAppInstall_Success(t *testing.T) {
	buf := &strings.Builder{}
	store := newMockStore()
	eng := successEngine(source.TypeAPT)
	app := newAppForTest(buf, eng, store, aptManagers())

	if err := app.install("git"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, found, _ := store.Get("git")
	if !found {
		t.Error("expected git to be recorded in local state")
	}
}

func TestAppInstall_NotFound(t *testing.T) {
	buf := &strings.Builder{}
	app := newAppForTest(buf, successEngine(source.TypeAPT), newMockStore(), aptManagers())
	err := app.install("xyznonexistent")
	if err == nil {
		t.Error("expected error for unknown application")
	}
}

func TestAppInstall_EmptyAppID(t *testing.T) {
	buf := &strings.Builder{}
	app := newAppWithStore(buf, newMockStore(), aptManagers())
	if err := app.install(""); err == nil {
		t.Error("expected error for empty appID")
	}
}

func TestAppInstall_NoCompatibleSource(t *testing.T) {
	buf := &strings.Builder{}
	// git only has APT source; Flatpak managers only — should fail
	app := newAppForTest(buf, successEngine(source.TypeFlatpak), newMockStore(), flatpakManagers())
	err := app.install("git")
	if err == nil {
		t.Error("expected error when no compatible source available for git on flatpak-only system")
	}
}

func TestAppInstall_HighRiskRequiresConfirmation_Confirmed(t *testing.T) {
	buf := &strings.Builder{}
	store := newMockStore()
	eng := successEngine(source.TypeAPT)
	app := newAppForTest(buf, eng, store, aptManagers())
	app.sources["git"] = []source.Source{
		{
			ApplicationID:    "git",
			SourceType:       source.TypeAPT,
			SourceIdentifier: "git",
			TrustLevel:       source.TrustVerified,
			RiskLevel:        source.RiskHigh,
		},
	}
	app.in = strings.NewReader("yes\n")

	if err := app.install("git"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, found, _ := store.Get("git"); !found {
		t.Fatal("expected git to be recorded in local state")
	}
}

func TestAppInstall_HighRiskRequiresConfirmation_Denied(t *testing.T) {
	buf := &strings.Builder{}
	store := newMockStore()
	eng := successEngine(source.TypeAPT)
	app := newAppForTest(buf, eng, store, aptManagers())
	app.sources["git"] = []source.Source{
		{
			ApplicationID:    "git",
			SourceType:       source.TypeAPT,
			SourceIdentifier: "git",
			TrustLevel:       source.TrustVerified,
			RiskLevel:        source.RiskHigh,
		},
	}
	app.in = strings.NewReader("no\n")

	err := app.install("git")
	if err == nil || !strings.Contains(err.Error(), "cancelled by user") {
		t.Fatalf("expected cancellation error, got: %v", err)
	}
	if _, found, _ := store.Get("git"); found {
		t.Fatal("did not expect git to be recorded in local state")
	}
}

func TestAppInstall_HighRiskRequiresInteractiveInput(t *testing.T) {
	buf := &strings.Builder{}
	store := newMockStore()
	eng := successEngine(source.TypeAPT)
	app := newAppForTest(buf, eng, store, aptManagers())
	app.sources["git"] = []source.Source{
		{
			ApplicationID:    "git",
			SourceType:       source.TypeAPT,
			SourceIdentifier: "git",
			TrustLevel:       source.TrustVerified,
			RiskLevel:        source.RiskHigh,
		},
	}
	app.in = nil

	err := app.install("git")
	if err == nil || !strings.Contains(err.Error(), "requires interactive confirmation") {
		t.Fatalf("expected interactive confirmation error, got: %v", err)
	}
}

func TestAppRemove_Success(t *testing.T) {
	buf := &strings.Builder{}
	store := newMockStore()
	// Pre-populate state.
	_ = store.Record(state.LocalInstallation{
		ApplicationID:      "git",
		SourceType:         source.TypeAPT,
		SourceIdentifier:   "git",
		InstallTimestamp:   time.Now(),
		InstallStatus:      state.StatusInstalled,
		VerificationStatus: state.VerificationPassed,
	})
	eng := successEngine(source.TypeAPT)
	app := newAppForTest(buf, eng, store, aptManagers())

	if err := app.remove("git"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec, found, _ := store.Get("git")
	if !found {
		t.Fatal("expected state record to still exist")
	}
	if rec.InstallStatus != state.StatusRemoved {
		t.Errorf("expected StatusRemoved, got %s", rec.InstallStatus)
	}
}

func TestAppRemove_EmptyAppID(t *testing.T) {
	buf := &strings.Builder{}
	app := newAppWithStore(buf, newMockStore(), aptManagers())
	if err := app.remove(""); err == nil {
		t.Error("expected error for empty appID")
	}
}

func TestAppRemove_RequiresInstalledStateRecord(t *testing.T) {
	buf := &strings.Builder{}
	app := newAppForTest(buf, successEngine(source.TypeAPT), newMockStore(), aptManagers())
	if err := app.remove("git"); err == nil {
		t.Error("expected error when app is not recorded as installed")
	}
}

func TestAppRemove_UsesRecordedSourceDetails(t *testing.T) {
	buf := &strings.Builder{}
	store := newMockStore()
	_ = store.Record(state.LocalInstallation{
		ApplicationID:      "obs-studio",
		SourceType:         source.TypeFlatpak,
		SourceIdentifier:   "com.obsproject.Studio",
		InstallTimestamp:   time.Now(),
		InstallStatus:      state.StatusInstalled,
		VerificationStatus: state.VerificationPassed,
	})

	aptAdapter := &mockAdapter{
		name:         "apt",
		available:    true,
		handledTypes: []source.Type{source.TypeAPT},
	}
	flatpakAdapter := &mockAdapter{
		name:         "flatpak",
		available:    true,
		handledTypes: []source.Type{source.TypeFlatpak},
	}
	eng := engine.NewDefaultEngine([]adapter.Adapter{aptAdapter, flatpakAdapter}, nil)
	app := newAppForTest(buf, eng, store, flatpakManagers())

	if err := app.remove("obs-studio"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(aptAdapter.removedIDs) != 0 {
		t.Fatalf("expected apt adapter not used, got %v", aptAdapter.removedIDs)
	}
	if len(flatpakAdapter.removedIDs) != 1 || flatpakAdapter.removedIDs[0] != "com.obsproject.Studio" {
		t.Fatalf("expected flatpak adapter remove identifier com.obsproject.Studio, got %v", flatpakAdapter.removedIDs)
	}
}

func TestAppExplain_Found(t *testing.T) {
	buf := &strings.Builder{}
	app := newAppWithStore(buf, newMockStore(), aptManagers())
	if err := app.explain("obs-studio"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "obs-studio") {
		t.Errorf("expected output to mention 'obs-studio', got:\n%s", out)
	}
	if !strings.Contains(out, "apt") {
		t.Errorf("expected output to mention 'apt', got:\n%s", out)
	}
}

func TestAppExplain_NotFound(t *testing.T) {
	buf := &strings.Builder{}
	app := newAppWithStore(buf, newMockStore(), aptManagers())
	err := app.explain("xyznonexistent")
	if err == nil {
		t.Error("expected error for unknown application")
	}
}

func TestAppExplain_EmptyAppID(t *testing.T) {
	buf := &strings.Builder{}
	app := newAppWithStore(buf, newMockStore(), aptManagers())
	if err := app.explain(""); err == nil {
		t.Error("expected error for empty appID")
	}
}

func TestAppExplain_JSONFound(t *testing.T) {
	buf := &strings.Builder{}
	app := newAppWithStore(buf, newMockStore(), aptManagers())
	if err := app.explainJSON("obs-studio"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out struct {
		ApplicationID         string `json:"application_id"`
		RecommendedSourceType string `json:"recommended_source_type"`
	}
	if err := json.Unmarshal([]byte(buf.String()), &out); err != nil {
		t.Fatalf("expected valid json output, got error: %v", err)
	}
	if out.ApplicationID != "obs-studio" {
		t.Fatalf("expected application_id obs-studio, got %q", out.ApplicationID)
	}
	if out.RecommendedSourceType == "" {
		t.Fatal("expected recommended source type to be set")
	}
}

func TestAppList_Empty(t *testing.T) {
	buf := &strings.Builder{}
	app := newAppWithStore(buf, newMockStore(), aptManagers())
	if err := app.list(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No applications") {
		t.Errorf("expected 'No applications' output, got:\n%s", buf.String())
	}
}

func TestAppList_WithInstallations(t *testing.T) {
	buf := &strings.Builder{}
	store := newMockStore()
	_ = store.Record(state.LocalInstallation{
		ApplicationID:      "git",
		SourceType:         source.TypeAPT,
		SourceIdentifier:   "git",
		InstallTimestamp:   time.Now(),
		InstallStatus:      state.StatusInstalled,
		VerificationStatus: state.VerificationPassed,
	})
	app := newAppWithStore(buf, store, aptManagers())
	if err := app.list(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "git") {
		t.Errorf("expected 'git' in list output, got:\n%s", buf.String())
	}
}

func TestAppList_HidesRemoved(t *testing.T) {
	buf := &strings.Builder{}
	store := newMockStore()
	_ = store.Record(state.LocalInstallation{
		ApplicationID:      "vlc",
		SourceType:         source.TypeAPT,
		SourceIdentifier:   "vlc",
		InstallTimestamp:   time.Now(),
		InstallStatus:      state.StatusRemoved,
		VerificationStatus: state.VerificationPassed,
	})
	app := newAppWithStore(buf, store, aptManagers())
	if err := app.list(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), "vlc") {
		t.Errorf("removed app should not appear in list, got:\n%s", buf.String())
	}
}

func TestAppList_JSONOutput(t *testing.T) {
	buf := &strings.Builder{}
	store := newMockStore()
	_ = store.Record(state.LocalInstallation{
		ApplicationID:      "git",
		SourceType:         source.TypeAPT,
		SourceIdentifier:   "git",
		InstallTimestamp:   time.Now(),
		InstallStatus:      state.StatusInstalled,
		VerificationStatus: state.VerificationPassed,
	})
	app := newAppWithStore(buf, store, aptManagers())
	if err := app.listJSON(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out struct {
		Applications []struct {
			ApplicationID string `json:"application_id"`
		} `json:"applications"`
	}
	if err := json.Unmarshal([]byte(buf.String()), &out); err != nil {
		t.Fatalf("expected valid json output, got error: %v", err)
	}
	if len(out.Applications) != 1 || out.Applications[0].ApplicationID != "git" {
		t.Fatalf("unexpected applications payload: %+v", out.Applications)
	}
}

func TestAppExplain_EnrichedTextOutputShowsTrustAndRisk(t *testing.T) {
buf := &strings.Builder{}
app := newAppWithStore(buf, newMockStore(), aptManagers())
if err := app.explain("obs-studio"); err != nil {
t.Fatalf("unexpected error: %v", err)
}
out := buf.String()
if !strings.Contains(out, "Trust:") {
t.Errorf("expected output to contain 'Trust:' label, got:\n%s", out)
}
if !strings.Contains(out, "Risk:") {
t.Errorf("expected output to contain 'Risk:' label, got:\n%s", out)
}
// obs-studio is an official APT source with low risk.
if !strings.Contains(out, "official") {
t.Errorf("expected trust level 'official' in output, got:\n%s", out)
}
if !strings.Contains(out, "low") {
t.Errorf("expected risk level 'low' in output, got:\n%s", out)
}
}

func TestAppExplain_EnrichedTextOutputShowsPrivilege(t *testing.T) {
buf := &strings.Builder{}
// APT requires privilege (sudo).
app := newAppWithStore(buf, newMockStore(), aptManagers())
if err := app.explain("obs-studio"); err != nil {
t.Fatalf("unexpected error: %v", err)
}
out := buf.String()
if !strings.Contains(out, "Privilege:") {
t.Errorf("expected output to contain 'Privilege:' label for APT, got:\n%s", out)
}
}

func TestAppExplain_EnrichedTextOutputShowsIdentifier(t *testing.T) {
buf := &strings.Builder{}
app := newAppWithStore(buf, newMockStore(), aptManagers())
if err := app.explain("obs-studio"); err != nil {
t.Fatalf("unexpected error: %v", err)
}
out := buf.String()
// The recommended line now shows "source_type (identifier)".
if !strings.Contains(out, "apt") {
t.Errorf("expected 'apt' in recommend line, got:\n%s", out)
}
if !strings.Contains(out, "obs-studio") {
t.Errorf("expected source identifier in recommend line, got:\n%s", out)
}
}

func TestAppExplain_JSONEnrichedOutput(t *testing.T) {
buf := &strings.Builder{}
app := newAppWithStore(buf, newMockStore(), aptManagers())
if err := app.explainJSON("obs-studio"); err != nil {
t.Fatalf("unexpected error: %v", err)
}
var out struct {
ApplicationID         string `json:"application_id"`
RecommendedSourceType string `json:"recommended_source_type"`
TrustLevel            string `json:"trust_level"`
RiskLevel             string `json:"risk_level"`
RequiresPrivilege     bool   `json:"requires_privilege"`
}
if err := json.Unmarshal([]byte(buf.String()), &out); err != nil {
t.Fatalf("expected valid JSON, got: %v\noutput: %s", err, buf.String())
}
if out.TrustLevel == "" {
t.Errorf("expected trust_level in JSON output, got empty\nfull output: %s", buf.String())
}
if out.RiskLevel == "" {
t.Errorf("expected risk_level in JSON output, got empty\nfull output: %s", buf.String())
}
// obs-studio via APT requires privilege.
if !out.RequiresPrivilege {
t.Errorf("expected requires_privilege=true for APT source\nfull output: %s", buf.String())
}
}

func TestAppExplain_ResolverExplanationContainsRiskAndTrust(t *testing.T) {
buf := &strings.Builder{}
app := newAppWithStore(buf, newMockStore(), aptManagers())
if err := app.explain("obs-studio"); err != nil {
t.Fatalf("unexpected error: %v", err)
}
out := buf.String()
// The enriched resolver explanation includes risk level wording.
if !strings.Contains(out, "risk") {
t.Errorf("expected resolver explanation to mention risk, got:\n%s", out)
}
}
