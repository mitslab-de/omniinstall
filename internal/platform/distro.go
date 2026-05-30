// Package platform provides system platform detection for OmniInstall.
//
// Platform detection reads /etc/os-release (or equivalent) to determine
// the Linux distribution family, which informs adapter availability checks
// and source resolution priority.
package platform

import (
	"bufio"
	"io"
	"strings"
)

// DistroFamily represents the family of the detected Linux distribution.
type DistroFamily string

const (
	// FamilyDebian covers Debian, Ubuntu, Mint, and derivatives.
	FamilyDebian DistroFamily = "debian"
	// FamilyFedora covers Fedora, RHEL, CentOS, AlmaLinux, Rocky, and derivatives.
	FamilyFedora DistroFamily = "fedora"
	// FamilyArch covers Arch Linux, Manjaro, and derivatives.
	FamilyArch DistroFamily = "arch"
	// FamilyOpenSUSE covers openSUSE Leap, Tumbleweed, and SLES.
	FamilyOpenSUSE DistroFamily = "opensuse"
	// FamilyUnknown is returned when the distribution cannot be identified.
	FamilyUnknown DistroFamily = "unknown"
)

// OsReleaseReader is a function that opens /etc/os-release for reading.
// It is a variable to allow injection of alternative readers in tests.
var OsReleaseReader = defaultOsReleaseReader

// defaultOsReleaseReader reads from /etc/os-release.
func defaultOsReleaseReader() (io.ReadCloser, error) {
	// Import os lazily to keep the package testable.
	return openOsRelease()
}

// Detect reads /etc/os-release and returns the distribution family.
// On any read error, FamilyUnknown is returned without an error — the
// caller should fall back to detecting available package managers instead.
func Detect() DistroFamily {
	r, err := OsReleaseReader()
	if err != nil {
		return FamilyUnknown
	}
	defer r.Close()
	return parseOsRelease(r)
}

// DetectFromReader parses os-release content from the given reader.
// This is the testable core of Detect().
func DetectFromReader(r io.Reader) DistroFamily {
	return parseOsRelease(r)
}

// parseOsRelease reads key=value pairs from os-release and classifies
// the distribution family from the ID and ID_LIKE fields.
func parseOsRelease(r io.Reader) DistroFamily {
	fields := make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		// Strip optional surrounding quotes.
		value = strings.Trim(value, `"'`)
		fields[strings.TrimSpace(key)] = value
	}
	return classifyFamily(fields["ID"], fields["ID_LIKE"])
}

// classifyFamily maps distro ID and ID_LIKE strings to a DistroFamily.
func classifyFamily(id, idLike string) DistroFamily {
	combined := strings.ToLower(id + " " + idLike)

	// Check for exact or substring matches in order of specificity.
	for _, marker := range []string{"arch", "manjaro"} {
		if strings.Contains(combined, marker) {
			return FamilyArch
		}
	}
	for _, marker := range []string{"opensuse", "suse", "sles"} {
		if strings.Contains(combined, marker) {
			return FamilyOpenSUSE
		}
	}
	for _, marker := range []string{"fedora", "rhel", "centos", "almalinux", "rocky", "ol"} {
		if strings.Contains(combined, marker) {
			return FamilyFedora
		}
	}
	for _, marker := range []string{"debian", "ubuntu", "mint", "pop", "kali", "raspbian"} {
		if strings.Contains(combined, marker) {
			return FamilyDebian
		}
	}
	return FamilyUnknown
}
