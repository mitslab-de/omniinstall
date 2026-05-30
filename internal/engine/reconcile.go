package engine

import (
	"fmt"

	adapter "github.com/mitslab-de/omniinstall/internal/adapters"
	"github.com/mitslab-de/omniinstall/internal/state"
)

// DriftKind describes why a local state record does not match the adapter.
type DriftKind string

const (
	// DriftNone means local state and adapter agree.
	DriftNone DriftKind = "none"

	// DriftMissingFromBackend means local state says installed but the
	// package manager does not know about the package. The application was
	// likely removed outside OmniInstall.
	DriftMissingFromBackend DriftKind = "missing_from_backend"

	// DriftPartialInBackend means local state says installed but the adapter
	// reports a partial/broken installation.
	DriftPartialInBackend DriftKind = "partial_in_backend"

	// DriftPresentButRemoved means local state says removed but the adapter
	// reports the package as still installed.
	DriftPresentButRemoved DriftKind = "present_but_removed"

	// DriftUnknown means the adapter could not determine the installed state
	// (e.g. backend unavailable or returned an error).
	DriftUnknown DriftKind = "unknown"
)

// DriftRecord holds the comparison result for a single local installation.
type DriftRecord struct {
	// Installation is the local state record that was checked.
	Installation state.LocalInstallation

	// DriftKind describes the type of drift detected, or DriftNone when
	// local state and backend agree.
	DriftKind DriftKind

	// BackendState is the adapter-reported installed state.
	BackendState adapter.InstalledState

	// Message is a human-readable description of the drift.
	Message string
}

// Reconcile checks each installation record against the installed state
// reported by the available adapters and returns a DriftRecord for each
// entry. Records whose adapter cannot be found produce DriftUnknown.
//
// Only installations with InstallStatus installed or removed are evaluated;
// entries with failed or pending status are skipped with DriftNone.
func (e *DefaultEngine) Reconcile(installs []state.LocalInstallation) []DriftRecord {
	results := make([]DriftRecord, 0, len(installs))
	for _, inst := range installs {
		results = append(results, e.reconcileOne(inst))
	}
	return results
}

// reconcileOne performs reconciliation for a single installation record.
func (e *DefaultEngine) reconcileOne(inst state.LocalInstallation) DriftRecord {
	// Only check installed and removed statuses; skip pending/failed.
	if inst.InstallStatus != state.StatusInstalled && inst.InstallStatus != state.StatusRemoved {
		return DriftRecord{
			Installation: inst,
			DriftKind:    DriftNone,
			Message:      fmt.Sprintf("%s: skipped (status %s)", inst.ApplicationID, inst.InstallStatus),
		}
	}

	a := e.selectRemovalAdapter(inst.SourceType)
	if a == nil {
		return DriftRecord{
			Installation: inst,
			DriftKind:    DriftUnknown,
			BackendState: adapter.StateUnknown,
			Message: fmt.Sprintf("%s: no available adapter for source type %s",
				inst.ApplicationID, inst.SourceType),
		}
	}

	backendState, err := a.CheckInstalled(inst.SourceIdentifier)
	if err != nil {
		return DriftRecord{
			Installation: inst,
			DriftKind:    DriftUnknown,
			BackendState: adapter.StateUnknown,
			Message: fmt.Sprintf("%s: adapter check failed: %s",
				inst.ApplicationID, err.Error()),
		}
	}

	switch inst.InstallStatus {
	case state.StatusInstalled:
		switch backendState {
		case adapter.StateInstalled:
			return DriftRecord{
				Installation: inst,
				DriftKind:    DriftNone,
				BackendState: backendState,
				Message:      fmt.Sprintf("%s: installed (confirmed)", inst.ApplicationID),
			}
		case adapter.StateNotInstalled:
			return DriftRecord{
				Installation: inst,
				DriftKind:    DriftMissingFromBackend,
				BackendState: backendState,
				Message: fmt.Sprintf("%s: recorded as installed but not found by adapter "+
					"(possibly removed outside OmniInstall)", inst.ApplicationID),
			}
		case adapter.StatePartial:
			return DriftRecord{
				Installation: inst,
				DriftKind:    DriftPartialInBackend,
				BackendState: backendState,
				Message: fmt.Sprintf("%s: partially installed according to adapter "+
					"(run reinstall or remove to fix)", inst.ApplicationID),
			}
		default:
			return DriftRecord{
				Installation: inst,
				DriftKind:    DriftUnknown,
				BackendState: backendState,
				Message:      fmt.Sprintf("%s: adapter returned unknown state %q", inst.ApplicationID, backendState),
			}
		}

	case state.StatusRemoved:
		if backendState == adapter.StateInstalled {
			return DriftRecord{
				Installation: inst,
				DriftKind:    DriftPresentButRemoved,
				BackendState: backendState,
				Message: fmt.Sprintf("%s: recorded as removed but still installed in backend "+
					"(possibly reinstalled outside OmniInstall)", inst.ApplicationID),
			}
		}
		return DriftRecord{
			Installation: inst,
			DriftKind:    DriftNone,
			BackendState: backendState,
			Message:      fmt.Sprintf("%s: removed (confirmed)", inst.ApplicationID),
		}
	}

	return DriftRecord{
		Installation: inst,
		DriftKind:    DriftNone,
		BackendState: backendState,
	}
}
