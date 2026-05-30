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
	if len(conflicts) != 2 {
		t.Errorf("expected 2 conflicts, got %d", len(conflicts))
	}
}

func TestDetectConflictsEmptySources(t *testing.T) {
	conflicts := resolver.DetectConflicts(nil, sysWithAll())
	if len(conflicts) != 0 {
		t.Errorf("expected 0 conflicts for empty sources, got %d", len(conflicts))
	}
}
