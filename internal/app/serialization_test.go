package app_test

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/mitslab-de/omniinstall/internal/app"
)

func fullApplication() *app.Application {
	return &app.Application{
		ID:          "obs-studio",
		DisplayName: "OBS Studio",
		Summary:     "Video recording and live streaming software.",
		Categories:  []string{"video", "streaming"},
		Aliases:     []string{"obs"},
		Homepage:    "https://obsproject.com",
		License:     "GPL-2.0-or-later",
		Publisher:   "OBS Project",
		Description: "Open Broadcaster Software Studio.",
	}
}

func TestApplication_JSON_RoundTrip(t *testing.T) {
	original := fullApplication()

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var restored app.Application
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if restored.ID != original.ID {
		t.Errorf("ID mismatch: got %q, want %q", restored.ID, original.ID)
	}
	if restored.DisplayName != original.DisplayName {
		t.Errorf("DisplayName mismatch: got %q, want %q", restored.DisplayName, original.DisplayName)
	}
	if restored.Summary != original.Summary {
		t.Errorf("Summary mismatch: got %q, want %q", restored.Summary, original.Summary)
	}
	if len(restored.Categories) != len(original.Categories) {
		t.Errorf("Categories length mismatch: got %d, want %d", len(restored.Categories), len(original.Categories))
	}
	if restored.Homepage != original.Homepage {
		t.Errorf("Homepage mismatch: got %q, want %q", restored.Homepage, original.Homepage)
	}
	if restored.License != original.License {
		t.Errorf("License mismatch: got %q, want %q", restored.License, original.License)
	}
}

func TestApplication_JSON_FieldNames(t *testing.T) {
	a := &app.Application{
		ID:          "vlc",
		DisplayName: "VLC",
		Summary:     "Media player.",
		Categories:  []string{"video"},
	}

	data, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json.Unmarshal to map failed: %v", err)
	}

	requiredKeys := []string{"id", "display_name", "summary", "categories"}
	for _, key := range requiredKeys {
		if _, ok := raw[key]; !ok {
			t.Errorf("expected JSON key %q, not found in: %s", key, string(data))
		}
	}

	// Optional fields must be omitted when empty.
	absentKeys := []string{"aliases", "homepage", "license", "publisher", "description"}
	for _, key := range absentKeys {
		if _, ok := raw[key]; ok {
			t.Errorf("expected JSON key %q to be omitted when empty, but it was present", key)
		}
	}
}

func TestApplication_YAML_RoundTrip(t *testing.T) {
	original := fullApplication()

	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("yaml.Marshal failed: %v", err)
	}

	var restored app.Application
	if err := yaml.Unmarshal(data, &restored); err != nil {
		t.Fatalf("yaml.Unmarshal failed: %v", err)
	}

	if restored.ID != original.ID {
		t.Errorf("ID mismatch: got %q, want %q", restored.ID, original.ID)
	}
	if restored.DisplayName != original.DisplayName {
		t.Errorf("DisplayName mismatch: got %q, want %q", restored.DisplayName, original.DisplayName)
	}
	if len(restored.Categories) != len(original.Categories) {
		t.Errorf("Categories length mismatch: got %d, want %d", len(restored.Categories), len(original.Categories))
	}
}

func TestApplication_YAML_FieldNames(t *testing.T) {
	a := &app.Application{
		ID:          "firefox",
		DisplayName: "Firefox",
		Summary:     "Web browser.",
		Categories:  []string{"browsers"},
	}

	data, err := yaml.Marshal(a)
	if err != nil {
		t.Fatalf("yaml.Marshal failed: %v", err)
	}

	yamlStr := string(data)
	requiredKeys := []string{"id:", "display_name:", "summary:", "categories:"}
	for _, key := range requiredKeys {
		found := false
		for i := 0; i <= len(yamlStr)-len(key); i++ {
			if yamlStr[i:i+len(key)] == key {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected YAML key %q in output:\n%s", key, yamlStr)
		}
	}
}
