# Task 0013 - External Local Catalog Loading

Status: Open
Priority: P1

## Goal

Allow discovery engine to load an optional local catalog file while preserving MVP embedded fallback catalog.

## Deliverables

- catalog loader for local YAML/JSON file
- fallback to embedded MVP catalog when file missing/invalid
- tests for success/fallback/error paths

## Acceptance Criteria

- search works with loaded local catalog
- missing catalog file does not break MVP behavior
- tests validate deterministic search behavior

## References

- docs/architecture/02-DISCOVERY_ENGINE.md
- docs/product/05-MVP_DEFINITION.md
