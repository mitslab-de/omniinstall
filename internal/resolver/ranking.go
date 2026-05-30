package resolver

import (
	"github.com/mitslab-de/omniinstall/internal/source"
)

// nativeTypes is the set of distribution-native package manager source types.
var nativeTypes = map[source.Type]bool{
	source.TypeAPT:    true,
	source.TypeDNF:    true,
	source.TypePacman: true,
	source.TypeZypper: true,
}

// baseTypePriority returns the baseline score for a source type.
// Higher is more preferred, per specs/05-RESOLVER_RULES.md.
func baseTypePriority(t source.Type) int {
	switch t {
	case source.TypeAPT, source.TypeDNF, source.TypePacman, source.TypeZypper:
		return 50 // native package managers
	case source.TypeFlatpak:
		return 40
	case source.TypeVendor:
		return 35
	case source.TypeSnap:
		return 30
	case source.TypeAppImage:
		return 25
	case source.TypeDirectDownload:
		return 20
	default:
		return 10
	}
}

// trustModifier returns the score adjustment for a trust level.
func trustModifier(t source.TrustLevel) int {
	switch t {
	case source.TrustOfficial:
		return 20
	case source.TrustVerified:
		return 15
	case source.TrustCommunity:
		return 5
	default:
		return 0
	}
}

// riskModifier returns the score adjustment for a risk level.
// Lower risk yields a higher bonus.
func riskModifier(r source.RiskLevel) int {
	switch r {
	case source.RiskLow:
		return 10
	case source.RiskMedium:
		return 5
	case source.RiskHigh:
		return 0
	case source.RiskCritical:
		return -20
	default:
		return 0
	}
}

// preferenceModifier returns additional score based on user preferences.
func preferenceModifier(t source.Type, prefs UserPreferences) int {
	if prefs.PreferNative && nativeTypes[t] {
		return 15
	}
	if prefs.PreferFlatpak && t == source.TypeFlatpak {
		return 15
	}
	return 0
}

// scoreSource computes the total ranking score for a source given the system
// context. Returns -1 if the source is disqualified (blocked or unavailable).
func scoreSource(s source.Source, ctx SystemContext) int {
	// Blocked sources are never recommended.
	if s.TrustLevel == source.TrustBlocked {
		return -1
	}

	// Source type must be available on the current system.
	available := false
	for _, m := range ctx.AvailableManagers {
		if m == s.SourceType {
			available = true
			break
		}
	}
	if !available {
		return -1
	}

	return baseTypePriority(s.SourceType) +
		trustModifier(s.TrustLevel) +
		riskModifier(s.RiskLevel) +
		preferenceModifier(s.SourceType, ctx.UserPreferences)
}
