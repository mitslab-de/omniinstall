package state_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mitslab-de/omniinstall/internal/source"
	"github.com/mitslab-de/omniinstall/internal/state"
)

// newTempStore creates a FileStore backed by a temp directory.
func newTempStore(t *testing.T) *state.FileStore {
	t.Helper()
	dir := t.TempDir()
	return state.NewFileStore(filepath.Join(dir, "state.json"))
}

// sampleRecord returns a valid LocalInstallation for testing.
func sampleRecord() state.LocalInstallation {
	return state.LocalInstallation{
		ApplicationID:      "obs-studio",
		SourceType:         source.TypeAPT,
		SourceIdentifier:   "obs-studio",
		InstallTimestamp:   time.Now(),
		InstallStatus:      state.StatusInstalled,
		VerificationStatus: state.VerificationPassed,
	}
}

func TestFileStoreRecordAndGet(t *testing.T) {
	s := newTempStore(t)
	rec := sampleRecord()

	if err := s.Record(rec); err != nil {
		t.Fatalf("Record failed: %v", err)
	}

	got, found, err := s.Get("obs-studio")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !found {
		t.Fatal("expected record to be found")
	}
	if got.ApplicationID != rec.ApplicationID {
		t.Errorf("expected ApplicationID=%s, got %s", rec.ApplicationID, got.ApplicationID)
	}
	if got.SourceType != rec.SourceType {
		t.Errorf("expected SourceType=%s, got %s", rec.SourceType, got.SourceType)
	}
}

func TestFileStoreGetNotFound(t *testing.T) {
	s := newTempStore(t)
	_, found, err := s.Get("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Error("expected found=false for nonexistent application")
	}
}

func TestFileStoreList(t *testing.T) {
	s := newTempStore(t)

	r1 := sampleRecord()
	r2 := sampleRecord()
	r2.ApplicationID = "vlc"
	r2.SourceIdentifier = "vlc"

	if err := s.Record(r1); err != nil {
		t.Fatalf("Record r1 failed: %v", err)
	}
	if err := s.Record(r2); err != nil {
		t.Fatalf("Record r2 failed: %v", err)
	}

	list, err := s.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 records, got %d", len(list))
	}
}

func TestFileStoreListEmpty(t *testing.T) {
	s := newTempStore(t)
	list, err := s.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d records", len(list))
	}
}

func TestFileStoreMarkRemoved(t *testing.T) {
	s := newTempStore(t)
	if err := s.Record(sampleRecord()); err != nil {
		t.Fatalf("Record failed: %v", err)
	}

	if err := s.MarkRemoved("obs-studio"); err != nil {
		t.Fatalf("MarkRemoved failed: %v", err)
	}

	got, found, err := s.Get("obs-studio")
	if err != nil || !found {
		t.Fatalf("expected record after MarkRemoved, found=%v err=%v", found, err)
	}
	if got.InstallStatus != state.StatusRemoved {
		t.Errorf("expected StatusRemoved, got %s", got.InstallStatus)
	}
}

func TestFileStoreMarkRemovedNotFound(t *testing.T) {
	s := newTempStore(t)
	err := s.MarkRemoved("nonexistent")
	if err == nil {
		t.Error("expected error when marking nonexistent application as removed")
	}
}

func TestFileStoreRecordOverwrites(t *testing.T) {
	s := newTempStore(t)
	rec := sampleRecord()
	if err := s.Record(rec); err != nil {
		t.Fatalf("first Record failed: %v", err)
	}

	rec.Version = "1.2.3"
	rec.VerificationStatus = state.VerificationFailed
	if err := s.Record(rec); err != nil {
		t.Fatalf("second Record failed: %v", err)
	}

	got, _, _ := s.Get("obs-studio")
	if got.Version != "1.2.3" {
		t.Errorf("expected Version=1.2.3, got %s", got.Version)
	}
	if got.VerificationStatus != state.VerificationFailed {
		t.Errorf("expected VerificationFailed, got %s", got.VerificationStatus)
	}
}

func TestFileStoreInvalidRecordRejected(t *testing.T) {
	s := newTempStore(t)
	err := s.Record(state.LocalInstallation{}) // missing required fields
	if err == nil {
		t.Error("expected error for invalid record, got nil")
	}
}

func TestFileStorePersistsAcrossInstances(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	s1 := state.NewFileStore(path)
	if err := s1.Record(sampleRecord()); err != nil {
		t.Fatalf("Record failed: %v", err)
	}

	// Open a second instance pointing at the same file.
	s2 := state.NewFileStore(path)
	_, found, err := s2.Get("obs-studio")
	if err != nil || !found {
		t.Fatalf("expected record in second instance, found=%v err=%v", found, err)
	}
}

func TestFileStoreStateFileContainsSchemaVersion(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")
	s := state.NewFileStore(path)
	if err := s.Record(sampleRecord()); err != nil {
		t.Fatalf("Record failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading state file: %v", err)
	}
	if !contains(string(data), "schema_version") {
		t.Error("expected state file to contain schema_version field")
	}
}

func TestFileStoreImplementsStore(t *testing.T) {
	var _ state.Store = state.NewFileStore("")
}

func TestDefaultStorePath(t *testing.T) {
	path := state.DefaultStorePath()
	if path == "" {
		t.Error("DefaultStorePath must return a non-empty string")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(substr) == 0 ||
		containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
