package resolver_test

import (
	"testing"

	"github.com/mitslab-de/omniinstall/internal/resolver"
	"github.com/mitslab-de/omniinstall/internal/source"
)

// aptSource is a well-trusted native APT source for test fixtures.
var aptSource = source.Source{
	ApplicationID:    "test-app",
	SourceType:       source.TypeAPT,
	SourceIdentifier: "test-app",
	TrustLevel:       source.TrustOfficial,
	RiskLevel:        source.RiskLow,
}

// flatpakSource is a verified Flatpak source for test fixtures.
var flatpakSource = source.Source{
	ApplicationID:    "test-app",
	SourceType:       source.TypeFlatpak,
	SourceIdentifier: "org.test.App",
	TrustLevel:       source.TrustVerified,
	RiskLevel:        source.RiskLow,
}

// snapSource is a snap source for test fixtures.
var snapSource = source.Source{
	ApplicationID:    "test-app",
	SourceType:       source.TypeSnap,
	SourceIdentifier: "test-app",
	TrustLevel:       source.TrustVerified,
	RiskLevel:        source.RiskLow,
}

// blockedSource is a blocked source for test fixtures.
var blockedSource = source.Source{
	ApplicationID:    "test-app",
	SourceType:       source.TypeAPT,
	SourceIdentifier: "bad-package",
	TrustLevel:       source.TrustBlocked,
	RiskLevel:        source.RiskCritical,
}

// sysWithAll returns a SystemContext with APT, Flatpak, and Snap available.
func sysWithAll() resolver.SystemContext {
	return resolver.SystemContext{
		Distribution:      "ubuntu",
		Architecture:      "amd64",
		AvailableManagers: []source.Type{source.TypeAPT, source.TypeFlatpak, source.TypeSnap},
	}
}

// sysAPTOnly returns a SystemContext with only APT available.
func sysAPTOnly() resolver.SystemContext {
	return resolver.SystemContext{
		Distribution:      "ubuntu",
		Architecture:      "amd64",
		AvailableManagers: []source.Type{source.TypeAPT},
	}
}

func TestDefaultResolverReturnsRecommendation(t *testing.T) {
	r := resolver.NewDefaultResolver()
	recs, err := r.Resolve("test-app", []source.Source{aptSource}, sysWithAll())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(recs))
	}
	if recs[0].Plan == nil {
		t.Fatal("expected non-nil plan")
	}
	if recs[0].Plan.SourceType != source.TypeAPT {
		t.Errorf("expected APT source, got %s", recs[0].Plan.SourceType)
	}
	if recs[0].Explanation == "" {
		t.Error("expected non-empty explanation")
	}
}

func TestDefaultResolverNativePreferredOverFlatpak(t *testing.T) {
	r := resolver.NewDefaultResolver()
	recs, err := r.Resolve("test-app", []source.Source{flatpakSource, aptSource}, sysWithAll())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) < 2 {
		t.Fatalf("expected ≥2 recommendations, got %d", len(recs))
	}
	if recs[0].Plan.SourceType != source.TypeAPT {
		t.Errorf("expected APT as first choice, got %s", recs[0].Plan.SourceType)
	}
	if recs[1].Plan.SourceType != source.TypeFlatpak {
		t.Errorf("expected Flatpak as second choice, got %s", recs[1].Plan.SourceType)
	}
}

func TestDefaultResolverFlatpakPreferredWhenNoNative(t *testing.T) {
	r := resolver.NewDefaultResolver()
	ctx := resolver.SystemContext{
		Distribution:      "ubuntu",
		Architecture:      "amd64",
		AvailableManagers: []source.Type{source.TypeFlatpak},
	}
	recs, err := r.Resolve("test-app", []source.Source{aptSource, flatpakSource}, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) != 1 {
		t.Fatalf("expected 1 recommendation (APT unavailable), got %d", len(recs))
	}
	if recs[0].Plan.SourceType != source.TypeFlatpak {
		t.Errorf("expected Flatpak, got %s", recs[0].Plan.SourceType)
	}
}

func TestDefaultResolverBlockedSourceExcluded(t *testing.T) {
	r := resolver.NewDefaultResolver()
	recs, err := r.Resolve("test-app", []source.Source{blockedSource, flatpakSource}, sysWithAll())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, rec := range recs {
		if rec.Plan.SourceIdentifier == "bad-package" {
			t.Error("blocked source should not appear in recommendations")
		}
	}
}

func TestDefaultResolverNoCompatibleSources(t *testing.T) {
	r := resolver.NewDefaultResolver()
	recs, err := r.Resolve("test-app", []source.Source{aptSource}, resolver.SystemContext{
		Distribution:      "ubuntu",
		AvailableManagers: []source.Type{source.TypeFlatpak}, // APT not available
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) != 0 {
		t.Errorf("expected 0 recommendations, got %d", len(recs))
	}
}

func TestDefaultResolverEmptyApplicationID(t *testing.T) {
	r := resolver.NewDefaultResolver()
	_, err := r.Resolve("", []source.Source{aptSource}, sysWithAll())
	if err == nil {
		t.Error("expected error for empty applicationID")
	}
}

func TestDefaultResolverDeterministic(t *testing.T) {
	r := resolver.NewDefaultResolver()
	sources := []source.Source{flatpakSource, aptSource, snapSource}
	ctx := sysWithAll()

	recs1, err := r.Resolve("test-app", sources, ctx)
	if err != nil {
		t.Fatalf("first resolve error: %v", err)
	}
	recs2, err := r.Resolve("test-app", sources, ctx)
	if err != nil {
		t.Fatalf("second resolve error: %v", err)
	}
	if len(recs1) != len(recs2) {
		t.Fatalf("non-deterministic: lengths differ (%d vs %d)", len(recs1), len(recs2))
	}
	for i := range recs1 {
		if recs1[i].Plan.SourceType != recs2[i].Plan.SourceType {
			t.Errorf("non-deterministic at index %d: %s vs %s",
				i, recs1[i].Plan.SourceType, recs2[i].Plan.SourceType)
		}
	}
}

func TestDefaultResolverPreferFlatpakPreference(t *testing.T) {
	r := resolver.NewDefaultResolver()
	ctx := sysWithAll()
	ctx.UserPreferences.PreferFlatpak = true

	recs, err := r.Resolve("test-app", []source.Source{aptSource, flatpakSource}, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) < 2 {
		t.Fatalf("expected ≥2 recommendations, got %d", len(recs))
	}
	if recs[0].Plan.SourceType != source.TypeFlatpak {
		t.Errorf("expected Flatpak as first choice when PreferFlatpak=true, got %s", recs[0].Plan.SourceType)
	}
}

func TestDefaultResolverPreferNativePreference(t *testing.T) {
	r := resolver.NewDefaultResolver()
	ctx := sysWithAll()
	ctx.UserPreferences.PreferNative = true

	recs, err := r.Resolve("test-app", []source.Source{flatpakSource, aptSource}, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) < 1 {
		t.Fatal("expected recommendations")
	}
	if recs[0].Plan.SourceType != source.TypeAPT {
		t.Errorf("expected APT as first choice when PreferNative=true, got %s", recs[0].Plan.SourceType)
	}
}

func TestDefaultResolverPlanIsValid(t *testing.T) {
	r := resolver.NewDefaultResolver()
	recs, err := r.Resolve("test-app", []source.Source{aptSource}, sysWithAll())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) == 0 {
		t.Fatal("expected at least one recommendation")
	}
	if err := recs[0].Plan.Validate(); err != nil {
		t.Errorf("plan validation failed: %v", err)
	}
}

func TestDefaultResolverHighRiskRequiresConfirmation(t *testing.T) {
	r := resolver.NewDefaultResolver()
	highRiskSrc := source.Source{
		ApplicationID:    "risky-app",
		SourceType:       source.TypeAPT,
		SourceIdentifier: "risky-pkg",
		TrustLevel:       source.TrustCommunity,
		RiskLevel:        source.RiskHigh,
	}
	recs, err := r.Resolve("risky-app", []source.Source{highRiskSrc}, sysAPTOnly())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) == 0 {
		t.Fatal("expected a recommendation")
	}
	if !recs[0].Plan.RequiresConfirmation {
		t.Error("high-risk source should require confirmation")
	}
}

func TestDefaultResolverLowRiskDoesNotRequireConfirmation(t *testing.T) {
	r := resolver.NewDefaultResolver()
	recs, err := r.Resolve("test-app", []source.Source{aptSource}, sysAPTOnly())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) == 0 {
		t.Fatal("expected a recommendation")
	}
	if recs[0].Plan.RequiresConfirmation {
		t.Error("low-risk official source should not require confirmation")
	}
}

func TestDefaultResolverNativeRequiresPrivilege(t *testing.T) {
	r := resolver.NewDefaultResolver()
	recs, err := r.Resolve("test-app", []source.Source{aptSource}, sysAPTOnly())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) == 0 {
		t.Fatal("expected a recommendation")
	}
	if !recs[0].Plan.RequiresPrivilege {
		t.Error("native APT install should require privilege")
	}
}

func TestDefaultResolverEmptySources(t *testing.T) {
	r := resolver.NewDefaultResolver()
	recs, err := r.Resolve("test-app", nil, sysWithAll())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) != 0 {
		t.Errorf("expected 0 recommendations for empty sources, got %d", len(recs))
	}
}

func TestDefaultResolverImplementsInterface(t *testing.T) {
	var _ resolver.Resolver = resolver.NewDefaultResolver()
}
