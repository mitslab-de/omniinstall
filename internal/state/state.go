// Package state defines the Local State model for OmniInstall.
//
// Local state tracks what OmniInstall has installed, from which source,
// and the verification status of each installation.
package state

import (
	"errors"
	"time"

	"github.com/mitslab-de/omniinstall/internal/source"
)

// InstallStatus describes the outcome of an installation attempt.
type InstallStatus string

const (
	StatusInstalled InstallStatus = "installed"
	StatusFailed    InstallStatus = "failed"
	StatusRemoved   InstallStatus = "removed"
	StatusPending   InstallStatus = "pending"
)

// VerificationStatus describes the outcome of a post-install verification.
type VerificationStatus string

const (
	VerificationPassed  VerificationStatus = "passed"
	VerificationFailed  VerificationStatus = "failed"
	VerificationSkipped VerificationStatus = "skipped"
	VerificationPending VerificationStatus = "pending"
)

// LocalInstallation is a single local state entry for an installed application.
type LocalInstallation struct {
	// ApplicationID is the canonical OmniInstall application identifier.
	ApplicationID string `json:"application_id" yaml:"application_id"`

	// SourceType identifies the backend used to install the application.
	SourceType source.Type `json:"source_type" yaml:"source_type"`

	// SourceIdentifier is the backend-specific package or application ID.
	SourceIdentifier string `json:"source_identifier" yaml:"source_identifier"`

	// InstallTimestamp is when the installation was performed.
	InstallTimestamp time.Time `json:"install_timestamp" yaml:"install_timestamp"`

	// InstallStatus is the outcome of the last install operation.
	InstallStatus InstallStatus `json:"install_status" yaml:"install_status"`

	// VerificationStatus is the outcome of the last verification check.
	VerificationStatus VerificationStatus `json:"verification_status" yaml:"verification_status"`

	// Version is the installed version (optional).
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
}

// Validate checks that the LocalInstallation has all required fields.
func (r *LocalInstallation) Validate() error {
	if r.ApplicationID == "" {
		return errors.New("application_id is required")
	}
	if r.SourceType == "" {
		return errors.New("source_type is required")
	}
	if r.SourceIdentifier == "" {
		return errors.New("source_identifier is required")
	}
	if r.InstallTimestamp.IsZero() {
		return errors.New("install_timestamp is required")
	}
	if r.InstallStatus == "" {
		return errors.New("install_status is required")
	}
	if r.VerificationStatus == "" {
		return errors.New("verification_status is required")
	}
	return nil
}
