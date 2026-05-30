// Package source defines the Source domain model for OmniInstall.
//
// A source connects an application identity to an installable artifact
// via a specific package manager or distribution mechanism.
package source

import "errors"

// Type represents a supported installation source type.
type Type string

const (
	TypeAPT            Type = "apt"
	TypeDNF            Type = "dnf"
	TypePacman         Type = "pacman"
	TypeZypper         Type = "zypper"
	TypeFlatpak        Type = "flatpak"
	TypeSnap           Type = "snap"
	TypeAppImage       Type = "appimage"
	TypeVendor         Type = "vendor"
	TypeDirectDownload Type = "direct-download"
)

// TrustLevel represents how trusted a source is.
type TrustLevel string

const (
	TrustOfficial  TrustLevel = "official"
	TrustVerified  TrustLevel = "verified"
	TrustCommunity TrustLevel = "community"
	TrustUnknown   TrustLevel = "unknown"
	TrustBlocked   TrustLevel = "blocked"
)

// RiskLevel represents the risk of installing from a source.
type RiskLevel string

const (
	RiskLow      RiskLevel = "low"
	RiskMedium   RiskLevel = "medium"
	RiskHigh     RiskLevel = "high"
	RiskCritical RiskLevel = "critical"
)

// Source maps an application identity to an installable artifact.
type Source struct {
	// ApplicationID is the canonical OmniInstall application identifier.
	ApplicationID string `json:"application_id" yaml:"application_id"`

	// SourceType identifies the backend mechanism.
	SourceType Type `json:"source_type" yaml:"source_type"`

	// SourceIdentifier is the backend-specific package or application ID.
	SourceIdentifier string `json:"source_identifier" yaml:"source_identifier"`

	// TrustLevel describes how trusted this source is.
	TrustLevel TrustLevel `json:"trust_level" yaml:"trust_level"`

	// RiskLevel describes the risk of using this source.
	RiskLevel RiskLevel `json:"risk_level" yaml:"risk_level"`
}

// validTypes is the set of known source type values.
var validTypes = map[Type]bool{
	TypeAPT:            true,
	TypeDNF:            true,
	TypePacman:         true,
	TypeZypper:         true,
	TypeFlatpak:        true,
	TypeSnap:           true,
	TypeAppImage:       true,
	TypeVendor:         true,
	TypeDirectDownload: true,
}

// validTrustLevels is the set of known trust level values.
var validTrustLevels = map[TrustLevel]bool{
	TrustOfficial:  true,
	TrustVerified:  true,
	TrustCommunity: true,
	TrustUnknown:   true,
	TrustBlocked:   true,
}

// validRiskLevels is the set of known risk level values.
var validRiskLevels = map[RiskLevel]bool{
	RiskLow:      true,
	RiskMedium:   true,
	RiskHigh:     true,
	RiskCritical: true,
}

// Validate checks that the Source has all required fields with valid values.
func (s *Source) Validate() error {
	if s.ApplicationID == "" {
		return errors.New("application_id is required")
	}
	if s.SourceType == "" {
		return errors.New("source_type is required")
	}
	if !validTypes[s.SourceType] {
		return errors.New("unknown source_type: " + string(s.SourceType))
	}
	if s.SourceIdentifier == "" {
		return errors.New("source_identifier is required")
	}
	if s.TrustLevel == "" {
		return errors.New("trust_level is required")
	}
	if !validTrustLevels[s.TrustLevel] {
		return errors.New("unknown trust_level: " + string(s.TrustLevel))
	}
	if s.RiskLevel == "" {
		return errors.New("risk_level is required")
	}
	if !validRiskLevels[s.RiskLevel] {
		return errors.New("unknown risk_level: " + string(s.RiskLevel))
	}
	return nil
}
