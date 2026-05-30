# Task 0026 - Discovery Search Relevance Tuning

Status: Open
Priority: P2

## Goal

Improve local catalog search ranking for exact id/name matches and aliases.

## Deliverables

- weighted ranking adjustments for exact vs partial matches
- alias normalization tests
- no regression in existing discovery behavior

## Acceptance Criteria

- exact app id/name matches rank highest consistently
- search remains case-insensitive and deterministic
- tests pass

## References

- docs/architecture/02-DISCOVERY_ENGINE.md
- docs/product/04-UX_PRINCIPLES.md
