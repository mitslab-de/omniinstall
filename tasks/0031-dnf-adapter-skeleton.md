# Task 0031 - DNF Adapter Skeleton

Status: Open
Priority: P2

## Goal

Create a skeleton DNF adapter (Fedora/RHEL family) following the Adapter interface, mirroring the APT adapter structure.

## Deliverables

- `internal/adapters/dnf/dnf.go` with Name, IsAvailable, CanHandle, Install, Remove, Verify, CheckInstalled
- unit tests for IsAvailable and CanHandle logic
- adapter registered in default engine setup

## Acceptance Criteria

- DNF adapter compiles and implements the full adapter.Adapter interface
- IsAvailable returns false when `dnf` is not on PATH (unit-testable via mock exec)
- tests pass

## References

- specs/04-ADAPTER_INTERFACE.md
- internal/adapters/apt/apt.go
