# OmniInstall UX Principles

Version: 1.0 Draft  
Status: Foundation Document  
Owner: MITSLab  
Document Priority: Critical  
Last Updated: 2026-05-30

---

# 1. Purpose

This document defines the user experience principles that govern OmniInstall.

Every screen, workflow, interaction, dialog, search result, installation flow, and future feature should be evaluated against these principles.

If a design decision conflicts with these principles, the principles take precedence.

---

# 2. Core UX Philosophy

OmniInstall exists to reduce cognitive load.

The user should think about software.

The system should think about package managers.

Bad experience:

- choose package manager
- choose repository
- choose package format
- choose installation method

Good experience:

- choose application
- click install

---

# 3. Users See Applications

Users install applications.

Users do not install package formats.

The primary object in the interface is always the application.

Not:

- Flatpak package
- RPM package
- DEB package

Instead:

- OBS Studio
- Steam
- Docker Desktop
- Firefox

Technical details should be available but secondary.

---

# 4. Simplicity First

The default path should be the simplest safe path.

The beginner should be able to complete common tasks without reading documentation.

A user should be able to:

- search
- install
- update
- remove

without understanding Linux packaging systems.

---

# 5. Progressive Disclosure

Complexity should never disappear.

Complexity should be hidden until needed.

Example:

Beginner View:

- Install
- Update
- Remove

Advanced View:

- Source details
- Version information
- Package mapping
- Repository details
- Verification details

Both audiences must be supported.

---

# 6. Trust Must Be Visible

Trust information should be understandable.

Examples:

- Official source
- Community source
- Verified source
- Unverified source

Users should not need security expertise to make reasonable decisions.

The system should communicate risk clearly.

---

# 7. Recommendations Over Decisions

Whenever possible OmniInstall should recommend rather than force users to decide.

Example:

Available Sources:

- Native Repository
- Flatpak
- AppImage

Recommended:

- Native Repository

The user may override the recommendation.

The user should not be forced to research the choice.

---

# 8. Search Is The Primary Workflow

Search is the heart of OmniInstall.

The primary interaction model is:

Search
→ Discover
→ Install

Navigation structures must never overshadow search.

---

# 9. Fast Feedback

Users should always know:

- what is happening
- why it is happening
- what happens next

Installation progress should be understandable.

Failures should be actionable.

---

# 10. Desktop First

The desktop experience is a first-class product surface.

The GUI is not a secondary wrapper around the CLI.

The GUI is part of the product vision.

The CLI and GUI should share the same core engine.

---

# 11. Accessibility

The product should be usable by as many people as possible.

Requirements:

- keyboard navigation
- screen reader compatibility
- clear typography
- meaningful color usage
- understandable language

---

# 12. Beginner Friendly Language

Language should prioritize clarity.

Prefer:

- Install
- Update
- Remove

Avoid exposing unnecessary technical terminology.

Technical terms should be explained when shown.

---

# 13. One Goal Per Screen

Screens should have a clear primary purpose.

Examples:

Application Page:

Goal:
Install application.

Update Screen:

Goal:
Review updates.

Stack Page:

Goal:
Install complete setup.

---

# 14. UX Success Criteria

A successful UX means:

- users install software without documentation
- users trust recommendations
- users understand warnings
- users recover from failures
- users recommend OmniInstall to others

---

# 15. UX Principle Summary

The best OmniInstall experience is one where users think about what they want to achieve, not about how Linux packaging works.

The software should absorb complexity so users can focus on outcomes.
