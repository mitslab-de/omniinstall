# Task 0028 - CLI Integration Fixtures

Status: Open
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
