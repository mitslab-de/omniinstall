package engine_test

import (
	"errors"
	"testing"

	adapter "github.com/mitslab-de/omniinstall/internal/adapters"
	"github.com/mitslab-de/omniinstall/internal/engine"
	"github.com/mitslab-de/omniinstall/internal/install"
	"github.com/mitslab-de/omniinstall/internal/logging"
	"github.com/mitslab-de/omniinstall/internal/source"
)

// mockAdapter is a configurable fake adapter for testing.
type mockAdapter struct {
	name          string
	available     bool
	handledTypes  []source.Type
	installResult *adapter.Result
	installErr    error
	removeResult  *adapter.Result
	removeErr     error
	verifyResult  *adapter.Result
	verifyErr     error
	checkState    adapter.InstalledState
	checkErr      error
	removedIDs    []string
}

func (m *mockAdapter) Name() string { return m.name }

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
	return m.checkState, m.checkErr
}

func (m *mockAdapter) Install(plan *install.Plan) (*adapter.Result, error) {
	return m.installResult, m.installErr
}

func (m *mockAdapter) Remove(id string) (*adapter.Result, error) {
	m.removedIDs = append(m.removedIDs, id)
	return m.removeResult, m.removeErr
}

func (m *mockAdapter) Verify(plan *install.Plan) (*adapter.VerificationResult, error) {
	if m.verifyErr != nil {
		return nil, m.verifyErr
	}
	if m.verifyResult != nil {
		return &adapter.VerificationResult{Verified: m.verifyResult.Success, Details: m.verifyResult.Message}, nil
	}
	return &adapter.VerificationResult{Verified: true, Details: "ok"}, nil
}

// goodPlan is a valid install plan for use in tests.
var goodPlan = &install.Plan{
	ApplicationID:    "test-app",
	SourceType:       source.TypeAPT,
	SourceIdentifier: "test-app",
	RiskLevel:        source.RiskLow,
	Verification:     []install.VerificationRule{{Command: "test-app"}},
}

// successAdapter returns a mockAdapter that succeeds on install/remove/verify.
func successAdapter() *mockAdapter {
	return &mockAdapter{
		name:         "mock-apt",
		available:    true,
		handledTypes: []source.Type{source.TypeAPT},
		installResult: &adapter.Result{
			Success: true,
			Message: "installed ok",
		},
		removeResult: &adapter.Result{
			Success: true,
			Message: "removed ok",
		},
	}
}

func TestDefaultEngineInstallSuccess(t *testing.T) {
	e := engine.NewDefaultEngine([]adapter.Adapter{successAdapter()}, nil)
	result, err := e.Install(goodPlan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected success, got message: %s", result.Message)
	}
	if result.ApplicationID != "test-app" {
		t.Errorf("expected application_id=test-app, got %s", result.ApplicationID)
	}
}

func TestDefaultEngineInstallNilPlan(t *testing.T) {
	e := engine.NewDefaultEngine([]adapter.Adapter{successAdapter()}, nil)
	_, err := e.Install(nil)
	if err == nil {
		t.Error("expected error for nil plan")
	}
}

func TestDefaultEngineInstallInvalidPlan(t *testing.T) {
	e := engine.NewDefaultEngine([]adapter.Adapter{successAdapter()}, nil)
	badPlan := &install.Plan{} // missing required fields
	result, err := e.Install(badPlan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected failure for invalid plan")
	}
	if result.ErrorCategory != "invalid_plan" {
		t.Errorf("expected error_category=invalid_plan, got %s", result.ErrorCategory)
	}
}

func TestDefaultEngineInstallNoAdapter(t *testing.T) {
	e := engine.NewDefaultEngine(nil, nil)
	result, err := e.Install(goodPlan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected failure when no adapters available")
	}
	if result.ErrorCategory != adapter.ErrBackendUnavailable {
		t.Errorf("expected error_category=%s, got %s", adapter.ErrBackendUnavailable, result.ErrorCategory)
	}
}

func TestDefaultEngineInstallAdapterError(t *testing.T) {
	a := successAdapter()
	a.installResult = &adapter.Result{Success: false, Message: "apt failed", ErrorCategory: adapter.ErrExecutionFailed}
	e := engine.NewDefaultEngine([]adapter.Adapter{a}, nil)
	result, err := e.Install(goodPlan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected failure when adapter reports failure")
	}
}

func TestDefaultEngineInstallPreflightBackendUnavailable(t *testing.T) {
	a := successAdapter()
	a.checkErr = errors.New("backend probe failed")
	e := engine.NewDefaultEngine([]adapter.Adapter{a}, nil)
	result, err := e.Install(goodPlan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Fatal("expected preflight failure")
	}
	if result.ErrorCategory != adapter.ErrBackendUnavailable {
		t.Fatalf("expected %s, got %s", adapter.ErrBackendUnavailable, result.ErrorCategory)
	}
}

func TestDefaultEngineInstallPreflightPermissionDenied(t *testing.T) {
	a := successAdapter()
	a.checkErr = errors.New("permission denied while checking backend")
	e := engine.NewDefaultEngine([]adapter.Adapter{a}, nil)
	result, err := e.Install(goodPlan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Fatal("expected preflight failure")
	}
	if result.ErrorCategory != adapter.ErrPermissionDenied {
		t.Fatalf("expected %s, got %s", adapter.ErrPermissionDenied, result.ErrorCategory)
	}
}

func TestDefaultEngineInstallVerificationHandoff(t *testing.T) {
	var events []engine.ProgressEvent
	handler := func(ev engine.ProgressEvent) { events = append(events, ev) }
	e := engine.NewDefaultEngine([]adapter.Adapter{successAdapter()}, handler)
	_, err := e.Install(goodPlan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var foundVerify bool
	for _, ev := range events {
		if ev.Kind == engine.EventVerifying {
			foundVerify = true
		}
	}
	if !foundVerify {
		t.Error("expected EventVerifying to be emitted")
	}
}

func TestDefaultEngineProgressEvents(t *testing.T) {
	var events []engine.ProgressEvent
	handler := func(ev engine.ProgressEvent) { events = append(events, ev) }
	e := engine.NewDefaultEngine([]adapter.Adapter{successAdapter()}, handler)
	if _, err := e.Install(goodPlan); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantKinds := []engine.EventKind{
		engine.EventStarted,
		engine.EventValidating,
		engine.EventSelectingAdapter,
		engine.EventValidating,
		engine.EventExecuting,
		engine.EventVerifying,
		engine.EventCompleted,
	}
	if len(events) < len(wantKinds) {
		t.Fatalf("expected ≥%d events, got %d", len(wantKinds), len(events))
	}
	for i, kind := range wantKinds {
		if events[i].Kind != kind {
			t.Errorf("event[%d]: expected %s, got %s", i, kind, events[i].Kind)
		}
		if events[i].ApplicationID != "test-app" {
			t.Errorf("event[%d]: expected applicationID=test-app, got %s", i, events[i].ApplicationID)
		}
	}
}

func TestDefaultEngineRemoveSuccess(t *testing.T) {
	e := engine.NewDefaultEngine([]adapter.Adapter{successAdapter()}, nil)
	result, err := e.Remove("test-app", source.TypeAPT, "test-app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected success, got: %s", result.Message)
	}
	if result.ApplicationID != "test-app" {
		t.Errorf("expected application_id=test-app, got %s", result.ApplicationID)
	}
}

func TestDefaultEngineRemoveEmptyID(t *testing.T) {
	e := engine.NewDefaultEngine([]adapter.Adapter{successAdapter()}, nil)
	_, err := e.Remove("", source.TypeAPT, "test-app")
	if err == nil {
		t.Error("expected error for empty applicationID")
	}
}

func TestDefaultEngineRemoveNoAdapter(t *testing.T) {
	e := engine.NewDefaultEngine(nil, nil)
	result, err := e.Remove("test-app", source.TypeAPT, "test-app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected failure when no adapters available")
	}
	if result.ErrorCategory != adapter.ErrBackendUnavailable {
		t.Errorf("expected error_category=%s, got %s", adapter.ErrBackendUnavailable, result.ErrorCategory)
	}
}

func TestDefaultEngineRemovePreflightFailure(t *testing.T) {
	a := successAdapter()
	a.checkErr = errors.New("backend readiness check failed")
	e := engine.NewDefaultEngine([]adapter.Adapter{a}, nil)
	result, err := e.Remove("test-app", source.TypeAPT, "test-app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Fatal("expected preflight failure")
	}
	if result.ErrorCategory != adapter.ErrBackendUnavailable {
		t.Fatalf("expected %s, got %s", adapter.ErrBackendUnavailable, result.ErrorCategory)
	}
}

func TestDefaultEngineRemoveUsesSourceIdentifierAndSourceType(t *testing.T) {
	apt := successAdapter()
	apt.handledTypes = []source.Type{source.TypeAPT}
	flatpak := successAdapter()
	flatpak.name = "flatpak"
	flatpak.handledTypes = []source.Type{source.TypeFlatpak}

	e := engine.NewDefaultEngine([]adapter.Adapter{apt, flatpak}, nil)
	result, err := e.Remove("obs-studio", source.TypeFlatpak, "com.obsproject.Studio")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Fatalf("expected success, got %s", result.Message)
	}
	if len(apt.removedIDs) != 0 {
		t.Fatalf("expected apt adapter not to be used, got removals: %v", apt.removedIDs)
	}
	if len(flatpak.removedIDs) != 1 || flatpak.removedIDs[0] != "com.obsproject.Studio" {
		t.Fatalf("expected flatpak removal with source identifier, got %v", flatpak.removedIDs)
	}
}

func TestDefaultEngineSkipsUnavailableAdapter(t *testing.T) {
	unavail := successAdapter()
	unavail.available = false

	avail := successAdapter()
	avail.name = "second-adapter"

	e := engine.NewDefaultEngine([]adapter.Adapter{unavail, avail}, nil)
	result, err := e.Install(goodPlan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected second adapter to be used, got: %s", result.Message)
	}
}

func TestDefaultEngineImplementsInterface(t *testing.T) {
	var _ engine.Engine = engine.NewDefaultEngine(nil, nil)
}

func TestProgressEventKindConstants(t *testing.T) {
	constants := []engine.EventKind{
		engine.EventStarted,
		engine.EventValidating,
		engine.EventSelectingAdapter,
		engine.EventExecuting,
		engine.EventVerifying,
		engine.EventCompleted,
		engine.EventFailed,
	}
	for _, c := range constants {
		if c == "" {
			t.Error("EventKind constant must not be empty")
		}
	}
}

func TestDefaultEngineInstallEmitsStructuredInstallAndVerifyLogs(t *testing.T) {
	e := engine.NewDefaultEngine([]adapter.Adapter{successAdapter()}, nil)
	var entries []logging.Entry
	e.WithLogger(logging.EmitFunc(func(entry logging.Entry) {
		entries = append(entries, entry)
	}))

	if _, err := e.Install(goodPlan); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 log entries (verify + install), got %d", len(entries))
	}
	if entries[0].Action != logging.ActionVerify {
		t.Fatalf("expected first action verify, got %q", entries[0].Action)
	}
	if entries[0].VerificationStatus != "verified" {
		t.Fatalf("expected verification status verified, got %q", entries[0].VerificationStatus)
	}
	if entries[1].Action != logging.ActionInstall {
		t.Fatalf("expected second action install, got %q", entries[1].Action)
	}
	if entries[1].Result != "success" {
		t.Fatalf("expected install success result, got %q", entries[1].Result)
	}
	if entries[1].Duration <= 0 {
		t.Fatalf("expected positive duration, got %v", entries[1].Duration)
	}
}

func TestDefaultEngineInstallEmitsFailureErrorCategory(t *testing.T) {
	e := engine.NewDefaultEngine([]adapter.Adapter{successAdapter()}, nil)
	var entries []logging.Entry
	e.WithLogger(logging.EmitFunc(func(entry logging.Entry) {
		entries = append(entries, entry)
	}))

	if _, err := e.Install(&install.Plan{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(entries))
	}
	if entries[0].Action != logging.ActionInstall {
		t.Fatalf("expected install action, got %q", entries[0].Action)
	}
	if entries[0].Result != "failure" {
		t.Fatalf("expected failure result, got %q", entries[0].Result)
	}
	if entries[0].ErrorCategory != "invalid_plan" {
		t.Fatalf("expected invalid_plan, got %q", entries[0].ErrorCategory)
	}
}

func TestDefaultEngineRemoveEmitsFailureErrorCategory(t *testing.T) {
	a := successAdapter()
	a.checkErr = errors.New("permission denied while checking backend")
	e := engine.NewDefaultEngine([]adapter.Adapter{a}, nil)
	var entries []logging.Entry
	e.WithLogger(logging.EmitFunc(func(entry logging.Entry) {
		entries = append(entries, entry)
	}))

	if _, err := e.Remove("test-app", source.TypeAPT, "test-app"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(entries))
	}
	if entries[0].Action != logging.ActionRemove {
		t.Fatalf("expected remove action, got %q", entries[0].Action)
	}
	if entries[0].Result != "failure" {
		t.Fatalf("expected failure result, got %q", entries[0].Result)
	}
	if entries[0].ErrorCategory != adapter.ErrPermissionDenied {
		t.Fatalf("expected %q, got %q", adapter.ErrPermissionDenied, entries[0].ErrorCategory)
	}
}
