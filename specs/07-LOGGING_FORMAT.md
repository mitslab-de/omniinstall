# OmniInstall Logging Format Specification

## Purpose

Defines structured logging requirements.

---

## Logging Goals

- troubleshooting
- diagnostics
- support
- auditing
- testing visibility

---

## Required Fields

- timestamp
- action
- application_id
- source_type
- result
- duration

Optional fields:

- risk_level
- verification_status
- error_category

---

## Rules

- logs must be structured
- logs must avoid secrets
- logs must support machine parsing

---

## MVP Goal

Provide consistent logging across resolver, install engine, and adapters.
