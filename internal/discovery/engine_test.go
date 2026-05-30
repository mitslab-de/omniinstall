package discovery_test

import (
	"testing"

	"github.com/mitslab-de/omniinstall/internal/app"
	"github.com/mitslab-de/omniinstall/internal/discovery"
	"github.com/mitslab-de/omniinstall/internal/logging"
)

// testCatalog builds a small catalog for engine tests.
func testCatalog(t *testing.T) *discovery.Catalog {
	t.Helper()
	c := discovery.NewCatalog()
	apps := []*app.Application{
		{
			ID:          "obs-studio",
			DisplayName: "OBS Studio",
			Summary:     "Video recording and live streaming software.",
			Categories:  []string{"video", "streaming"},
			Aliases:     []string{"obs", "open broadcaster software"},
		},
		{
			ID:          "vlc",
			DisplayName: "VLC",
			Summary:     "Versatile media player.",
			Categories:  []string{"video", "audio"},
			Aliases:     []string{"vlc media player", "videolan"},
		},
		{
			ID:          "visual-studio-code",
			DisplayName: "Visual Studio Code",
			Summary:     "Code editor.",
			Categories:  []string{"development"},
			Aliases:     []string{"vscode", "vs code", "code"},
		},
		{
			ID:          "firefox",
			DisplayName: "Firefox",
			Summary:     "Fast, private and secure web browser.",
			Categories:  []string{"browsers"},
			Aliases:     []string{"mozilla firefox", "mozilla"},
		},
		{
			ID:          "git",
			DisplayName: "Git",
			Summary:     "Distributed version control system.",
			Categories:  []string{"development"},
			Aliases:     []string{"git scm"},
		},
	}
	for _, a := range apps {
		if err := c.Add(a); err != nil {
			t.Fatalf("failed to add application %q to catalog: %v", a.ID, err)
		}
	}
	return c
}

func newEngine(t *testing.T) *discovery.LocalEngine {
	t.Helper()
	return discovery.NewLocalEngine(testCatalog(t))
}

// --- Search: exact matches ---

func TestSearch_ExactID(t *testing.T) {
	e := newEngine(t)
	results, err := e.Search("obs-studio")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one result, got none")
	}
	if results[0].Application.ID != "obs-studio" {
		t.Errorf("expected top result 'obs-studio', got %q", results[0].Application.ID)
	}
	if results[0].MatchScore != 100 {
		t.Errorf("expected score 100 for exact ID match, got %d", results[0].MatchScore)
	}
}

func TestSearch_ExactDisplayName(t *testing.T) {
	e := newEngine(t)
	results, err := e.Search("OBS Studio")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one result, got none")
	}
	if results[0].Application.ID != "obs-studio" {
		t.Errorf("expected top result 'obs-studio', got %q", results[0].Application.ID)
	}
	if results[0].MatchScore != 90 {
		t.Errorf("expected score 90 for exact display name match, got %d", results[0].MatchScore)
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	e := newEngine(t)
	results, err := e.Search("FIREFOX")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one result, got none")
	}
	if results[0].Application.ID != "firefox" {
		t.Errorf("expected top result 'firefox', got %q", results[0].Application.ID)
	}
}

// --- Search: alias matches ---

func TestSearch_ExactAlias(t *testing.T) {
	e := newEngine(t)
	results, err := e.Search("obs")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one result, got none")
	}
	if results[0].Application.ID != "obs-studio" {
		t.Errorf("expected top result 'obs-studio', got %q", results[0].Application.ID)
	}
	if results[0].MatchScore != 80 {
		t.Errorf("expected score 80 for alias match, got %d", results[0].MatchScore)
	}
}

func TestSearch_AliasVSCode(t *testing.T) {
	e := newEngine(t)
	results, err := e.Search("vscode")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one result, got none")
	}
	if results[0].Application.ID != "visual-studio-code" {
		t.Errorf("expected top result 'visual-studio-code', got %q", results[0].Application.ID)
	}
}

// --- Search: partial matches ---

func TestSearch_PartialID(t *testing.T) {
	e := newEngine(t)
	results, err := e.Search("obs")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected results for partial ID, got none")
	}
	// 'obs' is an alias for obs-studio (score 80), not a prefix of the ID "obs-studio" (which would be score 70).
	// The alias match wins.
	found := false
	for _, r := range results {
		if r.Application.ID == "obs-studio" {
			found = true
		}
	}
	if !found {
		t.Error("expected obs-studio in results for query 'obs'")
	}
}

func TestSearch_PartialDisplayName(t *testing.T) {
	e := newEngine(t)
	results, err := e.Search("visual")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected results for partial display name, got none")
	}
	if results[0].Application.ID != "visual-studio-code" {
		t.Errorf("expected top result 'visual-studio-code', got %q", results[0].Application.ID)
	}
}

func TestSearch_ContainsMatch(t *testing.T) {
	e := newEngine(t)
	results, err := e.Search("studio")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected results for contains match, got none")
	}
	// "studio" is in both "obs-studio" (ID contains) and "visual studio code" (display name contains).
	found := false
	for _, r := range results {
		if r.Application.ID == "obs-studio" {
			found = true
		}
	}
	if !found {
		t.Error("expected obs-studio in results for query 'studio'")
	}
}

// --- Search: ranking ---

func TestSearch_RankingIsDeterministic(t *testing.T) {
	e := newEngine(t)
	results1, err := e.Search("code")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	results2, err := e.Search("code")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results1) != len(results2) {
		t.Fatalf("non-deterministic result count: %d vs %d", len(results1), len(results2))
	}
	for i := range results1 {
		if results1[i].Application.ID != results2[i].Application.ID {
			t.Errorf("non-deterministic order at index %d: %q vs %q",
				i, results1[i].Application.ID, results2[i].Application.ID)
		}
	}
}

func TestSearch_HigherScoreFirst(t *testing.T) {
	e := newEngine(t)
	results, err := e.Search("vlc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) < 1 {
		t.Fatal("expected at least one result")
	}
	for i := 1; i < len(results); i++ {
		if results[i].MatchScore > results[i-1].MatchScore {
			t.Errorf("results not sorted by descending score: index %d (%d) > index %d (%d)",
				i, results[i].MatchScore, i-1, results[i-1].MatchScore)
		}
	}
}

// --- Search: unknown / empty queries ---

func TestSearch_UnknownQuery(t *testing.T) {
	e := newEngine(t)
	results, err := e.Search("completelymadeupname12345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty results for unknown query, got %d", len(results))
	}
}

func TestSearch_EmptyQuery(t *testing.T) {
	e := newEngine(t)
	_, err := e.Search("")
	if err == nil {
		t.Error("expected error for empty query, got nil")
	}
}

func TestSearch_WhitespaceOnlyQuery(t *testing.T) {
	e := newEngine(t)
	_, err := e.Search("   ")
	if err == nil {
		t.Error("expected error for whitespace-only query, got nil")
	}
}

// --- Lookup ---

func TestLookup_ExistingID(t *testing.T) {
	e := newEngine(t)
	a, err := e.Lookup("git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected application, got nil")
	}
	if a.ID != "git" {
		t.Errorf("expected application id 'git', got %q", a.ID)
	}
}

func TestLookup_UnknownID(t *testing.T) {
	e := newEngine(t)
	a, err := e.Lookup("does-not-exist")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a != nil {
		t.Errorf("expected nil for unknown id, got %q", a.ID)
	}
}

func TestSearch_EmitsStructuredLogOnSuccess(t *testing.T) {
	e := newEngine(t)
	var entries []logging.Entry
	e.WithLogger(logging.EmitFunc(func(entry logging.Entry) {
		entries = append(entries, entry)
	}))

	if _, err := e.Search("obs-studio"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(entries))
	}
	if entries[0].Action != logging.ActionSearch {
		t.Fatalf("expected action %q, got %q", logging.ActionSearch, entries[0].Action)
	}
	if entries[0].Result != "success" {
		t.Fatalf("expected success result, got %q", entries[0].Result)
	}
	if entries[0].Duration <= 0 {
		t.Fatalf("expected positive duration, got %v", entries[0].Duration)
	}
}

func TestSearch_EmitsStructuredLogOnFailure(t *testing.T) {
	e := newEngine(t)
	var entries []logging.Entry
	e.WithLogger(logging.EmitFunc(func(entry logging.Entry) {
		entries = append(entries, entry)
	}))

	if _, err := e.Search(" "); err == nil {
		t.Fatal("expected error for empty query")
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(entries))
	}
	if entries[0].Result != "failure" {
		t.Fatalf("expected failure result, got %q", entries[0].Result)
	}
	if entries[0].ErrorCategory != "invalid_query" {
		t.Fatalf("expected invalid_query, got %q", entries[0].ErrorCategory)
	}
}

// --- Catalog ---

func TestCatalog_Add_DuplicateID(t *testing.T) {
	c := discovery.NewCatalog()
	a := &app.Application{
		ID:          "vlc",
		DisplayName: "VLC",
		Summary:     "A player.",
		Categories:  []string{"video"},
	}
	if err := c.Add(a); err != nil {
		t.Fatalf("first Add failed: %v", err)
	}
	if err := c.Add(a); err == nil {
		t.Error("expected error for duplicate ID, got nil")
	}
}

func TestCatalog_Add_InvalidApplication(t *testing.T) {
	c := discovery.NewCatalog()
	invalid := &app.Application{}
	if err := c.Add(invalid); err == nil {
		t.Error("expected validation error for invalid application, got nil")
	}
}

func TestCatalog_Size(t *testing.T) {
	c := testCatalog(t)
	if c.Size() != 5 {
		t.Errorf("expected catalog size 5, got %d", c.Size())
	}
}

func TestCatalog_Get_Existing(t *testing.T) {
	c := testCatalog(t)
	a := c.Get("vlc")
	if a == nil {
		t.Fatal("expected application for 'vlc', got nil")
	}
	if a.ID != "vlc" {
		t.Errorf("expected id 'vlc', got %q", a.ID)
	}
}

func TestCatalog_Get_Unknown(t *testing.T) {
	c := testCatalog(t)
	if a := c.Get("unknown"); a != nil {
		t.Errorf("expected nil for unknown id, got %q", a.ID)
	}
}

// --- MVP seed catalog ---

func TestMVPCatalog_HasRequiredApplications(t *testing.T) {
	c := discovery.MVPCatalog()
	required := []string{"obs-studio", "vlc", "firefox", "bitwarden", "git", "docker", "visual-studio-code"}
	for _, id := range required {
		if a := c.Get(id); a == nil {
			t.Errorf("expected application %q in MVP catalog, not found", id)
		}
	}
}

func TestMVPCatalog_AllApplicationsAreValid(t *testing.T) {
	c := discovery.MVPCatalog()
	for _, a := range c.All() {
		if err := a.Validate(); err != nil {
			t.Errorf("application %q is invalid: %v", a.ID, err)
		}
	}
}
