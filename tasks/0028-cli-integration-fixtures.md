# Task 0028 - CLI Integration Fixtures

Status: Done
Priority: P2

## Goal

Create reusable CLI test fixtures for command routing, output capture, and state setup.

## Deliverables

- shared test helpers for command execution
- reduced duplication in cmd/omniinstall tests
- updated tests using fixtures

## Acceptance Criteria

- CLI tests remain readable and deterministic
- fixture helpers do not change command behavior
- go test ./... passes

## References

- specs/09-TESTING_STRATEGY.md
- docs/architecture/08-COMPONENT_BOUNDARIES.md

## Completion Notes

- Created `cmd/omniinstall/fixtures_test.go` with `CLIFixture` struct, `newFixture()`, `withAPT()`, `withFlatpak()`, `withInstalledAPT()`, `withInstalledFlatpak()`, `withRemovedAPT()`, `withVerifyResult()`, `mustRun()`, `runExpectError()`, `resetOutput()`, `Output()`.
- `CLIFixture.run()` dispatches to the correct App method by command name (search, install, remove, verify, list, explain).
- Included `fixtureAdapter` — a configurable adapter for fixture use that avoids dependence on the `mockAdapter` private fields in `app_test.go`.
- Created `cmd/omniinstall/fixtures_usage_test.go` with 11 tests demonstrating fixture helpers across search, list, verify, remove, flatpak, and output-reset flows.
- All 58 cmd/omniinstall tests pass; `go test ./...` green.
