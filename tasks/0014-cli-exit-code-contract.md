# Task 0014 - CLI Exit Code Contract

Status: Done
Priority: P1

## Goal

Define and implement stable CLI exit codes for common failure categories.

## Deliverables

- exit code mapping for usage, not-found, backend-unavailable, execution-failed
- tests for command exit behavior
- documented mapping in CLI docs/spec

## Acceptance Criteria

- commands exit non-zero on failures with stable code mapping
- usage errors return a distinct exit code
- tests cover success and failure paths

## References

- specs/08-CLI_SPEC.md
- docs/architecture/08-COMPONENT_BOUNDARIES.md

## Completion Notes

- Added stable CLI exit code contract in `cmd/omniinstall/main.go`:
  - `0` success
  - `2` usage error
  - `3` not found
  - `4` backend unavailable
  - `5` execution failed
- Updated CLI entrypoint to use `runCLI(args)` and map command errors via `exitCodeForError`.
- Added tests in `cmd/omniinstall/main_test.go` for:
  - success exit code (`version`)
  - usage exit code (`search` without query)
  - not-found exit code (`install` unknown app)
  - backend-unavailable mapping
  - execution-failed default mapping
- Documented exit code mapping in `specs/08-CLI_SPEC.md`.
- Verified with:
  - `gofmt -w $(find . -name '*.go' -not -path './.git/*')`
  - `go test ./...`
  - `go build ./cmd/omniinstall`
