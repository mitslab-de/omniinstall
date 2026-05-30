# Task 0036 - Resolver Source Priority Ordering

Status: Done
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

## Completion Notes

- Added `PriorityConfig` type (`map[source.Type]int`) and `DefaultPriorityConfig` var
  to `internal/resolver/ranking.go` — scores: APT/DNF/Pacman/Zypper=50, Flatpak=40,
  Vendor=35, Snap=30, AppImage=25, DirectDownload=20
- `PriorityConfig.priority()` has nil-safe fallback to DefaultPriorityConfig
- Added `WithPriorityConfig()` builder method on `DefaultResolver`
- `scoreSource` now accepts `PriorityConfig` instead of calling `baseTypePriority`
  directly, allowing per-call customisation
- Added 5 new tests: NativeBeforeFlatpak, FlatpakBeforeSnap, WithPriorityConfig_FlatpakOverNative,
  TieBreakByIdentifier, DefaultPriorityConfig_ValuesDefined
- All 35 resolver tests pass; `go test ./...` and `go build ./...` green
