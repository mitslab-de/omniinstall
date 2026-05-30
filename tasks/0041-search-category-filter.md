# Task 0041 - Search Category Filter

Status: Open
Priority: P2

## Goal

Add optional `--category <name>` flag to `omniinstall search` to filter results by application category.

## Deliverables

- `omniinstall search <query> --category development` filters results
- App.search() accepts optional category parameter
- main.go flag parsing for --category
- tests for filtered and unfiltered search

## Acceptance Criteria

- category filter is case-insensitive
- empty category argument returns all results
- filtered results still respect relevance scoring
- tests pass

## References

- specs/08-CLI_SPEC.md
- internal/discovery/engine.go
