# Task 0044 - Structured Log Level Control

Status: Open
Priority: P3

## Goal

Add `--verbose` and `--quiet` global flags to control structured log output level.

## Deliverables

- `--verbose` enables debug log output
- `--quiet` suppresses info logs (errors only)
- global flag parsing before command dispatch in main.go
- tests for flag parsing and log level selection

## Acceptance Criteria

- verbose/quiet flags affect logging without changing command output
- flags are accepted before and after command name
- tests pass

## References

- specs/08-CLI_SPEC.md
- internal/logging/logging.go
