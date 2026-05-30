// Package apt provides the APT adapter for OmniInstall.
//
// The APT adapter handles Debian/Ubuntu APT package manager operations.
// It implements the adapter.Adapter interface.
package apt

import (
	"github.com/mitslab-de/omniinstall/internal/source"
)

// Name is the human-readable name for the APT adapter.
const Name = "apt"

// CanHandle reports whether the APT adapter supports the given source type.
func CanHandle(sourceType source.Type) bool {
	return sourceType == source.TypeAPT
}
