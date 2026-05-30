# Task 0012 - Resolver Determinism Hardening

Status: Done
Priority: P1

## Goal

Strengthen deterministic resolver behavior for tie-break scenarios and manager-availability permutations.

## Deliverables

- expanded tie-break test matrix
- explicit deterministic ordering contract in resolver docs/tests
- no nondeterministic map iteration in ranking paths

## Acceptance Criteria

- resolver returns stable recommendations across repeated runs
- tests cover identical-score sources with predictable ordering
- go test ./... passes

## References

- specs/05-RESOLVER_RULES.md
- docs/architecture/03-SOURCE_RESOLVER.md

## Completion Notes

- Hardened deterministic tie-breaking in `internal/resolver/default_resolver.go` by ordering equal-score candidates by `SourceIdentifier` and then `SourceType`.
- Added deterministic coverage in `internal/resolver/default_resolver_test.go` for:
  - equal-score sources with identical source identifiers but different source types
  - stable recommendation ordering when available manager list order is permuted
- Preserved existing preference behavior (`PreferFlatpak`) while making tie handling independent of input ordering.
- Verified with `go test ./...` and `go test ./... -race`.
