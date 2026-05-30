package apt_test

import (
	"errors"
	"testing"

	adapter "github.com/mitslab-de/omniinstall/internal/adapters"
	adaptercontract "github.com/mitslab-de/omniinstall/internal/adapters"
	"github.com/mitslab-de/omniinstall/internal/adapters/apt"
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

// newTestAdapter creates an APT Adapter using a fakeExecutor (package-internal).
// We use the exported NewWithExecutorForTest to keep the executor unexported.
func newTestAdapter(e *fakeExecutor) *apt.Adapter {
	return apt.NewWithExecutorForTest(e)
}

// goodPlan is a minimal valid install plan.
var goodPlan = &install.Plan{
	ApplicationID:    "wget",
	SourceType:       source.TypeAPT,
	SourceIdentifier: "wget",
	RiskLevel:        source.RiskLow,
	Verification:     []install.VerificationRule{{Command: "wget"}},
}

func TestCanHandle(t *testing.T) {
	if !apt.CanHandle(source.TypeAPT) {
		t.Error("expected APT adapter to handle TypeAPT")
	}
	if apt.CanHandle(source.TypeFlatpak) {
		t.Error("expected APT adapter NOT to handle TypeFlatpak")
	}
	if apt.CanHandle(source.TypeDNF) {
		t.Error("expected APT adapter NOT to handle TypeDNF")
	}
}

func TestName(t *testing.T) {
	if apt.Name == "" {
		t.Error("expected non-empty adapter name")
	}
	a := newTestAdapter(&fakeExecutor{})
	if a.Name() == "" {
		t.Error("Adapter.Name() must return non-empty string")
	}
}

func TestIsAvailableWhenAptGetPresent(t *testing.T) {
	e := &fakeExecutor{
		lookPathFn: func(file string) (string, error) {
			if file == "apt-get" {
				return "/usr/bin/apt-get", nil
			}
			return "", errors.New("not found")
		},
	}
	a := newTestAdapter(e)
	if !a.IsAvailable() {
		t.Error("expected IsAvailable=true when apt-get is in PATH")
	}
}

func TestIsAvailableWhenAptGetMissing(t *testing.T) {
	e := &fakeExecutor{
		lookPathFn: func(file string) (string, error) {
			return "", errors.New("not found")
		},
	}
	a := newTestAdapter(e)
	if a.IsAvailable() {
		t.Error("expected IsAvailable=false when apt-get is not in PATH")
	}
}

func TestAdapterCanHandle(t *testing.T) {
	a := newTestAdapter(&fakeExecutor{})
	if !a.CanHandle(source.TypeAPT) {
		t.Error("Adapter.CanHandle must return true for TypeAPT")
	}
	if a.CanHandle(source.TypeFlatpak) {
		t.Error("Adapter.CanHandle must return false for TypeFlatpak")
	}
}

func TestCheckInstalledFound(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "install ok installed", 0, nil
		},
	}
	a := newTestAdapter(e)
	state, err := a.CheckInstalled("wget")
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
			return "", 1, errors.New("no packages found")
		},
	}
	a := newTestAdapter(e)
	state, err := a.CheckInstalled("nonexistent")
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
	state, err := a.CheckInstalled("wget")
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
			return "Reading package lists...\nDone.", 0, nil
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
	if result.Duration == 0 {
		t.Error("expected non-zero Duration")
	}
}

func TestInstallFailure(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "E: Unable to locate package nonexistent", 100, errors.New("exit status 100")
		},
	}
	a := newTestAdapter(e)
	plan := &install.Plan{
		ApplicationID:    "nonexistent",
		SourceType:       source.TypeAPT,
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
			return "Removing wget ...", 0, nil
		},
	}
	a := newTestAdapter(e)
	result, err := a.Remove("wget")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected success, got: %s", result.Message)
	}
}

func TestRemoveFailure(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "E: Package 'wget' is not installed", 1, errors.New("exit status 1")
		},
	}
	a := newTestAdapter(e)
	result, err := a.Remove("wget")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected failure")
	}
}

func TestVerifySuccess(t *testing.T) {
	e := &fakeExecutor{
		lookPathFn: func(file string) (string, error) {
			return "/usr/bin/" + file, nil
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
		lookPathFn: func(file string) (string, error) {
			return "", errors.New("not found")
		},
	}
	a := newTestAdapter(e)
	vr, err := a.Verify(goodPlan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vr.Verified {
		t.Error("expected Verified=false when command not in PATH")
	}
}

func TestAdapterImplementsInterface(t *testing.T) {
	var _ adaptercontract.Adapter = newTestAdapter(&fakeExecutor{})
}

func TestAdapterContract(t *testing.T) {
	// Use the contract test helper from the adapters package.
	// We instantiate a fake executor so APT need not be installed.
	a := newTestAdapter(&fakeExecutor{
		lookPathFn: func(file string) (string, error) {
			return "/usr/bin/" + file, nil
		},
	})
	// Manually run contract checks (avoiding import cycle with test file).
	if a.Name() == "" {
		t.Error("adapter Name() must return a non-empty string")
	}
	if a.CanHandle("unknown-source-type-xyz") {
		t.Error("CanHandle must return false for unknown source type")
	}
}

// ---------------------------------------------------------------------------
// CheckInstalled dpkg status format robustness tests (task 0024)
// ---------------------------------------------------------------------------

func TestCheckInstalled_FullDpkgFormat_Installed(t *testing.T) {
// Realistic dpkg-query output: "<want> <eflag> <status>"
e := &fakeExecutor{
runFn: func(name string, args ...string) (string, int, error) {
return "install ok installed", 0, nil
},
}
a := newTestAdapter(e)
state, err := a.CheckInstalled("git")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if state != adapter.StateInstalled {
t.Errorf("expected StateInstalled, got %q", state)
}
}

func TestCheckInstalled_ConfigFiles_NotInstalled(t *testing.T) {
// Package removed but config files remain.
e := &fakeExecutor{
runFn: func(name string, args ...string) (string, int, error) {
return "deinstall ok config-files", 0, nil
},
}
a := newTestAdapter(e)
state, err := a.CheckInstalled("git")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if state != adapter.StateNotInstalled {
t.Errorf("expected StateNotInstalled for config-files, got %q", state)
}
}

func TestCheckInstalled_Purge_NotInstalled(t *testing.T) {
e := &fakeExecutor{
runFn: func(name string, args ...string) (string, int, error) {
return "purge ok not-installed", 0, nil
},
}
a := newTestAdapter(e)
state, err := a.CheckInstalled("wget")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if state != adapter.StateNotInstalled {
t.Errorf("expected StateNotInstalled, got %q", state)
}
}

func TestCheckInstalled_HalfInstalled_Partial(t *testing.T) {
e := &fakeExecutor{
runFn: func(name string, args ...string) (string, int, error) {
return "install ok half-installed", 0, nil
},
}
a := newTestAdapter(e)
state, err := a.CheckInstalled("wget")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if state != adapter.StatePartial {
t.Errorf("expected StatePartial for half-installed, got %q", state)
}
}

func TestCheckInstalled_Unpacked_Partial(t *testing.T) {
e := &fakeExecutor{
runFn: func(name string, args ...string) (string, int, error) {
return "install ok unpacked", 0, nil
},
}
a := newTestAdapter(e)
state, err := a.CheckInstalled("wget")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if state != adapter.StatePartial {
t.Errorf("expected StatePartial for unpacked, got %q", state)
}
}

func TestCheckInstalled_HalfConfigured_Partial(t *testing.T) {
e := &fakeExecutor{
runFn: func(name string, args ...string) (string, int, error) {
return "install reinstreq half-configured", 0, nil
},
}
a := newTestAdapter(e)
state, err := a.CheckInstalled("wget")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if state != adapter.StatePartial {
t.Errorf("expected StatePartial for half-configured, got %q", state)
}
}

func TestCheckInstalled_TriggersAwaited_Partial(t *testing.T) {
e := &fakeExecutor{
runFn: func(name string, args ...string) (string, int, error) {
return "install ok triggers-awaited", 0, nil
},
}
a := newTestAdapter(e)
state, err := a.CheckInstalled("wget")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if state != adapter.StatePartial {
t.Errorf("expected StatePartial for triggers-awaited, got %q", state)
}
}

func TestCheckInstalled_TriggersPending_Partial(t *testing.T) {
e := &fakeExecutor{
runFn: func(name string, args ...string) (string, int, error) {
return "install ok triggers-pending", 0, nil
},
}
a := newTestAdapter(e)
state, err := a.CheckInstalled("wget")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if state != adapter.StatePartial {
t.Errorf("expected StatePartial for triggers-pending, got %q", state)
}
}

func TestCheckInstalled_EmptyOutput_NotInstalled(t *testing.T) {
// dpkg-query returns empty string with exit 0 for some dpkg versions
// when the package is not in the database.
e := &fakeExecutor{
runFn: func(name string, args ...string) (string, int, error) {
return "", 0, nil
},
}
a := newTestAdapter(e)
state, err := a.CheckInstalled("ghost-pkg")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if state != adapter.StateNotInstalled {
t.Errorf("expected StateNotInstalled for empty output, got %q", state)
}
}

func TestCheckInstalled_UnknownStatusField_Unknown(t *testing.T) {
e := &fakeExecutor{
runFn: func(name string, args ...string) (string, int, error) {
return "install ok future-status-field", 0, nil
},
}
a := newTestAdapter(e)
state, err := a.CheckInstalled("wget")
if err == nil {
t.Error("expected error for unknown status field")
}
if state != adapter.StateUnknown {
t.Errorf("expected StateUnknown, got %q", state)
}
}

func TestCheckInstalled_MalformedOutput_Unknown(t *testing.T) {
// Only one word — not a valid dpkg status line.
e := &fakeExecutor{
runFn: func(name string, args ...string) (string, int, error) {
return "garbage", 0, nil
},
}
a := newTestAdapter(e)
state, err := a.CheckInstalled("wget")
if err == nil {
t.Error("expected error for malformed dpkg output")
}
if state != adapter.StateUnknown {
t.Errorf("expected StateUnknown, got %q", state)
}
}

func TestCheckInstalled_ExitCode1_NotInstalled(t *testing.T) {
// Non-zero exit means package not in dpkg database.
e := &fakeExecutor{
runFn: func(name string, args ...string) (string, int, error) {
return "dpkg-query: no packages found matching nonexistent-pkg", 1, errors.New("exit status 1")
},
}
a := newTestAdapter(e)
state, err := a.CheckInstalled("nonexistent-pkg")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if state != adapter.StateNotInstalled {
t.Errorf("expected StateNotInstalled, got %q", state)
}
}

func TestCheckInstalled_TrailingWhitespace_Installed(t *testing.T) {
// Some dpkg versions add trailing newline/whitespace.
e := &fakeExecutor{
runFn: func(name string, args ...string) (string, int, error) {
return "install ok installed\n", 0, nil
},
}
a := newTestAdapter(e)
state, err := a.CheckInstalled("wget")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if state != adapter.StateInstalled {
t.Errorf("expected StateInstalled, got %q", state)
}
}
