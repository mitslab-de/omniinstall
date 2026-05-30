# Task 0004 - Source Resolver

Status: Done
Priority: P0

## Goal

Implement deterministic source recommendation logic.

## Deliverables

- ranking engine
- explanation generation
- source comparison model
- conflict detection foundation

## Acceptance Criteria

- resolver returns recommended source
- recommendation is explainable
- deterministic behavior verified by tests

## References

- docs/architecture/03-SOURCE_RESOLVER.md
- specs/05-RESOLVER_RULES.md

## Completion Notes

Implemented on 2026-05-30.

- `internal/resolver/ranking.go` — source scoring: type priority (native=50, flatpak=40, vendor=35, snap=30, appimage=25, direct=20) + trust modifier (official+20, verified+15, community+5) + risk modifier (low+10, medium+5, critical-20) + user preference boost (+15 for PreferNative/PreferFlatpak)
- `internal/resolver/explanation.go` — user-facing explanation and alternative explanation generation
- `internal/resolver/conflict.go` — ConflictInfo + DetectConflicts: surfaces blocked sources and unavailable managers
- `internal/resolver/default_resolver.go` — DefaultResolver implementing the Resolver interface; deterministic sort (score desc, identifier asc on tie); builds valid install.Plan per recommendation
- 19 tests covering: ranking, determinism, user preferences, blocked source exclusion, empty inputs, conflict kinds, plan validation
