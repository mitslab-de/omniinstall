package resolver

import (
	"errors"
	"sort"
	"time"

	"github.com/mitslab-de/omniinstall/internal/install"
	"github.com/mitslab-de/omniinstall/internal/logging"
	"github.com/mitslab-de/omniinstall/internal/source"
)

// DefaultResolver is the standard implementation of the Resolver interface.
// It ranks sources by type priority, trust level, risk level, and user
// preferences, producing a deterministic ordered recommendation list.
type DefaultResolver struct {
	logger         logging.Emitter
	now            func() time.Time
	priorityConfig PriorityConfig
}

// NewDefaultResolver creates a DefaultResolver with default settings.
func NewDefaultResolver() *DefaultResolver {
	return &DefaultResolver{now: time.Now}
}

// WithLogger configures structured action logging for resolver operations.
func (r *DefaultResolver) WithLogger(logger logging.Emitter) *DefaultResolver {
	r.logger = logger
	return r
}

// WithPriorityConfig replaces the default source type priority table.
// Use this to adjust the preference ordering per distribution or user policy.
func (r *DefaultResolver) WithPriorityConfig(cfg PriorityConfig) *DefaultResolver {
	r.priorityConfig = cfg
	return r
}

func (r *DefaultResolver) emitLog(entry logging.Entry) {
	if r.logger == nil {
		return
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = r.now()
	}
	r.logger.Emit(entry)
}

// Resolve ranks all available sources for the given applicationID and returns
// an ordered slice of Recommendations. The first entry is the top choice.
// Sources that are disqualified (blocked or manager unavailable) are omitted.
// Returns an empty slice (no error) when no compatible source exists.
func (r *DefaultResolver) Resolve(
	applicationID string,
	available []source.Source,
	ctx SystemContext,
) ([]Recommendation, error) {
	start := time.Now()
	if applicationID == "" {
		err := errors.New("applicationID must not be empty")
		r.emitLog(logging.Entry{
			Action:        logging.ActionResolve,
			ApplicationID: applicationID,
			Result:        "failure",
			Duration:      time.Since(start),
			ErrorCategory: "invalid_application_id",
		})
		return nil, err
	}

	type scored struct {
		src   source.Source
		score int
	}

	var candidates []scored
	for _, s := range available {
		sc := scoreSource(s, ctx, r.priorityConfig)
		if sc < 0 {
			continue // disqualified
		}
		candidates = append(candidates, scored{src: s, score: sc})
	}

	// Sort descending by score; ties broken by SourceIdentifier and then
	// SourceType for determinism independent of input ordering.
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].score != candidates[j].score {
			return candidates[i].score > candidates[j].score
		}
		if candidates[i].src.SourceIdentifier != candidates[j].src.SourceIdentifier {
			return candidates[i].src.SourceIdentifier < candidates[j].src.SourceIdentifier
		}
		return candidates[i].src.SourceType < candidates[j].src.SourceType
	})

	recommendations := make([]Recommendation, 0, len(candidates))
	for idx, c := range candidates {
		var exp string
		if idx == 0 {
			exp = explain(c.src, c.score, ctx)
		} else {
			exp = explainAlternative(c.src)
		}

		plan := buildPlan(applicationID, c.src, exp)
		recommendations = append(recommendations, Recommendation{
			Plan:        plan,
			Explanation: exp,
		})
	}

	entry := logging.Entry{
		Action:        logging.ActionResolve,
		ApplicationID: applicationID,
		Result:        "success",
		Duration:      time.Since(start),
	}
	if len(recommendations) > 0 && recommendations[0].Plan != nil {
		entry.SourceType = recommendations[0].Plan.SourceType
	}
	r.emitLog(entry)

	return recommendations, nil
}

// buildPlan creates an install.Plan from a source and its explanation.
func buildPlan(applicationID string, s source.Source, explanation string) *install.Plan {
	requiresPrivilege := nativeTypes[s.SourceType] || s.SourceType == source.TypeSnap
	requiresConfirmation := s.RiskLevel == source.RiskHigh || s.RiskLevel == source.RiskCritical

	plan := &install.Plan{
		ApplicationID:        applicationID,
		SourceType:           s.SourceType,
		SourceIdentifier:     s.SourceIdentifier,
		RiskLevel:            s.RiskLevel,
		RequiresPrivilege:    requiresPrivilege,
		RequiresConfirmation: requiresConfirmation,
		Explanation:          explanation,
		Verification:         []install.VerificationRule{{Command: s.SourceIdentifier}},
	}
	return plan
}
