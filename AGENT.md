# OmniInstall Agent Operating Instructions

Version: 1.0 Draft  
Status: Required Agent Instruction Document  
Owner: MITSLab  
Last Updated: 2026-05-30

---

## 1. Role

You are an AI development agent working on OmniInstall.

OmniInstall is a universal Linux application discovery and installation platform designed to make Linux software installation simple for normal users.

You are not building a generic package manager.

You are building a user-focused application layer above existing Linux software ecosystems.

---

## 2. Required Reading Before Coding

Before making code changes, read:

- `docs/product/01-PRODUCT.md`
- `docs/product/04-UX_PRINCIPLES.md`
- `docs/product/05-MVP_DEFINITION.md`
- `docs/architecture/01-SYSTEM_OVERVIEW.md`
- `docs/architecture/03-SOURCE_RESOLVER.md`
- `docs/architecture/04-INSTALL_ENGINE.md`
- `docs/architecture/06-SECURITY_MODEL.md`
- `docs/architecture/08-COMPONENT_BOUNDARIES.md`
- `specs/01-APPLICATION_MODEL.md`
- `specs/02-SOURCE_MODEL.md`
- `specs/03-INSTALL_PLAN.md`
- `specs/04-ADAPTER_INTERFACE.md`
- `specs/05-RESOLVER_RULES.md`
- `specs/08-CLI_SPEC.md`
- `specs/09-TESTING_STRATEGY.md`

If there is a conflict between documents, prefer this order:

1. Product documents
2. Architecture documents
3. Specs
4. Tasks
5. Existing code

---

## 3. Working Method

Always work task by task.

1. Open `tasks/`.
2. Select the first task marked `Status: Open`.
3. Work only on that task.
4. Implement the smallest complete solution.
5. Add or update tests.
6. Update documentation if needed.
7. Mark the task as `Status: Done` only when complete.
8. Add a completion note to the task file.
9. Commit with a clear message.
10. Continue to the next open task only after the current task is complete.

Do not skip tasks unless the task is blocked and you document why.

---

## 4. Core Product Guardrails

Never forget:

- Users install applications, not package formats.
- OmniInstall hides unnecessary complexity but does not hide meaningful risk.
- The GUI and CLI must use the same core engine.
- The Source Resolver chooses sources.
- The Install Engine executes plans.
- Adapters contain backend-specific logic.
- The Security Engine classifies risk.
- The project is open source first.

---

## 5. Architecture Rules

Do not put package-manager-specific logic in the product layer.

Do not put source ranking logic inside adapters.

Do not put installation execution logic inside the Discovery Engine.

Do not duplicate business logic between CLI and future GUI.

Do not silently execute high-risk actions.

Do not introduce a new package format unless an ADR exists.

---

## 6. MVP Scope

The initial MVP focuses on:

- application model
- source model
- install plan
- discovery engine
- source resolver
- adapter interface
- APT adapter
- Flatpak adapter
- CLI commands
- tests

Do not implement enterprise, cloud, accounts, marketplace, or complex GUI features until MVP foundation tasks are complete.

---

## 7. Testing Rules

Every core component needs tests.

Resolver logic must be deterministic.

Adapters should be tested with mocks first.

Live package-manager tests should be separated from unit tests.

Do not mark a task done without tests unless you clearly document why tests are not applicable.

---

## 8. Documentation Rules

If implementation changes architecture, specs, CLI behavior, or data models, update the relevant document.

Documentation is not optional.

Undocumented functionality is incomplete functionality.

---

## 9. Commit Rules

Use clear commit messages.

Examples:

- `Add core domain models`
- `Implement source resolver ranking`
- `Add adapter interface and tests`

Avoid vague messages such as:

- `update`
- `fix stuff`
- `changes`

---

## 10. Final Instruction

Build OmniInstall as a product, not as a collection of scripts.

Every change should make Linux software installation easier, safer, or more understandable.
