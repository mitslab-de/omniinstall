# OmniInstall Product Requirements Document (PRD)

Version: 1.0 Draft  
Status: Foundation Document  
Owner: MITSLab  
Document Priority: Critical  
Last Updated: 2026-05-30

---

## 1. Executive Summary

OmniInstall is a universal Linux software discovery, installation, update, and management platform.

Its purpose is to make software installation on Linux as simple and approachable as installing applications on Windows or macOS, while still respecting the technical diversity and openness of the Linux ecosystem.

OmniInstall is not intended to replace existing Linux package managers. It is designed to sit above them as a user-focused abstraction layer.

The user should not need to understand whether an application is available through APT, DNF, Pacman, Zypper, Flatpak, Snap, AppImage, a vendor repository, or a GitHub release.

The user should only need to answer one question:

> What do I want to install?

OmniInstall handles the discovery, source selection, installation, verification, update, and removal workflow.

The long-term goal is to become the universal Linux application layer: a consistent app-store-like experience that works across distributions and software sources.

---

## 2. Product Thesis

Linux is technically mature enough for many everyday users, developers, creators, students, and small organizations.

The limiting factor is often not Linux itself.

The limiting factor is the experience around finding, installing, updating, and trusting software.

Today, users are forced to understand concepts that should be implementation details:

- distribution-specific repositories
- package managers
- package formats
- third-party repositories
- portable application formats
- vendor-specific installation instructions
- terminal commands

This creates friction at exactly the moment where a new user expects confidence.

OmniInstall is based on the thesis that:

> Linux adoption is not limited by software availability alone. Linux adoption is limited by fragmented software discovery and installation complexity.

The product therefore focuses on user experience first and packaging mechanics second.

OmniInstall should make Linux feel less like a collection of technical systems and more like a coherent product experience.

---

## 3. The Linux Software Installation Problem

### 3.1 Fragmented Sources

A single Linux application may be available through several channels:

- native distribution repositories
- Flatpak
- Snap
- AppImage
- vendor repositories
- direct downloads
- GitHub releases
- source builds

Each source may have different tradeoffs around updates, sandboxing, integration, trust, and compatibility.

A normal user is not equipped to evaluate these tradeoffs.

### 3.2 Distribution-Specific Instructions

Installation instructions often look different for every distribution:

```bash
sudo apt install example
```

```bash
sudo dnf install example
```

```bash
sudo pacman -S example
```

```bash
flatpak install flathub org.example.App
```

The user is expected to know which command applies.

That is unacceptable for a product aimed at normal users.

### 3.3 Inconsistent Graphical Experiences

Existing graphical software centers are valuable, but they are often tied to a specific desktop environment, distribution, or source configuration.

They do not consistently answer the broader user question:

> What is the best way to install this software on my Linux system?

### 3.4 Trust and Safety Confusion

Users often cannot tell which source is official, maintained, secure, outdated, community-provided, sandboxed, or risky.

OmniInstall must make trust visible without overwhelming users.

---

## 4. Product Vision

OmniInstall should become the universal Linux App Store layer.

It should provide a single consistent experience for:

- discovering software
- installing software
- updating software
- removing software
- installing complete software stacks
- sharing curated setups

The desired user experience is simple:

1. Search for an application.
2. Read clear human-friendly information.
3. Click install.
4. Use the software.

The complexity of package formats and source selection should be handled by OmniInstall.

---

## 5. Target Users

### 5.1 Primary Users

#### Linux Beginners

Users switching from Windows or macOS who want Linux to feel approachable.

They need:

- clear search
- simple installation
- safe defaults
- no required terminal usage
- understandable warnings

#### Casual Desktop Users

Users who can use Linux but do not want to become system administrators.

They need:

- reliable installation
- easy updates
- clean removal
- minimal technical decisions

#### Students and Educational Users

Users in schools, universities, workshops, or learning environments.

They need:

- easy setup
- predictable software collections
- low support overhead

### 5.2 Secondary Users

#### Developers

Developers may already understand Linux, but they benefit from reproducible setup stacks.

Example stacks:

- Web Developer Setup
- PHP Developer Setup
- Node.js Developer Setup
- Shopify Developer Setup
- AI Workstation Setup

#### Creators

Content creators need tools like OBS, Audacity, Kdenlive, Discord, and graphics software without understanding packaging systems.

#### Homelab Users

Homelab users benefit from curated installation workflows for common tools and services.

#### Small Organizations

Small organizations need repeatable software deployment without enterprise complexity.

---

## 6. User Personas

### 6.1 Anna — Linux Beginner

Anna previously used Windows. She installs Linux because she wants a free, modern, privacy-friendly system.

She wants to install a browser, office software, a password manager, and communication tools.

Her pain points:

- terminal commands feel unsafe
- different download formats are confusing
- she does not know which source to trust

Success for Anna:

- she opens OmniInstall
- searches for software
- clicks install
- never needs to understand APT, Flatpak, or AppImage

### 6.2 Markus — Developer

Markus is comfortable with the terminal but wants faster machine setup.

He wants Git, VS Code, Docker, Node.js, PHP, database tools, and command-line utilities.

Success for Markus:

- he installs a complete development stack
- he can inspect what OmniInstall will do
- he can reuse the stack on another machine

### 6.3 Sarah — Creator

Sarah creates videos and streams.

She needs OBS Studio, Discord, Audacity, Kdenlive, and graphics tools.

Success for Sarah:

- she installs a Streamer Setup stack
- OmniInstall chooses suitable sources
- tools appear in her desktop menu afterward

### 6.4 Small Organization Admin

A small school or non-profit wants ten Linux laptops with the same software.

Success:

- a curated stack installs required software
- updates remain manageable
- support questions decrease

---

## 7. Jobs To Be Done

### Core Job

When I need software on Linux, help me find and install it safely without learning package systems.

### Functional Jobs

- find software
- compare available sources when needed
- install software
- update software
- remove software
- install groups of applications
- verify that installation worked

### Emotional Jobs

- feel safe
- feel confident
- avoid fear of breaking the system
- avoid confusion
- feel that Linux is approachable

---

## 8. Product Principles

### 8.1 Users Install Applications, Not Package Formats

Users should see applications, not packaging technologies.

Bad:

> Install Flatpak package org.example.App.

Good:

> Install Example App.

### 8.2 The Best Default Should Be Chosen Automatically

OmniInstall should recommend the best available source based on:

- distribution compatibility
- source trust
- update behavior
- maintenance status
- desktop integration
- user preference
- security characteristics

### 8.3 Transparency Must Be Available

Normal users should not be forced to understand details, but advanced users must be able to inspect them.

OmniInstall should support both:

- simple mode
- advanced details

### 8.4 Existing Standards Over Reinvention

OmniInstall should reuse existing Linux package managers and standards whenever possible.

It should not create a new package format unless there is a strong product and technical reason.

### 8.5 Open Source First

The core platform should remain open source.

The project should avoid artificial lock-in.

Commercial opportunities may exist around hosting, support, certification, consulting, managed services, and training, but the core product should remain open.

### 8.6 Safety Before Convenience

OmniInstall should not hide dangerous behavior.

When an action has risk, the user should receive clear guidance.

### 8.7 Desktop Experience Is First-Class

The CLI is important, but OmniInstall is ultimately meant to make Linux easier for normal users.

The desktop experience must not be an afterthought.

---

## 9. Core Product Areas

### 9.1 Discovery

Search across software sources and present applications in a unified catalog.

### 9.2 Source Resolution

Identify available installation sources and recommend the best one.

### 9.3 Installation

Install software through the appropriate backend.

### 9.4 Verification

Confirm that installation succeeded.

### 9.5 Updates

Provide a unified update experience where technically feasible.

### 9.6 Removal

Remove software cleanly and predictably.

### 9.7 Stacks

Install curated groups of applications and configuration steps.

Examples:

- Gaming Setup
- Streamer Setup
- Developer Setup
- Shopify Developer Setup
- Homelab Setup

---

## 10. MVP Definition

The first meaningful MVP should not attempt to solve every Linux packaging problem.

It should prove the core product promise:

> A user can search for a known application and install it through OmniInstall without knowing the underlying source.

### MVP Must Include

- CLI foundation
- core engine foundation
- source resolver foundation
- support for at least two native package managers
- Flatpak support
- basic application metadata
- install command
- remove command
- verification step
- simple GUI prototype or GUI-ready API design

### MVP Should Avoid

- accounts
- cloud sync
- enterprise features
- marketplace monetization
- complex social features
- broad plugin ecosystem

The MVP should be small, testable, and focused.

---

## 11. Open Source Strategy

OmniInstall should be open source first.

Recommended model:

- OmniInstall Core: open source
- CLI: open source
- GUI: open source
- recipe/catalog definitions: open source
- registry protocol: open source
- documentation: open source

Potential commercial offerings may include:

- hosted registry infrastructure
- managed mirrors
- certified stacks
- enterprise support
- consulting
- training
- sponsored listings with strict transparency rules

The open-source strategy must protect community trust.

Commercial features must not compromise the core promise.

---

## 12. Monetization Strategy

OmniInstall should not start with monetization as the primary driver.

The first goal is adoption and trust.

Possible later revenue streams:

- donations
- sponsorships
- paid support
- enterprise support contracts
- certified business stacks
- hosted services
- training material
- consulting around Linux deployments

Advertising should be treated carefully. If ever introduced, it must be transparent, clearly labeled, and never distort source recommendations in a way that harms users.

---

## 13. Success Metrics

Early success:

- users can install common applications successfully
- installation failure rate decreases over time
- documentation becomes easier to follow
- contributors can add catalog entries

Community success:

- external contributors submit recipes/catalog metadata
- multiple distributions are tested
- issues are triaged consistently
- users recommend OmniInstall as the easiest Linux installation path

Product success:

- OmniInstall becomes a common first stop for Linux software installation
- users stop searching distribution-specific installation commands for supported apps

---

## 14. Risks

### Technical Risks

- inconsistent package names across distributions
- different versions in different sources
- conflicting installations
- privilege handling
- rollback complexity

### Product Risks

- becoming too technical
- trying to support too many sources too early
- losing beginner focus
- building a package manager instead of a user product

### Community Risks

- low contribution quality
- unsafe recipes
- disputes over recommended sources
- trust issues around monetization

---

## 15. Non Goals

OmniInstall is not:

- a Linux distribution
- a replacement for APT
- a replacement for DNF
- a replacement for Pacman
- a replacement for Flatpak
- a replacement for Snap
- a configuration management system first
- an enterprise deployment product first

OmniInstall is a user-facing application discovery and installation layer.

---

## 16. Product Decision Log

### Decision 1: Universal Linux App Store Direction

OmniInstall is positioned as a universal Linux App Store layer, not merely a package manager.

### Decision 2: Open Source First

The project should be open source first, with commercial opportunities built around services rather than locking down the core.

### Decision 3: Desktop Experience Matters

The GUI must be considered a first-class product surface, not a later optional add-on.

### Decision 4: Stacks Are Additive

Application stacks are a major feature, but they must not replace the core single-application install experience.

---

## 17. Final Product Statement

OmniInstall exists to make installing software on Linux simple for everyone.

It hides package complexity without hiding important trust and safety information.

It respects the Linux ecosystem instead of replacing it.

It aims to become the universal application layer that Linux has always needed.
