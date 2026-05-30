# Task 0001 - Core Workspace

Status: Done
Priority: P0

## Goal

Create the initial project structure for the OmniInstall core.

## Deliverables

- source directory layout
- test directory layout
- package management setup
- linting setup
- formatting setup
- CI-ready structure

## Acceptance Criteria

- project builds successfully
- tests can run
- structure matches architecture documents

## References

- docs/architecture/*
- AGENT.md

## Completion Note

Implemented on 2026-05-30.

Created Go module `github.com/mitslab-de/omniinstall` with the following structure:

- `cmd/omniinstall/` — CLI entry point with basic command routing
- `internal/app/` — Application domain model with validation
- `internal/source/` — Source model with type, trust, and risk level constants
- `internal/install/` — Install Plan model with validation
- `internal/resolver/` — Source Resolver interface and SystemContext type
- `internal/engine/` — Install Engine interface and Result type
- `internal/discovery/` — Discovery Engine interface and Candidate type
- `internal/security/` — Security Engine interface and Assessment type
- `internal/adapters/` — Adapter interface with error category constants
- `internal/adapters/apt/` — APT adapter stub
- `internal/adapters/flatpak/` — Flatpak adapter stub
- `internal/state/` — Local state model with InstallRecord validation
- `internal/logging/` — Structured log entry type

All packages have tests. All tests pass. Formatting and CI workflow added.
