package snap_test

import (
	"errors"
	"strings"
	"testing"

	adapter "github.com/mitslab-de/omniinstall/internal/adapters"
	"github.com/mitslab-de/omniinstall/internal/adapters/snap"
	"github.com/mitslab-de/omniinstall/internal/install"
	"github.com/mitslab-de/omniinstall/internal/source"
)

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

func newTestAdapter(e *fakeExecutor) *snap.Adapter {
	return snap.NewWithExecutorForTest(e)
}

var goodPlan = &install.Plan{
	ApplicationID:    "vlc",
	SourceType:       source.TypeSnap,
	SourceIdentifier: "vlc",
	RiskLevel:        source.RiskLow,
	Verification:     []install.VerificationRule{{Command: "vlc"}},
}

func TestCanHandle(t *testing.T) {
	if !snap.CanHandle(source.TypeSnap) {
		t.Error("expected Snap adapter to handle TypeSnap")
	}
	if snap.CanHandle(source.TypeAPT) {
		t.Error("expected Snap adapter NOT to handle TypeAPT")
	}
	if snap.CanHandle(source.TypeFlatpak) {
		t.Error("expected Snap adapter NOT to handle TypeFlatpak")
	}
}

func TestAdapterCanHandle_Method(t *testing.T) {
	a := newTestAdapter(&fakeExecutor{})
	if !a.CanHandle(source.TypeSnap) {
		t.Error("Adapter.CanHandle should return true for TypeSnap")
	}
	if a.CanHandle(source.TypeAPT) {
		t.Error("Adapter.CanHandle should return false for TypeAPT")
	}
}

func TestName(t *testing.T) {
	if snap.AdapterName == "" {
		t.Error("expected non-empty adapter name constant")
	}
	a := newTestAdapter(&fakeExecutor{})
	if a.Name() != "snap" {
		t.Errorf("expected adapter name 'snap', got %q", a.Name())
	}
}

func TestIsAvailable_WhenSnapPresent(t *testing.T) {
	e := &fakeExecutor{
		lookPathFn: func(file string) (string, error) {
			if file == "snap" {
				return "/usr/bin/snap", nil
			}
			return "", errors.New("not found")
		},
	}
	a := newTestAdapter(e)
	if !a.IsAvailable() {
		t.Error("expected IsAvailable to return true when snap is on PATH")
	}
}

func TestIsAvailable_WhenSnapAbsent(t *testing.T) {
	e := &fakeExecutor{
		lookPathFn: func(_ string) (string, error) {
			return "", errors.New("not found")
		},
	}
	a := newTestAdapter(e)
	if a.IsAvailable() {
		t.Error("expected IsAvailable to return false when snap is not on PATH")
	}
}

func TestCheckInstalled_Installed(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "Name  Version  Rev  Tracking  Publisher  Notes\nvlc   3.0.18   2988 latest/stable  videolan  -\n", 0, nil
		},
	}
	a := newTestAdapter(e)
	state, err := a.CheckInstalled("vlc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state != adapter.StateInstalled {
		t.Errorf("expected StateInstalled, got %q", state)
	}
}

func TestCheckInstalled_NotInstalled_NotFound(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "error: snap \"vlc\" is not installed\n", 1, errors.New("exit status 1")
		},
	}
	a := newTestAdapter(e)
	state, err := a.CheckInstalled("vlc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state != adapter.StateNotInstalled {
		t.Errorf("expected StateNotInstalled, got %q", state)
	}
}

func TestCheckInstalled_NotInstalled_NoSnaps(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "error: no snaps are installed yet\n", 1, errors.New("exit status 1")
		},
	}
	a := newTestAdapter(e)
	state, err := a.CheckInstalled("vlc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state != adapter.StateNotInstalled {
		t.Errorf("expected StateNotInstalled for no-snaps error, got %q", state)
	}
}

func TestCheckInstalled_UnexpectedError(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "snapd is not running\n", 5, errors.New("exit status 5")
		},
	}
	a := newTestAdapter(e)
	state, err := a.CheckInstalled("vlc")
	if state != adapter.StateUnknown {
		t.Errorf("expected StateUnknown for unexpected output, got %q", state)
	}
	if err == nil {
		t.Error("expected an error for unexpected output")
	}
}

func TestCheckInstalled_EmptyOutput(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "", 0, nil
		},
	}
	a := newTestAdapter(e)
	state, err := a.CheckInstalled("vlc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state != adapter.StateNotInstalled {
		t.Errorf("expected StateNotInstalled for empty output, got %q", state)
	}
}

func TestCheckInstalled_LaunchError(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "", 0, errors.New("cannot exec snap")
		},
	}
	a := newTestAdapter(e)
	state, err := a.CheckInstalled("vlc")
	if state != adapter.StateUnknown {
		t.Errorf("expected StateUnknown for launch error, got %q", state)
	}
	if err == nil {
		t.Error("expected an error when command cannot launch")
	}
}

func TestInstall_Success(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "vlc 3.0.18 from VideoLAN installed\n", 0, nil
		},
	}
	a := newTestAdapter(e)
	result, err := a.Install(goodPlan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected install to succeed, got: %s", result.Message)
	}
	if !result.ChangedSystem {
		t.Error("expected ChangedSystem to be true on successful install")
	}
}

func TestInstall_SnapNotFound(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "error: snap not found\n", 1, errors.New("exit status 1")
		},
	}
	a := newTestAdapter(e)
	result, err := a.Install(goodPlan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected install to fail on non-zero exit")
	}
	if result.ErrorCategory != adapter.ErrPackageNotFound {
		t.Errorf("expected ErrPackageNotFound, got %q", result.ErrorCategory)
	}
}

func TestRemove_Success(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "vlc removed\n", 0, nil
		},
	}
	a := newTestAdapter(e)
	result, err := a.Remove("vlc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected remove to succeed, got: %s", result.Message)
	}
}

func TestVerify_CommandFound(t *testing.T) {
	e := &fakeExecutor{
		lookPathFn: func(file string) (string, error) {
			return "/snap/bin/" + file, nil
		},
	}
	a := newTestAdapter(e)
	result, err := a.Verify(goodPlan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Verified {
		t.Errorf("expected verify to pass, details: %s", result.Details)
	}
}

func TestVerify_CommandNotFound(t *testing.T) {
	e := &fakeExecutor{
		lookPathFn: func(_ string) (string, error) {
			return "", errors.New("not found")
		},
	}
	a := newTestAdapter(e)
	result, err := a.Verify(goodPlan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Verified {
		t.Error("expected verify to fail when command not in PATH")
	}
	if !strings.Contains(result.Details, "not found") {
		t.Errorf("expected details to mention not found, got: %q", result.Details)
	}
}
