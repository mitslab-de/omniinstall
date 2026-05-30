# Task 0049 - Shell Completion Support

Status: Open
Priority: P3

## Goal

Add `omniinstall completion <shell>` command that outputs shell completion scripts for bash and zsh.

## Deliverables

- `omniinstall completion bash` outputs a bash completion script
- `omniinstall completion zsh` outputs a zsh completion script
- completion covers: search, install, remove, verify, list, explain, stack, state commands
- tests for completion command output (non-empty, valid shebang)

## Acceptance Criteria

- completion scripts are valid and installable
- unsupported shell returns a clear error
- tests pass

## References

- specs/08-CLI_SPEC.md
