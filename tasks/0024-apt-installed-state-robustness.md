# Task 0024 - APT Installed-State Robustness

Status: Open
Priority: P2

## Goal

Harden APT adapter installed-state detection across common dpkg/apt output variations.

## Deliverables

- parser improvements for installed-state detection
- error handling for backend command failures
- expanded adapter tests with realistic samples

## Acceptance Criteria

- apt installed-state checks are deterministic on supported formats
- unexpected output paths return actionable errors
- tests pass

## References

- specs/04-ADAPTER_INTERFACE.md
- docs/architecture/04-INSTALL_ENGINE.md
