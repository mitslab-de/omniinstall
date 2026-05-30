# Task 0024 - APT Installed-State Robustness

Status: Done
Priority: P2

## Goal

Harden APT adapter installed-state detection across common dpkg/apt output variations.

## Deliverables

- parser improvements for installed-state detection
- error handling for backend command failures
- expanded adapter tests with realistic samples

## Acceptance Criteria

- apt installed-state checks are deterministic on supported formats
- unexpected output paths return actionable errors
- tests pass

## References

- specs/04-ADAPTER_INTERFACE.md
- docs/architecture/04-INSTALL_ENGINE.md

## Completion Notes

- Added `StatePartial InstalledState = "partial"` to `internal/adapters/adapter.go` for half-installed/unpacked/triggers states.
- Replaced `CheckInstalled` substring check with a proper 3-field dpkg status parser `parseDpkgStatus()` in `apt.go`.
- Parser maps: "installed" → StateInstalled; "not-installed"/"config-files" → StateNotInstalled; "half-installed"/"unpacked"/"half-configured"/"triggers-awaited"/"triggers-pending" → StatePartial; empty output → StateNotInstalled; <3 fields or unknown status field → StateUnknown + error.
- Non-zero dpkg-query exit code → StateNotInstalled (package not in dpkg database).
- Added 13 new tests covering all dpkg status variants, malformed output, trailing whitespace, and non-zero exit codes.
- All 14 packages pass `go test ./...`.
