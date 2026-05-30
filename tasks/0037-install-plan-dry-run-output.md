# Task 0037 - Install Plan Dry-Run Output Enrichment

Status: Open
Priority: P2

## Goal

Improve dry-run output to clearly show what would be installed, which adapter would be used, and what the estimated disk impact is.

## Deliverables

- App.install() dry-run output shows: adapter name, package identifier, estimated size (if known), risk level
- JSON dry-run output includes these fields
- unit tests for enriched dry-run output

## Acceptance Criteria

- dry-run output is human-readable and includes all required fields
- JSON output includes `adapter`, `package_id`, `risk_level`
- no actual installation occurs in dry-run mode

## References

- specs/03-INSTALL_PLAN.md
- specs/08-CLI_SPEC.md
- cmd/omniinstall/app.go
