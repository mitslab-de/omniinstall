# OmniInstall Data Model

## Purpose

Defines the core data entities used by OmniInstall.

### Primary Entities

- Application
- Source
- InstallPlan
- Stack
- VerificationRule
- LocalInstallation
- SecurityClassification

### Application

Canonical user-facing software identity.

Examples:

- obs-studio
- firefox
- docker

### Source

Maps an application to an installable source.

Examples:

- apt package
- dnf package
- flatpak id
- appimage source

### InstallPlan

Resolved execution plan produced by Source Resolver.

### Stack

Collection of applications and optional future actions.

### LocalInstallation

Records installation state tracked by OmniInstall.

### VerificationRule

Defines how installation success is validated.

### SecurityClassification

Defines trust level and risk level metadata.

The detailed schema will be defined later in specs/.
