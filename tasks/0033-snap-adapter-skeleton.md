# Task 0033 - Snap Adapter Skeleton

Status: Open
Priority: P2

## Goal

Create a skeleton Snap adapter following the Adapter interface.

## Deliverables

- `internal/adapters/snap/snap.go` with full adapter.Adapter implementation
- IsAvailable checks for `snap` command on PATH
- CheckInstalled uses `snap list <name>` and parses output
- unit tests for IsAvailable, CanHandle, CheckInstalled

## Acceptance Criteria

- Snap adapter compiles and implements adapter.Adapter
- CheckInstalled correctly distinguishes installed vs not-installed vs error
- tests pass

## References

- specs/04-ADAPTER_INTERFACE.md
- internal/adapters/flatpak/flatpak.go
