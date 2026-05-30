package discovery_test

import (
	"testing"

	"github.com/mitslab-de/omniinstall/internal/app"
	"github.com/mitslab-de/omniinstall/internal/discovery"
)

func TestCandidate_Fields(t *testing.T) {
	a := &app.Application{
		ID:          "obs-studio",
		DisplayName: "OBS Studio",
		Summary:     "Video recording and live streaming software.",
		Categories:  []string{"video"},
	}
	c := &discovery.Candidate{
		Application: a,
		MatchScore:  95,
	}
	if c.Application == nil {
		t.Error("expected non-nil Application")
	}
	if c.Application.ID != "obs-studio" {
		t.Errorf("expected application id 'obs-studio', got %q", c.Application.ID)
	}
	if c.MatchScore != 95 {
		t.Errorf("expected match score 95, got %d", c.MatchScore)
	}
}
