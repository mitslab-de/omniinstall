# OmniInstall Application Model Specification

Version: 1.0 Draft  
Status: Technical Specification  
Owner: MITSLab  
Document Priority: Critical  
Last Updated: 2026-05-30

---

## 1. Purpose

The Application Model defines how OmniInstall represents software at the product level.

An application is not a package.

An application is the user-facing identity of software.

Example:

- Application: OBS Studio
- APT package: obs-studio
- Flatpak ID: com.obsproject.Studio

The Application Model separates what the user wants from how the system installs it.

---

## 2. Canonical Application ID

Every application must have a canonical OmniInstall ID.

Requirements:

- lowercase
- ASCII
- kebab-case
- stable over time
- independent from package manager names where possible

Examples:

```text
obs-studio
firefox
vlc
docker
visual-studio-code
```

The canonical ID is used internally and in metadata references.

---

## 3. Required Fields

Minimum application metadata:

```yaml
id: obs-studio
display_name: OBS Studio
summary: Video recording and live streaming software.
categories:
  - video
  - streaming
```

Required fields:

- id
- display_name
- summary
- categories

---

## 4. Recommended Fields

Recommended metadata fields:

```yaml
homepage: https://obsproject.com
license: GPL-2.0-or-later
publisher: OBS Project
aliases:
  - obs
  - open broadcaster software
description: Full user-facing description.
```

Recommended fields improve discovery, trust, and presentation.

---

## 5. Optional Fields

Optional fields may include:

- icon
- screenshots
- documentation_url
- support_url
- source_code_url
- tags
- maturity
- popularity_hint
- localization

These are not required for MVP.

---

## 6. Categories

Categories should be user-facing and understandable.

Examples:

- browsers
- office
- development
- video
- audio
- streaming
- gaming
- communication
- security
- system-tools
- education

Categories must not expose repository-specific technical categories unless useful.

---

## 7. Aliases

Aliases improve search.

Examples:

```yaml
aliases:
  - vscode
  - code
```

Aliases should include:

- common short names
- former names
- common user spellings

Aliases must not be used to mislead search results.

---

## 8. Application Identity Rules

The application identity must remain stable even if package names differ.

Example:

```yaml
id: visual-studio-code
display_name: Visual Studio Code
aliases:
  - vscode
  - code
```

Package-specific identifiers belong in source metadata, not application metadata.

---

## 9. MVP Requirements

The MVP must include application metadata for at least 5-10 known applications.

Recommended MVP seed applications:

- obs-studio
- vlc
- firefox
- bitwarden
- git
- docker
- visual-studio-code

---

## 10. Validation Rules

Application metadata validation must check:

- required fields exist
- id format is valid
- display_name is non-empty
- summary is non-empty
- categories are non-empty
- aliases do not duplicate ID unnecessarily

Invalid metadata must fail tests.

---

## 11. Summary

The Application Model is the foundation of OmniInstall's user-facing catalog.

It ensures users interact with applications, not package formats.
