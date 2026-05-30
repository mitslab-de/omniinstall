package resolver_test

import (
	"testing"

	"github.com/mitslab-de/omniinstall/internal/logging"
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

func TestDefaultResolverDeterministicWithEqualScoreSources(t *testing.T) {
	r := resolver.NewDefaultResolver()
	apt := source.Source{
		ApplicationID:    "equal-app",
		SourceType:       source.TypeAPT,
		SourceIdentifier: "equal-app",
		TrustLevel:       source.TrustOfficial,
		RiskLevel:        source.RiskLow,
	}
	dnf := source.Source{
		ApplicationID:    "equal-app",
		SourceType:       source.TypeDNF,
		SourceIdentifier: "equal-app",
		TrustLevel:       source.TrustOfficial,
		RiskLevel:        source.RiskLow,
	}
	ctx := resolver.SystemContext{
		AvailableManagers: []source.Type{source.TypeDNF, source.TypeAPT},
	}

	first, err := r.Resolve("equal-app", []source.Source{dnf, apt}, ctx)
	if err != nil {
		t.Fatalf("first resolve error: %v", err)
	}
	second, err := r.Resolve("equal-app", []source.Source{apt, dnf}, ctx)
	if err != nil {
		t.Fatalf("second resolve error: %v", err)
	}
	if len(first) != 2 || len(second) != 2 {
		t.Fatalf("expected 2 recommendations in both runs, got %d and %d", len(first), len(second))
	}
	if first[0].Plan.SourceType != source.TypeAPT || second[0].Plan.SourceType != source.TypeAPT {
		t.Fatalf("expected deterministic APT-first tie break, got %s and %s", first[0].Plan.SourceType, second[0].Plan.SourceType)
	}
}

func TestDefaultResolverDeterministicAcrossManagerOrder(t *testing.T) {
	r := resolver.NewDefaultResolver()
	sources := []source.Source{flatpakSource, aptSource, snapSource}

	ctx1 := resolver.SystemContext{AvailableManagers: []source.Type{source.TypeAPT, source.TypeFlatpak, source.TypeSnap}}
	ctx2 := resolver.SystemContext{AvailableManagers: []source.Type{source.TypeSnap, source.TypeAPT, source.TypeFlatpak}}

	recs1, err := r.Resolve("test-app", sources, ctx1)
	if err != nil {
		t.Fatalf("first resolve error: %v", err)
	}
	recs2, err := r.Resolve("test-app", sources, ctx2)
	if err != nil {
		t.Fatalf("second resolve error: %v", err)
	}
	if len(recs1) != len(recs2) {
		t.Fatalf("expected same recommendation count, got %d and %d", len(recs1), len(recs2))
	}
	for i := range recs1 {
		if recs1[i].Plan.SourceType != recs2[i].Plan.SourceType {
			t.Fatalf("expected same source order, index %d differs: %s vs %s", i, recs1[i].Plan.SourceType, recs2[i].Plan.SourceType)
		}
	}
}

func TestDefaultResolverImplementsInterface(t *testing.T) {
	var _ resolver.Resolver = resolver.NewDefaultResolver()
}

func TestDefaultResolverEmitsStructuredLogOnSuccess(t *testing.T) {
	r := resolver.NewDefaultResolver()
	var entries []logging.Entry
	r.WithLogger(logging.EmitFunc(func(entry logging.Entry) {
		entries = append(entries, entry)
	}))

	if _, err := r.Resolve("test-app", []source.Source{aptSource}, sysWithAll()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(entries))
	}
	if entries[0].Action != logging.ActionResolve {
		t.Fatalf("expected action %q, got %q", logging.ActionResolve, entries[0].Action)
	}
	if entries[0].Result != "success" {
		t.Fatalf("expected success result, got %q", entries[0].Result)
	}
	if entries[0].SourceType != source.TypeAPT {
		t.Fatalf("expected source type %q, got %q", source.TypeAPT, entries[0].SourceType)
	}
	if entries[0].Duration <= 0 {
		t.Fatalf("expected positive duration, got %v", entries[0].Duration)
	}
}

func TestDefaultResolverEmitsStructuredLogOnFailure(t *testing.T) {
	r := resolver.NewDefaultResolver()
	var entries []logging.Entry
	r.WithLogger(logging.EmitFunc(func(entry logging.Entry) {
		entries = append(entries, entry)
	}))

	if _, err := r.Resolve("", []source.Source{aptSource}, sysWithAll()); err == nil {
		t.Fatal("expected error for empty applicationID")
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(entries))
	}
	if entries[0].Result != "failure" {
		t.Fatalf("expected failure result, got %q", entries[0].Result)
	}
	if entries[0].ErrorCategory != "invalid_application_id" {
		t.Fatalf("expected invalid_application_id, got %q", entries[0].ErrorCategory)
	}
}

// ---------------------------------------------------------------------------
// Priority ordering and PriorityConfig tests (task 0036)
// ---------------------------------------------------------------------------

func TestDefaultPriorityConfig_NativeBeforeFlatpak(t *testing.T) {
// APT (native, score 50) should beat Flatpak (score 40) when both have
// identical trust and risk levels.
apt := source.Source{
ApplicationID:    "git",
SourceType:       source.TypeAPT,
SourceIdentifier: "git",
TrustLevel:       source.TrustOfficial,
RiskLevel:        source.RiskLow,
}
flatpak := source.Source{
ApplicationID:    "git",
SourceType:       source.TypeFlatpak,
SourceIdentifier: "org.git.Git",
TrustLevel:       source.TrustOfficial,
RiskLevel:        source.RiskLow,
}
r := resolver.NewDefaultResolver()
recs, err := r.Resolve("git", []source.Source{flatpak, apt}, sysWithAll())
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(recs) < 2 {
t.Fatalf("expected 2 recommendations, got %d", len(recs))
}
if recs[0].Plan.SourceType != source.TypeAPT {
t.Errorf("expected APT first (native priority), got %q", recs[0].Plan.SourceType)
}
}

func TestDefaultPriorityConfig_FlatpakBeforeSnap(t *testing.T) {
flatpak := source.Source{
ApplicationID:    "vlc",
SourceType:       source.TypeFlatpak,
SourceIdentifier: "org.videolan.VLC",
TrustLevel:       source.TrustVerified,
RiskLevel:        source.RiskLow,
}
snap := source.Source{
ApplicationID:    "vlc",
SourceType:       source.TypeSnap,
SourceIdentifier: "vlc",
TrustLevel:       source.TrustVerified,
RiskLevel:        source.RiskLow,
}
r := resolver.NewDefaultResolver()
recs, err := r.Resolve("vlc", []source.Source{snap, flatpak}, sysWithAll())
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(recs) < 2 {
t.Fatalf("expected 2 recommendations, got %d", len(recs))
}
if recs[0].Plan.SourceType != source.TypeFlatpak {
t.Errorf("expected Flatpak first (priority over snap), got %q", recs[0].Plan.SourceType)
}
}

func TestWithPriorityConfig_FlatpakOverNative(t *testing.T) {
// Custom config that elevates Flatpak above native.
custom := resolver.PriorityConfig{
source.TypeAPT:     30,
source.TypeFlatpak: 60,
source.TypeSnap:    20,
}
apt := source.Source{
ApplicationID:    "git",
SourceType:       source.TypeAPT,
SourceIdentifier: "git",
TrustLevel:       source.TrustOfficial,
RiskLevel:        source.RiskLow,
}
flatpak := source.Source{
ApplicationID:    "git",
SourceType:       source.TypeFlatpak,
SourceIdentifier: "org.git.Git",
TrustLevel:       source.TrustOfficial,
RiskLevel:        source.RiskLow,
}
r := resolver.NewDefaultResolver().WithPriorityConfig(custom)
recs, err := r.Resolve("git", []source.Source{apt, flatpak}, sysWithAll())
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(recs) < 2 {
t.Fatalf("expected 2 recommendations, got %d", len(recs))
}
if recs[0].Plan.SourceType != source.TypeFlatpak {
t.Errorf("expected Flatpak first with custom config, got %q", recs[0].Plan.SourceType)
}
}

func TestPriorityConfig_TieBreakByIdentifier(t *testing.T) {
// Two APT sources with identical trust/risk — tie-break by SourceIdentifier.
apt1 := source.Source{
ApplicationID:    "curl",
SourceType:       source.TypeAPT,
SourceIdentifier: "curl",
TrustLevel:       source.TrustOfficial,
RiskLevel:        source.RiskLow,
}
apt2 := source.Source{
ApplicationID:    "curl",
SourceType:       source.TypeAPT,
SourceIdentifier: "libcurl",
TrustLevel:       source.TrustOfficial,
RiskLevel:        source.RiskLow,
}
r := resolver.NewDefaultResolver()
recs, err := r.Resolve("curl", []source.Source{apt2, apt1}, sysWithAll())
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(recs) < 2 {
t.Fatalf("expected 2 recommendations, got %d", len(recs))
}
// "curl" sorts before "libcurl".
if recs[0].Plan.SourceIdentifier != "curl" {
t.Errorf("expected tie-break by identifier: 'curl' before 'libcurl', got %q", recs[0].Plan.SourceIdentifier)
}
}

func TestDefaultPriorityConfig_ValuesDefined(t *testing.T) {
cfg := resolver.DefaultPriorityConfig
// Native managers should have the highest default score.
nativeScore := cfg[source.TypeAPT]
if cfg[source.TypeFlatpak] >= nativeScore {
t.Errorf("expected Flatpak score (%d) < APT score (%d)", cfg[source.TypeFlatpak], nativeScore)
}
if cfg[source.TypeSnap] >= cfg[source.TypeFlatpak] {
t.Errorf("expected Snap score (%d) < Flatpak score (%d)", cfg[source.TypeSnap], cfg[source.TypeFlatpak])
}
}
