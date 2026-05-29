# OmniInstall Source Resolver Architecture

Version: 1.0 Draft  
Status: Architecture Document  
Owner: MITSLab  
Document Priority: Critical  
Last Updated: 2026-05-30

---

## 1. Purpose

The Source Resolver is one of the most important components in OmniInstall.

It determines the best available installation source for a selected application on the user's current Linux system.

The user should not need to decide between APT, DNF, Pacman, Flatpak, Snap, AppImage, vendor repositories, or direct downloads.

The Source Resolver makes a recommendation based on system compatibility, trust, availability, update behavior, and user preferences.

---

## 2. Product Role

The Source Resolver enables OmniInstall's core product promise:

> The user chooses an application. OmniInstall chooses the best way to install it.

It is the bridge between human intent and Linux packaging complexity.

Discovery answers:

> What application does the user mean?

Source Resolver answers:

> What is the best installation source for this application on this system?

Install Engine answers:

> How do we execute that plan safely?

---

## 3. Responsibilities

The Source Resolver is responsible for:

- receiving a normalized application candidate
- detecting available source mappings
- evaluating current system compatibility
- ranking possible sources
- identifying the recommended source
- explaining the recommendation
- detecting conflicts and unsupported states
- producing an install plan candidate

It does not execute installation commands.

Execution belongs to the Install Engine and adapters.

---

## 4. Non-Responsibilities

The Source Resolver must not:

- install packages
- remove packages
- modify repositories directly
- run privileged commands directly
- download binaries directly
- make UI decisions
- hide security concerns

The resolver decides what should happen, not how execution is performed.

---

## 5. Input Model

Inputs include:

- canonical application ID
- application metadata
- source metadata
- detected distribution
- detected architecture
- installed package managers
- enabled sources
- user preferences
- policy configuration
- local installation state

Example conceptual input:

```yaml
application_id: obs-studio
system:
  distro: ubuntu
  version: "24.04"
  architecture: x86_64
available_sources:
  apt:
    package: obs-studio
  flatpak:
    id: com.obsproject.Studio
user_preferences:
  prefer_sandboxed_apps: false
```

---

## 6. Output Model

The resolver should output a ranked source resolution result.

Example conceptual output:

```yaml
application_id: obs-studio
recommended_source: apt
reason: Native package available for this distribution and marked as trusted.
alternatives:
  - source: flatpak
    reason: Cross-distribution package available.
    tradeoff: Larger runtime dependency.
risk_level: low
requires_confirmation: false
```

The output must be explainable.

A future GUI should be able to show:

> Recommended for your system: Native package.

and advanced details:

> APT package obs-studio is available on Ubuntu 24.04. Flatpak is also available as an alternative.

---

## 7. Source Types

Initial source types:

- native package manager
- Flatpak
- Snap
- AppImage
- vendor repository
- direct download
- GitHub release

MVP source types should be limited to:

- APT
- DNF or Pacman
- Flatpak

Additional source types should be added after the resolver rules are stable.

---

## 8. Ranking Factors

The resolver should evaluate sources using multiple weighted factors.

### 8.1 Compatibility

Can this source install on the current system?

Factors:

- distribution
- version
- architecture
- package manager availability
- runtime availability

### 8.2 Trust

How trustworthy is the source?

Potential trust levels:

- official upstream
- distribution repository
- verified Flatpak
- community maintained
- unknown
- untrusted

### 8.3 Update Behavior

How will the application receive updates?

Preferred sources should have clear update behavior.

### 8.4 Integration

How well does the app integrate with the desktop?

Examples:

- desktop entry
- icon
- MIME handlers
- command availability

### 8.5 Security

Does the source require risky actions?

Examples:

- adding repository keys
- executing scripts
- installing unsigned binaries
- downloading from unknown domains

### 8.6 User Preference

Users may prefer:

- native packages
- sandboxed apps
- newest versions
- stable versions
- open-source-only software

User preferences should influence ranking but not override critical safety rules silently.

---

## 9. Recommended Default Ranking

Default ranking should generally prefer:

1. trusted native package source
2. verified Flatpak
3. trusted vendor repository
4. Snap where supported and appropriate
5. AppImage from verified upstream
6. GitHub release from verified upstream
7. community source
8. unknown direct download

This ranking is not absolute.

The resolver must allow application-specific rules.

Example:

Some applications may be better maintained through Flatpak than through an outdated native repository.

---

## 10. Explanation Requirement

Every resolver decision must be explainable.

Bad:

> Selected apt.

Good:

> OmniInstall recommends the native Ubuntu package because it is available for your system, updates through your system package manager, and has a low risk level.

The resolver must produce machine-readable and user-readable explanations.

---

## 11. Conflict Handling

The resolver must detect possible conflicts.

Examples:

- application already installed from another source
- multiple versions installed
- package conflicts with existing package
- Flatpak and native package both present
- vendor repository already configured differently

Conflict results should not be hidden.

The user should receive clear options.

---

## 12. Installed State Awareness

The resolver should consider installed state.

If an application is already installed through Flatpak, OmniInstall should not blindly install a native package unless the user explicitly requests a source change.

Possible states:

- not installed
- installed through OmniInstall
- installed outside OmniInstall
- installed from same source
- installed from different source
- unknown state

---

## 13. Policy Integration

Future policy rules may affect source selection.

Examples:

- only open-source software
- prefer sandboxed apps
- block Snap
- block unverified sources
- organization-approved sources only

The MVP does not need a full policy engine, but resolver design must allow policy integration.

---

## 14. MVP Requirements

The MVP resolver must support:

- distribution detection input
- architecture detection input
- application source metadata
- at least APT
- one additional native source such as DNF or Pacman
- Flatpak
- ranking logic
- explanation output
- no-source-found handling
- deterministic tests

---

## 15. MVP Non-Goals

The MVP resolver should not yet support:

- complex enterprise policies
- paid source ranking
- full trust graph
- remote reputation scoring
- automatic vendor repository onboarding for many apps
- full AppImage lifecycle management
- source migration flows

These can be added later.

---

## 16. Testing Requirements

Resolver tests must cover:

- native package available
- Flatpak fallback
- no supported source
- unsupported distribution
- installed from different source
- user preference changes
- ranking stability
- explanation generation
- conflict detection

Tests should be deterministic and not require live package repositories for core logic.

Live integration tests may exist separately.

---

## 17. Failure Modes

Common failure modes:

- no known source
- source metadata invalid
- package manager unavailable
- distribution unsupported
- architecture unsupported
- source disabled
- policy blocks all sources

Failures must be actionable.

Example:

> No supported source was found for OBS Studio on your system. Flatpak support may provide an install path if enabled.

---

## 18. Future Expansion

Future resolver features may include:

- remote catalog scoring
- verified publisher metadata
- source health checks
- stale package detection
- CVE awareness
- user source preferences
- organization policies
- source migration assistant
- source comparison UI

---

## 19. Summary

The Source Resolver is the decision-making heart of OmniInstall.

It transforms a selected application into a recommended installation path.

It must be transparent, testable, policy-ready, and user-centered.

If built correctly, it allows OmniInstall to feel like one app store while still respecting the diversity of the Linux ecosystem.
