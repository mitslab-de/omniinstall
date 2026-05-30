# Task 0016 - Adapter Preflight Checks

Status: Done
Priority: P1

## Goal

Add standardized preflight checks for adapter availability, permissions, and basic backend readiness.

## Deliverables

- preflight validation before install/remove execution
- clear user-facing errors for failed checks
- unit tests for preflight failure categories

## Acceptance Criteria

- engine surfaces actionable preflight failures
- adapters return stable error categories
- tests pass across adapter and engine packages

## References

- docs/architecture/04-INSTALL_ENGINE.md
- specs/04-ADAPTER_INTERFACE.md

## Completion Notes

- Added standardized engine preflight checks in `internal/engine/default_engine.go` before install/remove execution:
  - validates adapter availability
  - runs backend readiness probe via `CheckInstalled`
- Added clear preflight failure surfacing with stable categories:
  - `permission_denied` for permission-related readiness failures
  - `backend_unavailable` for other readiness/backend failures
- Updated adapter readiness behavior:
  - `internal/adapters/apt/apt.go` now returns `StateUnknown` + error when dpkg readiness probe fails unexpectedly
  - `internal/adapters/flatpak/flatpak.go` now returns `StateUnknown` + error when flatpak readiness probe fails unexpectedly
- Added unit tests for preflight categories and readiness errors:
  - `internal/engine/default_engine_test.go`
  - `internal/adapters/apt/apt_test.go`
  - `internal/adapters/flatpak/flatpak_test.go`
- Verified with:
  - `gofmt -w $(find . -name '*.go' -not -path './.git/*')`
  - `go test ./...`
  - `go build ./cmd/omniinstall`
