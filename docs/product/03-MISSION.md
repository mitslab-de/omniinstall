# OmniInstall Mission

Version: 1.0 Draft  
Status: Foundation Document  
Owner: MITSLab  
Document Priority: Critical  
Last Updated: 2026-05-30

---

## 1. Mission Statement

OmniInstall's mission is to make software installation on Linux simple, safe, understandable, and accessible for normal users while preserving the openness and power of the Linux ecosystem.

The project exists to reduce the gap between what Linux can technically do and what users can comfortably experience.

A user should not need to become a package-management expert to install everyday software.

OmniInstall removes unnecessary complexity from the installation journey.

---

## 2. Daily Mission

Every product, design, documentation, and engineering decision should be measured against one question:

> Does this make Linux software installation easier, safer, or more understandable for users?

If the answer is no, the decision should be questioned.

If a feature adds power but also adds complexity, the complexity must be justified, hidden behind progressive disclosure, or moved out of the beginner path.

---

## 3. Product Mission

OmniInstall should provide one trusted place where users can:

- discover Linux applications
- install applications
- update applications
- remove applications
- understand source and trust information when needed
- install curated software stacks
- share reusable setup definitions

The product mission is not to expose every technical detail by default.

The product mission is to make the correct action obvious while keeping expert detail available.

---

## 4. Engineering Mission

The engineering mission is to build a reliable abstraction layer above existing Linux software distribution systems.

OmniInstall should orchestrate, not replace.

It should integrate with:

- APT
- DNF
- Pacman
- Zypper
- Flatpak
- Snap where supported
- AppImage workflows
- vendor repositories
- trusted direct downloads
- GitHub release metadata where appropriate

The system must be designed for:

- testability
- modularity
- source transparency
- safe execution
- predictable failure behavior
- long-term maintainability

---

## 5. User Mission

For users, OmniInstall should create confidence.

A user should feel:

- I can find the software I need.
- I understand what will happen before it happens.
- I do not need to know which package manager my distribution uses.
- I can trust the recommendation or inspect it.
- I can undo or remove software when needed.
- Linux feels approachable.

The emotional outcome matters as much as the technical outcome.

---

## 6. Community Mission

OmniInstall should make it possible for the Linux community to improve software accessibility together.

Contributors should be able to help by adding:

- application metadata
- source mappings
- verification rules
- translations
- documentation
- stack definitions
- compatibility reports
- safety reviews

The project should treat documentation and metadata as first-class contributions, not secondary work.

A contributor who improves clarity for beginners contributes as meaningfully as a contributor who writes code.

---

## 7. Open Source Mission

OmniInstall should remain open source first.

The core should be inspectable, forkable, auditable, and community-improvable.

The project should avoid artificial lock-in and should not depend on proprietary control over the user's ability to install software.

Commercial services may exist around the ecosystem, but they must not undermine the open core or the user's trust.

---

## 8. Documentation Mission

Documentation is not optional.

OmniInstall must be understandable by:

- normal users
- contributors
- maintainers
- AI coding agents
- future project maintainers

Every major feature should include:

- user-facing explanation
- technical specification
- test expectations
- security considerations
- failure behavior

If a feature cannot be explained, it is not ready.

---

## 9. Mission Guardrails

OmniInstall should not become:

- a package-manager replacement first
- an enterprise management system first
- a configuration-management tool first
- a Linux distribution first
- a tool only experts can love

Those directions may become adjacent opportunities later, but they must not distract from the primary mission.

The primary mission remains:

> Make Linux software installation easy for normal users.

---

## 10. Mission Summary

OmniInstall exists to make the Linux software ecosystem easier to access without making it less open.

It should turn fragmented installation paths into a coherent user experience.

It should help Linux feel like a product ordinary people can use, not a puzzle they must solve.
