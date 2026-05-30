// Package security defines the Security Engine interface for OmniInstall.
//
// The Security Engine classifies actions by risk and determines whether
// user confirmation is required before execution.
package security

import (
	"github.com/mitslab-de/omniinstall/internal/install"
	"github.com/mitslab-de/omniinstall/internal/source"
)

// Assessment is the result of a security risk classification.
type Assessment struct {
	// RiskLevel is the classified risk of the action.
	RiskLevel source.RiskLevel

	// RequiresConfirmation indicates whether the user must explicitly confirm.
	RequiresConfirmation bool

	// Explanation is a user-facing risk summary.
	Explanation string
}

// Engine classifies installation actions for safety and trust purposes.
type Engine interface {
	// Assess evaluates the risk of executing the given install plan.
	Assess(plan *install.Plan) (*Assessment, error)
}
