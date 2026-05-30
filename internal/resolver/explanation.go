package resolver

import (
	"fmt"
	"strings"

	"github.com/mitslab-de/omniinstall/internal/source"
)

// explain produces a user-facing explanation for why a source was recommended.
func explain(s source.Source, score int, ctx SystemContext) string {
	var parts []string

	switch s.SourceType {
	case source.TypeAPT:
		parts = append(parts, fmt.Sprintf("Native APT package %q is available", s.SourceIdentifier))
	case source.TypeDNF:
		parts = append(parts, fmt.Sprintf("Native DNF package %q is available", s.SourceIdentifier))
	case source.TypePacman:
		parts = append(parts, fmt.Sprintf("Native Pacman package %q is available", s.SourceIdentifier))
	case source.TypeZypper:
		parts = append(parts, fmt.Sprintf("Native Zypper package %q is available", s.SourceIdentifier))
	case source.TypeFlatpak:
		parts = append(parts, fmt.Sprintf("Flatpak application %q is available", s.SourceIdentifier))
	case source.TypeSnap:
		parts = append(parts, fmt.Sprintf("Snap package %q is available", s.SourceIdentifier))
	case source.TypeVendor:
		parts = append(parts, fmt.Sprintf("Vendor source %q is available", s.SourceIdentifier))
	case source.TypeAppImage:
		parts = append(parts, fmt.Sprintf("AppImage %q is available", s.SourceIdentifier))
	default:
		parts = append(parts, fmt.Sprintf("Source %q is available", s.SourceIdentifier))
	}

	switch s.TrustLevel {
	case source.TrustOfficial:
		parts = append(parts, "from an official source")
	case source.TrustVerified:
		parts = append(parts, "from a verified source")
	case source.TrustCommunity:
		parts = append(parts, "from a community source")
	}

	switch s.RiskLevel {
	case source.RiskLow:
		parts = append(parts, "with low installation risk")
	case source.RiskMedium:
		parts = append(parts, "with medium installation risk")
	case source.RiskHigh:
		parts = append(parts, "with high installation risk — review carefully")
	case source.RiskCritical:
		parts = append(parts, "with critical installation risk — not recommended")
	}

	if ctx.UserPreferences.PreferNative && nativeTypes[s.SourceType] {
		parts = append(parts, "(preferred: native package)")
	}
	if ctx.UserPreferences.PreferFlatpak && s.SourceType == source.TypeFlatpak {
		parts = append(parts, "(preferred: Flatpak)")
	}

	return strings.Join(parts, ", ") + "."
}

// explainAlternative produces a shorter explanation for an alternative source.
func explainAlternative(s source.Source) string {
	switch s.SourceType {
	case source.TypeAPT, source.TypeDNF, source.TypePacman, source.TypeZypper:
		return fmt.Sprintf("Native package %q (%s) is also available as an alternative.", s.SourceIdentifier, s.SourceType)
	case source.TypeFlatpak:
		return fmt.Sprintf("Flatpak %q is also available as a cross-distribution alternative.", s.SourceIdentifier)
	case source.TypeSnap:
		return fmt.Sprintf("Snap %q is also available.", s.SourceIdentifier)
	default:
		return fmt.Sprintf("Source %q (%s) is also available.", s.SourceIdentifier, s.SourceType)
	}
}
