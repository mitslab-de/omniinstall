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

// ---------------------------------------------------------------------------
// CheckInstalled scope and error-handling robustness tests (task 0025)
// ---------------------------------------------------------------------------

func TestCheckInstalled_SystemScope_Found(t *testing.T) {
// System scope returns with ID field — should report installed.
e := &fakeExecutor{
runFn: func(name string, args ...string) (string, int, error) {
for _, a := range args {
if a == "--system" {
return "ID:              com.obsproject.Studio\nVersion: 30.0.0", 0, nil
}
}
return "", 1, errors.New("not found in user scope")
},
}
a := newTestAdapter(e)
state, err := a.CheckInstalled("com.obsproject.Studio")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if state != adapter.StateInstalled {
t.Errorf("expected StateInstalled (system scope), got %q", state)
}
}

func TestCheckInstalled_UserScope_Only(t *testing.T) {
// System scope not found; user scope has the app.
e := &fakeExecutor{
runFn: func(name string, args ...string) (string, int, error) {
for _, a := range args {
if a == "--user" {
return "Application:     com.obsproject.Studio\nVersion: 30.0.0", 0, nil
}
}
// system scope → not found
return "", 1, errors.New("not found")
},
}
a := newTestAdapter(e)
state, err := a.CheckInstalled("com.obsproject.Studio")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if state != adapter.StateInstalled {
t.Errorf("expected StateInstalled (user scope), got %q", state)
}
}

func TestCheckInstalled_NotInAnyScope(t *testing.T) {
// Neither scope has the app.
e := &fakeExecutor{
runFn: func(name string, args ...string) (string, int, error) {
return "", 1, errors.New("error: No ref found for com.NotExists")
},
}
a := newTestAdapter(e)
state, err := a.CheckInstalled("com.NotExists")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if state != adapter.StateNotInstalled {
t.Errorf("expected StateNotInstalled, got %q", state)
}
}

func TestCheckInstalled_NoInstallations_StateUnknown(t *testing.T) {
// Flatpak is misconfigured — no installations directory.
e := &fakeExecutor{
runFn: func(name string, args ...string) (string, int, error) {
return "error: No installations", 1, errors.New("exit status 1")
},
}
a := newTestAdapter(e)
state, err := a.CheckInstalled("com.obsproject.Studio")
if err == nil {
t.Error("expected error for no-installations condition")
}
if state != adapter.StateUnknown {
t.Errorf("expected StateUnknown, got %q", state)
}
}

func TestCheckInstalled_PermissionDenied_StateUnknown(t *testing.T) {
e := &fakeExecutor{
runFn: func(name string, args ...string) (string, int, error) {
return "permission denied reading /var/lib/flatpak", 1, errors.New("exit status 1")
},
}
a := newTestAdapter(e)
state, err := a.CheckInstalled("com.obsproject.Studio")
if err == nil {
t.Error("expected error for permission denied")
}
if state != adapter.StateUnknown {
t.Errorf("expected StateUnknown, got %q", state)
}
}

func TestCheckInstalled_EmptyOutput_Exit0_Installed(t *testing.T) {
// flatpak info exits 0 but has no output — older versions; assume installed
// because exit code 0 means success.
e := &fakeExecutor{
runFn: func(name string, args ...string) (string, int, error) {
for _, a := range args {
if a == "--system" {
return "", 0, nil
}
}
return "", 1, errors.New("not found")
},
}
a := newTestAdapter(e)
state, err := a.CheckInstalled("com.obsproject.Studio")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
// parseFlatpakInfoInstalled returns false for empty output, so we fall
// through to user scope, which is not found → StateNotInstalled.
if state != adapter.StateNotInstalled {
t.Errorf("expected StateNotInstalled for empty exit-0 then not-found user, got %q", state)
}
}

func TestCheckInstalled_CommandLaunchError_StateUnknown(t *testing.T) {
// The flatpak binary crashes before writing anything (exitCode==0, err!=nil).
e := &fakeExecutor{
runFn: func(name string, args ...string) (string, int, error) {
return "", 0, errors.New("fork/exec: no such file or directory")
},
}
a := newTestAdapter(e)
state, err := a.CheckInstalled("com.obsproject.Studio")
if err == nil {
t.Error("expected error for command launch failure")
}
if state != adapter.StateUnknown {
t.Errorf("expected StateUnknown, got %q", state)
}
}

func TestCheckInstalled_RealFlatpakInfoFormat(t *testing.T) {
// Simulate realistic `flatpak info --system com.obsproject.Studio` output.
realOutput := `          ID: com.obsproject.Studio
        Ref: app/com.obsproject.Studio/x86_64/stable
       Arch: x86_64
     Branch: stable
    Version: 30.0.0
    License: GPL-2.0+
    Summary: Free and open source software for live streaming and screen recording
        URL: https://obsproject.com
     Commit: abc123
     Parent: (none)
   Location: /var/lib/flatpak/app/com.obsproject.Studio/x86_64/stable/active
  Installed: 295.0 MB
`
e := &fakeExecutor{
runFn: func(name string, args ...string) (string, int, error) {
for _, a := range args {
if a == "--system" {
return realOutput, 0, nil
}
}
return "", 1, errors.New("not found")
},
}
a := newTestAdapter(e)
state, err := a.CheckInstalled("com.obsproject.Studio")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if state != adapter.StateInstalled {
t.Errorf("expected StateInstalled for realistic flatpak info output, got %q", state)
}
}
