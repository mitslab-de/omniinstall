# Task 0032 - Pacman Adapter Skeleton

Status: Open
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
