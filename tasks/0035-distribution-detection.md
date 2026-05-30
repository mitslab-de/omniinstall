# Task 0035 - Distribution Detection

Status: Done
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

## Completion Notes

- Created `internal/platform/distro.go` with `DistroFamily` type: FamilyDebian, FamilyFedora, FamilyArch, FamilyOpenSUSE, FamilyUnknown.
- `Detect()` uses injectable `OsReleaseReader` var defaulting to `/etc/os-release`.
- `DetectFromReader(r io.Reader)` is the testable core: parses key=value pairs, handles comments, strips quotes, maps ID + ID_LIKE to family.
- `classifyFamily()` checks arch/manjaro → FamilyArch; opensuse/suse/sles → FamilyOpenSUSE; fedora/rhel/centos/almalinux/rocky → FamilyFedora; debian/ubuntu/mint/pop/kali/raspbian → FamilyDebian.
- Created `internal/platform/distro_test.go` with 13 tests: Ubuntu, Debian, Fedora, RHEL, AlmaLinux, Arch, Manjaro, OpenSUSE, Unknown, Empty, CommentLines, QuotedValues, ConstantValues.
- All 18 packages pass `go test ./...`.
