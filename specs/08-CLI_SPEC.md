# OmniInstall CLI Specification

## Purpose

Defines the command-line interface.

---

## Core Commands

```bash
omni search <query>
omni install <application>
omni remove <application>
omni explain <application>
omni list
```

---

## Future Commands

```bash
omni stack install <stack>
omni doctor
omni update
```

---

## Output Principles

- beginner friendly
- actionable
- consistent
- explainable

---

## Exit Codes

OmniInstall CLI commands must return stable non-zero exit codes on failures:

- `0`: success
- `2`: usage error (invalid/missing command arguments, unknown command)
- `3`: not found (unknown application or missing recorded installation)
- `4`: backend unavailable (no compatible/available adapter or source backend)
- `5`: execution failed (all other operational failures)

---

## MVP Goal

A small but stable CLI surface that maps directly to core engine capabilities.
