# Task 0019 - CLI JSON Output Mode

Status: Done
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

## Completion Notes

- Added optional `--json` mode to CLI commands `search`, `explain`, and `list` in `cmd/omniinstall/main.go`, including argument parsing and usage validation.
- Implemented stable JSON response payloads in `cmd/omniinstall/app.go`:
  - search: query + ordered result entries
  - explain: application metadata, recommended source, alternatives, and conflicts
  - list: sorted installed applications
- Kept existing human-readable text output as default behavior.
- Added tests for JSON mode in `cmd/omniinstall/app_test.go` and command-level flag handling in `cmd/omniinstall/main_test.go`.
- Verified with `go test ./...` and `go build ./...`.
