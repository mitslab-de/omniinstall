package source_test

import (
	"testing"

	"github.com/mitslab-de/omniinstall/internal/source"
)

func TestSource_Validate_Valid(t *testing.T) {
	s := &source.Source{
		ApplicationID:    "obs-studio",
		SourceType:       source.TypeFlatpak,
		SourceIdentifier: "com.obsproject.Studio",
		TrustLevel:       source.TrustVerified,
		RiskLevel:        source.RiskLow,
	}
	if err := s.Validate(); err != nil {
		t.Errorf("expected valid source, got error: %v", err)
	}
}

func TestSource_Validate_MissingApplicationID(t *testing.T) {
	s := &source.Source{
		SourceType:       source.TypeAPT,
		SourceIdentifier: "obs-studio",
		TrustLevel:       source.TrustOfficial,
		RiskLevel:        source.RiskLow,
	}
	if err := s.Validate(); err == nil {
		t.Error("expected error for missing application_id, got nil")
	}
}

func TestSource_Validate_InvalidSourceType(t *testing.T) {
	s := &source.Source{
		ApplicationID:    "obs-studio",
		SourceType:       source.Type("invalid"),
		SourceIdentifier: "obs-studio",
		TrustLevel:       source.TrustOfficial,
		RiskLevel:        source.RiskLow,
	}
	if err := s.Validate(); err == nil {
		t.Error("expected error for invalid source_type, got nil")
	}
}

func TestSource_Validate_MissingSourceIdentifier(t *testing.T) {
	s := &source.Source{
		ApplicationID: "obs-studio",
		SourceType:    source.TypeAPT,
		TrustLevel:    source.TrustOfficial,
		RiskLevel:     source.RiskLow,
	}
	if err := s.Validate(); err == nil {
		t.Error("expected error for missing source_identifier, got nil")
	}
}

func TestSource_Validate_InvalidTrustLevel(t *testing.T) {
	s := &source.Source{
		ApplicationID:    "obs-studio",
		SourceType:       source.TypeAPT,
		SourceIdentifier: "obs-studio",
		TrustLevel:       source.TrustLevel("bad"),
		RiskLevel:        source.RiskLow,
	}
	if err := s.Validate(); err == nil {
		t.Error("expected error for invalid trust_level, got nil")
	}
}

func TestSource_Validate_InvalidRiskLevel(t *testing.T) {
	s := &source.Source{
		ApplicationID:    "obs-studio",
		SourceType:       source.TypeAPT,
		SourceIdentifier: "obs-studio",
		TrustLevel:       source.TrustOfficial,
		RiskLevel:        source.RiskLevel("bad"),
	}
	if err := s.Validate(); err == nil {
		t.Error("expected error for invalid risk_level, got nil")
	}
}

func TestSourceTypeConstants(t *testing.T) {
	types := []source.Type{
		source.TypeAPT,
		source.TypeDNF,
		source.TypePacman,
		source.TypeZypper,
		source.TypeFlatpak,
		source.TypeSnap,
		source.TypeAppImage,
		source.TypeVendor,
		source.TypeDirectDownload,
	}
	for _, st := range types {
		s := &source.Source{
			ApplicationID:    "some-app",
			SourceType:       st,
			SourceIdentifier: "some-id",
			TrustLevel:       source.TrustUnknown,
			RiskLevel:        source.RiskLow,
		}
		if err := s.Validate(); err != nil {
			t.Errorf("expected valid source for type %q, got error: %v", st, err)
		}
	}
}
