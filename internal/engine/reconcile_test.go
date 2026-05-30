package engine_test

import (
	"errors"
	"testing"
	"time"

	adapter "github.com/mitslab-de/omniinstall/internal/adapters"
	"github.com/mitslab-de/omniinstall/internal/engine"
	"github.com/mitslab-de/omniinstall/internal/install"
	"github.com/mitslab-de/omniinstall/internal/source"
	"github.com/mitslab-de/omniinstall/internal/state"
)

// reconcileAdapter is a test adapter for reconciliation tests.
type reconcileAdapter struct {
	name          string
	available     bool
	handledTypes  []source.Type
	checkState    adapter.InstalledState
	checkErr      error
}

func (a *reconcileAdapter) Name() string      { return a.name }
func (a *reconcileAdapter) IsAvailable() bool { return a.available }
func (a *reconcileAdapter) CanHandle(t source.Type) bool {
	for _, ht := range a.handledTypes {
		if ht == t {
			return true
		}
	}
	return false
}
func (a *reconcileAdapter) CheckInstalled(_ string) (adapter.InstalledState, error) {
	return a.checkState, a.checkErr
}
func (a *reconcileAdapter) Install(_ *install.Plan) (*adapter.Result, error) {
	return &adapter.Result{Success: true}, nil
}
func (a *reconcileAdapter) Remove(_ string) (*adapter.Result, error) {
	return &adapter.Result{Success: true}, nil
}
func (a *reconcileAdapter) Verify(_ *install.Plan) (*adapter.VerificationResult, error) {
	return &adapter.VerificationResult{Verified: true}, nil
}

// helpers

func installedRec(id string, st source.Type, si string) state.LocalInstallation {
	return state.LocalInstallation{
		ApplicationID:      id,
		SourceType:         st,
		SourceIdentifier:   si,
		InstallTimestamp:   time.Now(),
		InstallStatus:      state.StatusInstalled,
		VerificationStatus: state.VerificationPassed,
	}
}

func removedRec(id string, st source.Type, si string) state.LocalInstallation {
	return state.LocalInstallation{
		ApplicationID:      id,
		SourceType:         st,
		SourceIdentifier:   si,
		InstallTimestamp:   time.Now(),
		InstallStatus:      state.StatusRemoved,
		VerificationStatus: state.VerificationSkipped,
	}
}

func aptAdapter(state adapter.InstalledState) *reconcileAdapter {
	return &reconcileAdapter{
		name:         "apt",
		available:    true,
		handledTypes: []source.Type{source.TypeAPT},
		checkState:   state,
	}
}

// ---------------------------------------------------------------------------

func TestReconcile_Installed_Confirmed(t *testing.T) {
	a := aptAdapter(adapter.StateInstalled)
	eng := engine.NewDefaultEngine([]adapter.Adapter{a}, nil)

	recs := eng.Reconcile([]state.LocalInstallation{
		installedRec("git", source.TypeAPT, "git"),
	})

	if len(recs) != 1 {
		t.Fatalf("expected 1 drift record, got %d", len(recs))
	}
	if recs[0].DriftKind != engine.DriftNone {
		t.Errorf("expected DriftNone, got %q", recs[0].DriftKind)
	}
	if recs[0].BackendState != adapter.StateInstalled {
		t.Errorf("expected BackendState installed, got %q", recs[0].BackendState)
	}
}

func TestReconcile_Installed_MissingFromBackend(t *testing.T) {
	a := aptAdapter(adapter.StateNotInstalled)
	eng := engine.NewDefaultEngine([]adapter.Adapter{a}, nil)

	recs := eng.Reconcile([]state.LocalInstallation{
		installedRec("wget", source.TypeAPT, "wget"),
	})

	if len(recs) != 1 {
		t.Fatalf("expected 1 drift record, got %d", len(recs))
	}
	if recs[0].DriftKind != engine.DriftMissingFromBackend {
		t.Errorf("expected DriftMissingFromBackend, got %q", recs[0].DriftKind)
	}
	if recs[0].BackendState != adapter.StateNotInstalled {
		t.Errorf("expected StateNotInstalled, got %q", recs[0].BackendState)
	}
	if recs[0].Message == "" {
		t.Error("expected non-empty drift message")
	}
}

func TestReconcile_Installed_PartialInBackend(t *testing.T) {
	a := aptAdapter(adapter.StatePartial)
	eng := engine.NewDefaultEngine([]adapter.Adapter{a}, nil)

	recs := eng.Reconcile([]state.LocalInstallation{
		installedRec("curl", source.TypeAPT, "curl"),
	})

	if len(recs) != 1 {
		t.Fatalf("expected 1 drift record, got %d", len(recs))
	}
	if recs[0].DriftKind != engine.DriftPartialInBackend {
		t.Errorf("expected DriftPartialInBackend, got %q", recs[0].DriftKind)
	}
}

func TestReconcile_Removed_Confirmed(t *testing.T) {
	a := aptAdapter(adapter.StateNotInstalled)
	eng := engine.NewDefaultEngine([]adapter.Adapter{a}, nil)

	recs := eng.Reconcile([]state.LocalInstallation{
		removedRec("wget", source.TypeAPT, "wget"),
	})

	if len(recs) != 1 {
		t.Fatalf("expected 1 drift record, got %d", len(recs))
	}
	if recs[0].DriftKind != engine.DriftNone {
		t.Errorf("expected DriftNone for confirmed removal, got %q", recs[0].DriftKind)
	}
}

func TestReconcile_Removed_PresentInBackend(t *testing.T) {
	a := aptAdapter(adapter.StateInstalled)
	eng := engine.NewDefaultEngine([]adapter.Adapter{a}, nil)

	recs := eng.Reconcile([]state.LocalInstallation{
		removedRec("wget", source.TypeAPT, "wget"),
	})

	if len(recs) != 1 {
		t.Fatalf("expected 1 drift record, got %d", len(recs))
	}
	if recs[0].DriftKind != engine.DriftPresentButRemoved {
		t.Errorf("expected DriftPresentButRemoved, got %q", recs[0].DriftKind)
	}
}

func TestReconcile_NoAdapter_DriftUnknown(t *testing.T) {
	// No adapter handles Flatpak when only APT adapter is configured.
	a := aptAdapter(adapter.StateInstalled)
	eng := engine.NewDefaultEngine([]adapter.Adapter{a}, nil)

	recs := eng.Reconcile([]state.LocalInstallation{
		installedRec("obs-studio", source.TypeFlatpak, "com.obsproject.Studio"),
	})

	if len(recs) != 1 {
		t.Fatalf("expected 1 drift record, got %d", len(recs))
	}
	if recs[0].DriftKind != engine.DriftUnknown {
		t.Errorf("expected DriftUnknown (no adapter), got %q", recs[0].DriftKind)
	}
}

func TestReconcile_AdapterError_DriftUnknown(t *testing.T) {
	a := &reconcileAdapter{
		name:         "apt",
		available:    true,
		handledTypes: []source.Type{source.TypeAPT},
		checkErr:     errors.New("backend unreachable"),
	}
	eng := engine.NewDefaultEngine([]adapter.Adapter{a}, nil)

	recs := eng.Reconcile([]state.LocalInstallation{
		installedRec("git", source.TypeAPT, "git"),
	})

	if len(recs) != 1 {
		t.Fatalf("expected 1 drift record, got %d", len(recs))
	}
	if recs[0].DriftKind != engine.DriftUnknown {
		t.Errorf("expected DriftUnknown on adapter error, got %q", recs[0].DriftKind)
	}
}

func TestReconcile_PendingStatus_Skipped(t *testing.T) {
	a := aptAdapter(adapter.StateNotInstalled)
	eng := engine.NewDefaultEngine([]adapter.Adapter{a}, nil)

	pending := state.LocalInstallation{
		ApplicationID:      "git",
		SourceType:         source.TypeAPT,
		SourceIdentifier:   "git",
		InstallTimestamp:   time.Now(),
		InstallStatus:      state.StatusPending,
		VerificationStatus: state.VerificationPending,
	}
	recs := eng.Reconcile([]state.LocalInstallation{pending})

	if len(recs) != 1 {
		t.Fatalf("expected 1 drift record, got %d", len(recs))
	}
	if recs[0].DriftKind != engine.DriftNone {
		t.Errorf("expected DriftNone for pending status, got %q", recs[0].DriftKind)
	}
}

func TestReconcile_FailedStatus_Skipped(t *testing.T) {
	a := aptAdapter(adapter.StateNotInstalled)
	eng := engine.NewDefaultEngine([]adapter.Adapter{a}, nil)

	failed := state.LocalInstallation{
		ApplicationID:      "git",
		SourceType:         source.TypeAPT,
		SourceIdentifier:   "git",
		InstallTimestamp:   time.Now(),
		InstallStatus:      state.StatusFailed,
		VerificationStatus: state.VerificationPending,
	}
	recs := eng.Reconcile([]state.LocalInstallation{failed})

	if len(recs) != 1 {
		t.Fatalf("expected 1 drift record, got %d", len(recs))
	}
	if recs[0].DriftKind != engine.DriftNone {
		t.Errorf("expected DriftNone for failed status, got %q", recs[0].DriftKind)
	}
}

func TestReconcile_EmptyList_NoResults(t *testing.T) {
	eng := engine.NewDefaultEngine(nil, nil)

	recs := eng.Reconcile(nil)
	if len(recs) != 0 {
		t.Errorf("expected empty results for nil input, got %d", len(recs))
	}
}

func TestReconcile_MultipleRecords(t *testing.T) {
	installed := aptAdapter(adapter.StateInstalled)
	eng := engine.NewDefaultEngine([]adapter.Adapter{installed}, nil)

	records := []state.LocalInstallation{
		installedRec("git", source.TypeAPT, "git"),
		installedRec("wget", source.TypeAPT, "wget"),
	}
	recs := eng.Reconcile(records)
	if len(recs) != 2 {
		t.Fatalf("expected 2 drift records, got %d", len(recs))
	}
	// Both confirm installed since adapter always returns StateInstalled.
	for _, r := range recs {
		if r.DriftKind != engine.DriftNone {
			t.Errorf("expected DriftNone for %s, got %q", r.Installation.ApplicationID, r.DriftKind)
		}
	}
}

func TestReconcile_MessageIsNonEmpty(t *testing.T) {
	a := aptAdapter(adapter.StateNotInstalled)
	eng := engine.NewDefaultEngine([]adapter.Adapter{a}, nil)

	recs := eng.Reconcile([]state.LocalInstallation{
		installedRec("git", source.TypeAPT, "git"),
	})

	if recs[0].Message == "" {
		t.Error("expected non-empty drift message for missing-from-backend drift")
	}
}
