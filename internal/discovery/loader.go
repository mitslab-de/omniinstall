package discovery

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitslab-de/omniinstall/internal/app"
	"gopkg.in/yaml.v3"
)

type catalogEnvelope struct {
	Applications []*app.Application `json:"applications" yaml:"applications"`
}

// LoadCatalogFile loads a catalog from a YAML or JSON file.
// Supported structures are either:
//   - a top-level array of applications
//   - an object with an "applications" field
func LoadCatalogFile(path string) (*Catalog, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("catalog path must not be empty")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read catalog file: %w", err)
	}

	apps, err := parseCatalogApplications(content)
	if err != nil {
		return nil, err
	}

	catalog := NewCatalog()
	for i, a := range apps {
		if err := catalog.Add(a); err != nil {
			return nil, fmt.Errorf("invalid application at index %d: %w", i, err)
		}
	}
	return catalog, nil
}

// LoadCatalogWithFallback loads a catalog file when configured, and falls back
// to the embedded MVP catalog if loading fails.
func LoadCatalogWithFallback(path string) (*Catalog, error) {
	if strings.TrimSpace(path) == "" {
		return MVPCatalog(), nil
	}

	catalog, err := LoadCatalogFile(path)
	if err != nil {
		return MVPCatalog(), fmt.Errorf("failed to load local catalog %q, using embedded MVP catalog: %w", path, err)
	}

	return catalog, nil
}

func parseCatalogApplications(content []byte) ([]*app.Application, error) {
	var list []*app.Application
	if err := yaml.Unmarshal(content, &list); err == nil && len(list) > 0 {
		return list, nil
	}

	var env catalogEnvelope
	if err := yaml.Unmarshal(content, &env); err == nil && len(env.Applications) > 0 {
		return env.Applications, nil
	}

	return nil, fmt.Errorf("catalog must contain at least one application")
}
