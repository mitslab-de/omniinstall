package resolver_test

import (
	"testing"

	"github.com/mitslab-de/omniinstall/internal/resolver"
	"github.com/mitslab-de/omniinstall/internal/source"
)

func TestDetectConflictsBlockedSource(t *testing.T) {
	ctx := sysWithAll()
	conflicts := resolver.DetectConflicts([]source.Source{blockedSource}, ctx)
	if len(conflicts) == 0 {
		t.Fatal("expected conflict for blocked source")
	}
	if conflicts[0].Kind != resolver.ConflictBlockedSource {
		t.Errorf("expected ConflictBlockedSource, got %s", conflicts[0].Kind)
	}
	if conflicts[0].Message == "" {
		t.Error("expected non-empty conflict message")
	}
}

func TestDetectConflictsUnavailableManager(t *testing.T) {
	ctx := sysAPTOnly()
	conflicts := resolver.DetectConflicts([]source.Source{flatpakSource}, ctx)
	if len(conflicts) == 0 {
		t.Fatal("expected conflict for unavailable Flatpak manager")
	}
	if conflicts[0].Kind != resolver.ConflictUnavailableManager {
		t.Errorf("expected ConflictUnavailableManager, got %s", conflicts[0].Kind)
	}
}

func TestDetectConflictsNoConflict(t *testing.T) {
	ctx := sysAPTOnly()
	conflicts := resolver.DetectConflicts([]source.Source{aptSource}, ctx)
	if len(conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(conflicts))
	}
}

func TestDetectConflictsMultiple(t *testing.T) {
	ctx := sysAPTOnly()
	conflicts := resolver.DetectConflicts([]source.Source{blockedSource, flatpakSource}, ctx)
	// Expect 3: one per excluded source + the ConflictNoCompatibleSource summary.
	if len(conflicts) != 3 {
		t.Errorf("expected 3 conflicts (2 per-source + 1 summary), got %d", len(conflicts))
	}
	last := conflicts[len(conflicts)-1]
	if last.Kind != resolver.ConflictNoCompatibleSource {
		t.Errorf("expected last conflict to be ConflictNoCompatibleSource, got %s", last.Kind)
	}
}

func TestDetectConflictsEmptySources(t *testing.T) {
	conflicts := resolver.DetectConflicts(nil, sysWithAll())
	if len(conflicts) != 0 {
		t.Errorf("expected 0 conflicts for empty sources, got %d", len(conflicts))
	}
}

func TestDetectConflictsNoCompatibleSourceSummary(t *testing.T) {
// When every source is excluded (manager unavailable), a summary conflict is appended.
ctx := sysAPTOnly()
conflicts := resolver.DetectConflicts([]source.Source{flatpakSource}, ctx)
// Expect 2: one ConflictUnavailableManager + one ConflictNoCompatibleSource.
if len(conflicts) != 2 {
t.Fatalf("expected 2 conflicts, got %d", len(conflicts))
}
last := conflicts[len(conflicts)-1]
if last.Kind != resolver.ConflictNoCompatibleSource {
t.Fatalf("expected ConflictNoCompatibleSource as summary, got %s", last.Kind)
}
if last.Message == "" {
t.Error("expected non-empty message for ConflictNoCompatibleSource")
}
if last.Suggestion == "" {
t.Error("expected actionable suggestion for ConflictNoCompatibleSource")
}
}

func TestDetectConflicts_SuggestionsPopulated(t *testing.T) {
// Each per-source conflict should carry an actionable suggestion.
ctx := sysAPTOnly()
conflicts := resolver.DetectConflicts([]source.Source{flatpakSource}, ctx)
for _, c := range conflicts {
if c.Kind == resolver.ConflictNoCompatibleSource {
continue // summary suggestion checked separately
}
if c.Suggestion == "" {
t.Errorf("expected suggestion for conflict kind %s, got empty", c.Kind)
}
}
}

func TestDetectConflicts_BlockedSourceSuggestion(t *testing.T) {
ctx := sysWithAll()
conflicts := resolver.DetectConflicts([]source.Source{blockedSource}, ctx)
if len(conflicts) == 0 {
t.Fatal("expected at least one conflict for blocked source")
}
if conflicts[0].Suggestion == "" {
t.Error("expected actionable suggestion for blocked source conflict")
}
}

func TestDetectConflicts_SomeSources_NoSummaryConflict(t *testing.T) {
// When at least one source is usable, no ConflictNoCompatibleSource is appended.
ctx := sysAPTOnly()
// One usable (APT), one not (Flatpak).
conflicts := resolver.DetectConflicts([]source.Source{aptSource, flatpakSource}, ctx)
for _, c := range conflicts {
if c.Kind == resolver.ConflictNoCompatibleSource {
t.Error("should not emit ConflictNoCompatibleSource when some sources are usable")
}
}
}
