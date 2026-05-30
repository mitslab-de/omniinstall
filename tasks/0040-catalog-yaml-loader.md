# Task 0040 - Catalog YAML Loader Validation

Status: Open
Priority: P2

## Goal

Strengthen catalog loading from YAML files to validate all fields and produce clear errors for invalid entries.

## Deliverables

- `internal/discovery/loader.go` validates each loaded application entry
- invalid entries produce itemised error messages, not panics
- tests for catalogs with valid, invalid, and mixed entries

## Acceptance Criteria

- load errors include the application ID and offending field
- catalog with any invalid entry returns error, not partial catalog
- tests pass

## References

- internal/discovery/loader.go
- internal/app/app.go Validate()
