package app_test

import (
	"testing"

	"github.com/mitslab-de/omniinstall/internal/app"
)

func TestApplication_Validate_Valid(t *testing.T) {
	a := &app.Application{
		ID:          "obs-studio",
		DisplayName: "OBS Studio",
		Summary:     "Video recording and live streaming software.",
		Categories:  []string{"video", "streaming"},
	}
	if err := a.Validate(); err != nil {
		t.Errorf("expected valid application, got error: %v", err)
	}
}

func TestApplication_Validate_MissingID(t *testing.T) {
	a := &app.Application{
		DisplayName: "OBS Studio",
		Summary:     "Video recording and live streaming software.",
		Categories:  []string{"video"},
	}
	if err := a.Validate(); err == nil {
		t.Error("expected error for missing id, got nil")
	}
}

func TestApplication_Validate_InvalidID(t *testing.T) {
	cases := []string{
		"OBS_Studio",
		"obs studio",
		"obs.studio",
		"ObsStudio",
		"-obs",
		"obs-",
	}
	for _, id := range cases {
		a := &app.Application{
			ID:          id,
			DisplayName: "OBS Studio",
			Summary:     "A summary.",
			Categories:  []string{"video"},
		}
		if err := a.Validate(); err == nil {
			t.Errorf("expected error for invalid id %q, got nil", id)
		}
	}
}

func TestApplication_Validate_MissingDisplayName(t *testing.T) {
	a := &app.Application{
		ID:         "obs-studio",
		Summary:    "A summary.",
		Categories: []string{"video"},
	}
	if err := a.Validate(); err == nil {
		t.Error("expected error for missing display_name, got nil")
	}
}

func TestApplication_Validate_MissingSummary(t *testing.T) {
	a := &app.Application{
		ID:          "obs-studio",
		DisplayName: "OBS Studio",
		Categories:  []string{"video"},
	}
	if err := a.Validate(); err == nil {
		t.Error("expected error for missing summary, got nil")
	}
}

func TestApplication_Validate_MissingCategories(t *testing.T) {
	a := &app.Application{
		ID:          "obs-studio",
		DisplayName: "OBS Studio",
		Summary:     "A summary.",
	}
	if err := a.Validate(); err == nil {
		t.Error("expected error for missing categories, got nil")
	}
}
