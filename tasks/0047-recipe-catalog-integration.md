# Task 0047 - Recipe Catalog Integration

Status: Open
Priority: P2

## Goal

Load recipe files from `examples/recipes/` at startup and merge them into the discovery catalog, providing source mappings for catalog applications.

## Deliverables

- recipe loader reads YAML files from a configured recipes directory
- recipe source entries are merged into the catalog application's sources
- App uses merged catalog for search and explain
- tests for recipe loading and catalog merge

## Acceptance Criteria

- recipe sources do not override catalog metadata (ID, DisplayName, etc.)
- duplicate recipe sources are deduplicated
- tests pass

## References

- internal/recipes/recipes.go
- internal/discovery/loader.go
- examples/recipes/
