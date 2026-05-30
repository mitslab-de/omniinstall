package resolver

import "github.com/mitslab-de/omniinstall/internal/source"

// ConflictKind classifies the type of source conflict detected.
type ConflictKind string

const (
	// ConflictUnavailableManager is reported when the required package manager
	// is not installed on the current system.
	ConflictUnavailableManager ConflictKind = "unavailable_manager"

	// ConflictBlockedSource is reported when a source has been blocked by policy.
	ConflictBlockedSource ConflictKind = "blocked_source"

	// ConflictNoCompatibleSource is reported when no sources are compatible
	// with the current system context.
	ConflictNoCompatibleSource ConflictKind = "no_compatible_source"
)

// ConflictInfo describes a conflict that the resolver detected.
// Conflicts are surfaced rather than silently suppressed.
type ConflictInfo struct {
	// Kind classifies the conflict.
	Kind ConflictKind

	// Source is the source that triggered the conflict (may be zero value).
	Source source.Source

	// Message is a user-facing description of the conflict.
	Message string

	// Suggestion is an actionable next step the user can take to resolve the conflict.
	Suggestion string
}

// DetectConflicts inspects all candidate sources and returns any conflicts
// found. A non-empty slice means at least one source cannot be used.
// Conflicts do not prevent the resolver from recommending other sources.
//
// A ConflictNoCompatibleSource entry is appended when every source is
// excluded and no recommendation is possible.
func DetectConflicts(available []source.Source, ctx SystemContext) []ConflictInfo {
	var conflicts []ConflictInfo

	managerSet := make(map[source.Type]bool, len(ctx.AvailableManagers))
	for _, m := range ctx.AvailableManagers {
		managerSet[m] = true
	}

	usable := 0
	for _, s := range available {
		if s.TrustLevel == source.TrustBlocked {
			conflicts = append(conflicts, ConflictInfo{
				Kind:       ConflictBlockedSource,
				Source:     s,
				Message:    "Source " + s.SourceIdentifier + " (" + string(s.SourceType) + ") is blocked by policy.",
				Suggestion: "Contact your administrator or choose a different source.",
			})
			continue
		}
		if !managerSet[s.SourceType] {
			conflicts = append(conflicts, ConflictInfo{
				Kind:       ConflictUnavailableManager,
				Source:     s,
				Message:    "Package manager " + string(s.SourceType) + " is not available on this system.",
				Suggestion: "Install " + string(s.SourceType) + " to enable this source, or try a different source type.",
			})
			continue
		}
		usable++
	}

	// If every source was excluded, emit a summary conflict.
	if len(available) > 0 && usable == 0 {
		conflicts = append(conflicts, ConflictInfo{
			Kind:       ConflictNoCompatibleSource,
			Message:    "No compatible source is available for this application on your system.",
			Suggestion: "Check that a supported package manager (apt, flatpak, snap) is installed, or search for an alternative application.",
		})
	}

	return conflicts
}
