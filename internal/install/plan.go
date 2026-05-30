// Package install defines the Install Plan model for OmniInstall.
//
// An install plan is an immutable execution instruction produced by the
// Source Resolver and consumed by the Install Engine.
package install

import (
	"errors"

	"github.com/mitslab-de/omniinstall/internal/source"
)

// VerificationCheck describes a post-installation verification step.
type VerificationCheck struct {
	// Command is a binary name that should be findable after installation.
	Command string
}

// Plan is an immutable, explicit instruction for installing one application.
// The Install Engine must reject any plan that does not pass Validate().
type Plan struct {
	// ApplicationID is the canonical OmniInstall application identifier.
	ApplicationID string

	// SourceType identifies the backend to use.
	SourceType source.Type

	// SourceIdentifier is the backend-specific package or application ID.
	SourceIdentifier string

	// RiskLevel reflects the assessed risk of this installation.
	RiskLevel source.RiskLevel

	// RequiresPrivilege indicates whether elevated privileges are needed.
	RequiresPrivilege bool

	// RequiresConfirmation indicates whether the user must explicitly confirm.
	RequiresConfirmation bool

	// Explanation is a user-facing description of why this plan was created.
	Explanation string

	// Verification is the list of checks to perform after installation.
	Verification []VerificationCheck

	// RepositoryChanges describes any repository or PPA additions required.
	RepositoryChanges []string

	// ExpectedVersion is the version expected after installation (optional).
	ExpectedVersion string
}

// validPlanTypes is the set of source types accepted in an install plan.
var validPlanTypes = map[source.Type]bool{
	source.TypeAPT:            true,
	source.TypeDNF:            true,
	source.TypePacman:         true,
	source.TypeZypper:         true,
	source.TypeFlatpak:        true,
	source.TypeSnap:           true,
	source.TypeAppImage:       true,
	source.TypeVendor:         true,
	source.TypeDirectDownload: true,
}

// validRiskLevels is the set of accepted risk level values.
var validRiskLevels = map[source.RiskLevel]bool{
	source.RiskLow:      true,
	source.RiskMedium:   true,
	source.RiskHigh:     true,
	source.RiskCritical: true,
}

// Validate ensures the Plan contains all required fields with valid values.
// The Install Engine must not execute a plan that does not pass this check.
func (p *Plan) Validate() error {
	if p.ApplicationID == "" {
		return errors.New("application_id is required")
	}
	if p.SourceType == "" {
		return errors.New("source_type is required")
	}
	if !validPlanTypes[p.SourceType] {
		return errors.New("unknown source_type: " + string(p.SourceType))
	}
	if p.SourceIdentifier == "" {
		return errors.New("source_identifier is required")
	}
	if p.RiskLevel == "" {
		return errors.New("risk_level is required")
	}
	if !validRiskLevels[p.RiskLevel] {
		return errors.New("unknown risk_level: " + string(p.RiskLevel))
	}
	if len(p.Verification) == 0 {
		return errors.New("at least one verification check is required")
	}
	return nil
}
