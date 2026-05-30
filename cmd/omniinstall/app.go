package main

import (
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
	flatpakadapter "github.com/mitslab-de/omniinstall/internal/adapters/flatpak"
	"github.com/mitslab-de/omniinstall/internal/discovery"
	"github.com/mitslab-de/omniinstall/internal/engine"
	"github.com/mitslab-de/omniinstall/internal/resolver"
	"github.com/mitslab-de/omniinstall/internal/source"
	"github.com/mitslab-de/omniinstall/internal/state"
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
	out       io.Writer
}

// newApp creates a production App, detecting available package managers and
// wiring all services together.
func newApp(out io.Writer) *App {
	managers := detectAvailableManagers()

	adapters := []adapter.Adapter{
		aptadapter.New(),
		flatpakadapter.New(),
	}

	storePath := state.DefaultStorePath()

	return &App{
		discovery: discovery.NewLocalEngine(discovery.MVPCatalog()),
		resolver:  resolver.NewDefaultResolver(),
		engine:    engine.NewDefaultEngine(adapters, nil),
		store:     state.NewFileStore(storePath),
		sources:   discovery.MVPSources(),
		ctx: resolver.SystemContext{
			AvailableManagers: managers,
		},
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
		return fmt.Errorf("no compatible source found for %q on this system", resolvedID)
	}

	rec := recs[0]
	fmt.Fprintf(a.out, "Installing %s...\n", app.DisplayName)
	fmt.Fprintf(a.out, "  Source: %s\n", rec.Plan.SourceType)
	fmt.Fprintf(a.out, "  Reason: %s\n", rec.Explanation)

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
		return nil
	}

	fmt.Fprintf(a.out, "  Recommended: %s\n", recs[0].Plan.SourceType)
	fmt.Fprintf(a.out, "  %s\n", recs[0].Explanation)

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

// handleProgress prints progress events to the output writer.
func handleProgress(out io.Writer) engine.ProgressHandler {
	return func(ev engine.ProgressEvent) {
		switch ev.Kind {
		case engine.EventStarted, engine.EventCompleted, engine.EventFailed:
			fmt.Fprintf(out, "  [%s] %s\n", ev.Kind, ev.Message)
		}
	}
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
		out: out,
	}
}

// Ensure App uses os.Stdout for the production path.
var _ = os.Stdout
