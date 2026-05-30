# Task 0019 - CLI JSON Output Mode

Status: Open
Priority: P1

## Goal

Add optional JSON output for search, explain, and list commands for scriptable usage.

## Deliverables

- --json flag support for selected commands
- stable output structs and field names
- tests for human vs json output modes

## Acceptance Criteria

- json output is valid and deterministic
- default output remains beginner-friendly text
- tests verify backward compatibility

## References

- specs/08-CLI_SPEC.md
- docs/product/04-UX_PRINCIPLES.md
