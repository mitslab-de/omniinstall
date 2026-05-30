package discovery_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mitslab-de/omniinstall/internal/discovery"
)

func TestLoadCatalogFile_YAMLArray(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "catalog.yaml")
	content := `
- id: test-browser
  display_name: Test Browser
  summary: Browser app for loader test.
  categories: [browsers]
  aliases: [tbrowser]
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	catalog, err := discovery.LoadCatalogFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if catalog.Get("test-browser") == nil {
		t.Fatal("expected test-browser in loaded catalog")
	}
}

func TestLoadCatalogFile_JSONEnvelope(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "catalog.json")
	content := `{
  "applications": [
    {
      "id": "test-editor",
      "display_name": "Test Editor",
      "summary": "Editor app for loader test.",
      "categories": ["development"],
      "aliases": ["teditor"]
    }
  ]
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	catalog, err := discovery.LoadCatalogFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if catalog.Get("test-editor") == nil {
		t.Fatal("expected test-editor in loaded catalog")
	}
}

func TestLoadCatalogWithFallback_MissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.yaml")
	catalog, err := discovery.LoadCatalogWithFallback(path)
	if err == nil {
		t.Fatal("expected fallback warning error for missing catalog file")
	}
	if !strings.Contains(err.Error(), "using embedded MVP catalog") {
		t.Fatalf("expected fallback warning in error, got: %v", err)
	}
	if catalog.Get("obs-studio") == nil {
		t.Fatal("expected embedded MVP catalog after fallback")
	}
}

func TestLoadCatalogWithFallback_InvalidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "invalid.yaml")
	if err := os.WriteFile(path, []byte(`invalid: true`), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	catalog, err := discovery.LoadCatalogWithFallback(path)
	if err == nil {
		t.Fatal("expected fallback warning error for invalid catalog file")
	}
	if catalog.Get("obs-studio") == nil {
		t.Fatal("expected embedded MVP catalog after fallback")
	}
}

func TestLoadCatalogFile_InvalidContentReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "invalid.yaml")
	if err := os.WriteFile(path, []byte(`applications: []`), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if _, err := discovery.LoadCatalogFile(path); err == nil {
		t.Fatal("expected error for empty catalog")
	}
}

func TestSearch_DeterministicWithLoadedCatalog(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "catalog.yaml")
	content := `
- id: alpha-tool
  display_name: Alpha Tool
  summary: Alpha helper.
  categories: [development]
  aliases: [tool]
- id: beta-tool
  display_name: Beta Tool
  summary: Beta helper.
  categories: [development]
  aliases: [tool]
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	catalog, err := discovery.LoadCatalogFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	engine := discovery.NewLocalEngine(catalog)

	first, err := engine.Search("tool")
	if err != nil {
		t.Fatalf("unexpected search error: %v", err)
	}
	second, err := engine.Search("tool")
	if err != nil {
		t.Fatalf("unexpected search error: %v", err)
	}
	if len(first) != len(second) {
		t.Fatalf("non-deterministic count: %d vs %d", len(first), len(second))
	}
	for i := range first {
		if first[i].Application.ID != second[i].Application.ID {
			t.Fatalf("non-deterministic order at index %d: %q vs %q", i, first[i].Application.ID, second[i].Application.ID)
		}
	}
}
