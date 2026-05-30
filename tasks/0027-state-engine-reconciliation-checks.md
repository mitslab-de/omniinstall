# Task 0027 - State-Engine Reconciliation Checks

Status: Done
Priority: P2

## Goal

Add reconciliation checks between local state and adapter-reported installed state.

## Deliverables

- reconciliation helper for list/verify paths
- status transitions for drifted records
- tests for stale and missing backend installations

## Acceptance Criteria

- local state drift can be detected and surfaced clearly
- state updates remain safe and explicit
- tests pass

## References

- specs/06-LOCAL_STATE.md
- docs/architecture/04-INSTALL_ENGINE.md

## Completion Notes

- Created `internal/engine/reconcile.go` with `DriftKind` type (DriftNone, DriftMissingFromBackend, DriftPartialInBackend, DriftPresentButRemoved, DriftUnknown) and `DriftRecord` struct.
- `DefaultEngine.Reconcile([]state.LocalInstallation) []DriftRecord` checks each record via the appropriate adapter's `CheckInstalled`, handles no-adapter, adapter-error, and status-based skipping (pending/failed → DriftNone).
- Status transitions: installed+not-installed→DriftMissingFromBackend; installed+partial→DriftPartialInBackend; removed+installed→DriftPresentButRemoved.
- Created `internal/engine/reconcile_test.go` with 12 tests: Confirmed, MissingFromBackend, PartialInBackend, Removed_Confirmed, Removed_PresentInBackend, NoAdapter, AdapterError, PendingSkipped, FailedSkipped, EmptyList, MultipleRecords, MessageNonEmpty.
- All 14 packages pass `go test ./...`.
