# Task 0011 - Source-Aware Removal Routing

Status: Done
Priority: P0

## Goal

Ensure remove operations target the installed source identifier and source type from local state, not only the application ID.

## Deliverables

- engine remove API supports source-aware removal input
- CLI remove resolves installation record before calling engine
- tests for mixed-source removal selection and identifier handling

## Acceptance Criteria

- removing an app installed from flatpak does not route through apt
- removal uses recorded source_identifier when it differs from application id
- all related unit tests pass

## References

- specs/04-ADAPTER_INTERFACE.md
- specs/06-LOCAL_STATE.md
- docs/architecture/04-INSTALL_ENGINE.md

## Completion Notes

- Updated `internal/engine.Engine` and `internal/engine.DefaultEngine` removal API to require `applicationID`, `sourceType`, and `sourceIdentifier`.
- Removal now selects an adapter by recorded `sourceType` and invokes adapter removal with recorded `sourceIdentifier`, preventing incorrect backend routing.
- Updated `cmd/omniinstall/app.go` remove flow to read local installation state first and reject remove requests for apps not currently recorded as installed.
- Added tests in `internal/engine/default_engine_test.go` and `cmd/omniinstall/app_test.go` for:
  - source-type-aware adapter selection (Flatpak removal does not route through APT)
  - source-identifier pass-through (e.g., `com.obsproject.Studio`)
  - missing-installed-state validation on remove
- Verified with `go test ./...` and `go test ./... -race`.
