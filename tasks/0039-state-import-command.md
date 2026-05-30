# Task 0039 - State Import Command

Status: Open
Priority: P3

## Goal

Add `omniinstall state import <file>` command to restore state from a previously exported file.

## Deliverables

- `omniinstall state import <file>` reads JSON/YAML and merges into local state
- validation: reject malformed or schema-incompatible files
- dry-run mode shows what would be imported without writing
- tests for import with valid, malformed, and conflicting state files

## Acceptance Criteria

- import validates schema_version before writing
- conflicting IDs are handled (existing record wins unless --force)
- tests pass

## References

- specs/06-LOCAL_STATE.md
- internal/state/file_store.go
