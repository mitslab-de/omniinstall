# OmniInstall System Overview

Version: 1.0 Draft  
Status: Architecture Foundation Document  
Owner: MITSLab  
Document Priority: Critical  
Last Updated: 2026-05-30

---

## 1. Purpose

This document defines the high-level system architecture for OmniInstall.

It translates the product vision into a buildable software architecture.

OmniInstall is not a single package manager. It is a user-facing application discovery and installation platform that orchestrates existing Linux software sources through a consistent experience.

---

## 2. Architectural Goal

The architecture must support one central product promise:

> Users can find and install Linux software without knowing package managers, package formats, or distribution-specific commands.

The system must therefore separate user intent from installation mechanics.

User intent:

- install OBS Studio
- remove VLC
- update applications
- install a Gaming Setup stack

Installation mechanics:

- apt install
- dnf install
- pacman -S
- flatpak install
- AppImage handling
- vendor repository setup

OmniInstall connects these layers safely and transparently.

---

## 3. High-Level Architecture

OmniInstall is composed of the following major components:

```text
User Interface Layer
    ├── GUI
    └── CLI

Core Application Layer
    ├── Discovery Engine
    ├── Source Resolver
    ├── Install Engine
    ├── Update Engine
    ├── Removal Engine
    ├── Verification Engine
    ├── Stack Engine
    └── Security Engine

Integration Layer
    ├── Package Manager Adapters
    ├── Flatpak Adapter
    ├── Snap Adapter
    ├── AppImage Adapter
    ├── Vendor Source Adapter
    └── Registry Adapter

Data Layer
    ├── Application Metadata
    ├── Source Metadata
    ├── Stack Definitions
    ├── Local State
    ├── Logs
    └── Policy Configuration
```

---

## 4. Core Architectural Principle

The GUI and CLI must use the same core engine.

The CLI must not contain business logic that the GUI cannot reuse.

The GUI must not implement installation behavior independently.

All user interfaces call into the same core services.

This ensures:

- consistent behavior
- testability
- future desktop readiness
- predictable automation
- fewer duplicated bugs

---

## 5. User Interface Layer

### 5.1 GUI

The GUI is a first-class product surface.

It should provide:

- application search
- application detail pages
- install button
- update view
- remove flow
- stack browsing
- progress feedback
- trust information
- advanced details when requested

The GUI must be beginner-friendly by default.

### 5.2 CLI

The CLI is required for:

- development
- automation
- testing
- power users
- CI environments
- server use cases

Example commands:

```bash
omni search obs
omni install obs-studio
omni remove obs-studio
omni explain obs-studio
omni stack install gaming-setup
```

The CLI should use beginner-friendly language where possible, but may expose advanced flags.

---

## 6. Core Application Layer

### 6.1 Discovery Engine

Responsible for finding applications and presenting unified results.

Input:

- search query
- category browsing
- stack references

Output:

- normalized application candidates
- metadata
- availability hints

### 6.2 Source Resolver

Responsible for choosing the best installation source for the current system.

It evaluates:

- distribution
- architecture
- package availability
- trust level
- update behavior
- user preference
- installed state
- source quality

### 6.3 Install Engine

Responsible for turning a resolved installation plan into actual system changes.

It must support:

- planning
- confirmation
- execution
- progress reporting
- error handling
- verification handoff

### 6.4 Verification Engine

Responsible for confirming whether installation succeeded.

Verification may include:

- package manager state
- command existence
- desktop entry existence
- application binary availability
- optional version checks

### 6.5 Stack Engine

Responsible for installing curated sets of applications.

Stacks are not the MVP's primary focus, but the architecture must allow them.

Examples:

- Gaming Setup
- Streamer Setup
- Developer Setup
- Shopify Developer Setup
- Homelab Setup

### 6.6 Security Engine

Responsible for risk analysis and safety decisions.

It must classify actions such as:

- installing from trusted repositories
- adding repositories
- executing scripts
- downloading binaries
- modifying system files

---

## 7. Integration Layer

Package manager and source-specific behavior must live behind adapters.

Adapters should expose a common interface to the core engine.

Examples:

- AptAdapter
- DnfAdapter
- PacmanAdapter
- ZypperAdapter
- FlatpakAdapter
- SnapAdapter
- AppImageAdapter

Adapters should not decide product-level policy.

They report capabilities and perform actions.

The Source Resolver decides which adapter should be used.

---

## 8. Data Layer

OmniInstall requires structured data to function reliably.

### 8.1 Application Metadata

Describes applications in user-facing terms.

Examples:

- name
- description
- icon
- website
- categories
- aliases

### 8.2 Source Metadata

Maps applications to available installation sources.

Examples:

- apt package name
- dnf package name
- flatpak ID
- snap name
- AppImage URL
- vendor repository instructions

### 8.3 Stack Definitions

Describe collections of applications and optional setup steps.

### 8.4 Local State

Tracks what OmniInstall installed, when, from which source, and how it was verified.

### 8.5 Logs

Logs must support debugging, user support, and future auditability.

---

## 9. Installation Flow

High-level installation flow:

```text
User selects application
    ↓
Discovery Engine resolves application identity
    ↓
Source Resolver evaluates available sources
    ↓
Security Engine assesses risk
    ↓
Install Engine creates installation plan
    ↓
User confirms if needed
    ↓
Adapter executes installation
    ↓
Verification Engine checks result
    ↓
Local State and logs are updated
```

---

## 10. Failure Handling

Failure handling must be part of the architecture from the beginning.

Failures should include:

- no source found
- source unavailable
- permission denied
- package manager error
- network failure
- verification failure
- partial install
- unsupported distribution

Errors must be:

- understandable
- actionable
- logged
- testable

Bad error:

> backend exited with code 100

Good error:

> OBS Studio could not be installed because the package source for your system is currently unavailable. Try again later or use the Flatpak source.

---

## 11. MVP Architecture Boundary

The MVP should include:

- CLI
- core engine
- application metadata format
- source resolver
- APT adapter
- one additional native adapter, preferably DNF or Pacman
- Flatpak adapter
- install flow
- remove flow
- verification flow
- structured logging

The MVP should not include:

- user accounts
- cloud sync
- enterprise policies
- paid marketplace
- complex plugin marketplace
- full rollback system

---

## 12. Architectural Non-Goals

OmniInstall should not become:

- a replacement package manager
- a new package format first
- a configuration management tool first
- a Linux distribution first
- a closed app store

The architecture should support future expansion without losing the core identity.

---

## 13. Future Architecture Expansion

Future components may include:

- hosted registry
- private catalogs
- certified stacks
- policy engine
- enterprise audit layer
- GUI app store client
- plugin SDK
- remote metadata sync
- OpenStore OS integration

These should be built on top of the core architecture, not mixed into the MVP core.

---

## 14. System Overview Summary

OmniInstall is architected as a user-facing application layer above existing Linux software ecosystems.

The system separates:

- what the user wants
- what sources exist
- which source is best
- how installation is executed
- how success is verified

This separation is essential for building a product that is simple for beginners, useful for advanced users, and maintainable for contributors.
