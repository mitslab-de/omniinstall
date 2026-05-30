# Task 0034 - Source Type Registration for New Adapters

Status: Done
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

## Completion Notes

- Source type constants TypeDNF, TypePacman, TypeSnap (and TypeZypper, TypeAppImage, TypeVendor, TypeDirectDownload) were already present in `internal/source/source.go` prior to this task.
- IsKnownType (via validTypes map) already included all these types.
- No code changes needed; task marked Done to reflect existing state.
