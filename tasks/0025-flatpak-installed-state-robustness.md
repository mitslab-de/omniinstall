# Task 0025 - Flatpak Installed-State Robustness

Status: Open
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
