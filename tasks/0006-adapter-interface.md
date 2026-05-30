# Task 0006 - Adapter Interface

Status: Done
Priority: P0

## Goal

Implement the adapter contract used by backend integrations.

## Deliverables

- adapter interface
- adapter result model
- verification result model
- adapter contract tests

## Acceptance Criteria

- install engine depends on adapter interface only
- adapters are swappable

## References

- specs/04-ADAPTER_INTERFACE.md

## Completion Notes

Implemented on 2026-05-30.

- `internal/adapters/adapter.go` — Adapter interface (Name, IsAvailable, CanHandle, CheckInstalled, Install, Remove, Verify); Result model updated with `Duration time.Duration` field per spec; VerificationResult; error category constants; InstalledState constants
- `internal/adapters/contract.go` — documents behavioral guarantees via ContractRequirements comment
- `internal/adapters/contract_test.go` — RunAdapterContract() reusable contract test helper; TestAdapterContractWithMock proves the helper works; Install Engine already depends only on the Adapter interface
