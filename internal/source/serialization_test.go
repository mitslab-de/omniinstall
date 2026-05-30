package source_test

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/mitslab-de/omniinstall/internal/source"
)

func fullSource() *source.Source {
	return &source.Source{
		ApplicationID:    "obs-studio",
		SourceType:       source.TypeFlatpak,
		SourceIdentifier: "com.obsproject.Studio",
		TrustLevel:       source.TrustVerified,
		RiskLevel:        source.RiskLow,
	}
}

func TestSource_JSON_RoundTrip(t *testing.T) {
	original := fullSource()

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var restored source.Source
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if restored.ApplicationID != original.ApplicationID {
		t.Errorf("ApplicationID mismatch: got %q, want %q", restored.ApplicationID, original.ApplicationID)
	}
	if restored.SourceType != original.SourceType {
		t.Errorf("SourceType mismatch: got %q, want %q", restored.SourceType, original.SourceType)
	}
	if restored.SourceIdentifier != original.SourceIdentifier {
		t.Errorf("SourceIdentifier mismatch: got %q, want %q", restored.SourceIdentifier, original.SourceIdentifier)
	}
	if restored.TrustLevel != original.TrustLevel {
		t.Errorf("TrustLevel mismatch: got %q, want %q", restored.TrustLevel, original.TrustLevel)
	}
	if restored.RiskLevel != original.RiskLevel {
		t.Errorf("RiskLevel mismatch: got %q, want %q", restored.RiskLevel, original.RiskLevel)
	}
}

func TestSource_JSON_FieldNames(t *testing.T) {
	s := fullSource()

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json.Unmarshal to map failed: %v", err)
	}

	requiredKeys := []string{"application_id", "source_type", "source_identifier", "trust_level", "risk_level"}
	for _, key := range requiredKeys {
		if _, ok := raw[key]; !ok {
			t.Errorf("expected JSON key %q, not found in: %s", key, string(data))
		}
	}
}

func TestSource_YAML_RoundTrip(t *testing.T) {
	original := fullSource()

	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("yaml.Marshal failed: %v", err)
	}

	var restored source.Source
	if err := yaml.Unmarshal(data, &restored); err != nil {
		t.Fatalf("yaml.Unmarshal failed: %v", err)
	}

	if restored.ApplicationID != original.ApplicationID {
		t.Errorf("ApplicationID mismatch: got %q, want %q", restored.ApplicationID, original.ApplicationID)
	}
	if restored.SourceType != original.SourceType {
		t.Errorf("SourceType mismatch: got %q, want %q", restored.SourceType, original.SourceType)
	}
	if restored.TrustLevel != original.TrustLevel {
		t.Errorf("TrustLevel mismatch: got %q, want %q", restored.TrustLevel, original.TrustLevel)
	}
	if restored.RiskLevel != original.RiskLevel {
		t.Errorf("RiskLevel mismatch: got %q, want %q", restored.RiskLevel, original.RiskLevel)
	}
}
