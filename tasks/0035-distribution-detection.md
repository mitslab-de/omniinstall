# Task 0035 - Distribution Detection

Status: Open
Priority: P2

## Goal

Implement distribution detection by reading /etc/os-release and mapping to known families (debian, fedora, arch, unknown).

## Deliverables

- `internal/platform/distro.go` with Detect() returning DistroFamily
- DistroFamily type: FamilyDebian, FamilyFedora, FamilyArch, FamilyUnknown
- unit tests using a mock /etc/os-release reader
- integration with adapter availability checks

## Acceptance Criteria

- Detect() parses ID and ID_LIKE fields from /etc/os-release
- Unknown distributions return FamilyUnknown without error
- tests pass

## References

- docs/architecture/01-SYSTEM_OVERVIEW.md
- specs/04-ADAPTER_INTERFACE.md
