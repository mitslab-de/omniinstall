package resolver_test

import (
	"testing"

	"github.com/mitslab-de/omniinstall/internal/resolver"
	"github.com/mitslab-de/omniinstall/internal/source"
)

// Ensure the exported types are usable; this test validates the package compiles
// and that the SystemContext fields are accessible.
func TestSystemContext_Fields(t *testing.T) {
	ctx := resolver.SystemContext{
		Distribution:      "ubuntu",
		Architecture:      "amd64",
		AvailableManagers: []source.Type{source.TypeAPT, source.TypeFlatpak},
		UserPreferences: resolver.UserPreferences{
			PreferNative: true,
		},
	}

	if ctx.Distribution != "ubuntu" {
		t.Errorf("expected distribution 'ubuntu', got %q", ctx.Distribution)
	}
	if ctx.Architecture != "amd64" {
		t.Errorf("expected architecture 'amd64', got %q", ctx.Architecture)
	}
	if len(ctx.AvailableManagers) != 2 {
		t.Errorf("expected 2 available managers, got %d", len(ctx.AvailableManagers))
	}
	if !ctx.UserPreferences.PreferNative {
		t.Error("expected PreferNative to be true")
	}
}

func TestRecommendation_Fields(t *testing.T) {
	r := resolver.Recommendation{
		Explanation: "Native APT package is available and trusted.",
	}
	if r.Explanation == "" {
		t.Error("expected non-empty explanation")
	}
	if r.Plan != nil {
		t.Error("expected nil plan when not set")
	}
}
