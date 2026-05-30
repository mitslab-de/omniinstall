# Task 0027 - State-Engine Reconciliation Checks

Status: Open
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
