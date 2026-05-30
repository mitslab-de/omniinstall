package security_test

import (
	"testing"

	"github.com/mitslab-de/omniinstall/internal/security"
	"github.com/mitslab-de/omniinstall/internal/source"
)

func TestAssessment_Fields(t *testing.T) {
	a := &security.Assessment{
		RiskLevel:            source.RiskLow,
		RequiresConfirmation: false,
		Explanation:          "Installing from the official system repository is low risk.",
	}
	if a.RiskLevel != source.RiskLow {
		t.Errorf("expected risk level %q, got %q", source.RiskLow, a.RiskLevel)
	}
	if a.RequiresConfirmation {
		t.Error("expected RequiresConfirmation to be false")
	}
	if a.Explanation == "" {
		t.Error("expected non-empty explanation")
	}
}

func TestAssessment_HighRisk(t *testing.T) {
	a := &security.Assessment{
		RiskLevel:            source.RiskHigh,
		RequiresConfirmation: true,
		Explanation:          "This installation requires adding an external repository.",
	}
	if a.RiskLevel != source.RiskHigh {
		t.Errorf("expected risk level %q, got %q", source.RiskHigh, a.RiskLevel)
	}
	if !a.RequiresConfirmation {
		t.Error("expected RequiresConfirmation to be true for high risk")
	}
}
