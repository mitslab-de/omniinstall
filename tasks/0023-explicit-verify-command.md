# Task 0023 - Explicit Verify Command

Status: Done
Priority: P2

## Goal

Expose a dedicated CLI command to re-run verification checks for installed applications.

## Deliverables

- omni verify <application> command path
- engine/adapter verification call integration for existing installs
- tests for verified/failed/not-installed cases

## Acceptance Criteria

- users can verify installed apps without reinstalling
- verification result is reflected in output and local state
- tests pass

## References

- specs/08-CLI_SPEC.md
- docs/architecture/04-INSTALL_ENGINE.md

## Completion Notes

- Added `Verify(applicationID, sourceType, sourceIdentifier)` method to `engine.Engine` interface and implemented in `DefaultEngine` (`internal/engine/default_engine.go`). Emits EventStarted, EventVerifying, EventCompleted/EventFailed and logs with `logging.ActionVerify`.
- Added `App.verify(appID)` in `cmd/omniinstall/app.go`: looks up local state, rejects not-found/removed apps, delegates to engine, updates `VerificationStatus` (VerificationPassed/VerificationFailed) in state, prints ✔/✘ output.
- Added `case "verify":` dispatcher in `cmd/omniinstall/main.go` and updated `printUsage()`.
- Extended `mockAdapter` in `app_test.go` with `verifyResult`/`verifyErr` fields.
- Added 6 tests: Success (state→passed, ✔), Failed (state→failed, ✘), EmptyAppID, NotInstalled, RemovedApp, UpdatesStateOnFailure.
- All 14 packages pass `go test ./...`.
