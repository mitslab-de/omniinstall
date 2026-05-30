package pacman_test

import (
	"errors"
	"strings"
	"testing"

	adapter "github.com/mitslab-de/omniinstall/internal/adapters"
	"github.com/mitslab-de/omniinstall/internal/adapters/pacman"
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

func newTestAdapter(e *fakeExecutor) *pacman.Adapter {
	return pacman.NewWithExecutorForTest(e)
}

var goodPlan = &install.Plan{
	ApplicationID:    "git",
	SourceType:       source.TypePacman,
	SourceIdentifier: "git",
	RiskLevel:        source.RiskLow,
	Verification:     []install.VerificationRule{{Command: "git"}},
}

func TestCanHandle(t *testing.T) {
	if !pacman.CanHandle(source.TypePacman) {
		t.Error("expected Pacman adapter to handle TypePacman")
	}
	if pacman.CanHandle(source.TypeAPT) {
		t.Error("expected Pacman adapter NOT to handle TypeAPT")
	}
	if pacman.CanHandle(source.TypeFlatpak) {
		t.Error("expected Pacman adapter NOT to handle TypeFlatpak")
	}
}

func TestAdapterCanHandle_Method(t *testing.T) {
	a := newTestAdapter(&fakeExecutor{})
	if !a.CanHandle(source.TypePacman) {
		t.Error("Adapter.CanHandle should return true for TypePacman")
	}
	if a.CanHandle(source.TypeAPT) {
		t.Error("Adapter.CanHandle should return false for TypeAPT")
	}
}

func TestName(t *testing.T) {
	if pacman.AdapterName == "" {
		t.Error("expected non-empty adapter name constant")
	}
	a := newTestAdapter(&fakeExecutor{})
	if a.Name() != "pacman" {
		t.Errorf("expected adapter name 'pacman', got %q", a.Name())
	}
}

func TestIsAvailable_WhenPacmanPresent(t *testing.T) {
	e := &fakeExecutor{
		lookPathFn: func(file string) (string, error) {
			if file == "pacman" {
				return "/usr/bin/pacman", nil
			}
			return "", errors.New("not found")
		},
	}
	a := newTestAdapter(e)
	if !a.IsAvailable() {
		t.Error("expected IsAvailable to return true when pacman is on PATH")
	}
}

func TestIsAvailable_WhenPacmanAbsent(t *testing.T) {
	e := &fakeExecutor{
		lookPathFn: func(_ string) (string, error) {
			return "", errors.New("not found")
		},
	}
	a := newTestAdapter(e)
	if a.IsAvailable() {
		t.Error("expected IsAvailable to return false when pacman is not on PATH")
	}
}

func TestCheckInstalled_Installed(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "git 2.40.1-1\n", 0, nil
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

func TestCheckInstalled_NotInstalled(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "error: package 'git' was not found\n", 1, errors.New("exit status 1")
		},
	}
	a := newTestAdapter(e)
	state, err := a.CheckInstalled("git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state != adapter.StateNotInstalled {
		t.Errorf("expected StateNotInstalled, got %q", state)
	}
}

func TestCheckInstalled_UnexpectedError(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "some unexpected failure\n", 5, errors.New("exit status 5")
		},
	}
	a := newTestAdapter(e)
	state, err := a.CheckInstalled("git")
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
	state, err := a.CheckInstalled("git")
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
			return "", 0, errors.New("cannot exec pacman")
		},
	}
	a := newTestAdapter(e)
	state, err := a.CheckInstalled("git")
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
			return "Packages (1) git-2.40.1-1\nTotal Installed Size: 48.59 MiB\n", 0, nil
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

func TestInstall_NonZeroExit(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "error: target not found: git\n", 1, errors.New("exit status 1")
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
			return "removing git...\n", 0, nil
		},
	}
	a := newTestAdapter(e)
	result, err := a.Remove("git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected remove to succeed, got: %s", result.Message)
	}
}

func TestRemove_NonZeroExit(t *testing.T) {
	e := &fakeExecutor{
		runFn: func(name string, args ...string) (string, int, error) {
			return "error: target not found: git\n", 1, errors.New("exit status 1")
		},
	}
	a := newTestAdapter(e)
	result, err := a.Remove("git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected remove to fail on non-zero exit")
	}
}

func TestVerify_CommandFound(t *testing.T) {
	e := &fakeExecutor{
		lookPathFn: func(file string) (string, error) {
			return "/usr/bin/" + file, nil
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
