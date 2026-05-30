# Task 0050 - End-to-End Smoke Test Harness

Status: Open
Priority: P2

## Goal

Create an end-to-end smoke test harness that runs the compiled `omniinstall` binary through core flows (search, explain, install dry-run, list) using a file-system-isolated test environment.

## Deliverables

- `internal/e2e/smoke_test.go` with subprocess-based binary execution tests
- test binary built with `go build` in TestMain
- tests for: search returns results, explain shows metadata, install dry-run succeeds, list returns empty state
- CI-compatible (no root required, no real installs)

## Acceptance Criteria

- tests invoke the binary via exec.Command and assert exit code and stdout
- tests are isolated: each uses a separate temp directory for state
- tests pass on CI without root

## References

- specs/09-TESTING_STRATEGY.md
- docs/product/10-V01-RELEASE-READINESS.md
