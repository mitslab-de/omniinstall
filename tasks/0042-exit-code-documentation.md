# Task 0042 - Exit Code Documentation and Test Coverage

Status: Open
Priority: P3

## Goal

Document all CLI exit codes in a spec file and add tests to verify each exit code is returned under the expected conditions.

## Deliverables

- `docs/product/11-EXIT_CODES.md` listing all exit codes with conditions
- tests for each exit code: 0 (success), 1 (error), 2 (usage), 3 (backend-unavailable), 4 (not-found)
- verify exit codes match current main.go implementation

## Acceptance Criteria

- all exit codes documented with triggering conditions
- tests cover each exit code path
- tests pass

## References

- specs/08-CLI_SPEC.md
- cmd/omniinstall/main.go exitCodeForError()
