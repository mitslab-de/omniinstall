package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	adapter "github.com/mitslab-de/omniinstall/internal/adapters"
	aptadapter "github.com/mitslab-de/omniinstall/internal/adapters/apt"
	dnfadapter "github.com/mitslab-de/omniinstall/internal/adapters/dnf"
	flatpakadapter "github.com/mitslab-de/omniinstall/internal/adapters/flatpak"
	pacmanadapter "github.com/mitslab-de/omniinstall/internal/adapters/pacman"
	snapadapter "github.com/mitslab-de/omniinstall/internal/adapters/snap"
	"github.com/mitslab-de/omniinstall/internal/discovery"
	"github.com/mitslab-de/omniinstall/internal/engine"
	"github.com/mitslab-de/omniinstall/internal/resolver"
	"github.com/mitslab-de/omniinstall/internal/source"
	"github.com/mitslab-de/omniinstall/internal/state"
	"gopkg.in/yaml.v3"
)

// App wires together all core services and implements the CLI commands.
// It holds no business logic — all decisions are delegated to service layers.
type App struct {
	discovery *discovery.LocalEngine
	resolver  *resolver.DefaultResolver
	engine    *engine.DefaultEngine
	store     state.Store
	sources   map[string][]source.Source
	ctx       resolver.SystemContext
	in        io.Reader
	out       io.Writer
}

type searchJSONResult struct {
	ApplicationID string `json:"application_id"`
	DisplayName   string `json:"display_name"`
	Summary       string `json:"summary"`
	MatchScore    int    `json:"match_score"`
}

type searchJSONOutput struct {
	Query   string             `json:"query"`
	Results []searchJSONResult `json:"results"`
}

type explainJSONAlternative struct {
	SourceType  source.Type `json:"source_type"`
	Explanation string      `json:"explanation"`
}

type explainJSONConflict struct {
	Kind       resolver.ConflictKind `json:"kind"`
	Message    string                `json:"message"`
	Suggestion string                `json:"suggestion,omitempty"`
}

type explainJSONOutput struct {
	ApplicationID          string                   `json:"application_id"`
	DisplayName            string                   `json:"display_name"`
	RecommendedSourceType  source.Type              `json:"recommended_source_type"`
	RecommendedExplanation string                   `json:"recommended_explanation"`
	TrustLevel             source.TrustLevel        `json:"trust_level,omitempty"`
	RiskLevel              source.RiskLevel         `json:"risk_level,omitempty"`
	RequiresPrivilege      bool                     `json:"requires_privilege,omitempty"`
	Alternatives           []explainJSONAlternative `json:"alternatives"`
	Conflicts              []explainJSONConflict    `json:"conflicts"`
}

type installDryRunJSONOutput struct {
	ApplicationID    string           `json:"application_id"`
	DisplayName      string           `json:"display_name"`
	Adapter          source.Type      `json:"adapter"`
	PackageID        string           `json:"package_id"`
	RiskLevel        source.RiskLevel `json:"risk_level"`
	RequiresPrivilege bool            `json:"requires_privilege,omitempty"`
	Explanation      string           `json:"explanation,omitempty"`
}

type listJSONApplication struct {
	ApplicationID string              `json:"application_id"`
	SourceType    source.Type         `json:"source_type"`
	Status        state.InstallStatus `json:"status"`
}

type listJSONOutput struct {
	Applications []listJSONApplication `json:"applications"`
}

// newApp creates a production App, detecting available package managers and
// wiring all services together.
func newApp(out io.Writer) *App {
	managers := detectAvailableManagers()
	catalog, catalogErr := discovery.LoadCatalogWithFallback(os.Getenv("OMNIINSTALL_CATALOG_PATH"))
	if catalogErr != nil {
		fmt.Fprintf(out, "Warning: %v\n", catalogErr)
	}

	adapters := []adapter.Adapter{
		aptadapter.New(),
		dnfadapter.New(),
		pacmanadapter.New(),
		snapadapter.New(),
		flatpakadapter.New(),
	}

	storePath := state.DefaultStorePath()

	return &App{
		discovery: discovery.NewLocalEngine(catalog),
		resolver:  resolver.NewDefaultResolver(),
		engine:    engine.NewDefaultEngine(adapters, nil),
		store:     state.NewFileStore(storePath),
		sources:   discovery.MVPSources(),
		ctx: resolver.SystemContext{
			AvailableManagers: managers,
		},
		in:  os.Stdin,
		out: out,
	}
}

// detectAvailableManagers returns the source types whose package managers are
// found in PATH on the current system.
func detectAvailableManagers() []source.Type {
	candidates := []struct {
		binary string
		t      source.Type
	}{
		{"apt-get", source.TypeAPT},
		{"dnf", source.TypeDNF},
		{"pacman", source.TypePacman},
		{"zypper", source.TypeZypper},
		{"flatpak", source.TypeFlatpak},
		{"snap", source.TypeSnap},
	}
	var available []source.Type
	for _, c := range candidates {
		if _, err := exec.LookPath(c.binary); err == nil {
			available = append(available, c.t)
		}
	}
	return available
}

// search handles `omni search <query>`.
func (a *App) search(query string) error {
	if query == "" {
		return errors.New("usage: omni search <query>")
	}

	candidates, err := a.discovery.Search(query)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}
	if len(candidates) == 0 {
		fmt.Fprintf(a.out, "No applications found for %q.\n", query)
		return nil
	}
	fmt.Fprintf(a.out, "Search results for %q:\n\n", query)
	for _, c := range candidates {
		fmt.Fprintf(a.out, "  %-25s %s\n", c.Application.ID, c.Application.Summary)
	}
	return nil
}

func (a *App) searchJSON(query string) error {
	if query == "" {
		return errors.New("usage: omni search <query>")
	}

	candidates, err := a.discovery.Search(query)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	out := searchJSONOutput{
		Query:   query,
		Results: make([]searchJSONResult, 0, len(candidates)),
	}
	for _, c := range candidates {
		out.Results = append(out.Results, searchJSONResult{
			ApplicationID: c.Application.ID,
			DisplayName:   c.Application.DisplayName,
			Summary:       c.Application.Summary,
			MatchScore:    c.MatchScore,
		})
	}
	return writeJSON(a.out, out)
}

// install handles `omni install <app>`.
func (a *App) install(appID string) error {
	if appID == "" {
		return errors.New("usage: omni install <application>")
	}

	// Discover the application.
	candidates, err := a.discovery.Search(appID)
	if err != nil {
		return fmt.Errorf("discovery failed: %w", err)
	}
	if len(candidates) == 0 {
		return fmt.Errorf("application %q not found", appID)
	}
	// Use the top-ranked candidate.
	app := candidates[0].Application
	resolvedID := app.ID

	// Resolve sources.
	srcs, ok := a.sources[resolvedID]
	if !ok || len(srcs) == 0 {
		return fmt.Errorf("no sources available for %q", resolvedID)
	}

	recs, err := a.resolver.Resolve(resolvedID, srcs, a.ctx)
	if err != nil {
		return fmt.Errorf("source resolution failed: %w", err)
	}
	if len(recs) == 0 {
		conflicts := resolver.DetectConflicts(srcs, a.ctx)
		if len(conflicts) > 0 {
			fmt.Fprintf(a.out, "Cannot install %q: no compatible source found.\n", resolvedID)
			for _, c := range conflicts {
				fmt.Fprintf(a.out, "  - %s\n", c.Message)
				if c.Suggestion != "" {
					fmt.Fprintf(a.out, "    Next step: %s\n", c.Suggestion)
				}
			}
		}
		return fmt.Errorf("no compatible source found for %q on this system", resolvedID)
	}

	rec := recs[0]
	fmt.Fprintf(a.out, "Installing %s...\n", app.DisplayName)
	fmt.Fprintf(a.out, "  Source: %s\n", rec.Plan.SourceType)
	fmt.Fprintf(a.out, "  Reason: %s\n", rec.Explanation)
	if rec.Plan.RequiresConfirmation {
		confirmed, confirmErr := a.confirmHighRiskInstall(app.DisplayName, rec.Plan.RiskLevel)
		if confirmErr != nil {
			return confirmErr
		}
		if !confirmed {
			return fmt.Errorf("installation cancelled by user")
		}
	}

	result, err := a.engine.Install(rec.Plan)
	if err != nil {
		return fmt.Errorf("install engine error: %w", err)
	}
	if !result.Success {
		return fmt.Errorf("installation failed: %s", result.Message)
	}

	// Record local state.
	_ = a.store.Record(state.LocalInstallation{
		ApplicationID:      resolvedID,
		SourceType:         rec.Plan.SourceType,
		SourceIdentifier:   rec.Plan.SourceIdentifier,
		InstallTimestamp:   time.Now(),
		InstallStatus:      state.StatusInstalled,
		VerificationStatus: state.VerificationPassed,
	})

	fmt.Fprintf(a.out, "\n%s\n", result.Message)
	return nil
}

// installDryRun handles `omni install <app> --dry-run`.
//
// It resolves the best source for the application and prints what would be
// installed (adapter, package identifier, risk level) without executing the
// installation.
func (a *App) installDryRun(appID string) error {
	if appID == "" {
		return errors.New("usage: omni install <application> --dry-run")
	}

	candidates, err := a.discovery.Search(appID)
	if err != nil {
		return fmt.Errorf("discovery failed: %w", err)
	}
	if len(candidates) == 0 {
		return fmt.Errorf("application %q not found", appID)
	}
	app := candidates[0].Application
	resolvedID := app.ID

	srcs, ok := a.sources[resolvedID]
	if !ok || len(srcs) == 0 {
		return fmt.Errorf("no sources available for %q", resolvedID)
	}

	recs, err := a.resolver.Resolve(resolvedID, srcs, a.ctx)
	if err != nil {
		return fmt.Errorf("source resolution failed: %w", err)
	}
	if len(recs) == 0 {
		return fmt.Errorf("no compatible source found for %q on this system", resolvedID)
	}

	rec := recs[0]
	fmt.Fprintf(a.out, "Dry-run: would install %s\n", app.DisplayName)
	fmt.Fprintf(a.out, "  Adapter:    %s\n", rec.Plan.SourceType)
	fmt.Fprintf(a.out, "  Package ID: %s\n", rec.Plan.SourceIdentifier)
	fmt.Fprintf(a.out, "  Risk level: %s\n", rec.Plan.RiskLevel)
	if rec.Plan.RequiresPrivilege {
		fmt.Fprintf(a.out, "  Requires elevated privileges: yes\n")
	}
	fmt.Fprintf(a.out, "  Reason: %s\n", rec.Explanation)
	fmt.Fprintf(a.out, "\nNo changes were made (dry-run).\n")
	return nil
}

// installDryRunJSON handles `omni install <app> --dry-run --json`.
func (a *App) installDryRunJSON(appID string) error {
	if appID == "" {
		return errors.New("usage: omni install <application> --dry-run --json")
	}

	candidates, err := a.discovery.Search(appID)
	if err != nil {
		return fmt.Errorf("discovery failed: %w", err)
	}
	if len(candidates) == 0 {
		return fmt.Errorf("application %q not found", appID)
	}
	app := candidates[0].Application
	resolvedID := app.ID

	srcs, ok := a.sources[resolvedID]
	if !ok || len(srcs) == 0 {
		return fmt.Errorf("no sources available for %q", resolvedID)
	}

	recs, err := a.resolver.Resolve(resolvedID, srcs, a.ctx)
	if err != nil {
		return fmt.Errorf("source resolution failed: %w", err)
	}
	if len(recs) == 0 {
		return fmt.Errorf("no compatible source found for %q on this system", resolvedID)
	}

	rec := recs[0]
	out := installDryRunJSONOutput{
		ApplicationID:     resolvedID,
		DisplayName:       app.DisplayName,
		Adapter:           rec.Plan.SourceType,
		PackageID:         rec.Plan.SourceIdentifier,
		RiskLevel:         rec.Plan.RiskLevel,
		RequiresPrivilege: rec.Plan.RequiresPrivilege,
		Explanation:       rec.Explanation,
	}
	return writeJSON(a.out, out)
}
//
// It re-runs post-install verification for an already-installed application
// without reinstalling it, and updates the verification status in local state.
func (a *App) verify(appID string) error {
	if appID == "" {
		return errors.New("usage: omni verify <application>")
	}

	rec, ok, err := a.store.Get(appID)
	if err != nil {
		return fmt.Errorf("could not read local state: %w", err)
	}
	if !ok {
		return fmt.Errorf("application %q is not recorded as installed", appID)
	}
	if rec.InstallStatus == state.StatusRemoved {
		return fmt.Errorf("application %q has been removed", appID)
	}

	result, err := a.engine.Verify(appID, rec.SourceType, rec.SourceIdentifier)
	if err != nil {
		return fmt.Errorf("verify engine error: %w", err)
	}

	// Update verification status in local state.
	newStatus := state.VerificationPassed
	if !result.Success {
		newStatus = state.VerificationFailed
	}
	rec.VerificationStatus = newStatus
	if storeErr := a.store.Record(rec); storeErr != nil {
		return fmt.Errorf("could not update verification status: %w", storeErr)
	}

	if result.Success {
		fmt.Fprintf(a.out, "✔ %s\n", result.Message)
	} else {
		fmt.Fprintf(a.out, "✘ %s\n", result.Message)
	}
	return nil
}

// remove handles `omni remove <app>`.
func (a *App) remove(appID string) error {
	if appID == "" {
		return errors.New("usage: omni remove <application>")
	}

	record, found, err := a.store.Get(appID)
	if err != nil {
		return fmt.Errorf("failed to read local state: %w", err)
	}
	if !found || record.InstallStatus == state.StatusRemoved {
		return fmt.Errorf("application %q is not recorded as installed", appID)
	}

	result, err := a.engine.Remove(appID, record.SourceType, record.SourceIdentifier)
	if err != nil {
		return fmt.Errorf("remove engine error: %w", err)
	}
	if !result.Success {
		return fmt.Errorf("removal failed: %s", result.Message)
	}

	// Update local state.
	if markErr := a.store.MarkRemoved(appID); markErr != nil && !errors.Is(markErr, state.ErrNotFound) {
		// Non-fatal: state update failure does not fail the removal.
		fmt.Fprintf(a.out, "Warning: could not update local state: %v\n", markErr)
	}

	fmt.Fprintf(a.out, "%s\n", result.Message)
	return nil
}

// explain handles `omni explain <app>`.
func (a *App) explain(appID string) error {
	if appID == "" {
		return errors.New("usage: omni explain <application>")
	}

	candidates, err := a.discovery.Search(appID)
	if err != nil {
		return fmt.Errorf("discovery failed: %w", err)
	}
	if len(candidates) == 0 {
		return fmt.Errorf("application %q not found", appID)
	}
	app := candidates[0].Application
	resolvedID := app.ID

	srcs, ok := a.sources[resolvedID]
	if !ok || len(srcs) == 0 {
		fmt.Fprintf(a.out, "No sources available for %q.\n", resolvedID)
		return nil
	}

	recs, err := a.resolver.Resolve(resolvedID, srcs, a.ctx)
	if err != nil {
		return fmt.Errorf("source resolution failed: %w", err)
	}

	fmt.Fprintf(a.out, "Source selection for %s (%s):\n\n", app.DisplayName, resolvedID)

	if len(recs) == 0 {
		fmt.Fprintf(a.out, "  No compatible source found on this system.\n")
		conflicts := resolver.DetectConflicts(srcs, a.ctx)
		if len(conflicts) > 0 {
			fmt.Fprintf(a.out, "\n  Reasons:\n")
			for _, c := range conflicts {
				fmt.Fprintf(a.out, "    - %s\n", c.Message)
				if c.Suggestion != "" {
					fmt.Fprintf(a.out, "      Next step: %s\n", c.Suggestion)
				}
			}
		}
		return nil
	}

	fmt.Fprintf(a.out, "  Recommended: %s (%s)\n", recs[0].Plan.SourceType, recs[0].Plan.SourceIdentifier)

	// Show enriched metadata for the top recommendation.
	if topSrc, ok := findSource(srcs, recs[0].Plan.SourceType, recs[0].Plan.SourceIdentifier); ok {
		fmt.Fprintf(a.out, "  Trust:       %s\n", topSrc.TrustLevel)
		fmt.Fprintf(a.out, "  Risk:        %s\n", topSrc.RiskLevel)
	}
	if recs[0].Plan.RequiresPrivilege {
		fmt.Fprintf(a.out, "  Privilege:   required (sudo)\n")
	}
	fmt.Fprintf(a.out, "\n  %s\n", recs[0].Explanation)

	if len(recs) > 1 {
		fmt.Fprintf(a.out, "\n  Alternatives:\n")
		for _, alt := range recs[1:] {
			fmt.Fprintf(a.out, "    - %s: %s\n", alt.Plan.SourceType, alt.Explanation)
		}
	}

	conflicts := resolver.DetectConflicts(srcs, a.ctx)
	if len(conflicts) > 0 {
		fmt.Fprintf(a.out, "\n  Conflicts:\n")
		for _, c := range conflicts {
			fmt.Fprintf(a.out, "    - %s\n", c.Message)
		}
	}

	return nil
}

func (a *App) explainJSON(appID string) error {
	if appID == "" {
		return errors.New("usage: omni explain <application>")
	}

	candidates, err := a.discovery.Search(appID)
	if err != nil {
		return fmt.Errorf("discovery failed: %w", err)
	}
	if len(candidates) == 0 {
		return fmt.Errorf("application %q not found", appID)
	}
	app := candidates[0].Application
	resolvedID := app.ID

	srcs, ok := a.sources[resolvedID]
	if !ok {
		srcs = nil
	}

	recs, err := a.resolver.Resolve(resolvedID, srcs, a.ctx)
	if err != nil {
		return fmt.Errorf("source resolution failed: %w", err)
	}

	out := explainJSONOutput{
		ApplicationID: resolvedID,
		DisplayName:   app.DisplayName,
		Alternatives:  []explainJSONAlternative{},
		Conflicts:     []explainJSONConflict{},
	}

	if len(recs) > 0 && recs[0].Plan != nil {
		out.RecommendedSourceType = recs[0].Plan.SourceType
		out.RecommendedExplanation = recs[0].Explanation
		out.RiskLevel = recs[0].Plan.RiskLevel
		out.RequiresPrivilege = recs[0].Plan.RequiresPrivilege
		if topSrc, ok := findSource(srcs, recs[0].Plan.SourceType, recs[0].Plan.SourceIdentifier); ok {
			out.TrustLevel = topSrc.TrustLevel
		}
	}
	if len(recs) > 1 {
		for _, alt := range recs[1:] {
			out.Alternatives = append(out.Alternatives, explainJSONAlternative{
				SourceType:  alt.Plan.SourceType,
				Explanation: alt.Explanation,
			})
		}
	}

	conflicts := resolver.DetectConflicts(srcs, a.ctx)
	for _, c := range conflicts {
		out.Conflicts = append(out.Conflicts, explainJSONConflict{
			Kind:       c.Kind,
			Message:    c.Message,
			Suggestion: c.Suggestion,
		})
	}

	return writeJSON(a.out, out)
}

// list handles `omni list`.
func (a *App) list() error {
	records, err := a.store.List()
	if err != nil {
		return fmt.Errorf("could not read local state: %w", err)
	}

	active := make([]state.LocalInstallation, 0, len(records))
	for _, r := range records {
		if r.InstallStatus != state.StatusRemoved {
			active = append(active, r)
		}
	}

	if len(active) == 0 {
		fmt.Fprintf(a.out, "No applications installed via OmniInstall.\n")
		return nil
	}

	// Sort by application ID for stable output.
	sort.Slice(active, func(i, j int) bool {
		return active[i].ApplicationID < active[j].ApplicationID
	})

	fmt.Fprintf(a.out, "Installed applications:\n\n")
	fmt.Fprintf(a.out, "  %-25s %-15s %s\n", "APPLICATION", "SOURCE", "STATUS")
	fmt.Fprintf(a.out, "  %s\n", strings.Repeat("-", 60))
	for _, r := range active {
		fmt.Fprintf(a.out, "  %-25s %-15s %s\n",
			r.ApplicationID,
			r.SourceType,
			r.InstallStatus,
		)
	}
	return nil
}

func (a *App) listJSON() error {
	records, err := a.store.List()
	if err != nil {
		return fmt.Errorf("could not read local state: %w", err)
	}

	active := make([]state.LocalInstallation, 0, len(records))
	for _, r := range records {
		if r.InstallStatus != state.StatusRemoved {
			active = append(active, r)
		}
	}

	sort.Slice(active, func(i, j int) bool {
		return active[i].ApplicationID < active[j].ApplicationID
	})

	out := listJSONOutput{Applications: make([]listJSONApplication, 0, len(active))}
	for _, r := range active {
		out.Applications = append(out.Applications, listJSONApplication{
			ApplicationID: r.ApplicationID,
			SourceType:    r.SourceType,
			Status:        r.InstallStatus,
		})
	}
	return writeJSON(a.out, out)
}

// stateExportFormat is the schema used for state export/import files.
type stateExportFormat struct {
	SchemaVersion string                   `json:"schema_version" yaml:"schema_version"`
	ExportedAt    string                   `json:"exported_at"    yaml:"exported_at"`
	Installations []state.LocalInstallation `json:"installations"  yaml:"installations"`
}

// stateExport handles `omniinstall state export [--yaml]`.
//
// Outputs all non-removed installations as JSON (default) or YAML.
func (a *App) stateExport(yamlOutput bool) error {
	records, err := a.store.List()
	if err != nil {
		return fmt.Errorf("could not read local state: %w", err)
	}

	active := make([]state.LocalInstallation, 0, len(records))
	for _, r := range records {
		if r.InstallStatus != state.StatusRemoved {
			active = append(active, r)
		}
	}
	sort.Slice(active, func(i, j int) bool {
		return active[i].ApplicationID < active[j].ApplicationID
	})

	export := stateExportFormat{
		SchemaVersion: "1.0",
		ExportedAt:    time.Now().UTC().Format(time.RFC3339),
		Installations: active,
	}

	if yamlOutput {
		enc := yaml.NewEncoder(a.out)
		enc.SetIndent(2)
		if err := enc.Encode(export); err != nil {
			return fmt.Errorf("yaml encode failed: %w", err)
		}
		return enc.Close()
	}
	return writeJSON(a.out, export)
}

// stateImport handles `omniinstall state import <file> [--force] [--dry-run]`.
//
// Reads a JSON/YAML export file and records any non-conflicting installations
// into local state. With --force, existing records are overwritten. With
// --dry-run, changes are shown but not applied.
func (a *App) stateImport(filePath string, force, dryRun bool) error {
	if filePath == "" {
		return errors.New("usage: omniinstall state import <file> [--force] [--dry-run]")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("could not read file %q: %w", filePath, err)
	}

	var export stateExportFormat
	// Try JSON first, fall back to YAML.
	if jsonErr := json.Unmarshal(data, &export); jsonErr != nil {
		if yamlErr := yaml.Unmarshal(data, &export); yamlErr != nil {
			return fmt.Errorf("file is not valid JSON or YAML: %v; %v", jsonErr, yamlErr)
		}
	}

	// Validate schema_version.
	if export.SchemaVersion == "" {
		return errors.New("import file is missing schema_version")
	}
	parts := strings.SplitN(export.SchemaVersion, ".", 2)
	if parts[0] != "1" {
		return fmt.Errorf("unsupported schema_version %q: only major version 1 is supported", export.SchemaVersion)
	}

	if len(export.Installations) == 0 {
		fmt.Fprintf(a.out, "No installations found in import file.\n")
		return nil
	}

	var toImport []state.LocalInstallation
	for _, rec := range export.Installations {
		if err := rec.Validate(); err != nil {
			fmt.Fprintf(a.out, "  Skipping %q: invalid record (%v)\n", rec.ApplicationID, err)
			continue
		}
		existing, found, lookupErr := a.store.Get(rec.ApplicationID)
		if lookupErr != nil {
			return fmt.Errorf("state lookup failed: %w", lookupErr)
		}
		if found && existing.InstallStatus != state.StatusRemoved && !force {
			fmt.Fprintf(a.out, "  Skipping %q: already installed (use --force to overwrite)\n", rec.ApplicationID)
			continue
		}
		toImport = append(toImport, rec)
	}

	if dryRun {
		fmt.Fprintf(a.out, "Dry-run: would import %d installation(s):\n", len(toImport))
		for _, rec := range toImport {
			fmt.Fprintf(a.out, "  %s (%s %s)\n", rec.ApplicationID, rec.SourceType, rec.SourceIdentifier)
		}
		fmt.Fprintf(a.out, "\nNo changes were made (dry-run).\n")
		return nil
	}

	imported := 0
	for _, rec := range toImport {
		if err := a.store.Record(rec); err != nil {
			fmt.Fprintf(a.out, "  Warning: could not import %q: %v\n", rec.ApplicationID, err)
			continue
		}
		fmt.Fprintf(a.out, "  Imported: %s\n", rec.ApplicationID)
		imported++
	}
	fmt.Fprintf(a.out, "\n%d installation(s) imported.\n", imported)
	return nil
}

// handleProgress prints progress events to the output writer.
func handleProgress(out io.Writer) engine.ProgressHandler {
	return func(ev engine.ProgressEvent) {
		switch ev.Kind {
		case engine.EventStarted, engine.EventCompleted, engine.EventFailed:
			fmt.Fprintf(out, "  [%s] %s\n", ev.Kind, ev.Message)
		}
	}
}

func (a *App) confirmHighRiskInstall(displayName string, riskLevel source.RiskLevel) (bool, error) {
	if !isInteractiveInput(a.in) {
		return false, fmt.Errorf("high-risk installation requires interactive confirmation")
	}

	fmt.Fprintf(a.out, "  Warning: This action is marked as %s risk.\n", riskLevel)
	fmt.Fprintf(a.out, "  Type 'yes' to continue installing %s: ", displayName)

	reader := bufio.NewReader(a.in)
	answer, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false, fmt.Errorf("failed to read confirmation input: %w", err)
	}
	return strings.EqualFold(strings.TrimSpace(answer), "yes"), nil
}

func isInteractiveInput(in io.Reader) bool {
	if in == nil {
		return false
	}
	file, ok := in.(*os.File)
	if !ok {
		return true
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}

// newAppWithStore creates an App with a custom store for testing.
func newAppWithStore(out io.Writer, store state.Store, managers []source.Type) *App {
	adapters := []adapter.Adapter{
		aptadapter.New(),
		flatpakadapter.New(),
	}
	return &App{
		discovery: discovery.NewLocalEngine(discovery.MVPCatalog()),
		resolver:  resolver.NewDefaultResolver(),
		engine:    engine.NewDefaultEngine(adapters, nil),
		store:     store,
		sources:   discovery.MVPSources(),
		ctx: resolver.SystemContext{
			AvailableManagers: managers,
		},
		in:  os.Stdin,
		out: out,
	}
}

// newAppForTest creates an App with an in-memory store and a custom engine for testing.
func newAppForTest(out io.Writer, eng *engine.DefaultEngine, store state.Store, managers []source.Type) *App {
	return &App{
		discovery: discovery.NewLocalEngine(discovery.MVPCatalog()),
		resolver:  resolver.NewDefaultResolver(),
		engine:    eng,
		store:     store,
		sources:   discovery.MVPSources(),
		ctx: resolver.SystemContext{
			AvailableManagers: managers,
		},
		in:  os.Stdin,
		out: out,
	}
}

// Ensure App uses os.Stdout for the production path.
var _ = os.Stdout

// findSource returns the first source from srcs that matches the given type
// and identifier, along with a boolean indicating whether it was found.
func findSource(srcs []source.Source, t source.Type, identifier string) (source.Source, bool) {
	for _, s := range srcs {
		if s.SourceType == t && s.SourceIdentifier == identifier {
			return s, true
		}
	}
	return source.Source{}, false
}

func writeJSON(out io.Writer, v any) error {
	enc := json.NewEncoder(out)
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}
