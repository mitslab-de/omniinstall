# Task 0020 - Engine Integration Test Harness

Status: Open
Priority: P1

## Goal

Create integration-style tests for install/remove/verify flows using fake adapters and temporary state.

## Deliverables

- test harness exercising discovery→resolver→engine paths
- fixtures for success and common failure scenarios
- coverage for install then remove lifecycle

## Acceptance Criteria

- integration tests validate end-to-end MVP flow contracts
- tests remain isolated and deterministic
- go test ./... passes

## References

- specs/09-TESTING_STRATEGY.md
- docs/architecture/01-SYSTEM_OVERVIEW.md
