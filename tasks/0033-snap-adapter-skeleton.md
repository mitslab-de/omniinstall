# Task 0033 - Snap Adapter Skeleton

Status: Done
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

## Completion Notes

- Created `internal/adapters/snap/snap.go`: full Adapter interface, CheckInstalled uses `snap list <name>` with exit-code and output parsing for "is not installed", "no snaps are installed", and "not found" patterns.
- Created `internal/adapters/snap/snap_test.go` with 15 tests: CanHandle, IsAvailable (present/absent), CheckInstalled (installed/not-installed variants/unexpected/empty/launch-error), Install (success/not-found), Remove (success), Verify (found/not-found).
- Registered `snapadapter.New()` in `cmd/omniinstall/app.go`.
- All 17 packages pass `go test ./...`.
