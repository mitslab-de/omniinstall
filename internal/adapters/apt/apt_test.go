package apt_test

import (
	"testing"

	"github.com/mitslab-de/omniinstall/internal/adapters/apt"
	"github.com/mitslab-de/omniinstall/internal/source"
)

func TestCanHandle(t *testing.T) {
	if !apt.CanHandle(source.TypeAPT) {
		t.Error("expected APT adapter to handle TypeAPT")
	}
	if apt.CanHandle(source.TypeFlatpak) {
		t.Error("expected APT adapter NOT to handle TypeFlatpak")
	}
	if apt.CanHandle(source.TypeDNF) {
		t.Error("expected APT adapter NOT to handle TypeDNF")
	}
}

func TestName(t *testing.T) {
	if apt.Name == "" {
		t.Error("expected non-empty adapter name")
	}
}
