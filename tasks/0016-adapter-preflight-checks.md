# Task 0016 - Adapter Preflight Checks

Status: Open
Priority: P1

## Goal

Add standardized preflight checks for adapter availability, permissions, and basic backend readiness.

## Deliverables

- preflight validation before install/remove execution
- clear user-facing errors for failed checks
- unit tests for preflight failure categories

## Acceptance Criteria

- engine surfaces actionable preflight failures
- adapters return stable error categories
- tests pass across adapter and engine packages

## References

- docs/architecture/04-INSTALL_ENGINE.md
- specs/04-ADAPTER_INTERFACE.md
