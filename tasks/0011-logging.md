# Task 0011 - Logging

Status: Open
Priority: P0

## Goal

Implement structured logging across the MVP core.

## Deliverables

- logging model
- structured event output
- resolver logging
- install engine logging
- adapter logging hooks

## Acceptance Criteria

- logs include required fields from specs/07-LOGGING_FORMAT.md
- logs do not leak secrets
- tests cover log event creation

## References

- specs/07-LOGGING_FORMAT.md
- docs/architecture/06-SECURITY_MODEL.md
