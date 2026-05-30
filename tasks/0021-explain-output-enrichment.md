# Task 0021 - Explain Output Enrichment

Status: Done
Priority: P2

## Goal

Enhance explain output with trust level, risk level, and compatibility rationale for top candidates.

## Deliverables

- explain output includes trust/risk and compatibility signals
- format remains clear for beginners
- tests for enriched output content

## Acceptance Criteria

- explain command communicates why a source is recommended
- output remains deterministic and readable
- tests pass

## References

- docs/product/04-UX_PRINCIPLES.md
- docs/architecture/03-SOURCE_RESOLVER.md

## Completion Notes

### Changes

**`internal/resolver/explanation.go`**
- Enriched `explain()` to include risk level wording ("with low installation risk", etc.)
- Joined explanation parts with `, ` instead of ` ` for clarity
- High/critical risk levels now include explicit warning phrases

**`cmd/omniinstall/app.go`**
- `explain()` text output now shows labeled `Trust:`, `Risk:`, and `Privilege:` fields for the top recommendation
- Recommended line now shows `source_type (identifier)` instead of just `source_type`
- Added `findSource()` helper to look up the source for metadata display
- `explainJSONOutput` now includes `trust_level`, `risk_level`, and `requires_privilege` fields
- `explainJSON()` populates all three new fields from the resolved plan and matching source

### Tests

Added 5 new tests in `cmd/omniinstall/app_test.go`:
- `TestAppExplain_EnrichedTextOutputShowsTrustAndRisk`
- `TestAppExplain_EnrichedTextOutputShowsPrivilege`
- `TestAppExplain_EnrichedTextOutputShowsIdentifier`
- `TestAppExplain_JSONEnrichedOutput`
- `TestAppExplain_ResolverExplanationContainsRiskAndTrust`

All tests pass; `go test ./...` and `go build ./...` succeed.
