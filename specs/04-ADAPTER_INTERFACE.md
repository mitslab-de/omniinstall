# OmniInstall Adapter Interface Specification

Version: 1.0 Draft  
Status: Technical Specification  
Owner: MITSLab  
Document Priority: Critical  
Last Updated: 2026-05-30

---

## 1. Purpose

Adapters isolate backend-specific package and source behavior from the OmniInstall core.

The core should not directly execute APT, DNF, Pacman, Flatpak, Snap, or AppImage logic.

Instead, the Install Engine calls adapters through a stable interface.

---

## 2. Adapter Responsibilities

Adapters are responsible for:

- checking source availability
- validating backend availability
- planning backend-specific execution details
- executing install/remove/update operations
- returning structured results
- exposing backend capabilities

Adapters must not make product-level decisions.

---

## 3. Required Adapter Methods

Conceptual interface:

```text
name() -> string
is_available(system_context) -> bool
can_handle(source_type) -> bool
check_installed(source_identifier) -> InstalledState
install(install_plan) -> AdapterResult
remove(source_identifier) -> AdapterResult
verify(source_identifier) -> VerificationResult
```

Exact implementation language may define this as trait, interface, protocol, abstract class, or equivalent.

---

## 4. Adapter Result

Adapter operations must return structured results.

Required fields:

- success
- exit_code if applicable
- message
- technical_details
- duration
- changed_system

Adapters must not return raw terminal output as the only result.

---

## 5. Error Handling

Adapters should normalize backend errors into stable error categories:

- backend_unavailable
- package_not_found
- permission_denied
- network_error
- dependency_conflict
- execution_failed
- unsupported_operation

---

## 6. MVP Adapters

MVP should include:

- AptAdapter
- one of DnfAdapter or PacmanAdapter
- FlatpakAdapter

---

## 7. Testing

Adapters should be testable through:

- unit tests with mocked command execution
- integration tests in containers where possible
- contract tests against the adapter interface

---

## 8. Summary

Adapters protect the core from backend-specific complexity and make OmniInstall extensible across Linux ecosystems.
