package install_test

import (
	"testing"

	"github.com/mitslab-de/omniinstall/internal/install"
	"github.com/mitslab-de/omniinstall/internal/source"
)

func TestPlan_Validate_Valid(t *testing.T) {
	p := &install.Plan{
		ApplicationID:     "obs-studio",
		SourceType:        source.TypeAPT,
		SourceIdentifier:  "obs-studio",
		RiskLevel:         source.RiskLow,
		RequiresPrivilege: true,
		Verification:      []install.VerificationRule{{Command: "obs"}},
	}
	if err := p.Validate(); err != nil {
		t.Errorf("expected valid plan, got error: %v", err)
	}
}

func TestPlan_Validate_MissingApplicationID(t *testing.T) {
	p := &install.Plan{
		SourceType:       source.TypeAPT,
		SourceIdentifier: "obs-studio",
		RiskLevel:        source.RiskLow,
		Verification:     []install.VerificationRule{{Command: "obs"}},
	}
	if err := p.Validate(); err == nil {
		t.Error("expected error for missing application_id, got nil")
	}
}

func TestPlan_Validate_MissingSourceType(t *testing.T) {
	p := &install.Plan{
		ApplicationID:    "obs-studio",
		SourceIdentifier: "obs-studio",
		RiskLevel:        source.RiskLow,
		Verification:     []install.VerificationRule{{Command: "obs"}},
	}
	if err := p.Validate(); err == nil {
		t.Error("expected error for missing source_type, got nil")
	}
}

func TestPlan_Validate_MissingSourceIdentifier(t *testing.T) {
	p := &install.Plan{
		ApplicationID: "obs-studio",
		SourceType:    source.TypeAPT,
		RiskLevel:     source.RiskLow,
		Verification:  []install.VerificationRule{{Command: "obs"}},
	}
	if err := p.Validate(); err == nil {
		t.Error("expected error for missing source_identifier, got nil")
	}
}

func TestPlan_Validate_MissingRiskLevel(t *testing.T) {
	p := &install.Plan{
		ApplicationID:    "obs-studio",
		SourceType:       source.TypeAPT,
		SourceIdentifier: "obs-studio",
		Verification:     []install.VerificationRule{{Command: "obs"}},
	}
	if err := p.Validate(); err == nil {
		t.Error("expected error for missing risk_level, got nil")
	}
}

func TestPlan_Validate_MissingVerification(t *testing.T) {
	p := &install.Plan{
		ApplicationID:    "obs-studio",
		SourceType:       source.TypeAPT,
		SourceIdentifier: "obs-studio",
		RiskLevel:        source.RiskLow,
	}
	if err := p.Validate(); err == nil {
		t.Error("expected error for missing verification checks, got nil")
	}
}
