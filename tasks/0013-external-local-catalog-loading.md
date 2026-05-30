# Task 0013 - External Local Catalog Loading

Status: Done
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

## Completion Notes

- Added `internal/discovery/loader.go` with:
  - `LoadCatalogFile(path)` for strict local YAML/JSON catalog loading.
  - `LoadCatalogWithFallback(path)` to preserve embedded `MVPCatalog()` behavior when local loading fails.
- Supported both local catalog shapes:
  - top-level application list
  - `{ applications: [...] }` envelope
- Wired CLI bootstrap (`cmd/omniinstall/app.go`) to optional `OMNIINSTALL_CATALOG_PATH` and fallback warning output while keeping MVP defaults.
- Added `internal/discovery/loader_test.go` coverage for:
  - YAML and JSON load success
  - fallback on missing file
  - fallback on invalid file
  - strict loader error path
  - deterministic search behavior on loaded local catalog
- Verified with:
  - `gofmt -w $(find . -name '*.go' -not -path './.git/*')`
  - `go test ./...`
  - `go build ./cmd/omniinstall`
