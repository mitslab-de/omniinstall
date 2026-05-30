# Task 0015 - High-Risk Install Confirmation

Status: Done
Priority: P1

## Goal

Require explicit user confirmation before executing high-risk install plans in CLI execute mode.

## Deliverables

- risk-aware confirmation prompt path in CLI
- non-interactive safety behavior for high-risk plans
- tests for confirm/deny/skip flows

## Acceptance Criteria

- high-risk actions are never silently executed
- user can decline and exit safely
- tests verify guardrail behavior

## References

- docs/architecture/06-SECURITY_MODEL.md
- docs/product/04-UX_PRINCIPLES.md

## Completion Notes

- Added risk-aware confirmation flow to CLI install path in `cmd/omniinstall/app.go`:
  - high-risk/critical plans (`RequiresConfirmation`) now prompt explicitly before execution
  - install proceeds only when user types `yes`
- Added non-interactive safety guard:
  - high-risk installs are blocked when no interactive input is available
- Added tests in `cmd/omniinstall/app_test.go` for:
  - confirmed high-risk install path
  - denied high-risk install path
  - non-interactive high-risk guard path
- Verified with:
  - `gofmt -w $(find . -name '*.go' -not -path './.git/*')`
  - `go test ./...`
  - `go build ./cmd/omniinstall`
