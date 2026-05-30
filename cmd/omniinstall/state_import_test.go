package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mitslab-de/omniinstall/internal/source"
	"github.com/mitslab-de/omniinstall/internal/state"
)

// writeExportFile writes an export JSON file to a temp path and returns it.
func writeExportFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "export.json")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("writeExportFile: %v", err)
	}
	return path
}

// validExportJSON returns a minimal valid export JSON with one installation.
func validExportJSON() string {
	rec := state.LocalInstallation{
		ApplicationID:    "curl",
		SourceType:       source.TypeAPT,
		SourceIdentifier: "curl",
		InstallTimestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		InstallStatus:    state.StatusInstalled,
		VerificationStatus: state.VerificationPassed,
	}
	type exportFmt struct {
		SchemaVersion string                   `json:"schema_version"`
		ExportedAt    string                   `json:"exported_at"`
		Installations []state.LocalInstallation `json:"installations"`
	}
	b, _ := json.Marshal(exportFmt{
		SchemaVersion: "1.0",
		ExportedAt:    "2024-01-01T00:00:00Z",
		Installations: []state.LocalInstallation{rec},
	})
	return string(b)
}

func TestStateImport_ValidFile_ImportsRecord(t *testing.T) {
	f := newFixture(t).withAPT()
	path := writeExportFile(t, validExportJSON())
	a := f.app()
	if err := a.stateImport(path, false, false); err != nil {
		t.Fatalf("stateImport: unexpected error: %v", err)
	}
	f.resetOutput()
	out := f.mustRun(t, "list")
	if !strings.Contains(out, "curl") {
		t.Errorf("expected 'curl' in list after import, got: %s", out)
	}
}

func TestStateImport_EmptyFilePath_ReturnsError(t *testing.T) {
	f := newFixture(t).withAPT()
	a := f.app()
	if err := a.stateImport("", false, false); err == nil {
		t.Fatal("expected error for empty file path")
	}
}

func TestStateImport_MissingFile_ReturnsError(t *testing.T) {
	f := newFixture(t).withAPT()
	a := f.app()
	if err := a.stateImport("/nonexistent/path.json", false, false); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestStateImport_MalformedFile_ReturnsError(t *testing.T) {
	f := newFixture(t).withAPT()
	path := writeExportFile(t, `{not valid json or yaml`)
	a := f.app()
	if err := a.stateImport(path, false, false); err == nil {
		t.Fatal("expected error for malformed file")
	}
}

func TestStateImport_MissingSchemaVersion_ReturnsError(t *testing.T) {
	f := newFixture(t).withAPT()
	path := writeExportFile(t, `{"installations":[]}`)
	a := f.app()
	if err := a.stateImport(path, false, false); err == nil {
		t.Fatal("expected error for missing schema_version")
	}
}

func TestStateImport_UnsupportedMajorVersion_ReturnsError(t *testing.T) {
	f := newFixture(t).withAPT()
	path := writeExportFile(t, `{"schema_version":"2.0","installations":[]}`)
	a := f.app()
	if err := a.stateImport(path, false, false); err == nil {
		t.Fatal("expected error for unsupported schema_version 2.0")
	}
}

func TestStateImport_ConflictWithoutForce_SkipsExisting(t *testing.T) {
	f := newFixture(t).withAPT().withInstalledAPT("curl", "curl")
	path := writeExportFile(t, validExportJSON())
	a := f.app()
	var buf strings.Builder
	a.out = &buf
	if err := a.stateImport(path, false, false); err != nil {
		t.Fatalf("stateImport: unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "Skipping") {
		t.Errorf("expected 'Skipping' message for conflict, got: %s", buf.String())
	}
}

func TestStateImport_DryRun_DoesNotWriteState(t *testing.T) {
	f := newFixture(t).withAPT()
	path := writeExportFile(t, validExportJSON())
	a := f.app()
	if err := a.stateImport(path, false, true); err != nil {
		t.Fatalf("stateImport dry-run: unexpected error: %v", err)
	}
	// State should be empty — dry-run should not persist.
	f.resetOutput()
	listOut := f.mustRun(t, "list")
	if strings.Contains(listOut, "curl") {
		t.Errorf("dry-run import should not write state, but 'curl' appeared in list: %s", listOut)
	}
}
