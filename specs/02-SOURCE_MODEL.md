# OmniInstall Source Model Specification

## Purpose

Defines how installable sources are represented.

Sources connect application identities to installable artifacts.

---

## Source Types

Supported source types:

- apt
- dnf
- pacman
- zypper
- flatpak
- snap
- appimage
- vendor
- direct-download

---

## Example

```yaml
application_id: obs-studio
source_type: flatpak
source_identifier: com.obsproject.Studio
trust_level: verified
risk_level: low
```

---

## Required Fields

- application_id
- source_type
- source_identifier
- trust_level
- risk_level

---

## Trust Levels

- official
- verified
- community
- unknown
- blocked

---

## Risk Levels

- low
- medium
- high
- critical

---

## MVP Requirements

MVP source model must support:

- apt
- one additional native package manager
- flatpak

Additional source types can be introduced later.
