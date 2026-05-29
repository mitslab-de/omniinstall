# Task 0003 - Discovery Engine

Status: Open
Priority: P0

## Goal

Implement the first version of the Discovery Engine.

## Deliverables

- local metadata catalog loader
- application search by ID, display name, and alias
- deterministic ranking for MVP
- tests for exact, alias, partial, and unknown searches

## Acceptance Criteria

- search returns normalized Application results
- unknown queries fail with a useful result state
- tests cover the MVP behavior

## References

- docs/architecture/02-DISCOVERY_ENGINE.md
- specs/01-APPLICATION_MODEL.md
