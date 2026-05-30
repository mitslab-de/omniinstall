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
	required := []string{"obs-studio", "vlc", "firefox", "bitwarden", "git", "docker", "visual-studio-code", "curl", "htop", "neovim"}
	for _, id := range required {
		if a := c.Get(id); a == nil {
			t.Errorf("expected application %q in MVP catalog, not found", id)
		}
	}
}

func TestMVPCatalog_HasMinimumApplicationCount(t *testing.T) {
	c := discovery.MVPCatalog()
	const minApps = 10
	if got := len(c.All()); got < minApps {
		t.Errorf("MVP catalog should have at least %d applications, got %d", minApps, got)
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

// ---------------------------------------------------------------------------
// Alias normalisation and word-segment scoring tests (task 0026)
// ---------------------------------------------------------------------------

func TestSearch_NormalisedID_HyphenAsSpace(t *testing.T) {
// Query "obs-studio" with underscores should match ID "obs-studio" via
// normalised exact ID match (score 85) when there is no exact ID match.
// We use a catalog entry where the display name does NOT match the query.
c := discovery.NewCatalog()
_ = c.Add(&app.Application{
ID:          "obs-studio",
DisplayName: "Open Broadcaster Software",
Summary:     "Streaming.",
Categories:  []string{"video"},
})
e := discovery.NewLocalEngine(c)
// Query uses underscores instead of hyphens.
results, err := e.Search("obs_studio")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(results) == 0 {
t.Fatal("expected results for normalised query, got none")
}
if results[0].Application.ID != "obs-studio" {
t.Errorf("expected top result 'obs-studio', got %q", results[0].Application.ID)
}
if results[0].MatchScore != 85 {
t.Errorf("expected score 85 for normalised ID match, got %d", results[0].MatchScore)
}
}

func TestSearch_NormalisedAlias_SpaceVsHyphen(t *testing.T) {
// "vs-code" query should match alias "vs code" via normalised alias match (score 75).
c := discovery.NewCatalog()
_ = c.Add(&app.Application{
ID:          "visual-studio-code",
DisplayName: "Visual Studio Code",
Summary:     "Code editor.",
Categories:  []string{"development"},
Aliases:     []string{"vs code"},
})
e := discovery.NewLocalEngine(c)
results, err := e.Search("vs-code")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(results) == 0 {
t.Fatal("expected results for normalised alias query, got none")
}
if results[0].Application.ID != "visual-studio-code" {
t.Errorf("expected 'visual-studio-code', got %q", results[0].Application.ID)
}
if results[0].MatchScore != 75 {
t.Errorf("expected score 75 for normalised alias match, got %d", results[0].MatchScore)
}
}

func TestSearch_WordSegment_Studio(t *testing.T) {
// "studio" matches both "obs-studio" (word segment, score 45) and
// "visual-studio-code" (word segment, score 45) and possibly display name contains.
// obs-studio gets score 45 (word segment); visual-studio-code also gets 45.
// visual studio code display name contains "studio" → score 30.
// Both IDs have "studio" as a segment → score 45 each. obs-studio ID also
// contains "studio" → score 40. Since word-segment check (45) fires before
// contains (40), obs-studio gets 45.
e := newEngine(t)
results, err := e.Search("studio")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(results) == 0 {
t.Fatal("expected results for word-segment query, got none")
}
foundOBS, foundVSC := false, false
for _, r := range results {
if r.Application.ID == "obs-studio" {
foundOBS = true
if r.MatchScore < 40 {
t.Errorf("obs-studio score too low for 'studio' query: %d", r.MatchScore)
}
}
if r.Application.ID == "visual-studio-code" {
foundVSC = true
if r.MatchScore < 30 {
t.Errorf("visual-studio-code score too low for 'studio' query: %d", r.MatchScore)
}
}
}
if !foundOBS {
t.Error("expected obs-studio in results for query 'studio'")
}
if !foundVSC {
t.Error("expected visual-studio-code in results for query 'studio'")
}
}

func TestSearch_NormalisedID_BeatsPartialMatch(t *testing.T) {
// Normalised exact match (85) must outrank alias prefix (50) and
// prefix ID match (70). We use a display name that does NOT match the
// query to avoid the display name exact match (90) interfering.
c := discovery.NewCatalog()
_ = c.Add(&app.Application{
ID:          "obs-studio",
DisplayName: "Open Broadcaster Studio",
Summary:     "Streaming.",
Categories:  []string{"video"},
Aliases:     []string{"obs recorder"},
})
_ = c.Add(&app.Application{
ID:          "obs-remote",
DisplayName: "Open Broadcaster Remote",
Summary:     "Control OBS.",
Categories:  []string{"video"},
})
e := discovery.NewLocalEngine(c)
// Query "obs_studio" normalises to "obs studio" which equals normalised ID "obs studio".
results, err := e.Search("obs_studio")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(results) == 0 {
t.Fatal("expected results")
}
if results[0].Application.ID != "obs-studio" {
t.Errorf("expected 'obs-studio' to rank first, got %q", results[0].Application.ID)
}
if results[0].MatchScore != 85 {
t.Errorf("expected score 85 (normalised exact ID), got %d", results[0].MatchScore)
}
}

func TestSearch_ExactID_BeatsNormalised(t *testing.T) {
// An exact display name match (90) must always outrank a normalised ID match (85).
// "obs studio" (query) matches display name "OBS Studio" exactly (90)
// and also matches normalised ID "obs-studio" (85). The display name wins.
e := newEngine(t)
results, err := e.Search("obs studio")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(results) == 0 {
t.Fatal("expected results")
}
if results[0].Application.ID != "obs-studio" {
t.Errorf("expected 'obs-studio' first, got %q", results[0].Application.ID)
}
// Display name "OBS Studio" exactly matches query "obs studio" → score 90.
if results[0].MatchScore != 90 {
t.Errorf("expected score 90 (exact display name beats normalised ID), got %d", results[0].MatchScore)
}
}

func TestNormalizeIdentifier_Various(t *testing.T) {
cases := []struct {
input string
want  string
}{
{"obs-studio", "obs studio"},
{"obs_studio", "obs studio"},
{"obs studio", "obs studio"},
{"OBS-Studio", "obs studio"},
{"  obs  studio  ", "obs studio"},
{"visual-studio-code", "visual studio code"},
{"", ""},
{"git", "git"},
}
for _, tc := range cases {
got := discovery.NormalizeIdentifier(tc.input)
if got != tc.want {
t.Errorf("NormalizeIdentifier(%q) = %q, want %q", tc.input, got, tc.want)
}
}
}
