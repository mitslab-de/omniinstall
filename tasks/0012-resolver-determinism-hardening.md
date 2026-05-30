# Task 0012 - Resolver Determinism Hardening

Status: Open
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
