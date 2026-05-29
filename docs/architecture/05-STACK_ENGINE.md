# OmniInstall Stack Engine Architecture

Version: 1.0 Draft  
Status: Architecture Document  
Owner: MITSLab  
Document Priority: High  
Last Updated: 2026-05-30

---

## 1. Purpose

The Stack Engine enables OmniInstall to install curated groups of applications and optional setup steps as one coherent workflow.

A stack represents a user outcome rather than a single application.

Examples:

- Gaming Setup
- Streamer Setup
- Developer Setup
- Shopify Developer Setup
- Homelab Setup
- School Laptop Setup

The Stack Engine is not required for the smallest MVP, but the architecture must support it from the beginning.

---

## 2. Product Role

Single-application installation solves the first core problem:

> How do I install this app?

Stacks solve a larger problem:

> How do I prepare my Linux system for this activity?

This makes OmniInstall more than an app installer.

It becomes a setup assistant for real-world workflows.

---

## 3. Responsibilities

The Stack Engine is responsible for:

- reading stack definitions
- validating stack metadata
- resolving stack items into applications
- planning multi-application installation
- ordering dependencies where needed
- presenting the stack plan to the user
- delegating each application to the Source Resolver and Install Engine
- tracking stack-level progress
- recording stack-level local state
- supporting future stack sharing

---

## 4. Non-Responsibilities

The Stack Engine must not:

- directly install packages
- bypass the Source Resolver
- bypass the Install Engine
- hide risk from users
- run arbitrary scripts without security classification
- become a general configuration-management system too early

Stacks coordinate installation.

They do not replace the core installation pipeline.

---

## 5. Stack Definition Model

A stack definition describes a named setup.

Conceptual example:

```yaml
id: streamer-setup
display_name: Streamer Setup
summary: Install common tools for streaming and content creation.
items:
  - application: obs-studio
    required: true
  - application: discord
    required: false
  - application: audacity
    required: false
  - application: kdenlive
    required: false
```

Future versions may support:

- optional groups
- configuration steps
- verification steps
- post-install guidance
- distro-specific notes
- maintainer metadata
- trust metadata

---

## 6. Stack Item Types

Initial stack item type:

- application reference

Future item types may include:

- configuration step
- service enablement
- repository requirement
- command-line tool
- desktop preference
- documentation link
- manual follow-up step

The MVP should focus on application references only.

---

## 7. Stack Planning Flow

High-level flow:

```text
User selects stack
    ↓
Stack Engine loads definition
    ↓
Validate stack metadata
    ↓
Resolve stack items to applications
    ↓
Source Resolver evaluates each application
    ↓
Security Engine reviews combined risk
    ↓
Install plan is presented to user
    ↓
Install Engine installs each selected item
    ↓
Verification runs per item
    ↓
Stack result is recorded
```

---

## 8. User Choice

Stacks should not force unnecessary installations.

The user should be able to review what will be installed.

Stack items may be:

- required
- recommended
- optional

Example:

Gaming Setup:

Required:

- Steam
- Mesa utilities

Recommended:

- Lutris
- Heroic Games Launcher

Optional:

- Discord
- OBS Studio

---

## 9. Risk Handling

The risk of a stack is the combined risk of its items.

If one item requires a high-risk action, the whole stack flow must communicate that risk clearly.

The Stack Engine must not hide risky actions inside a friendly stack name.

Example:

Bad:

> Install Developer Setup

while silently adding external repositories and running scripts.

Good:

> This setup installs 7 applications. 1 action requires adding an external repository. Review before continuing.

---

## 10. Failure Handling

Stack installation may partially succeed.

The system must clearly report:

- installed items
- failed items
- skipped items
- cancelled items
- verification failures

The Stack Engine must not mark a stack fully installed unless all required items succeed.

Optional failures should be reported but may not fail the full stack depending on policy.

---

## 11. Local State

Stack state should record:

- stack ID
- stack version
- selected items
- installed items
- failed items
- source decisions
- timestamp
- verification result

This enables future:

- stack updates
- stack repair
- stack removal assistance
- diagnostics

---

## 12. MVP Position

The Stack Engine is not part of the minimal first implementation unless the core app installation pipeline is already stable.

However, the MVP architecture must not block stacks.

Recommended MVP approach:

- define stack file format draft
- support internal stack planning tests
- defer full user-facing stack installation until after single-app installation works

---

## 13. Future Expansion

Future Stack Engine capabilities may include:

- community stack registry
- verified stacks
- certified MITSLab stacks
- organization-approved stacks
- localized stack descriptions
- stack ratings
- stack dependency graphs
- stack updates
- stack migration
- OpenStore OS onboarding stacks

---

## 14. Testing Requirements

Stack Engine tests should cover:

- valid stack definition
- invalid stack definition
- required item failure
- optional item failure
- source resolution per item
- combined risk summary
- partial installation reporting
- local state recording

---

## 15. Summary

The Stack Engine allows OmniInstall to install outcomes, not just apps.

It should remain layered above Discovery, Source Resolver, Install Engine, Verification Engine, and Security Engine.

Stacks are a major long-term differentiator, but they must be implemented carefully so they do not turn OmniInstall into an unsafe script runner or overly complex configuration-management system.
