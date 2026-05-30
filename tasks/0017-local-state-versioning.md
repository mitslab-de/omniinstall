# Task 0017 - Local State Versioning

Status: Open
Priority: P1

## Goal

Introduce explicit versioning for local state file schema to support future migrations safely.

## Deliverables

- state schema_version field and parser support
- backward-compatible load behavior for current files
- tests for version parsing and unknown-version handling

## Acceptance Criteria

- state store can read current and versioned state files
- unknown incompatible versions fail with clear errors
- tests cover migration boundary behavior

## References

- specs/06-LOCAL_STATE.md
- docs/architecture/07-DATA_MODEL.md
