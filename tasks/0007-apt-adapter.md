# Task 0007 - APT Adapter

Status: Done
Priority: P0

## Goal

Implement the first native package manager adapter.

## Deliverables

- availability detection
- install support
- remove support
- verification support

## Acceptance Criteria

- adapter passes contract tests
- adapter can execute MVP install plans

## References

- specs/04-ADAPTER_INTERFACE.md

## Completion Notes

Implemented on 2026-05-30.

- `internal/adapters/apt/apt.go` — Full `Adapter` struct implementing `adapter.Adapter`: IsAvailable (LookPath apt-get), CanHandle (TypeAPT), CheckInstalled (dpkg-query), Install/Remove (apt-get with timing), Verify (LookPath per verification rule), error categorization (exit 100→package not found, keyword-based others)
- `NewWithExecutorForTest` exported for test injection of fake executor
- 15 tests covering: availability, CanHandle, CheckInstalled, Install success/failure, Remove success/failure, Verify success/failure, interface compliance, contract checks
