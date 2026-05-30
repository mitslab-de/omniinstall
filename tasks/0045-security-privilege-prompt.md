# Task 0045 - Security Privilege Escalation Prompt

Status: Open
Priority: P2

## Goal

When an install plan requires privilege escalation (requires_root = true), display a clear warning and require interactive confirmation, matching the high-risk confirmation pattern.

## Deliverables

- App.install() checks plan.RequiresRoot and prompts before executing
- prompt message explains why privilege is needed
- tests for prompt shown, confirmed, and denied

## Acceptance Criteria

- RequiresRoot = true triggers the prompt in --execute mode
- dry-run mode shows warning but does not prompt
- prompt can be bypassed with --yes flag (non-interactive use)
- tests pass

## References

- specs/03-INSTALL_PLAN.md
- docs/architecture/06-SECURITY_MODEL.md
- cmd/omniinstall/app.go
