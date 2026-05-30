package state_test

import (
	"encoding/json"
	"testing"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/mitslab-de/omniinstall/internal/source"
	"github.com/mitslab-de/omniinstall/internal/state"
)

func fullLocalInstallation() *state.LocalInstallation {
	return &state.LocalInstallation{
		ApplicationID:      "obs-studio",
		SourceType:         source.TypeAPT,
		SourceIdentifier:   "obs-studio",
		InstallTimestamp:   time.Date(2026, 5, 30, 0, 0, 0, 0, time.UTC),
		InstallStatus:      state.StatusInstalled,
		VerificationStatus: state.VerificationPassed,
		Version:            "30.0.0",
	}
}

func TestLocalInstallation_JSON_RoundTrip(t *testing.T) {
	original := fullLocalInstallation()

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var restored state.LocalInstallation
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if restored.ApplicationID != original.ApplicationID {
		t.Errorf("ApplicationID mismatch: got %q, want %q", restored.ApplicationID, original.ApplicationID)
	}
	if restored.SourceType != original.SourceType {
		t.Errorf("SourceType mismatch: got %q, want %q", restored.SourceType, original.SourceType)
	}
	if restored.InstallStatus != original.InstallStatus {
		t.Errorf("InstallStatus mismatch: got %q, want %q", restored.InstallStatus, original.InstallStatus)
	}
	if restored.VerificationStatus != original.VerificationStatus {
		t.Errorf("VerificationStatus mismatch: got %q, want %q", restored.VerificationStatus, original.VerificationStatus)
	}
	if restored.Version != original.Version {
		t.Errorf("Version mismatch: got %q, want %q", restored.Version, original.Version)
	}
	if !restored.InstallTimestamp.Equal(original.InstallTimestamp) {
		t.Errorf("InstallTimestamp mismatch: got %v, want %v", restored.InstallTimestamp, original.InstallTimestamp)
	}
}

func TestLocalInstallation_JSON_FieldNames(t *testing.T) {
	r := fullLocalInstallation()

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json.Unmarshal to map failed: %v", err)
	}

	requiredKeys := []string{
		"application_id", "source_type", "source_identifier",
		"install_timestamp", "install_status", "verification_status",
	}
	for _, key := range requiredKeys {
		if _, ok := raw[key]; !ok {
			t.Errorf("expected JSON key %q, not found in: %s", key, string(data))
		}
	}
}

func TestLocalInstallation_YAML_RoundTrip(t *testing.T) {
	original := fullLocalInstallation()

	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("yaml.Marshal failed: %v", err)
	}

	var restored state.LocalInstallation
	if err := yaml.Unmarshal(data, &restored); err != nil {
		t.Fatalf("yaml.Unmarshal failed: %v", err)
	}

	if restored.ApplicationID != original.ApplicationID {
		t.Errorf("ApplicationID mismatch: got %q, want %q", restored.ApplicationID, original.ApplicationID)
	}
	if restored.InstallStatus != original.InstallStatus {
		t.Errorf("InstallStatus mismatch: got %q, want %q", restored.InstallStatus, original.InstallStatus)
	}
	if restored.VerificationStatus != original.VerificationStatus {
		t.Errorf("VerificationStatus mismatch: got %q, want %q", restored.VerificationStatus, original.VerificationStatus)
	}
}
