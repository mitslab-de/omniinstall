# Task 0023 - Explicit Verify Command

Status: Open
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
