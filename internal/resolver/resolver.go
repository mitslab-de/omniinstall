// Package resolver defines the Source Resolver interface for OmniInstall.
//
// The Source Resolver evaluates available sources for an application and
// produces an ordered list of recommendations. It must behave deterministically.
package resolver

import (
	"github.com/mitslab-de/omniinstall/internal/install"
	"github.com/mitslab-de/omniinstall/internal/source"
)

// Recommendation is a resolver output that pairs an install plan with a
// user-facing explanation of why this source was selected.
type Recommendation struct {
	// Plan is the install plan for this recommendation.
	Plan *install.Plan

	// Explanation is a user-facing reason for this recommendation.
	Explanation string
}

// Resolver evaluates available sources for an application and returns
// an ordered list of install plan recommendations.
//
// Implementations must be deterministic: given the same SystemContext
// and available sources, they must always return the same ordered result.
type Resolver interface {
	// Resolve returns an ordered list of recommendations for the application.
	// The first recommendation is the preferred choice.
	// Returns an empty slice if no sources are available.
	Resolve(applicationID string, available []source.Source, ctx SystemContext) ([]Recommendation, error)
}

// SystemContext captures the state of the current Linux system that the
// resolver uses to rank and filter sources.
type SystemContext struct {
	// Distribution is the detected Linux distribution ID (e.g. "ubuntu", "fedora").
	Distribution string

	// Architecture is the system CPU architecture (e.g. "amd64", "arm64").
	Architecture string

	// AvailableManagers is the set of package managers currently installed.
	AvailableManagers []source.Type

	// UserPreferences captures any user-supplied source preferences.
	UserPreferences UserPreferences
}

// UserPreferences holds user-supplied hints that may influence source ranking.
type UserPreferences struct {
	// PreferFlatpak instructs the resolver to prefer Flatpak sources.
	PreferFlatpak bool

	// PreferNative instructs the resolver to prefer the native package manager.
	PreferNative bool
}
