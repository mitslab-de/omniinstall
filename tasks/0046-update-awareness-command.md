# Task 0046 - Update Awareness Command

Status: Open
Priority: P2

## Goal

Add `omniinstall outdated` command that checks installed applications against available versions and reports which are outdated.

## Deliverables

- `omniinstall outdated` calls adapter's GetVersion() if available and compares with installed version
- App.outdated() method
- main.go dispatcher for "outdated"
- output shows app ID, installed version, available version
- tests for outdated detection

## Acceptance Criteria

- outdated works for adapters that support version checking
- apps with no version info are listed as "unknown" and not reported as outdated
- tests pass

## References

- docs/product/05-MVP_DEFINITION.md (basic update awareness)
- specs/04-ADAPTER_INTERFACE.md
