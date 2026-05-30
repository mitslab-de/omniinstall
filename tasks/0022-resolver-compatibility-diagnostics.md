# Task 0022 - Resolver Compatibility Diagnostics

Status: Open
Priority: P2

## Goal

Improve resolver diagnostics when no compatible source exists.

## Deliverables

- structured reasons for filtered-out sources
- CLI explain/install surfaces actionable next steps
- tests for no-compatible-source scenarios

## Acceptance Criteria

- users receive clear cause and next action when install cannot proceed
- diagnostics remain deterministic
- tests pass

## References

- docs/architecture/03-SOURCE_RESOLVER.md
- docs/architecture/04-INSTALL_ENGINE.md
