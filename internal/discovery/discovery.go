// Package discovery defines the Discovery Engine interface for OmniInstall.
//
// The Discovery Engine is responsible for resolving a user's search query
// into candidate applications with metadata. It does not install software.
package discovery

import (
	"github.com/mitslab-de/omniinstall/internal/app"
)

// Candidate is a matching application returned by the discovery engine.
type Candidate struct {
	// Application is the matching application metadata.
	Application *app.Application

	// MatchScore is an internal relevance score (higher is more relevant).
	MatchScore int
}

// Engine discovers applications based on user queries.
type Engine interface {
	// Search returns candidates matching the given query.
	// Returns an empty slice when no matches are found.
	Search(query string) ([]*Candidate, error)

	// Lookup returns the application with the exact canonical ID, or nil.
	Lookup(applicationID string) (*app.Application, error)
}
