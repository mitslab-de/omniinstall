# Task 0003 - Discovery Engine

Status: Done
Priority: P0

## Goal

Implement the first version of the Discovery Engine.

## Deliverables

- local metadata catalog loader
- application search by ID, display name, and alias
- deterministic ranking for MVP
- tests for exact, alias, partial, and unknown searches

## Acceptance Criteria

- search returns normalized Application results
- unknown queries fail with a useful result state
- tests cover the MVP behavior

## References

- docs/architecture/02-DISCOVERY_ENGINE.md
- specs/01-APPLICATION_MODEL.md

## Completion Notes

Implemented on 2026-05-30.

- `internal/discovery/catalog.go` — thread-safe in-memory Catalog with Add/Get/All/Size
- `internal/discovery/engine.go` — LocalEngine implementing the Engine interface with deterministic scoring
- `internal/discovery/seed.go` — MVPCatalog() with 7 seed applications (obs-studio, vlc, firefox, bitwarden, git, docker, visual-studio-code)

Search scoring tiers: exact ID (100) → exact display name (90) → exact alias (80) → ID prefix (70) → display name prefix (60) → alias prefix (50) → ID contains (40) → display name contains (30) → alias contains (20). Tied scores break by ID alphabetically for determinism. Empty/whitespace queries return an error. Unknown queries return an empty slice. All 23 tests pass.
