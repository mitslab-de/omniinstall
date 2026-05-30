package discovery

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mitslab-de/omniinstall/internal/app"
	"github.com/mitslab-de/omniinstall/internal/logging"
)

// LocalEngine is a Discovery Engine backed by a local Catalog.
// It implements the Engine interface.
type LocalEngine struct {
	catalog *Catalog
	logger  logging.Emitter
	now     func() time.Time
}

// NewLocalEngine creates a new LocalEngine backed by the given Catalog.
func NewLocalEngine(catalog *Catalog) *LocalEngine {
	return &LocalEngine{
		catalog: catalog,
		now:     time.Now,
	}
}

// WithLogger configures structured action logging for the engine.
func (e *LocalEngine) WithLogger(logger logging.Emitter) *LocalEngine {
	e.logger = logger
	return e
}

func (e *LocalEngine) emitLog(entry logging.Entry) {
	if e.logger == nil {
		return
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = e.now()
	}
	e.logger.Emit(entry)
}

// Search returns candidates from the catalog that match the query.
//
// Matching is case-insensitive. The results are ordered by descending match
// score, making the engine deterministic for identical inputs and catalogs.
//
// Returns an empty slice (never nil) when no matches are found.
// Returns a descriptive error when the query is empty.
func (e *LocalEngine) Search(query string) ([]*Candidate, error) {
	start := time.Now()
	trimmedQuery := strings.TrimSpace(query)

	if strings.TrimSpace(query) == "" {
		err := fmt.Errorf("search query must not be empty")
		e.emitLog(logging.Entry{
			Action:        logging.ActionSearch,
			ApplicationID: trimmedQuery,
			Result:        "failure",
			Duration:      time.Since(start),
			ErrorCategory: "invalid_query",
		})
		return nil, err
	}

	q := strings.ToLower(trimmedQuery)
	var results []*Candidate

	for _, a := range e.catalog.All() {
		score := scoreApplication(a, q)
		if score > 0 {
			results = append(results, &Candidate{
				Application: a,
				MatchScore:  score,
			})
		}
	}

	// Sort by descending score, then by ID for determinism on equal scores.
	sort.Slice(results, func(i, j int) bool {
		if results[i].MatchScore != results[j].MatchScore {
			return results[i].MatchScore > results[j].MatchScore
		}
		return results[i].Application.ID < results[j].Application.ID
	})

	if results == nil {
		e.emitLog(logging.Entry{
			Action:        logging.ActionSearch,
			ApplicationID: trimmedQuery,
			Result:        "success",
			Duration:      time.Since(start),
		})
		return []*Candidate{}, nil
	}
	e.emitLog(logging.Entry{
		Action:        logging.ActionSearch,
		ApplicationID: trimmedQuery,
		Result:        "success",
		Duration:      time.Since(start),
	})
	return results, nil
}

// Lookup returns the application with the exact canonical ID, or nil if not found.
func (e *LocalEngine) Lookup(applicationID string) (*app.Application, error) {
	a := e.catalog.Get(applicationID)
	return a, nil
}

// scoreApplication returns a relevance score for the application against the query.
//
// Scoring tiers (higher is better):
//
//	100  exact ID match
//	 90  exact display name match (case-insensitive)
//	 85  normalised exact ID match (hyphens/underscores treated as spaces)
//	 80  exact alias match
//	 75  normalised exact alias match
//	 70  ID prefix match
//	 65  normalised ID prefix match or word-segment prefix
//	 60  display name prefix match
//	 50  alias prefix match
//	 45  word-segment ID contains (query == a full hyphenated word in the ID)
//	 40  ID contains query
//	 30  display name contains query
//	 20  alias contains query
//
// Returns 0 if there is no match.
func scoreApplication(a *app.Application, q string) int {
	id := strings.ToLower(a.ID)
	name := strings.ToLower(a.DisplayName)
	qNorm := normalizeIdentifier(q)
	idNorm := normalizeIdentifier(id)

	if id == q {
		return 100
	}
	if name == q {
		return 90
	}
	if idNorm == qNorm && idNorm != "" {
		return 85
	}
	for _, alias := range a.Aliases {
		al := strings.ToLower(alias)
		if al == q {
			return 80
		}
		if normalizeIdentifier(al) == qNorm && qNorm != "" {
			return 75
		}
	}
	if strings.HasPrefix(id, q) {
		return 70
	}
	if strings.HasPrefix(idNorm, qNorm) && qNorm != "" {
		return 65
	}
	if strings.HasPrefix(name, q) {
		return 60
	}
	for _, alias := range a.Aliases {
		if strings.HasPrefix(strings.ToLower(alias), q) {
			return 50
		}
	}
	// Word-segment match: query equals one of the hyphen-separated words in the ID.
	for _, seg := range strings.Split(id, "-") {
		if seg == q {
			return 45
		}
	}
	if strings.Contains(id, q) {
		return 40
	}
	if strings.Contains(name, q) {
		return 30
	}
	for _, alias := range a.Aliases {
		if strings.Contains(strings.ToLower(alias), q) {
			return 20
		}
	}
	return 0
}

// NormalizeIdentifier converts a query or identifier to a canonical form for
// fuzzy comparison by lowercasing and collapsing hyphens, underscores, and
// spaces into a single space, then trimming.
func NormalizeIdentifier(s string) string {
	return normalizeIdentifier(s)
}

// normalizeIdentifier is the unexported implementation.
func normalizeIdentifier(s string) string {
	var b strings.Builder
	prev := ' '
	for _, ch := range strings.ToLower(s) {
		if ch == '-' || ch == '_' || ch == ' ' {
			ch = ' '
		}
		if ch == ' ' && prev == ' ' {
			continue
		}
		b.WriteRune(ch)
		prev = ch
	}
	return strings.TrimSpace(b.String())
}
