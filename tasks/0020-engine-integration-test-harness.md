# Task 0020 - Engine Integration Test Harness

Status: Done
Priority: P1

## Goal

Create integration-style tests for install/remove/verify flows using fake adapters and temporary state.

## Deliverables

- test harness exercising discovery→resolver→engine paths
- fixtures for success and common failure scenarios
- coverage for install then remove lifecycle

## Acceptance Criteria

- integration tests validate end-to-end MVP flow contracts
- tests remain isolated and deterministic
- go test ./... passes

## References

- specs/09-TESTING_STRATEGY.md
- docs/architecture/01-SYSTEM_OVERVIEW.md

## Completion Notes

Implemented in `internal/integration/harness_test.go`.

The harness wires together all three pipeline stages using lightweight fakes:

- **fakeAdapter** — configurable in-memory adapter satisfying the `adapter.Adapter` interface; records install/remove calls for assertion.
- **memStore** — in-memory `state.Store` for isolated, deterministic state verification.
- **buildPipeline** helper — creates a `discovery.LocalEngine`, `resolver.DefaultResolver`, and `engine.DefaultEngine` from fixture data.

Tests cover:
- Full APT install flow: discover → resolve → install → state record
- Full Flatpak install flow (Flatpak-only manager context)
- Install-then-remove lifecycle with state transitions
- Resolver source selection: APT preferred over Flatpak by default
- Resolver respects `PreferFlatpak` user preference (Flatpak wins when its combined score exceeds APT)
- No adapter available → `backend_unavailable` result
- Adapter execution error → `execution_failed` result
- Adapter reports package-not-found → correct error category propagated
- Backend preflight failure → `backend_unavailable` result
- No sources for app → empty recommendations, no error
- Lookup miss → nil result, no error
- Progress events: `EventStarted`, `EventExecuting`, `EventCompleted` emitted during full install

All 12 tests pass; `go test ./...` and `go build ./...` succeed.
