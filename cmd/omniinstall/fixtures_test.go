package main

// CLI test fixtures — reusable helpers for command routing, output capture,
// and state setup in cmd/omniinstall tests.
//
// Usage:
//
//	f := newFixture(t).withAPT()
//	f.mustRun(t, "search", "git")
//	if !strings.Contains(f.Output(), "Git") { t.Error(...) }

import (
	"strings"
	"testing"
	"time"

	adapter "github.com/mitslab-de/omniinstall/internal/adapters"
	"github.com/mitslab-de/omniinstall/internal/engine"
	"github.com/mitslab-de/omniinstall/internal/install"
	"github.com/mitslab-de/omniinstall/internal/source"
	"github.com/mitslab-de/omniinstall/internal/state"
)

// CLIFixture is a reusable test harness for cmd/omniinstall App methods.
//
// It bundles an output buffer, an in-memory state store, and a
// configurable engine so tests can focus on behaviour rather than setup.
type CLIFixture struct {
	buf     strings.Builder
	store   *mockStore
	engine  *engine.DefaultEngine
	managers []source.Type
}

// newFixture creates a CLIFixture with APT adapter and APT manager by default.
func newFixture(_ *testing.T) *CLIFixture {
	f := &CLIFixture{
		store:    newMockStore(),
		managers: aptManagers(),
	}
	f.engine = engine.NewDefaultEngine([]adapter.Adapter{&fixtureAdapter{
		name:         "apt",
		available:    true,
		handledTypes: []source.Type{source.TypeAPT},
	}}, nil)
	return f
}

// withAPT configures the fixture with a successful APT adapter.
func (f *CLIFixture) withAPT() *CLIFixture {
	f.managers = aptManagers()
	f.engine = engine.NewDefaultEngine([]adapter.Adapter{&fixtureAdapter{
		name:         "apt",
		available:    true,
		handledTypes: []source.Type{source.TypeAPT},
	}}, nil)
	return f
}

// withFlatpak configures the fixture with a successful Flatpak adapter.
func (f *CLIFixture) withFlatpak() *CLIFixture {
	f.managers = flatpakManagers()
	f.engine = engine.NewDefaultEngine([]adapter.Adapter{&fixtureAdapter{
		name:         "flatpak",
		available:    true,
		handledTypes: []source.Type{source.TypeFlatpak},
	}}, nil)
	return f
}

// withInstalledAPT adds an installed APT record to the fixture state store.
func (f *CLIFixture) withInstalledAPT(appID, pkg string) *CLIFixture {
	_ = f.store.Record(state.LocalInstallation{
		ApplicationID:      appID,
		SourceType:         source.TypeAPT,
		SourceIdentifier:   pkg,
		InstallTimestamp:   time.Now(),
		InstallStatus:      state.StatusInstalled,
		VerificationStatus: state.VerificationPassed,
	})
	return f
}

// withInstalledFlatpak adds an installed Flatpak record to the fixture state store.
func (f *CLIFixture) withInstalledFlatpak(appID, flatpakID string) *CLIFixture {
	_ = f.store.Record(state.LocalInstallation{
		ApplicationID:      appID,
		SourceType:         source.TypeFlatpak,
		SourceIdentifier:   flatpakID,
		InstallTimestamp:   time.Now(),
		InstallStatus:      state.StatusInstalled,
		VerificationStatus: state.VerificationPassed,
	})
	return f
}

// withRemovedAPT adds a removed APT record to the fixture state store.
func (f *CLIFixture) withRemovedAPT(appID, pkg string) *CLIFixture {
	_ = f.store.Record(state.LocalInstallation{
		ApplicationID:      appID,
		SourceType:         source.TypeAPT,
		SourceIdentifier:   pkg,
		InstallTimestamp:   time.Now(),
		InstallStatus:      state.StatusRemoved,
		VerificationStatus: state.VerificationSkipped,
	})
	return f
}

// withVerifyResult sets the adapter's verify outcome (true = pass, false = fail).
func (f *CLIFixture) withVerifyResult(verified bool, details string) *CLIFixture {
	f.engine = engine.NewDefaultEngine([]adapter.Adapter{&fixtureAdapter{
		name:          "apt",
		available:     true,
		handledTypes:  []source.Type{source.TypeAPT},
		verifyResult:  &adapter.VerificationResult{Verified: verified, Details: details},
	}}, nil)
	return f
}

// app returns a configured App for the fixture.
func (f *CLIFixture) app() *App {
	return newAppForTest(&f.buf, f.engine, f.store, f.managers)
}

// Output returns captured output from the last command run.
func (f *CLIFixture) Output() string {
	return f.buf.String()
}

// resetOutput clears the output buffer.
func (f *CLIFixture) resetOutput() {
	f.buf.Reset()
}

// mustRun dispatches to the appropriate App method based on the command name
// and fails the test if an error is returned.
//
// Supported commands: search, install, remove, verify, list, explain.
func (f *CLIFixture) mustRun(t *testing.T, cmd string, args ...string) string {
	t.Helper()
	err := f.run(cmd, args...)
	if err != nil {
		t.Fatalf("mustRun(%q %v): unexpected error: %v", cmd, args, err)
	}
	return f.buf.String()
}

// runExpectError dispatches to the appropriate App method and fails the test
// if NO error is returned.
func (f *CLIFixture) runExpectError(t *testing.T, cmd string, args ...string) error {
	t.Helper()
	err := f.run(cmd, args...)
	if err == nil {
		t.Fatalf("runExpectError(%q %v): expected an error, got nil", cmd, args)
	}
	return err
}

// run dispatches to the App method for the given command.
func (f *CLIFixture) run(cmd string, args ...string) error {
	a := f.app()
	arg0 := ""
	if len(args) > 0 {
		arg0 = args[0]
	}
	switch cmd {
	case "search":
		return a.search(arg0)
	case "install":
		return a.install(arg0)
	case "install-dry-run":
		return a.installDryRun(arg0)
	case "install-dry-run-json":
		return a.installDryRunJSON(arg0)
	case "state-export":
		return a.stateExport(false)
	case "state-export-yaml":
		return a.stateExport(true)
	case "remove":
		return a.remove(arg0)
	case "verify":
		return a.verify(arg0)
	case "list":
		return a.list()
	case "explain":
		return a.explain(arg0)
	default:
		return nil
	}
}

// fixtureAdapter is a configurable adapter for CLI fixtures.
type fixtureAdapter struct {
	name         string
	available    bool
	handledTypes []source.Type
	verifyResult *adapter.VerificationResult
}

func (a *fixtureAdapter) Name() string      { return a.name }
func (a *fixtureAdapter) IsAvailable() bool { return a.available }
func (a *fixtureAdapter) CanHandle(t source.Type) bool {
	for _, ht := range a.handledTypes {
		if ht == t {
			return true
		}
	}
	return false
}
func (a *fixtureAdapter) CheckInstalled(_ string) (adapter.InstalledState, error) {
	return adapter.StateInstalled, nil
}
func (a *fixtureAdapter) Install(_ *install.Plan) (*adapter.Result, error) {
	return &adapter.Result{Success: true, Message: "installed", ChangedSystem: true}, nil
}
func (a *fixtureAdapter) Remove(_ string) (*adapter.Result, error) {
	return &adapter.Result{Success: true, Message: "removed", ChangedSystem: true}, nil
}
func (a *fixtureAdapter) Verify(_ *install.Plan) (*adapter.VerificationResult, error) {
	if a.verifyResult != nil {
		return a.verifyResult, nil
	}
	return &adapter.VerificationResult{Verified: true, Details: "ok"}, nil
}
