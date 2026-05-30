package flatpak_test

import (
	"errors"
	"testing"

	adapter "github.com/mitslab-de/omniinstall/internal/adapters"
	"github.com/mitslab-de/omniinstall/internal/adapters/flatpak"
	"github.com/mitslab-de/omniinstall/internal/install"
	"github.com/mitslab-de/omniinstall/internal/source"
)

// fakeExecutor is a configurable test double for OS command execution.
type fakeExecutor struct {
	lookPathFn func(string) (string, error)
	runFn      func(name string, args ...string) (string, int, error)
}

func (f *fakeExecutor) LookPath(file string) (string, error) {
	if f.lookPathFn != nil {
		return f.lookPathFn(file)
	}
	return "/usr/bin/" + file, nil
}

func (f *fakeExecutor) Run(name string, args ...string) (string, int, error) {
	if f.runFn != nil {
		return f.runFn(name, args...)
	}
	return "", 0, nil
}

func newTestAdapter(e *fakeExecutor) *flatpak.Adapter {
	return flatpak.NewWithExecutorForTest(e)
}

var goodPlan = &install.Plan{
	ApplicationID:    "com.obsproject.Studio",
	SourceType:       source.TypeFlatpak,
	SourceIdentifier: "com.obsproject.Studio",
	RiskLevel:        source.RiskLow,
	Verification:     []install.VerificationRule{{Command: "com.obsproject.Studio"}},
}

func TestCanHandle(t *testing.T) {
	if !flatpak.CanHandle(source.TypeFlatpak) {
		t.Error("expected Flatpak adapter to handle TypeFlatpak")
	}
	if flatpak.CanHandle(source.TypeAPT) {
		t.Error("expected Flatpak adapter NOT to handle TypeAPT")
	}
	if flatpak.CanHandle(source.TypeDNF) {
		t.Error("expected Flatpak adapter NOT to handle TypeDNF")
	}
}

func TestName(t *testing.T) {
	if flatpak.Name == "" {
		t.Error("expected non-empty adapter name")
	}
	a := newTestAdapter(&fakeExecutor{})
	if a.Name() == "" {
		t.Error("Adapter.Name() must return non-empty string")
	}
}

func TestIsAvailableWhenFlatpakPresent(t *testing.T) {
	e := &fakeExecutor{
		lookPathFn: func(file string) (string, error) {
			if file == "flatpak" {
				return "/usr/bin/flatpak", nil
			}
			return "", errors.New("not found")
		},
	}
	a := newTestAdapter(e)
	if !a.IsAvailable() {
		t.Error("expected IsAvailable=true when flatpak is in PATH")
	}
}

func TestIsAvailableWhenFlatpakMissing(t *testing.T) {
	e := &fakeExecutor{
		lookPathFn: func(file string) (string, error) {
			return "", errors.New("not found")
		},
	}
	a := newTestAdapter(e)
	if a.IsAvailable() {
		t.Error("expected IsAvailable=false when flatpak is not in PATH")
	}
}

func TestAdapterCanHandle(t *testing.T) {
	a := newTestAdapter(&fakeExecutor{})
	if !a.CanHandle(source.TypeFlatpak) {
		t.Error("Adapter.CanHandle must return true for TypeFlatpak")
	}
	if a.CanHandle(source.TypeAPT) {
		t.Error("Adapter.CanHandle must return false for TypeAPT")
	}
}

func TestCheckInstalledFound(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "com.obsproject.Studio\nVersion: 30.0.0", 0, nil
		},
	}
	a := newTestAdapter(e)
	state, err := a.CheckInstalled("com.obsproject.Studio")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state != adapter.StateInstalled {
		t.Errorf("expected StateInstalled, got %s", state)
	}
}

func TestCheckInstalledNotFound(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "", 1, errors.New("error: No ref found")
		},
	}
	a := newTestAdapter(e)
	state, err := a.CheckInstalled("com.obsproject.Studio")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state != adapter.StateNotInstalled {
		t.Errorf("expected StateNotInstalled, got %s", state)
	}
}

func TestCheckInstalledReadinessError(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "", 0, errors.New("permission denied")
		},
	}
	a := newTestAdapter(e)
	state, err := a.CheckInstalled("com.obsproject.Studio")
	if err == nil {
		t.Fatal("expected readiness error")
	}
	if state != adapter.StateUnknown {
		t.Fatalf("expected StateUnknown, got %s", state)
	}
}

func TestInstallSuccess(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "Installing com.obsproject.Studio...\nDone.", 0, nil
		},
	}
	a := newTestAdapter(e)
	result, err := a.Install(goodPlan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected success, got: %s", result.Message)
	}
	if !result.ChangedSystem {
		t.Error("expected ChangedSystem=true on successful install")
	}
}

func TestInstallFailure(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "error: No ref found for 'nonexistent'", 1, errors.New("exit status 1")
		},
	}
	a := newTestAdapter(e)
	plan := &install.Plan{
		ApplicationID:    "nonexistent",
		SourceType:       source.TypeFlatpak,
		SourceIdentifier: "nonexistent",
		RiskLevel:        source.RiskLow,
		Verification:     []install.VerificationRule{{Command: "nonexistent"}},
	}
	result, err := a.Install(plan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected failure")
	}
	if result.ErrorCategory != adapter.ErrPackageNotFound {
		t.Errorf("expected ErrPackageNotFound, got %s", result.ErrorCategory)
	}
}

func TestRemoveSuccess(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "Removing com.obsproject.Studio...\nDone.", 0, nil
		},
	}
	a := newTestAdapter(e)
	result, err := a.Remove("com.obsproject.Studio")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected success, got: %s", result.Message)
	}
}

func TestVerifySuccess(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			// flatpak info returns output for an installed app
			return "com.obsproject.Studio\nVersion: 30.0.0", 0, nil
		},
	}
	a := newTestAdapter(e)
	vr, err := a.Verify(goodPlan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !vr.Verified {
		t.Errorf("expected Verified=true, got details: %s", vr.Details)
	}
}

func TestVerifyFailure(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "", 1, errors.New("not installed")
		},
	}
	a := newTestAdapter(e)
	vr, err := a.Verify(goodPlan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vr.Verified {
		t.Error("expected Verified=false when app not found")
	}
}

func TestAdapterImplementsInterface(t *testing.T) {
	var _ adapter.Adapter = newTestAdapter(&fakeExecutor{})
}

func TestAdapterContract(t *testing.T) {
	a := newTestAdapter(&fakeExecutor{})
	if a.Name() == "" {
		t.Error("adapter Name() must return a non-empty string")
	}
	if a.CanHandle("unknown-source-type-xyz") {
		t.Error("CanHandle must return false for unknown source type")
	}
}
