# Task 0029 - MVP Catalog Expansion

Status: Done
Priority: P2

## Goal

Expand embedded MVP seed catalog to cover at least 10 representative applications.

## Deliverables

- additional validated application entries
- corresponding source mappings for supported managers
- tests ensuring seed catalog validity

## Acceptance Criteria

- catalog includes minimum representative app set per spec
- seed validation remains strict and deterministic
- tests pass

## References

- specs/01-APPLICATION_MODEL.md
- docs/product/05-MVP_DEFINITION.md

## Completion Notes

- Added 3 new applications to `internal/discovery/seed.go`: curl, htop, neovim — bringing total to 10.
- All entries have: ID (kebab-case), DisplayName, Summary, Categories, Aliases, Homepage, License, Publisher, Description.
- Updated `TestMVPCatalog_HasRequiredApplications` to require all 10 IDs.
- Added `TestMVPCatalog_HasMinimumApplicationCount` asserting catalog.All() >= 10.
- All tests pass (`go test ./...`).
