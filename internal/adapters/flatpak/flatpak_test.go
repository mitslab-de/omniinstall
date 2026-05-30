package flatpak_test

import (
	"testing"

	"github.com/mitslab-de/omniinstall/internal/adapters/flatpak"
	"github.com/mitslab-de/omniinstall/internal/source"
)

func TestCanHandle(t *testing.T) {
	if !flatpak.CanHandle(source.TypeFlatpak) {
		t.Error("expected Flatpak adapter to handle TypeFlatpak")
	}
	if flatpak.CanHandle(source.TypeAPT) {
		t.Error("expected Flatpak adapter NOT to handle TypeAPT")
	}
	if flatpak.CanHandle(source.TypeDNF) {
		t.Error("expected Flatpak adapter NOT to handle TypeDNF")
	}
}

func TestName(t *testing.T) {
	if flatpak.Name == "" {
		t.Error("expected non-empty adapter name")
	}
}
