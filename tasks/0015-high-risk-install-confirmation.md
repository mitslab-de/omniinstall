# Task 0015 - High-Risk Install Confirmation

Status: Open
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
