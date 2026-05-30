# Task 0048 - Stack Install Command

Status: Open
Priority: P2

## Goal

Add `omniinstall stack install <stack-file>` command that installs all applications defined in a stack YAML file.

## Deliverables

- App.stackInstall(path string) reads a stack file and installs each recipe
- main.go dispatcher for `stack install <file>`
- dry-run and --execute modes
- tests for successful stack install, missing file, and invalid stack

## Acceptance Criteria

- stack install installs applications in order
- any failure stops the stack install and reports which app failed
- tests pass

## References

- docs/stacks.md
- examples/stacks/
- cmd/omniinstall/main.go
