# Task 0038 - State Export Command

Status: Open
Priority: P2

## Goal

Add `omniinstall state export` command that outputs the current local state as JSON or YAML for backup and migration.

## Deliverables

- `omniinstall state export` outputs JSON by default, YAML with --yaml flag
- App.stateExport() method
- main.go dispatcher for `state` subcommand
- tests for export output format

## Acceptance Criteria

- export includes all non-removed installations
- output is valid JSON/YAML
- tests pass

## References

- specs/06-LOCAL_STATE.md
- cmd/omniinstall/app.go
