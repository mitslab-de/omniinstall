package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// stateFileVersion is the schema version written to and expected in state files.
const stateFileVersion = "1.0"

// stateFile is the on-disk format for the local state store.
type stateFile struct {
	// SchemaVersion enables future migrations.
	SchemaVersion string                       `json:"schema_version"`
	Installations map[string]LocalInstallation `json:"installations"`
}

// FileStore persists local installation state to a JSON file on disk.
// It is safe for concurrent use.
type FileStore struct {
	mu   sync.RWMutex
	path string
}

// NewFileStore returns a FileStore that reads and writes to path.
// The file and its parent directories are created on first write.
func NewFileStore(path string) *FileStore {
	return &FileStore{path: path}
}

// DefaultStorePath returns the platform-appropriate default path for the state
// file, following XDG Base Directory conventions.
func DefaultStorePath() string {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			dataHome = filepath.Join(home, ".local", "share")
		}
	}
	if dataHome != "" {
		return filepath.Join(dataHome, "omniinstall", "state.json")
	}
	return filepath.Join(".omniinstall", "state.json")
}

// Record adds or updates the installation record for an application.
func (s *FileStore) Record(installation LocalInstallation) error {
	if err := installation.Validate(); err != nil {
		return fmt.Errorf("invalid installation record: %w", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	sf, err := s.load()
	if err != nil {
		return err
	}
	sf.Installations[installation.ApplicationID] = installation
	return s.save(sf)
}

// Get returns the installation record for the given application ID.
func (s *FileStore) Get(applicationID string) (LocalInstallation, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sf, err := s.load()
	if err != nil {
		return LocalInstallation{}, false, err
	}
	rec, ok := sf.Installations[applicationID]
	return rec, ok, nil
}

// List returns all stored installation records.
func (s *FileStore) List() ([]LocalInstallation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sf, err := s.load()
	if err != nil {
		return nil, err
	}
	records := make([]LocalInstallation, 0, len(sf.Installations))
	for _, v := range sf.Installations {
		records = append(records, v)
	}
	return records, nil
}

// MarkRemoved updates the InstallStatus of the given application to StatusRemoved.
func (s *FileStore) MarkRemoved(applicationID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sf, err := s.load()
	if err != nil {
		return err
	}
	rec, ok := sf.Installations[applicationID]
	if !ok {
		return ErrNotFound
	}
	rec.InstallStatus = StatusRemoved
	sf.Installations[applicationID] = rec
	return s.save(sf)
}

// load reads the state file from disk. If the file does not exist, it returns
// an empty stateFile. Caller must hold at least a read lock.
func (s *FileStore) load() (*stateFile, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return &stateFile{
			SchemaVersion: stateFileVersion,
			Installations: make(map[string]LocalInstallation),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading state file %s: %w", s.path, err)
	}
	var sf stateFile
	if err := json.Unmarshal(data, &sf); err != nil {
		return nil, fmt.Errorf("parsing state file %s: %w", s.path, err)
	}
	if sf.Installations == nil {
		sf.Installations = make(map[string]LocalInstallation)
	}
	return &sf, nil
}

// save writes the state file to disk atomically (via a temp file + rename).
// Caller must hold the write lock.
func (s *FileStore) save(sf *stateFile) error {
	sf.SchemaVersion = stateFileVersion

	data, err := json.MarshalIndent(sf, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding state: %w", err)
	}

	// Ensure parent directory exists.
	if err := os.MkdirAll(filepath.Dir(s.path), 0o750); err != nil {
		return fmt.Errorf("creating state directory: %w", err)
	}

	// Write to a temp file then rename for atomicity.
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return fmt.Errorf("writing state temp file: %w", err)
	}
	if err := os.Rename(tmp, s.path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("committing state file: %w", err)
	}
	return nil
}
