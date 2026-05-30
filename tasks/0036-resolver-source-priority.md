# Task 0036 - Resolver Source Priority Ordering

Status: Open
Priority: P2

## Goal

Implement deterministic source priority ordering in the resolver so that when multiple sources are available for an application, the preferred source is selected consistently.

## Deliverables

- priority ordering: native > flatpak > snap (configurable per distribution)
- Resolver.Resolve() uses priority when multiple sources match
- unit tests for priority ordering and tie-breaking

## Acceptance Criteria

- priority is deterministic and documented
- existing resolver tests remain passing
- priority ordering can be overridden per source type

## References

- specs/05-RESOLVER_RULES.md
- internal/resolver/resolver.go
