# Task 0026 - Discovery Search Relevance Tuning

Status: Done
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

## Completion Notes

- Added `normalizeIdentifier()` (+ exported `NormalizeIdentifier()`) to collapse hyphens, underscores, spaces into canonical form for fuzzy comparison.
- Added scoring tiers 85 (normalised exact ID), 75 (normalised exact alias), 65 (normalised ID prefix), 45 (word-segment match for hyphenated IDs).
- All existing scoring tiers preserved with no regressions.
- Added 7 new tests: NormalisedID_HyphenAsSpace, NormalisedAlias_SpaceVsHyphen, WordSegment_Studio, NormalisedID_BeatsPartialMatch, ExactID_BeatsNormalised, NormalizeIdentifier_Various, NormalisedID_BeatsPartialMatch.
- All 14 packages pass `go test ./...`.
