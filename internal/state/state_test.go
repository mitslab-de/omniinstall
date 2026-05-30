package state_test

import (
	"testing"
	"time"

	"github.com/mitslab-de/omniinstall/internal/source"
	"github.com/mitslab-de/omniinstall/internal/state"
)

func TestLocalInstallation_Validate_Valid(t *testing.T) {
	r := &state.LocalInstallation{
		ApplicationID:      "obs-studio",
		SourceType:         source.TypeAPT,
		SourceIdentifier:   "obs-studio",
		InstallTimestamp:   time.Now(),
		InstallStatus:      state.StatusInstalled,
		VerificationStatus: state.VerificationPassed,
	}
	if err := r.Validate(); err != nil {
		t.Errorf("expected valid record, got error: %v", err)
	}
}

func TestLocalInstallation_Validate_MissingApplicationID(t *testing.T) {
	r := &state.LocalInstallation{
		SourceType:         source.TypeAPT,
		SourceIdentifier:   "obs-studio",
		InstallTimestamp:   time.Now(),
		InstallStatus:      state.StatusInstalled,
		VerificationStatus: state.VerificationPassed,
	}
	if err := r.Validate(); err == nil {
		t.Error("expected error for missing application_id, got nil")
	}
}

func TestLocalInstallation_Validate_MissingTimestamp(t *testing.T) {
	r := &state.LocalInstallation{
		ApplicationID:      "obs-studio",
		SourceType:         source.TypeAPT,
		SourceIdentifier:   "obs-studio",
		InstallStatus:      state.StatusInstalled,
		VerificationStatus: state.VerificationPassed,
	}
	if err := r.Validate(); err == nil {
		t.Error("expected error for missing install_timestamp, got nil")
	}
}

func TestInstallStatusConstants(t *testing.T) {
	statuses := []state.InstallStatus{
		state.StatusInstalled,
		state.StatusFailed,
		state.StatusRemoved,
		state.StatusPending,
	}
	for _, s := range statuses {
		if s == "" {
			t.Error("expected non-empty install status constant")
		}
	}
}

func TestVerificationStatusConstants(t *testing.T) {
	statuses := []state.VerificationStatus{
		state.VerificationPassed,
		state.VerificationFailed,
		state.VerificationSkipped,
		state.VerificationPending,
	}
	for _, s := range statuses {
		if s == "" {
			t.Error("expected non-empty verification status constant")
		}
	}
}
