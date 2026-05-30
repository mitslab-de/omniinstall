# OmniInstall Local State Specification

## Purpose

Defines data persisted locally by OmniInstall.

---

## Goals

- track installations
- support removal
- support updates
- support diagnostics
- support future migration and rollback features

---

## Local Installation Record

Required fields:

- application_id
- source_type
- source_identifier
- install_timestamp
- install_status
- verification_status

Optional fields:

- version
- install_plan_reference
- stack_reference

---

## Storage Principles

- human inspectable where practical
- versioned format
- migration friendly
- corruption resistant

### File Schema Versioning

- top-level `schema_version` field is required for newly written state files
- current writer emits `schema_version: "1.0"`
- loaders must remain backward-compatible with legacy files that do not have `schema_version`
- incompatible major schema versions must fail with clear, user-facing errors

---

## MVP Goal

Provide enough local state to support install, verify, remove, and troubleshoot workflows.
