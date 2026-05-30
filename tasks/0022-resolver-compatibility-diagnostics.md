# Task 0022 - Resolver Compatibility Diagnostics

Status: Done
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

## Completion Notes

### Changes

**`internal/resolver/conflict.go`**
- Added `Suggestion` field to `ConflictInfo` with actionable next-step text
- `ConflictUnavailableManager`: suggests installing the required package manager
- `ConflictBlockedSource`: suggests contacting an administrator
- When ALL sources are excluded and at least one source exists, a `ConflictNoCompatibleSource` summary entry is appended (with suggestion to check installed managers or try an alternative app)

**`cmd/omniinstall/app.go`**
- `explain()`: when no recommendation is found, shows a "Reasons:" section listing each conflict message and "Next step:" suggestion
- `install()`: when no recommendation is found, prints the conflicts and suggestions before returning the error
- `explainJSONConflict` struct gains a `suggestion` field (omitempty)
- `explainJSON()` populates the `suggestion` field from `ConflictInfo.Suggestion`

**`internal/resolver/conflict_test.go`**
- Updated `TestDetectConflictsMultiple` to expect 3 conflicts (2 per-source + 1 summary)
- Added 4 new tests: `NoCompatibleSourceSummary`, `SuggestionsPopulated`, `BlockedSourceSuggestion`, `SomeSources_NoSummaryConflict`

**`cmd/omniinstall/app_test.go`**
- Added `TestAppExplain_NoCompatibleSourceShowsReasons`
- Added `TestAppExplain_JSONNoCompatibleSourceIncludesSuggestion`

All tests pass; `go test ./...` and `go build ./...` succeed.
