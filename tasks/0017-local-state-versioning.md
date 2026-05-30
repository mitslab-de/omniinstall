# Task 0017 - Local State Versioning

Status: Done
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

## Completion Notes

- Added explicit schema version validation in `internal/state/file_store.go`:
  - accepts legacy files without `schema_version` (backward compatibility)
  - accepts compatible `1.x` schema versions
  - rejects invalid formats and unsupported major versions with clear errors
- Added migration-boundary and version parsing tests in `internal/state/file_store_test.go`:
  - legacy unversioned file loading
  - compatible minor version loading (`1.5`)
  - unsupported major version rejection (`2.0`)
  - invalid version format rejection (`v1`)
- Updated `specs/06-LOCAL_STATE.md` with explicit `schema_version` behavior for writers/loaders.
