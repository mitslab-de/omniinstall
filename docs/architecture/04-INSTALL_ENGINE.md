# OmniInstall Install Engine Architecture

Version: 1.0 Draft  
Status: Architecture Document  
Owner: MITSLab  
Document Priority: Critical  
Last Updated: 2026-05-30

---

## 1. Purpose

The Install Engine is responsible for turning a resolved installation decision into a safe, observable, and verifiable system change.

It receives an installation plan from the Source Resolver and coordinates execution through the appropriate adapter.

The Install Engine must never guess what to install.

It executes a plan.

---

## 2. Product Role

The Install Engine is the component that makes OmniInstall real for the user.

Discovery helps the user find software.

The Source Resolver chooses the best available source.

The Install Engine performs the installation workflow.

It must preserve the product promise:

> The user can install Linux software without understanding the underlying package manager.

---

## 3. Responsibilities

The Install Engine is responsible for:

- receiving an install plan
- validating the plan
- checking required permissions
- coordinating user confirmation
- calling the correct adapter
- reporting progress
- handling failures
- handing off to verification
- recording local state
- writing structured logs

---

## 4. Non-Responsibilities

The Install Engine must not:

- perform application discovery
- choose the best source
- directly implement package-manager-specific logic
- directly implement GUI behavior
- ignore security warnings
- silently execute high-risk actions

Package-manager-specific behavior belongs to adapters.

Source choice belongs to the Source Resolver.

Risk classification belongs to the Security Engine.

---

## 5. Install Plan

The Install Engine consumes an install plan.

A conceptual install plan may include:

```yaml
application_id: obs-studio
display_name: OBS Studio
source:
  type: apt
  package: obs-studio
risk_level: low
requires_privilege: true
requires_confirmation: false
verification:
  - type: command
    value: obs
  - type: desktop_entry
    value: com.obsproject.Studio.desktop
```

The plan must be explicit.

The Install Engine should reject incomplete or ambiguous plans.

---

## 6. Installation Flow

High-level flow:

```text
Receive install plan
    ↓
Validate plan
    ↓
Check security classification
    ↓
Check permissions
    ↓
Request confirmation if needed
    ↓
Execute through adapter
    ↓
Track progress
    ↓
Handle adapter result
    ↓
Run verification
    ↓
Update local state
    ↓
Write structured logs
    ↓
Return user-facing result
```

---

## 7. Permission Handling

OmniInstall must treat privilege escalation carefully.

The Install Engine should know whether an operation requires elevated privileges before execution begins.

Examples:

- native package installation usually requires privileges
- Flatpak user installation may not require root
- system-wide Flatpak installation may require privileges
- AppImage download may not require privileges
- repository modification requires elevated privileges

The user must not be surprised by privilege prompts.

---

## 8. Adapter Execution

The Install Engine executes through adapters.

Examples:

- AptAdapter
- DnfAdapter
- PacmanAdapter
- ZypperAdapter
- FlatpakAdapter
- SnapAdapter
- AppImageAdapter

Adapters expose a common execution interface.

The Install Engine should not contain raw package-manager command details except for adapter coordination.

---

## 9. Progress Reporting

Installation progress must be understandable.

The Install Engine should emit structured progress events such as:

- preparing
- checking permissions
- downloading
- installing
- verifying
- completed
- failed

These events should be usable by both CLI and GUI.

CLI example:

> Installing OBS Studio using the recommended source for your system.

GUI example:

> Preparing installation...
> Installing...
> Verifying...
> Done.

---

## 10. Verification Handoff

The Install Engine must not assume success only because an adapter command exited successfully.

After execution, it should hand off to the Verification Engine.

Verification may check:

- package manager state
- command availability
- desktop entry availability
- installed version
- file existence

A successful installation requires execution success and verification success where verification is defined.

---

## 11. Local State Recording

OmniInstall should record what it installs.

Local state may include:

- application ID
- source used
- package name or source identifier
- version if known
- install timestamp
- installed through OmniInstall or detected externally
- verification result

This enables:

- removal
- updates
- conflict detection
- support diagnostics
- future migration flows

---

## 12. Error Handling

Errors must be clear, actionable, and logged.

Bad error:

> apt failed with code 100

Better error:

> OBS Studio could not be installed because the package manager reported a dependency problem. You can try updating package lists or use the Flatpak source if available.

Errors should include:

- user-facing message
- technical details for logs
- adapter result
- suggested next action when possible

---

## 13. Failure Categories

The Install Engine should categorize failures.

Examples:

- permission denied
- network unavailable
- package not found
- dependency conflict
- source unavailable
- adapter failure
- verification failure
- user cancelled
- unsupported operation

Failure categories should be stable so CLI, GUI, tests, and documentation can rely on them.

---

## 14. Cancellation

Where possible, long-running installations should support cancellation.

Cancellation behavior depends on backend capabilities.

If cancellation is unsafe or unsupported, the user should be informed.

The system must avoid leaving local state marked as successfully installed when installation was cancelled.

---

## 15. Rollback Preparation

The MVP does not require full rollback.

However, the Install Engine should be designed so rollback can be added later.

To prepare for rollback, the engine should track:

- planned operations
- completed operations
- package/source used
- repository changes
- local files written
- verification failures

Rollback should not be faked.

If rollback is unavailable, the user should be told clearly.

---

## 16. MVP Requirements

The MVP Install Engine must support:

- install plan validation
- adapter execution
- progress events
- structured errors
- verification handoff
- local state write
- APT install flow
- one additional native package manager flow
- Flatpak install flow
- removal-ready state tracking

---

## 17. MVP Non-Goals

The MVP should not include:

- full rollback
- complex transaction system
- enterprise policy enforcement
- GUI-specific implementation
- remote execution
- unattended organization deployment
- complex multi-step stack installation

These features may come later.

---

## 18. Testing Requirements

Install Engine tests should cover:

- valid install plan
- invalid install plan
- low-risk install
- high-risk confirmation requirement
- adapter success
- adapter failure
- verification success
- verification failure
- local state write
- progress event sequence

Core tests should use mocked adapters.

Live package-manager tests should be separated as integration tests.

---

## 19. Security Considerations

The Install Engine must respect Security Engine decisions.

It must never silently execute high-risk actions.

Examples of high-risk actions:

- executing remote scripts
- adding repositories
- installing unsigned binaries
- modifying system files outside package manager control

The engine should support policy blocks in future versions.

---

## 20. Summary

The Install Engine is responsible for safe execution.

It must be plan-driven, adapter-based, observable, verifiable, and failure-aware.

A successful Install Engine allows OmniInstall to provide beginner-friendly installation while preserving the reliability expected from serious Linux tooling.
