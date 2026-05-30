# Task 0014 - CLI Exit Code Contract

Status: Open
Priority: P1

## Goal

Define and implement stable CLI exit codes for common failure categories.

## Deliverables

- exit code mapping for usage, not-found, backend-unavailable, execution-failed
- tests for command exit behavior
- documented mapping in CLI docs/spec

## Acceptance Criteria

- commands exit non-zero on failures with stable code mapping
- usage errors return a distinct exit code
- tests cover success and failure paths

## References

- specs/08-CLI_SPEC.md
- docs/architecture/08-COMPONENT_BOUNDARIES.md
