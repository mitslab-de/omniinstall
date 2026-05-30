package adapter_test

import (
	"testing"

	adapter "github.com/mitslab-de/omniinstall/internal/adapters"
)

func TestResult_Success(t *testing.T) {
	r := &adapter.Result{
		Success:       true,
		Message:       "Package installed successfully.",
		ChangedSystem: true,
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
	if !r.ChangedSystem {
		t.Error("expected ChangedSystem to be true")
	}
}

func TestResult_Failure(t *testing.T) {
	r := &adapter.Result{
		Success:       false,
		ExitCode:      100,
		Message:       "Package not found.",
		ErrorCategory: adapter.ErrPackageNotFound,
		ChangedSystem: false,
	}
	if r.Success {
		t.Error("expected Success to be false")
	}
	if r.ErrorCategory != adapter.ErrPackageNotFound {
		t.Errorf("expected error category %q, got %q", adapter.ErrPackageNotFound, r.ErrorCategory)
	}
}

func TestVerificationResult_Fields(t *testing.T) {
	v := &adapter.VerificationResult{
		Verified: true,
		Details:  "Command 'obs' found at /usr/bin/obs.",
	}
	if !v.Verified {
		t.Error("expected Verified to be true")
	}
}

func TestInstalledStateConstants(t *testing.T) {
	states := []adapter.InstalledState{
		adapter.StateInstalled,
		adapter.StateNotInstalled,
		adapter.StateUnknown,
	}
	if len(states) != 3 {
		t.Errorf("expected 3 installed state constants, got %d", len(states))
	}
}

func TestErrorCategoryConstants(t *testing.T) {
	cats := []string{
		adapter.ErrBackendUnavailable,
		adapter.ErrPackageNotFound,
		adapter.ErrPermissionDenied,
		adapter.ErrNetworkError,
		adapter.ErrDependencyConflict,
		adapter.ErrExecutionFailed,
		adapter.ErrUnsupportedOperation,
	}
	for _, c := range cats {
		if c == "" {
			t.Error("expected non-empty error category constant")
		}
	}
}
