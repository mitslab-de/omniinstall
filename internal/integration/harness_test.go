// Package integration provides integration-style tests for the OmniInstall
// discovery→resolver→engine pipeline.
//
// Tests in this package exercise end-to-end MVP flow contracts using fake
// adapters and temporary in-memory state. All tests are isolated and
// deterministic.
package integration_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	adapter "github.com/mitslab-de/omniinstall/internal/adapters"
	"github.com/mitslab-de/omniinstall/internal/app"
	"github.com/mitslab-de/omniinstall/internal/discovery"
	"github.com/mitslab-de/omniinstall/internal/engine"
	"github.com/mitslab-de/omniinstall/internal/install"
	"github.com/mitslab-de/omniinstall/internal/resolver"
	"github.com/mitslab-de/omniinstall/internal/source"
	"github.com/mitslab-de/omniinstall/internal/state"
)

// ---------------------------------------------------------------------------
// Fake adapter
// ---------------------------------------------------------------------------

// fakeAdapter is a configurable in-memory adapter for integration tests.
type fakeAdapter struct {
	mu            sync.Mutex
	name          string
	available     bool
	handledTypes  []source.Type
	installResult *adapter.Result
	installErr    error
	removeResult  *adapter.Result
	removeErr     error
	checkState    adapter.InstalledState
	checkErr      error
	installedIDs  []string
	removedIDs    []string
}

func (f *fakeAdapter) Name() string { return f.name }

func (f *fakeAdapter) IsAvailable() bool { return f.available }

func (f *fakeAdapter) CanHandle(t source.Type) bool {
	for _, ht := range f.handledTypes {
		if ht == t {
			return true
		}
	}
	return false
}

func (f *fakeAdapter) CheckInstalled(id string) (adapter.InstalledState, error) {
	return f.checkState, f.checkErr
}

func (f *fakeAdapter) Install(plan *install.Plan) (*adapter.Result, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.installErr != nil {
		return nil, f.installErr
	}
	if f.installResult != nil && f.installResult.Success {
		f.installedIDs = append(f.installedIDs, plan.SourceIdentifier)
	}
	return f.installResult, nil
}

func (f *fakeAdapter) Remove(id string) (*adapter.Result, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.removeErr != nil {
		return nil, f.removeErr
	}
	f.removedIDs = append(f.removedIDs, id)
	return f.removeResult, nil
}

func (f *fakeAdapter) Verify(plan *install.Plan) (*adapter.VerificationResult, error) {
	return &adapter.VerificationResult{Verified: true, Details: "fake verify ok"}, nil
}

// ---------------------------------------------------------------------------
// In-memory state store
// ---------------------------------------------------------------------------

// memStore is a simple in-memory state.Store for integration tests.
type memStore struct {
	mu    sync.RWMutex
	items map[string]state.LocalInstallation
}

func newMemStore() *memStore {
	return &memStore{items: make(map[string]state.LocalInstallation)}
}

func (m *memStore) Record(inst state.LocalInstallation) error {
	if err := inst.Validate(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[inst.ApplicationID] = inst
	return nil
}

func (m *memStore) Get(applicationID string) (state.LocalInstallation, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.items[applicationID]
	return r, ok, nil
}

func (m *memStore) List() ([]state.LocalInstallation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]state.LocalInstallation, 0, len(m.items))
	for _, v := range m.items {
		result = append(result, v)
	}
	return result, nil
}

func (m *memStore) MarkRemoved(applicationID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	r, ok := m.items[applicationID]
	if !ok {
		return state.ErrNotFound
	}
	r.InstallStatus = state.StatusRemoved
	m.items[applicationID] = r
	return nil
}

// ---------------------------------------------------------------------------
// Test fixtures
// ---------------------------------------------------------------------------

// testApp returns a minimal valid Application fixture.
func testApp(id, name string) *app.Application {
	return &app.Application{
		ID:          id,
		DisplayName: name,
		Summary:     name + " – integration test fixture",
		Categories:  []string{"utilities"},
	}
}

// aptSource returns a low-risk APT source fixture for the given application.
func aptSource(appID, pkg string) source.Source {
	return source.Source{
		ApplicationID:    appID,
		SourceType:       source.TypeAPT,
		SourceIdentifier: pkg,
		TrustLevel:       source.TrustOfficial,
		RiskLevel:        source.RiskLow,
	}
}

// flatpakSource returns a low-risk Flatpak source fixture.
func flatpakSource(appID, ref string) source.Source {
	return source.Source{
		ApplicationID:    appID,
		SourceType:       source.TypeFlatpak,
		SourceIdentifier: ref,
		TrustLevel:       source.TrustVerified,
		RiskLevel:        source.RiskLow,
	}
}

// successAptAdapter returns a fakeAdapter that succeeds for APT operations.
func successAptAdapter() *fakeAdapter {
	return &fakeAdapter{
		name:         "fake-apt",
		available:    true,
		handledTypes: []source.Type{source.TypeAPT},
		installResult: &adapter.Result{
			Success: true,
			Message: "installed via fake-apt",
		},
		removeResult: &adapter.Result{
			Success: true,
			Message: "removed via fake-apt",
		},
		checkState: adapter.StateNotInstalled,
	}
}

// successFlatpakAdapter returns a fakeAdapter that succeeds for Flatpak operations.
func successFlatpakAdapter() *fakeAdapter {
	return &fakeAdapter{
		name:         "fake-flatpak",
		available:    true,
		handledTypes: []source.Type{source.TypeFlatpak},
		installResult: &adapter.Result{
			Success: true,
			Message: "installed via fake-flatpak",
		},
		removeResult: &adapter.Result{
			Success: true,
			Message: "removed via fake-flatpak",
		},
		checkState: adapter.StateNotInstalled,
	}
}

// sysCtxWithAPT returns a SystemContext that has APT available.
func sysCtxWithAPT() resolver.SystemContext {
	return resolver.SystemContext{
		Distribution:      "ubuntu",
		Architecture:      "amd64",
		AvailableManagers: []source.Type{source.TypeAPT},
	}
}

// sysCtxWithBoth returns a SystemContext with both APT and Flatpak available.
func sysCtxWithBoth() resolver.SystemContext {
	return resolver.SystemContext{
		Distribution:      "ubuntu",
		Architecture:      "amd64",
		AvailableManagers: []source.Type{source.TypeAPT, source.TypeFlatpak},
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// buildPipeline sets up a catalog, engine and resolver for integration tests.
func buildPipeline(t *testing.T, apps []*app.Application, adapters []adapter.Adapter) (
	*discovery.LocalEngine,
	*resolver.DefaultResolver,
	*engine.DefaultEngine,
) {
	t.Helper()
	cat := discovery.NewCatalog()
	for _, a := range apps {
		if err := cat.Add(a); err != nil {
			t.Fatalf("catalog.Add(%s): %v", a.ID, err)
		}
	}
	de := discovery.NewLocalEngine(cat)
	res := resolver.NewDefaultResolver()
	eng := engine.NewDefaultEngine(adapters, nil)
	return de, res, eng
}

// ---------------------------------------------------------------------------
// Integration tests: install flow
// ---------------------------------------------------------------------------

// TestIntegration_DiscoverResolveInstall_APT exercises the full
// discovery→resolver→engine install path for a single APT application.
func TestIntegration_DiscoverResolveInstall_APT(t *testing.T) {
	obsApp := testApp("obs-studio", "OBS Studio")
	obsSrc := aptSource("obs-studio", "obs-studio")

	aptAdapter := successAptAdapter()
	de, res, eng := buildPipeline(t, []*app.Application{obsApp}, []adapter.Adapter{aptAdapter})
	st := newMemStore()

	// 1. Discover
	candidates, err := de.Search("obs")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(candidates) == 0 {
		t.Fatal("expected at least one candidate")
	}
	found := candidates[0].Application
	if found.ID != "obs-studio" {
		t.Fatalf("expected obs-studio, got %s", found.ID)
	}

	// 2. Resolve
	recs, err := res.Resolve(found.ID, []source.Source{obsSrc}, sysCtxWithAPT())
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if len(recs) == 0 {
		t.Fatal("expected at least one recommendation")
	}
	plan := recs[0].Plan
	if plan.SourceType != source.TypeAPT {
		t.Fatalf("expected APT plan, got %s", plan.SourceType)
	}

	// 3. Install
	result, err := eng.Install(plan)
	if err != nil {
		t.Fatalf("Install: %v", err)
	}
	if !result.Success {
		t.Fatalf("expected install success, got: %s", result.Message)
	}
	if result.ApplicationID != "obs-studio" {
		t.Fatalf("expected application_id=obs-studio, got %s", result.ApplicationID)
	}

	// 4. Record state
	if err := st.Record(state.LocalInstallation{
		ApplicationID:      found.ID,
		SourceType:         plan.SourceType,
		SourceIdentifier:   plan.SourceIdentifier,
		InstallTimestamp:   time.Now(),
		InstallStatus:      state.StatusInstalled,
		VerificationStatus: state.VerificationPassed,
	}); err != nil {
		t.Fatalf("Record: %v", err)
	}

	rec, ok, err := st.Get("obs-studio")
	if err != nil || !ok {
		t.Fatalf("expected state record to exist: err=%v ok=%v", err, ok)
	}
	if rec.InstallStatus != state.StatusInstalled {
		t.Fatalf("expected installed status, got %s", rec.InstallStatus)
	}
}

// TestIntegration_DiscoverResolveInstall_Flatpak exercises the full
// discovery→resolver→engine install path for a Flatpak application.
func TestIntegration_DiscoverResolveInstall_Flatpak(t *testing.T) {
	vlcApp := testApp("vlc", "VLC Media Player")
	vlcSrc := flatpakSource("vlc", "org.videolan.VLC")

	flatpakAdapter := successFlatpakAdapter()
	de, res, eng := buildPipeline(t, []*app.Application{vlcApp}, []adapter.Adapter{flatpakAdapter})

	// 1. Discover
	candidates, err := de.Search("vlc")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(candidates) == 0 {
		t.Fatal("expected at least one candidate")
	}

	// 2. Resolve (only Flatpak available)
	ctx := resolver.SystemContext{
		Distribution:      "ubuntu",
		Architecture:      "amd64",
		AvailableManagers: []source.Type{source.TypeFlatpak},
	}
	recs, err := res.Resolve("vlc", []source.Source{vlcSrc}, ctx)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if len(recs) == 0 {
		t.Fatal("expected at least one recommendation")
	}
	plan := recs[0].Plan
	if plan.SourceType != source.TypeFlatpak {
		t.Fatalf("expected Flatpak plan, got %s", plan.SourceType)
	}

	// 3. Install
	result, err := eng.Install(plan)
	if err != nil {
		t.Fatalf("Install: %v", err)
	}
	if !result.Success {
		t.Fatalf("expected install success, got: %s", result.Message)
	}
}

// ---------------------------------------------------------------------------
// Integration test: install then remove lifecycle
// ---------------------------------------------------------------------------

// TestIntegration_InstallThenRemoveLifecycle verifies the full install→state
// record→remove→state updated lifecycle.
func TestIntegration_InstallThenRemoveLifecycle(t *testing.T) {
	gimpApp := testApp("gimp", "GIMP")
	gimpSrc := aptSource("gimp", "gimp")

	aptAdapter := successAptAdapter()
	de, res, eng := buildPipeline(t, []*app.Application{gimpApp}, []adapter.Adapter{aptAdapter})
	st := newMemStore()

	// --- Install phase ---
	candidates, err := de.Search("gimp")
	if err != nil || len(candidates) == 0 {
		t.Fatalf("Search failed: err=%v candidates=%d", err, len(candidates))
	}

	recs, err := res.Resolve("gimp", []source.Source{gimpSrc}, sysCtxWithAPT())
	if err != nil || len(recs) == 0 {
		t.Fatalf("Resolve failed: err=%v recs=%d", err, len(recs))
	}
	plan := recs[0].Plan

	installResult, err := eng.Install(plan)
	if err != nil {
		t.Fatalf("Install: %v", err)
	}
	if !installResult.Success {
		t.Fatalf("expected install success, got: %s", installResult.Message)
	}

	if err := st.Record(state.LocalInstallation{
		ApplicationID:      "gimp",
		SourceType:         plan.SourceType,
		SourceIdentifier:   plan.SourceIdentifier,
		InstallTimestamp:   time.Now(),
		InstallStatus:      state.StatusInstalled,
		VerificationStatus: state.VerificationPassed,
	}); err != nil {
		t.Fatalf("Record: %v", err)
	}

	// Confirm installed
	rec, ok, _ := st.Get("gimp")
	if !ok || rec.InstallStatus != state.StatusInstalled {
		t.Fatal("expected gimp to be marked installed in state")
	}

	// --- Remove phase ---
	removeResult, err := eng.Remove("gimp", rec.SourceType, rec.SourceIdentifier)
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if !removeResult.Success {
		t.Fatalf("expected remove success, got: %s", removeResult.Message)
	}

	if err := st.MarkRemoved("gimp"); err != nil {
		t.Fatalf("MarkRemoved: %v", err)
	}

	// Confirm removed
	rec, ok, _ = st.Get("gimp")
	if !ok || rec.InstallStatus != state.StatusRemoved {
		t.Fatalf("expected gimp to be marked removed, got ok=%v status=%s", ok, rec.InstallStatus)
	}

	// Adapter should have been called with the correct source identifier
	if len(aptAdapter.removedIDs) == 0 || aptAdapter.removedIDs[0] != "gimp" {
		t.Fatalf("expected adapter remove called with 'gimp', got %v", aptAdapter.removedIDs)
	}
}

// ---------------------------------------------------------------------------
// Integration tests: source selection with multiple sources
// ---------------------------------------------------------------------------

// TestIntegration_ResolverSelectsPreferredSource verifies that the resolver
// selects APT over Flatpak when both are available and no preference is set.
func TestIntegration_ResolverSelectsPreferredSource(t *testing.T) {
	inkscapeApp := testApp("inkscape", "Inkscape")
	aptSrc := aptSource("inkscape", "inkscape")
	fpSrc := flatpakSource("inkscape", "org.inkscape.Inkscape")

	aptAdapter := successAptAdapter()
	fpAdapter := successFlatpakAdapter()
	_, res, eng := buildPipeline(t,
		[]*app.Application{inkscapeApp},
		[]adapter.Adapter{aptAdapter, fpAdapter},
	)

	recs, err := res.Resolve("inkscape", []source.Source{aptSrc, fpSrc}, sysCtxWithBoth())
	if err != nil || len(recs) == 0 {
		t.Fatalf("Resolve: err=%v recs=%d", err, len(recs))
	}

	// Install with top recommendation
	result, err := eng.Install(recs[0].Plan)
	if err != nil {
		t.Fatalf("Install: %v", err)
	}
	if !result.Success {
		t.Fatalf("expected success, got: %s", result.Message)
	}

	// APT is natively ranked higher, so APT adapter should have been used
	if len(aptAdapter.installedIDs) == 0 {
		t.Error("expected APT adapter to be used for installation")
	}
	if len(fpAdapter.installedIDs) != 0 {
		t.Error("expected Flatpak adapter NOT to be used for top recommendation")
	}
}

// TestIntegration_ResolverPrefersUserFlatpakPreference verifies that user
// preference for Flatpak can push it above a lower-trust APT source.
// APT (community trust) scores 50+5+10=65; Flatpak (verified) + preference
// scores 40+15+10+15=80, so Flatpak wins.
func TestIntegration_ResolverPrefersUserFlatpakPreference(t *testing.T) {
	ffApp := testApp("firefox", "Firefox")
	// Use community-trust APT so Flatpak+preference outscores it.
	aptSrc := source.Source{
		ApplicationID:    "firefox",
		SourceType:       source.TypeAPT,
		SourceIdentifier: "firefox",
		TrustLevel:       source.TrustCommunity,
		RiskLevel:        source.RiskLow,
	}
	fpSrc := flatpakSource("firefox", "org.mozilla.firefox")

	fpAdapter := successFlatpakAdapter()
	aptAdapter := successAptAdapter()
	_, res, eng := buildPipeline(t,
		[]*app.Application{ffApp},
		[]adapter.Adapter{aptAdapter, fpAdapter},
	)

	ctx := resolver.SystemContext{
		Distribution:      "ubuntu",
		Architecture:      "amd64",
		AvailableManagers: []source.Type{source.TypeAPT, source.TypeFlatpak},
		UserPreferences:   resolver.UserPreferences{PreferFlatpak: true},
	}

	recs, err := res.Resolve("firefox", []source.Source{aptSrc, fpSrc}, ctx)
	if err != nil || len(recs) == 0 {
		t.Fatalf("Resolve: err=%v recs=%d", err, len(recs))
	}

	topPlan := recs[0].Plan
	if topPlan.SourceType != source.TypeFlatpak {
		t.Fatalf("expected Flatpak as top recommendation when user prefers Flatpak, got %s", topPlan.SourceType)
	}

	result, err := eng.Install(topPlan)
	if err != nil {
		t.Fatalf("Install: %v", err)
	}
	if !result.Success {
		t.Fatalf("expected install success, got: %s", result.Message)
	}
	if len(fpAdapter.installedIDs) == 0 {
		t.Error("expected Flatpak adapter to be used")
	}
}

// ---------------------------------------------------------------------------
// Integration tests: failure scenarios
// ---------------------------------------------------------------------------

// TestIntegration_NoAdapterAvailable verifies that Install fails gracefully
// when no adapter supports the resolved source type.
func TestIntegration_NoAdapterAvailable(t *testing.T) {
	keeApp := testApp("keepassxc", "KeePassXC")
	keepSrc := aptSource("keepassxc", "keepassxc")

	// Engine has no adapters configured
	_, res, eng := buildPipeline(t, []*app.Application{keeApp}, nil)

	recs, err := res.Resolve("keepassxc", []source.Source{keepSrc}, sysCtxWithAPT())
	if err != nil || len(recs) == 0 {
		t.Fatalf("Resolve: err=%v recs=%d", err, len(recs))
	}

	result, err := eng.Install(recs[0].Plan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Fatal("expected failure when no adapter is available")
	}
	if result.ErrorCategory != adapter.ErrBackendUnavailable {
		t.Fatalf("expected %s, got %s", adapter.ErrBackendUnavailable, result.ErrorCategory)
	}
}

// TestIntegration_AdapterInstallError verifies that an adapter execution error
// is surfaced as a failed engine result.
func TestIntegration_AdapterInstallError(t *testing.T) {
	codeApp := testApp("vscode", "VS Code")
	codeSrc := aptSource("vscode", "code")

	brokenAdapter := successAptAdapter()
	brokenAdapter.installErr = errors.New("dpkg lock held by another process")

	_, res, eng := buildPipeline(t, []*app.Application{codeApp}, []adapter.Adapter{brokenAdapter})

	recs, err := res.Resolve("vscode", []source.Source{codeSrc}, sysCtxWithAPT())
	if err != nil || len(recs) == 0 {
		t.Fatalf("Resolve: err=%v recs=%d", err, len(recs))
	}

	result, err := eng.Install(recs[0].Plan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Fatal("expected failure when adapter returns an error")
	}
	if result.ErrorCategory != adapter.ErrExecutionFailed {
		t.Fatalf("expected %s, got %s", adapter.ErrExecutionFailed, result.ErrorCategory)
	}
}

// TestIntegration_AdapterReportsFailure verifies that a non-success adapter
// result is propagated correctly through the engine.
func TestIntegration_AdapterReportsFailure(t *testing.T) {
	htopApp := testApp("htop", "htop")
	htopSrc := aptSource("htop", "htop")

	failingAdapter := successAptAdapter()
	failingAdapter.installResult = &adapter.Result{
		Success:       false,
		Message:       "package htop not found in repository",
		ErrorCategory: adapter.ErrPackageNotFound,
	}

	_, res, eng := buildPipeline(t, []*app.Application{htopApp}, []adapter.Adapter{failingAdapter})

	recs, err := res.Resolve("htop", []source.Source{htopSrc}, sysCtxWithAPT())
	if err != nil || len(recs) == 0 {
		t.Fatalf("Resolve: err=%v recs=%d", err, len(recs))
	}

	result, err := eng.Install(recs[0].Plan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Fatal("expected failure when adapter reports package not found")
	}
	if result.ErrorCategory != adapter.ErrPackageNotFound {
		t.Fatalf("expected %s, got %s", adapter.ErrPackageNotFound, result.ErrorCategory)
	}
}

// TestIntegration_PreflightFailurePropagated verifies that a backend
// unavailability detected during preflight is surfaced as a failed install.
func TestIntegration_PreflightFailurePropagated(t *testing.T) {
	nanoApp := testApp("nano", "nano")
	nanoSrc := aptSource("nano", "nano")

	unavailAdapter := successAptAdapter()
	unavailAdapter.checkErr = errors.New("apt binary not found")

	_, res, eng := buildPipeline(t, []*app.Application{nanoApp}, []adapter.Adapter{unavailAdapter})

	recs, err := res.Resolve("nano", []source.Source{nanoSrc}, sysCtxWithAPT())
	if err != nil || len(recs) == 0 {
		t.Fatalf("Resolve: err=%v recs=%d", err, len(recs))
	}

	result, err := eng.Install(recs[0].Plan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Fatal("expected preflight failure to propagate as install failure")
	}
	if result.ErrorCategory != adapter.ErrBackendUnavailable {
		t.Fatalf("expected %s, got %s", adapter.ErrBackendUnavailable, result.ErrorCategory)
	}
}

// TestIntegration_NoSourcesForApp verifies that Resolve returns an empty slice
// (no error) when the application has no available sources.
func TestIntegration_NoSourcesForApp(t *testing.T) {
	testAppObj := testApp("unknown-app", "Unknown App")
	_, res, _ := buildPipeline(t, []*app.Application{testAppObj}, nil)

	recs, err := res.Resolve("unknown-app", nil, sysCtxWithAPT())
	if err != nil {
		t.Fatalf("Resolve with no sources should not error, got: %v", err)
	}
	if len(recs) != 0 {
		t.Fatalf("expected empty recommendations, got %d", len(recs))
	}
}

// TestIntegration_LookupMiss verifies that Lookup returns nil for an app not
// in the catalog without an error.
func TestIntegration_LookupMiss(t *testing.T) {
	de, _, _ := buildPipeline(t, nil, nil)

	a, err := de.Lookup("does-not-exist")
	if err != nil {
		t.Fatalf("Lookup should not error on miss, got: %v", err)
	}
	if a != nil {
		t.Fatalf("expected nil for unknown app, got %+v", a)
	}
}

// ---------------------------------------------------------------------------
// Integration test: progress events across full pipeline
// ---------------------------------------------------------------------------

// TestIntegration_ProgressEventsEmittedDuringInstall verifies that all
// expected progress events are emitted during a full install flow.
func TestIntegration_ProgressEventsEmittedDuringInstall(t *testing.T) {
	tApp := testApp("tree", "tree")
	tSrc := aptSource("tree", "tree")

	var events []engine.ProgressEvent
	handler := func(ev engine.ProgressEvent) { events = append(events, ev) }

	cat := discovery.NewCatalog()
	if err := cat.Add(tApp); err != nil {
		t.Fatalf("catalog.Add: %v", err)
	}
	de := discovery.NewLocalEngine(cat)
	res := resolver.NewDefaultResolver()
	eng := engine.NewDefaultEngine([]adapter.Adapter{successAptAdapter()}, handler)

	candidates, err := de.Search("tree")
	if err != nil || len(candidates) == 0 {
		t.Fatalf("Search: err=%v candidates=%d", err, len(candidates))
	}

	recs, err := res.Resolve("tree", []source.Source{tSrc}, sysCtxWithAPT())
	if err != nil || len(recs) == 0 {
		t.Fatalf("Resolve: err=%v recs=%d", err, len(recs))
	}

	if _, err := eng.Install(recs[0].Plan); err != nil {
		t.Fatalf("Install: %v", err)
	}

	// Verify that at minimum started, executing, and completed were emitted.
	kindSet := make(map[engine.EventKind]bool)
	for _, ev := range events {
		kindSet[ev.Kind] = true
	}
	for _, required := range []engine.EventKind{
		engine.EventStarted,
		engine.EventExecuting,
		engine.EventCompleted,
	} {
		if !kindSet[required] {
			t.Errorf("expected progress event %s to be emitted", required)
		}
	}
}
