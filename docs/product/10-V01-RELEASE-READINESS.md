# OmniInstall v0.1 Release Readiness Checklist

Version: 1.0  
Status: Active  
Owner: MITSLab  
Last Updated: 2026-05-30

---

## Purpose

This checklist defines the quality gates that must pass before the OmniInstall v0.1 MVP
release is considered ready.

All checks are reproducible locally using the commands listed in each section.

---

## 1. Build and Test Gates

### 1.1 Full Test Suite Passes

```sh
go test ./...
```

Expected: all packages report `ok`.

### 1.2 Binary Builds Without Errors

```sh
go build ./...
```

Expected: exits with code 0, no compiler errors or warnings.

### 1.3 Code Formatting Clean

```sh
gofmt -w $(find . -name '*.go' -not -path './.git/*')
git diff --exit-code
```

Expected: no diff after formatting.

---

## 2. Core User-Flow Validation

Each of the following flows must work end-to-end on a supported Linux distribution
with at least one supported package manager (APT or Flatpak).

### 2.1 Search

- [ ] `omniinstall search git` returns at least one result.
- [ ] `omniinstall search git --json` returns valid JSON with a `results` array.
- [ ] Searching for an unknown term produces a clear "not found" message, not a crash.

### 2.2 Explain

- [ ] `omniinstall explain git` shows DisplayName, Summary, Categories, trust level, and risk level.
- [ ] `omniinstall explain git --json` returns enriched JSON with `trust_level` and `risk_level`.
- [ ] Explaining an unknown ID produces a clear "not found" message.

### 2.3 Install

- [ ] `omniinstall install git` (dry-run default) prints a plan without making changes.
- [ ] `omniinstall install git --execute` performs the real installation (requires supported system).
- [ ] Post-install state is recorded in local state file.
- [ ] Attempting to install an unknown application returns exit code 4 (not-found).

### 2.4 Verify

- [ ] `omniinstall verify git` after successful install prints `✔ git: verification passed`.
- [ ] `omniinstall verify git` for a not-installed app returns a clear error.
- [ ] Verify result is recorded in state as `passed` or `failed`.

### 2.5 Remove

- [ ] `omniinstall remove git` after install succeeds and updates state to `removed`.
- [ ] Attempting to remove an app that was never installed returns a clear error.
- [ ] `omniinstall list` after removal does not show the removed app.

### 2.6 List

- [ ] `omniinstall list` shows all installed applications.
- [ ] `omniinstall list --json` returns valid JSON.
- [ ] Removed applications do not appear in list output.

---

## 3. Error and Edge-Case Checks

- [ ] Running `omniinstall` with no arguments prints usage without crashing.
- [ ] Running `omniinstall --help` prints usage.
- [ ] Unknown commands return exit code 2.
- [ ] High-risk applications prompt for confirmation before install.
- [ ] CLI exits with code 0 on success, non-zero on error (see exit code table in `main.go`).

---

## 4. Adapter Readiness

### 4.1 APT Adapter

- [ ] `CheckInstalled` correctly parses dpkg status for installed, not-installed, partial, and unknown states.
- [ ] Remove uses `apt-get remove -y`.
- [ ] Verify calls `dpkg -s <package>` and checks status field.

### 4.2 Flatpak Adapter

- [ ] `CheckInstalled` checks both `--system` and `--user` scopes.
- [ ] Error classification distinguishes "not installed" from "backend unavailable".

---

## 5. State and Reconciliation

- [ ] Local state file is created at `$XDG_STATE_HOME/omniinstall/state.yaml` (fallback `~/.omniinstall/state.yaml`).
- [ ] State file includes `schema_version` field.
- [ ] `DefaultEngine.Reconcile()` detects `missing_from_backend`, `partial_in_backend`, `present_but_removed`, and `unknown` drift kinds.

---

## 6. Catalog Quality

- [ ] Embedded MVP catalog contains at least 10 validated applications.
- [ ] All catalog entries pass `Application.Validate()`.
- [ ] `omniinstall search <any catalog ID>` returns that application in results.

Local validation command:

```sh
go test ./internal/discovery/... -run TestMVPCatalog
```

---

## 7. Documentation

- [ ] `README.md` reflects current CLI commands and options.
- [ ] `docs/product/05-MVP_DEFINITION.md` is up to date.
- [ ] `CONTRIBUTING.md` describes how to run tests and format code.
- [ ] Architecture docs in `docs/architecture/` describe all active components.

---

## 8. CI Gate

- [ ] GitHub Actions CI passes on the release branch.
- [ ] CI runs `gofmt`, `go build ./...`, and `go test ./...`.
- [ ] No unresolved code scanning alerts at critical or high severity.

---

## Sign-off

All items above must be checked before tagging `v0.1`.

| Gate                    | Status  | Notes |
|-------------------------|---------|-------|
| Build and test          | Pending |       |
| Core user flows         | Pending |       |
| Error handling          | Pending |       |
| Adapter readiness       | Pending |       |
| State and reconciliation| Pending |       |
| Catalog quality         | Pending |       |
| Documentation           | Pending |       |
| CI gate                 | Pending |       |
