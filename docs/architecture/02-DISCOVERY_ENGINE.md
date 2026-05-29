# OmniInstall Discovery Engine Architecture

Version: 1.0 Draft  
Status: Architecture Document  
Owner: MITSLab  
Document Priority: Critical  
Last Updated: 2026-05-30

---

## 1. Purpose

The Discovery Engine is responsible for helping users find software without requiring knowledge of Linux package names, package managers, or distribution-specific repositories.

It turns user intent into normalized application candidates.

A user searches for:

> OBS

The Discovery Engine should understand that this may refer to:

> OBS Studio

and return a user-friendly application result that can later be resolved into installable sources by the Source Resolver.

---

## 2. Product Role

Discovery is the first major user-facing workflow in OmniInstall.

If discovery fails, installation does not matter.

OmniInstall is not only an install command. It is a software discovery and installation product.

Therefore, the Discovery Engine must be designed as a first-class system component.

---

## 3. Responsibilities

The Discovery Engine is responsible for:

- accepting user queries
- searching local and remote metadata
- handling aliases and common names
- normalizing application identity
- ranking results
- returning application metadata
- exposing availability hints
- supporting categories
- supporting future stack discovery

It is not responsible for choosing the final installation source.

That responsibility belongs to the Source Resolver.

---

## 4. Non-Responsibilities

The Discovery Engine must not:

- execute installation commands
- decide package manager strategy
- add repositories
- download binaries
- perform privilege escalation
- mutate system state
- decide final trust policy

Discovery must remain read-oriented and side-effect free.

---

## 5. Input Model

Discovery inputs may include:

- free-text search query
- category selection
- application alias
- stack reference
- direct application ID
- future web catalog selection

Example inputs:

```text
obs
OBS Studio
video recording
streaming
code editor
shopify developer setup
```

---

## 6. Output Model

The Discovery Engine should return normalized application candidates.

A candidate should include:

- canonical application ID
- display name
- short description
- long description if available
- icon reference if available
- categories
- aliases
- homepage
- license if known
- publisher or maintainer if known
- availability hints
- confidence score

Example conceptual result:

```yaml
id: obs-studio
display_name: OBS Studio
summary: Video recording and live streaming software.
categories:
  - video
  - streaming
aliases:
  - obs
  - open broadcaster software
homepage: https://obsproject.com
availability_hint:
  flatpak: true
  apt: true
  dnf: true
confidence: 0.98
```

---

## 7. Application Identity

Application identity must be independent from package names.

The canonical application ID is an OmniInstall-level identifier.

Example:

```text
obs-studio
```

This may map to:

- apt package: obs-studio
- flatpak ID: com.obsproject.Studio
- Fedora package: obs-studio
- Arch package: obs-studio

The user sees the application.

The system maps the sources.

---

## 8. Metadata Sources

Discovery may use multiple metadata sources over time.

### 8.1 Local Catalog

A bundled or installed metadata catalog for supported applications.

Used for:

- offline search
- predictable MVP behavior
- tests

### 8.2 Remote Catalog

A future hosted or community-maintained catalog.

Used for:

- broader application coverage
- metadata updates
- community contributions

### 8.3 System Sources

Existing package managers may expose searchable package data.

This can help discover availability but should not replace normalized application metadata.

### 8.4 Community Metadata

Community-submitted metadata may improve aliases, descriptions, categories, icons, and stack relationships.

---

## 9. Ranking

Discovery ranking should prioritize:

1. exact name matches
2. official application identity
3. common aliases
4. category relevance
5. popularity if available
6. metadata completeness
7. source availability
8. trust quality

Ranking must not be silently manipulated by monetization.

Sponsored or promoted content, if ever introduced, must be clearly labeled and separated from organic trust-based results.

---

## 10. Search Behavior

Search should tolerate:

- lowercase input
- uppercase input
- partial names
- common abbreviations
- spelling variations where possible
- old names or aliases

Example:

Search:

```text
vscode
```

May return:

```text
Visual Studio Code
VSCodium
```

The ranking should make the likely intended result clear while preserving alternatives.

---

## 11. Categories

Categories support browsing and non-search discovery.

Initial categories may include:

- Browsers
- Office
- Development
- Graphics
- Audio
- Video
- Streaming
- Gaming
- Communication
- Security
- System Tools
- Education

Categories should remain user-facing and understandable.

They should not mirror technical package repository categories unless those categories are useful to users.

---

## 12. Integration With Source Resolver

The Discovery Engine returns application candidates.

The Source Resolver receives a selected candidate and determines installable sources for the current system.

Discovery answers:

> What application does the user mean?

Source Resolver answers:

> What is the best way to install it here?

This separation is essential.

---

## 13. MVP Requirements

For the MVP, the Discovery Engine should support:

- local metadata catalog
- exact and partial name search
- aliases
- category metadata
- application ID lookup
- at least 5-10 known applications
- testable search results

The MVP should not require a remote service.

---

## 14. Failure Handling

Discovery failures should be helpful.

Bad:

> No result.

Better:

> No supported application was found for this search. Try another name or check whether OmniInstall supports this application yet.

Future versions may allow:

- request application support
- search external sources
- submit metadata

---

## 15. Testing Requirements

Discovery tests should cover:

- exact matches
- alias matches
- partial matches
- unknown queries
- ranking behavior
- category filtering
- metadata parsing

Discovery must be deterministic in tests.

---

## 16. Future Expansion

Future Discovery Engine capabilities may include:

- remote registry search
- fuzzy matching
- popularity signals
- localization
- screenshots
- user ratings
- verified publisher metadata
- stack recommendations
- compatibility badges

These must not compromise the simple MVP search experience.

---

## 17. Summary

The Discovery Engine is the user's entry point into OmniInstall.

It transforms human search intent into structured application identity.

It must remain user-centered, source-aware, side-effect free, and separated from installation decisions.
