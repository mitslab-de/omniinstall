# Task 0037 - Install Plan Dry-Run Output Enrichment

Status: Done
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

## Completion Notes

- Added `--dry-run` flag to `install` command in `main.go`; supports `--dry-run` and
  `--dry-run --json` combinations
- Added `App.installDryRun(appID)` — human-readable output showing Adapter, Package ID,
  Risk level, RequiresPrivilege, and Reason without executing any installation
- Added `App.installDryRunJSON(appID)` — JSON output with `application_id`, `display_name`,
  `adapter`, `package_id`, `risk_level`, `requires_privilege`, `explanation`
- Added `installDryRunJSONOutput` struct in `app.go`
- Extended `CLIFixture.run()` dispatch to handle `install-dry-run` and `install-dry-run-json`
- Added 6 tests: ShowsAdapterAndPackageID, ShowsRiskLevel, NoStateChange,
  EmptyAppID_ReturnsError, ContainsRequiredFields, JSON_NoActualInstall
- All 19 packages pass `go test ./...`; `go build ./...` succeeds
