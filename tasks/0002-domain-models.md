# Task 0002 - Domain Models

Status: Done
Priority: P0

## Goal

Implement core domain models.

## Models

- Application
- Source
- InstallPlan
- VerificationRule
- LocalInstallation

## Acceptance Criteria

- models implemented
- serialization supported
- validation implemented
- tests added

## References

- specs/01-APPLICATION_MODEL.md
- specs/02-SOURCE_MODEL.md
- specs/03-INSTALL_PLAN.md

## Completion Note

Implemented on 2026-05-30.

All five domain models are fully implemented with validation and serialization:

- `Application` (`internal/app/app.go`) — user-facing application identity with required field validation
- `Source` (`internal/source/source.go`) — maps application to installable artifact with type/trust/risk validation
- `Plan` / `VerificationRule` (`internal/install/plan.go`) — immutable install plan with verification rules
- `LocalInstallation` (`internal/state/state.go`) — local state record tracking what was installed and verified

All models have JSON and YAML struct tags with snake_case field names matching the specs. Optional fields use `omitempty`. Serialization round-trip tests added for all models. All tests pass.
