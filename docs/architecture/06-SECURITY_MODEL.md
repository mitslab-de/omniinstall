# OmniInstall Security Model Architecture

Version: 1.0 Draft  
Status: Architecture Document  
Owner: MITSLab  
Document Priority: Critical  
Last Updated: 2026-05-30

---

## 1. Purpose

OmniInstall performs actions that can change a Linux system.

It may install packages, remove packages, add repositories, download binaries, modify local state, and later execute stack workflows.

Therefore, security is not an optional feature.

Security must be built into the architecture from the beginning.

---

## 2. Security Philosophy

OmniInstall should make Linux easier without making unsafe behavior invisible.

The product should reduce user confusion, but it must not hide meaningful risk.

Security model principle:

> Hide unnecessary complexity. Never hide meaningful risk.

---

## 3. Core Security Goals

The security model must ensure:

- users understand risky actions before they happen
- untrusted sources are clearly marked
- privileged operations are explicit
- source recommendations are not silently manipulated
- install plans are inspectable
- logs support troubleshooting and auditability
- dangerous actions require confirmation
- future policies can restrict behavior

---

## 4. Risk Levels

OmniInstall should classify actions into risk levels.

### LOW

Examples:

- installing from trusted distribution repository
- installing verified Flatpak from trusted remote
- reading local metadata
- searching catalog

### MEDIUM

Examples:

- installing from community-maintained source
- enabling user-level Flatpak remote
- downloading AppImage from verified upstream
- removing an application installed by OmniInstall

### HIGH

Examples:

- adding external repository
- adding repository signing key
- installing binary from direct download
- running post-install commands
- modifying system services

### CRITICAL

Examples:

- executing remote scripts
- modifying system files outside controlled adapters
- disabling security features
- deleting broad filesystem paths
- installing from unknown or unverifiable source

Critical actions should be blocked by default unless explicitly allowed by advanced settings or policy.

---

## 5. Trust Levels

Sources should be classified independently from action risk.

Suggested trust levels:

- official distribution repository
- official upstream source
- verified publisher
- verified community source
- unverified community source
- unknown source
- blocked source

Trust level influences source ranking but does not automatically determine final risk.

Example:

A trusted vendor repository may still require a high-risk action if it adds a new repository and signing key.

---

## 6. User Confirmation Rules

Suggested confirmation behavior:

- LOW: no extra confirmation beyond install action
- MEDIUM: normal confirmation if meaningful
- HIGH: explicit warning and confirmation
- CRITICAL: blocked by default

Warnings must be understandable.

Bad warning:

> Risk level HIGH.

Good warning:

> This installation needs to add an external software repository to your system. Only continue if you trust the publisher.

---

## 7. Privilege Handling

OmniInstall must never request elevated privileges earlier than necessary.

The user should know why privileges are requested.

Examples:

- Native package install requires system privileges.
- User-level Flatpak install may not require root.
- Adding a repository requires elevated privileges.

Privilege prompts must be tied to a clear action.

---

## 8. Script Execution Policy

Script execution is one of the highest-risk behaviors.

OmniInstall should avoid arbitrary script execution in the MVP.

If script execution is added later, it must require:

- explicit metadata declaration
- risk classification
- user confirmation
- logging
- preferably sandboxing or restricted execution
- review for community catalog inclusion

Remote curl-pipe-shell patterns must not be treated as normal installation methods.

---

## 9. Repository Modification Policy

Adding repositories changes system trust boundaries.

Repository modifications must be treated as high-risk actions.

The system must show:

- repository URL
- publisher
- reason
- signing key information if available
- removal implications

Repository additions should not happen silently inside stacks.

---

## 10. Downloaded Binary Policy

Direct downloads and AppImages require trust evaluation.

The system should prefer:

- official publisher URLs
- checksums
- signatures
- release metadata
- reproducible source references where available

Unknown binary downloads should be discouraged or blocked by default depending on policy.

---

## 11. Logging And Auditability

Security-relevant actions must be logged.

Logs should include:

- timestamp
- application ID
- source selected
- risk level
- user confirmation status
- adapter used
- result
- error details if failure occurs

Logs must not leak secrets.

---

## 12. Policy Readiness

The MVP does not require a full policy engine, but the architecture must allow policies later.

Future policies may include:

- block unverified sources
- prefer sandboxed apps
- require official sources
- block external repositories
- allow only organization-approved stacks
- require signed metadata

---

## 13. Security MVP Requirements

The MVP must include:

- basic action risk classification
- source trust metadata field
- high-risk confirmation support
- structured security logs
- no arbitrary remote script execution
- no silent repository addition
- clear user-facing warnings

---

## 14. Security Non-Goals For MVP

The MVP should not include:

- full sandboxing engine
- full enterprise policy system
- remote reputation scoring
- automated malware scanning
- complete rollback security model
- cryptographic signing for every metadata item

These may be added later.

---

## 15. Summary

OmniInstall must be simple, but it must not be careless.

The security model exists to protect users while preserving a beginner-friendly experience.

The core rule is:

> Make safe choices easy. Make risky choices visible. Block dangerous choices by default.
