# Task 0034 - Source Type Registration for New Adapters

Status: Open
Priority: P2

## Goal

Add source.TypeDNF, source.TypePacman, and source.TypeSnap to the source type registry so new adapters have valid source type constants.

## Deliverables

- new constants in `internal/source/source.go`
- IsKnownType updated to include new types
- existing tests remain passing

## Acceptance Criteria

- source.TypeDNF, source.TypePacman, source.TypeSnap compile and pass IsKnownType
- no regression in existing source type tests

## References

- specs/02-SOURCE_MODEL.md
- internal/source/source.go
