package engine_test

import (
	"testing"

	"github.com/mitslab-de/omniinstall/internal/engine"
)

func TestResult_Fields(t *testing.T) {
	r := &engine.Result{
		Success:       true,
		ApplicationID: "obs-studio",
		Message:       "OBS Studio was installed successfully.",
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
	if r.ApplicationID != "obs-studio" {
		t.Errorf("expected application_id 'obs-studio', got %q", r.ApplicationID)
	}
	if r.Message == "" {
		t.Error("expected non-empty message")
	}
}

func TestResult_FailureFields(t *testing.T) {
	r := &engine.Result{
		Success:       false,
		ApplicationID: "obs-studio",
		Message:       "Installation failed: package source unavailable.",
		ErrorCategory: "backend_unavailable",
	}
	if r.Success {
		t.Error("expected Success to be false")
	}
	if r.ErrorCategory == "" {
		t.Error("expected non-empty error category for failures")
	}
}
