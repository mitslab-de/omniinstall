# Task 0011 - Source-Aware Removal Routing

Status: Open
Priority: P0

## Goal

Ensure remove operations target the installed source identifier and source type from local state, not only the application ID.

## Deliverables

- engine remove API supports source-aware removal input
- CLI remove resolves installation record before calling engine
- tests for mixed-source removal selection and identifier handling

## Acceptance Criteria

- removing an app installed from flatpak does not route through apt
- removal uses recorded source_identifier when it differs from application id
- all related unit tests pass

## References

- specs/04-ADAPTER_INTERFACE.md
- specs/06-LOCAL_STATE.md
- docs/architecture/04-INSTALL_ENGINE.md
