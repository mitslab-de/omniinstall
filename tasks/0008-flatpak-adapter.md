# Task 0008 - Flatpak Adapter

Status: Done
Priority: P0

## Goal

Implement Flatpak support.

## Deliverables

- Flatpak adapter
- install support
- remove support
- verification support

## Acceptance Criteria

- adapter passes contract tests
- resolver can target Flatpak sources

## References

- specs/04-ADAPTER_INTERFACE.md

## Completion Notes

Implemented on 2026-05-30.

- `internal/adapters/flatpak/flatpak.go` — Full `Adapter` struct implementing `adapter.Adapter`: IsAvailable (LookPath flatpak), CanHandle (TypeFlatpak), CheckInstalled (flatpak info), Install/Remove (flatpak install/remove --noninteractive -y, with timing), Verify (re-checks flatpak info), error categorization
- `NewWithExecutorForTest` exported for test injection
- 14 tests covering: availability, CanHandle, CheckInstalled, Install success/failure, Remove, Verify success/failure, interface compliance, contract checks
