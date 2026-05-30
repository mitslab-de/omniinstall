package discovery

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mitslab-de/omniinstall/internal/app"
)

// LocalEngine is a Discovery Engine backed by a local Catalog.
// It implements the Engine interface.
type LocalEngine struct {
	catalog *Catalog
}

// NewLocalEngine creates a new LocalEngine backed by the given Catalog.
func NewLocalEngine(catalog *Catalog) *LocalEngine {
	return &LocalEngine{catalog: catalog}
}

// Search returns candidates from the catalog that match the query.
//
// Matching is case-insensitive. The results are ordered by descending match
// score, making the engine deterministic for identical inputs and catalogs.
//
// Returns an empty slice (never nil) when no matches are found.
// Returns a descriptive error when the query is empty.
func (e *LocalEngine) Search(query string) ([]*Candidate, error) {
	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("search query must not be empty")
	}

	q := strings.ToLower(strings.TrimSpace(query))
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
		return []*Candidate{}, nil
	}
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
//	 80  exact alias match
//	 70  ID prefix match
//	 60  display name prefix match
//	 50  alias prefix match
//	 40  ID contains query
//	 30  display name contains query
//	 20  alias contains query
//
// Returns 0 if there is no match.
func scoreApplication(a *app.Application, q string) int {
	id := strings.ToLower(a.ID)
	name := strings.ToLower(a.DisplayName)

	if id == q {
		return 100
	}
	if name == q {
		return 90
	}
	for _, alias := range a.Aliases {
		if strings.ToLower(alias) == q {
			return 80
		}
	}
	if strings.HasPrefix(id, q) {
		return 70
	}
	if strings.HasPrefix(name, q) {
		return 60
	}
	for _, alias := range a.Aliases {
		if strings.HasPrefix(strings.ToLower(alias), q) {
			return 50
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
