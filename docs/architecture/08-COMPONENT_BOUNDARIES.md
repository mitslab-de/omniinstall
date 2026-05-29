# OmniInstall Component Boundaries

## Purpose

Defines ownership and responsibilities between major components.

## Discovery Engine

Owns:

- application search
- aliases
- metadata lookup

Must not:

- install software
- select final source

## Source Resolver

Owns:

- source ranking
- recommendation generation
- install plan creation

Must not:

- execute installation

## Install Engine

Owns:

- execution
- progress
- logging
- verification handoff

Must not:

- choose source

## Security Engine

Owns:

- risk classification
- trust evaluation
- confirmation requirements

## Stack Engine

Owns:

- stack parsing
- stack planning
- stack progress

Must not:

- bypass resolver or install engine

## Adapters

Own:

- backend-specific behavior

Examples:

- apt
- dnf
- pacman
- flatpak

Must not:

- make product decisions

## GUI

Owns:

- presentation
- interaction

Must not:

- implement installation logic

## CLI

Owns:

- command interface

Must not:

- duplicate core business logic
