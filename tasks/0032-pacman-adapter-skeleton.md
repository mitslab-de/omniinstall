# Task 0032 - Pacman Adapter Skeleton

Status: Done
Priority: P2

## Goal

Create a skeleton Pacman adapter (Arch Linux family) following the Adapter interface.

## Deliverables

- `internal/adapters/pacman/pacman.go` with full adapter.Adapter implementation
- IsAvailable check via `which pacman` / PATH detection
- unit tests for IsAvailable, CanHandle, and CheckInstalled parsing

## Acceptance Criteria

- Pacman adapter compiles and passes the Adapter interface contract
- CheckInstalled uses `pacman -Q <pkg>` and parses exit code / output
- tests pass

## References

- specs/04-ADAPTER_INTERFACE.md
- internal/adapters/apt/apt.go

## Completion Notes

- Created `internal/adapters/pacman/pacman.go`: full Adapter interface, CheckInstalled uses `pacman -Q` with exit-code and "not found" output parsing.
- Created `internal/adapters/pacman/pacman_test.go` with 15 tests: CanHandle, IsAvailable (present/absent), CheckInstalled (installed/not-installed/unexpected/empty/launch-error), Install (success/not-found), Remove (success/fail), Verify (found/not-found).
- Registered `pacmanadapter.New()` in `cmd/omniinstall/app.go`.
- All 17 packages pass `go test ./...`.
