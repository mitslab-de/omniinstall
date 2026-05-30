package state

import "errors"

// Store defines the interface for persisting and querying local installation state.
//
// Implementations must be safe for concurrent use.
type Store interface {
	// Record adds or updates the installation record for an application.
	// If a record with the same ApplicationID already exists, it is replaced.
	Record(installation LocalInstallation) error

	// Get returns the installation record for the given application ID.
	// Returns (record, true, nil) if found, (zero, false, nil) if not found.
	Get(applicationID string) (LocalInstallation, bool, error)

	// List returns all stored installation records in undefined order.
	List() ([]LocalInstallation, error)

	// MarkRemoved updates the InstallStatus of the record identified by
	// applicationID to StatusRemoved. Returns an error if no record exists.
	MarkRemoved(applicationID string) error
}

// ErrNotFound is returned when an application ID is not found in the store.
var ErrNotFound = errors.New("application not found in state store")
