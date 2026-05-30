# Task 0005 - Install Engine

Status: Done
Priority: P0

## Goal

Implement the Install Engine foundation.

## Deliverables

- install plan validation
- execution orchestration
- verification handoff
- progress event model

## Acceptance Criteria

- install plans can be executed through adapters
- progress events are emitted
- verification integration exists

## References

- docs/architecture/04-INSTALL_ENGINE.md
- specs/03-INSTALL_PLAN.md

## Completion Notes

Implemented on 2026-05-30.

- `internal/engine/progress.go` — ProgressEvent + EventKind constants (started, validating, selecting_adapter, executing, verifying, completed, failed) + ProgressHandler type
- `internal/engine/default_engine.go` — DefaultEngine implementing the Engine interface; validates plan → selects first available+capable adapter → executes → verifies; emits progress events at each stage
- 15 tests covering: success path, nil/invalid plan, no adapter, adapter errors, progress event sequence, verification handoff, skipping unavailable adapters, remove operations
