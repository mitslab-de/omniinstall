# Task 0043 - Progress Event Streaming to Output

Status: Open
Priority: P2

## Goal

Ensure install/remove/verify progress events (EventStarted, EventCompleted, EventFailed) are surfaced to the CLI output so users see real-time feedback.

## Deliverables

- App.install(), App.remove(), App.verify() consume the progress channel and print status lines
- progress lines are suppressed in --json mode
- tests verifying progress output appears in human-readable mode

## Acceptance Criteria

- human-readable output includes at least "Installing..." and "Installed" lines
- JSON output does not include progress lines
- tests pass

## References

- internal/engine/progress.go
- cmd/omniinstall/app.go
