// Package flatpak provides the Flatpak adapter for OmniInstall.
//
// The Flatpak adapter handles Flatpak application operations.
// It implements the adapter.Adapter interface.
package flatpak

import (
	"github.com/mitslab-de/omniinstall/internal/source"
)

// Name is the human-readable name for the Flatpak adapter.
const Name = "flatpak"

// CanHandle reports whether the Flatpak adapter supports the given source type.
func CanHandle(sourceType source.Type) bool {
	return sourceType == source.TypeFlatpak
}
