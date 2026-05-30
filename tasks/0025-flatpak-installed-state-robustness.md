# Task 0025 - Flatpak Installed-State Robustness

Status: Done
Priority: P2

## Goal

Harden Flatpak adapter installed-state detection across user/system installations.

## Deliverables

- flatpak state detection for user and system scopes
- clear handling when flatpak is unavailable
- expanded adapter tests with sample outputs

## Acceptance Criteria

- flatpak installed-state checks are deterministic and accurate
- error categories are stable
- tests pass

## References

- specs/04-ADAPTER_INTERFACE.md
- docs/architecture/04-INSTALL_ENGINE.md

## Completion Notes

- Replaced single `flatpak info <id>` with dual-scope calls (`--system` then `--user`) in `CheckInstalled`.
- Added `parseFlatpakInfoInstalled()` to validate output by checking for "ID:" / "Application:" header lines; falls back to non-empty output for older flatpak versions.
- Added `classifyFlatpakInfoError()` to surface backend errors ("No installations", "permission denied") as StateUnknown, distinguishing them from normal not-found exits.
- Added 8 new tests: SystemScope_Found, UserScope_Only, NotInAnyScope, NoInstallations, PermissionDenied, EmptyOutput, CommandLaunchError, RealFlatpakInfoFormat.
- All 14 packages pass `go test ./...`.
