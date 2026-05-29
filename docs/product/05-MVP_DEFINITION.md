# OmniInstall MVP Definition

Version: 1.0 Draft  
Status: Foundation Document  
Owner: MITSLab  
Document Priority: Critical  
Last Updated: 2026-05-30

---

## 1. Purpose

This document defines the first meaningful Minimum Viable Product for OmniInstall.

The MVP must prove the core product promise:

> A Linux user can find and install software through OmniInstall without needing to understand package managers, package formats, or distribution-specific commands.

The MVP is not the final product.

It is the smallest version that demonstrates the essential value.

---

## 2. MVP Product Goal

The MVP must enable a user to:

1. Search for a known application.
2. See a clear result.
3. Install the application through OmniInstall.
4. Verify that installation completed.
5. Remove the application again.

This flow must work on more than one Linux distribution and through more than one installation source.

---

## 3. MVP User Promise

The MVP promise is:

> Search. Install. Use.

The user does not need to know whether the application is installed through APT, DNF, Pacman, or Flatpak.

OmniInstall chooses the best available supported path for the current system.

---

## 4. MVP Scope

### 4.1 Required Product Capabilities

The MVP must include:

- application search
- application metadata display
- source resolution
- installation
- removal
- basic update awareness
- verification after installation
- clear error messages
- CLI interface
- GUI-ready core architecture

### 4.2 Required Technical Capabilities

The MVP must include:

- distribution detection
- package manager detection
- source availability detection
- native package manager integration for at least two package managers
- Flatpak support
- structured logging
- basic test suite
- documented architecture

### 4.3 Required Documentation

The MVP must include:

- product definition
- architecture overview
- source resolver specification
- install engine specification
- contribution guide
- initial user guide

---

## 5. MVP Supported Sources

The MVP should support a narrow but meaningful source set.

Required:

- APT
- DNF or Pacman
- Flatpak

Optional if practical:

- AppImage metadata support
- Snap detection

The MVP should not attempt to fully support every ecosystem immediately.

Depth is more important than superficial breadth.

---

## 6. MVP Supported Use Cases

### Use Case 1: Install a Common Desktop Application

Example:

- OBS Studio
- VLC
- Firefox
- Bitwarden

Expected result:

The user searches, selects, installs, and sees confirmation.

### Use Case 2: Install a Developer Tool

Example:

- Git
- VS Code
- Docker

Expected result:

OmniInstall resolves the best supported source and installs it.

### Use Case 3: Remove an Application

Expected result:

The user can remove software installed through OmniInstall.

### Use Case 4: Explain Source Choice

Expected result:

Advanced details show why a source was selected.

---

## 7. MVP Non-Goals

The MVP must avoid scope creep.

The MVP should not include:

- user accounts
- cloud sync
- enterprise policy management
- paid marketplace features
- social features
- broad plugin marketplace
- complete GUI polish
- full rollback system
- every Linux distribution
- every package format
- full stack marketplace

Stacks may be represented as a future-facing concept, but the MVP should focus primarily on single-application installation.

---

## 8. MVP User Experience Requirements

The MVP must still respect the product's UX principles.

Even if the first version is CLI-heavy, the experience must remain beginner-oriented in language and structure.

Required qualities:

- clear application names
- clear install confirmation
- clear progress output
- clear failure reasons
- no unexplained package-manager jargon in beginner-facing output

Example good output:

> Installing OBS Studio using the recommended source for your system.

Example bad output:

> Running backend resolver apt candidate obs-studio package.

---

## 9. MVP Architecture Requirements

The MVP must not be a throwaway prototype.

It must establish the final architecture direction:

- core engine separated from CLI
- GUI can call the same core engine later
- source resolver is modular
- package manager integrations are isolated
- application metadata is structured
- tests are part of the workflow

The MVP may be small, but it must not be architecturally misleading.

---

## 10. MVP Success Criteria

The MVP is successful when:

- at least 5 common applications can be installed through OmniInstall
- at least 2 Linux package ecosystems are supported
- Flatpak installation works for supported apps
- installation results are verified
- users can remove installed applications
- errors are understandable
- contributors can add new application metadata
- the architecture is documented
- automated tests exist for resolver and install planning logic

---

## 11. MVP Failure Criteria

The MVP fails if:

- it only becomes a wrapper around one package manager
- users still need to know package names and package managers
- architecture prevents a future GUI
- source selection is undocumented
- installation failures are unclear
- tests are missing
- documentation is outdated at release

---

## 12. MVP Release Name

Recommended MVP release name:

OmniInstall 0.1 "First Install"

Purpose:

Prove that universal Linux application installation can be simplified without replacing existing package managers.
