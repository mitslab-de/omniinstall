// Package engine defines the Install Engine interface for OmniInstall.
//
// The Install Engine takes a validated install plan and executes it through
// the appropriate adapter. It handles progress reporting, logging, and
// verification handoff.
package engine

import (
	"github.com/mitslab-de/omniinstall/internal/install"
)

// Result is the structured outcome of an engine execution.
type Result struct {
	// Success indicates whether the operation completed successfully.
	Success bool

	// ApplicationID is the canonical ID of the application acted on.
	ApplicationID string

	// Message is a user-facing description of what happened.
	Message string

	// ErrorCategory is a stable machine-readable error category, if applicable.
	ErrorCategory string
}

// Engine executes install plans through the appropriate adapter.
type Engine interface {
	// Install executes the given plan and returns a structured result.
	Install(plan *install.Plan) (*Result, error)

	// Remove removes an application by ID and source type.
	Remove(applicationID string) (*Result, error)
}
