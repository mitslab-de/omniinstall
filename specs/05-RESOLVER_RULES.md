# OmniInstall Resolver Rules Specification

## Purpose

Defines deterministic source selection rules.

---

## Default Preference Order

Suggested baseline:

1. trusted native package source
2. verified Flatpak
3. trusted vendor source
4. Snap where supported
5. verified AppImage
6. verified GitHub release
7. community source
8. unknown source

---

## Ranking Inputs

- compatibility
- trust level
- risk level
- update behavior
- integration quality
- user preferences

---

## Determinism

Given identical metadata and system context, the resolver must produce identical recommendations.

---

## Explanation Requirement

Every recommendation must be explainable in both machine-readable and user-facing form.

---

## Conflict Handling

Conflicts must be surfaced rather than silently ignored.

Examples:

- already installed from another source
- unsupported architecture
- unavailable package manager

---

## MVP Goal

Reliable and predictable source recommendations across supported package ecosystems.
