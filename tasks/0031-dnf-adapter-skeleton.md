# Task 0031 - DNF Adapter Skeleton

Status: Done
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

## Completion Notes

- Created `internal/adapters/dnf/dnf.go` with full adapter.Adapter implementation: Name, IsAvailable (dnf in PATH), CanHandle (TypeDNF), CheckInstalled (rpm -q + output parsing), Install (dnf install -y), Remove (dnf remove -y), Verify (LookPath for verification commands).
- executor interface and fakeExecutor test double mirror APT adapter pattern.
- Created `internal/adapters/dnf/dnf_test.go` with 15 tests covering all methods.
- Registered DNF adapter in `cmd/omniinstall/app.go` alongside APT and Flatpak.
- All 15 packages pass `go test ./...`.
