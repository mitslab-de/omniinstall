package install_test

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/mitslab-de/omniinstall/internal/install"
	"github.com/mitslab-de/omniinstall/internal/source"
)

func fullPlan() *install.Plan {
	return &install.Plan{
		ApplicationID:        "obs-studio",
		SourceType:           source.TypeAPT,
		SourceIdentifier:     "obs-studio",
		RiskLevel:            source.RiskLow,
		RequiresPrivilege:    true,
		RequiresConfirmation: false,
		Explanation:          "Native APT package is available.",
		Verification:         []install.VerificationRule{{Command: "obs"}},
		ExpectedVersion:      "30.0.0",
	}
}

func TestPlan_JSON_RoundTrip(t *testing.T) {
	original := fullPlan()

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var restored install.Plan
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if restored.ApplicationID != original.ApplicationID {
		t.Errorf("ApplicationID mismatch: got %q, want %q", restored.ApplicationID, original.ApplicationID)
	}
	if restored.SourceType != original.SourceType {
		t.Errorf("SourceType mismatch: got %q, want %q", restored.SourceType, original.SourceType)
	}
	if restored.RiskLevel != original.RiskLevel {
		t.Errorf("RiskLevel mismatch: got %q, want %q", restored.RiskLevel, original.RiskLevel)
	}
	if len(restored.Verification) != len(original.Verification) {
		t.Errorf("Verification length mismatch: got %d, want %d", len(restored.Verification), len(original.Verification))
	}
	if len(restored.Verification) > 0 && restored.Verification[0].Command != original.Verification[0].Command {
		t.Errorf("Verification[0].Command mismatch: got %q, want %q",
			restored.Verification[0].Command, original.Verification[0].Command)
	}
}

func TestPlan_JSON_FieldNames(t *testing.T) {
	p := fullPlan()

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json.Unmarshal to map failed: %v", err)
	}

	requiredKeys := []string{"application_id", "source_type", "source_identifier", "risk_level", "verification"}
	for _, key := range requiredKeys {
		if _, ok := raw[key]; !ok {
			t.Errorf("expected JSON key %q, not found in: %s", key, string(data))
		}
	}
}

func TestPlan_YAML_RoundTrip(t *testing.T) {
	original := fullPlan()

	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("yaml.Marshal failed: %v", err)
	}

	var restored install.Plan
	if err := yaml.Unmarshal(data, &restored); err != nil {
		t.Fatalf("yaml.Unmarshal failed: %v", err)
	}

	if restored.ApplicationID != original.ApplicationID {
		t.Errorf("ApplicationID mismatch: got %q, want %q", restored.ApplicationID, original.ApplicationID)
	}
	if restored.SourceType != original.SourceType {
		t.Errorf("SourceType mismatch: got %q, want %q", restored.SourceType, original.SourceType)
	}
	if len(restored.Verification) != len(original.Verification) {
		t.Errorf("Verification length mismatch: got %d, want %d", len(restored.Verification), len(original.Verification))
	}
}

func TestVerificationRule_JSON(t *testing.T) {
	r := install.VerificationRule{Command: "obs"}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var restored install.VerificationRule
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if restored.Command != r.Command {
		t.Errorf("Command mismatch: got %q, want %q", restored.Command, r.Command)
	}
}
