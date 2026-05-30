// Package flatpak provides the Flatpak adapter for OmniInstall.
//
// The Flatpak adapter handles Flatpak application operations.
// It implements the adapter.Adapter interface.
package flatpak

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	adapter "github.com/mitslab-de/omniinstall/internal/adapters"
	"github.com/mitslab-de/omniinstall/internal/install"
	"github.com/mitslab-de/omniinstall/internal/source"
)

// AdapterName is the human-readable name for the Flatpak adapter.
const AdapterName = "flatpak"

// Name is an alias kept for backwards compatibility.
const Name = AdapterName

// CanHandle reports whether the Flatpak adapter supports the given source type.
func CanHandle(sourceType source.Type) bool {
	return sourceType == source.TypeFlatpak
}

// executor abstracts command execution for testability.
type executor interface {
	LookPath(file string) (string, error)
	Run(name string, args ...string) (string, int, error)
}

// defaultExecutor runs real OS commands.
type defaultExecutor struct{}

func (d *defaultExecutor) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

func (d *defaultExecutor) Run(name string, args ...string) (string, int, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	exitCode := 0
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}
	return string(out), exitCode, err
}

// Adapter is the Flatpak implementation of the adapter.Adapter interface.
type Adapter struct {
	exec executor
}

// New returns a production-ready Flatpak Adapter.
func New() *Adapter {
	return &Adapter{exec: &defaultExecutor{}}
}

// NewWithExecutorForTest exposes a custom executor for package-level tests.
func NewWithExecutorForTest(e interface {
	LookPath(string) (string, error)
	Run(string, ...string) (string, int, error)
}) *Adapter {
	return &Adapter{exec: e}
}

// Name returns the adapter's human-readable name.
func (a *Adapter) Name() string { return AdapterName }

// IsAvailable returns true if the flatpak binary is present on the system.
func (a *Adapter) IsAvailable() bool {
	_, err := a.exec.LookPath("flatpak")
	return err == nil
}

// CanHandle returns true for the Flatpak source type.
func (a *Adapter) CanHandle(sourceType source.Type) bool {
	return CanHandle(sourceType)
}

// CheckInstalled returns the installation state for the given Flatpak application ID.
func (a *Adapter) CheckInstalled(sourceIdentifier string) (adapter.InstalledState, error) {
	out, exitCode, _ := a.exec.Run("flatpak", "info", sourceIdentifier)
	if exitCode != 0 || strings.TrimSpace(out) == "" {
		return adapter.StateNotInstalled, nil
	}
	return adapter.StateInstalled, nil
}

// Install installs the Flatpak application described in the plan.
func (a *Adapter) Install(plan *install.Plan) (*adapter.Result, error) {
	start := time.Now()
	out, exitCode, err := a.exec.Run("flatpak", "install", "--noninteractive", "-y", plan.SourceIdentifier)
	duration := time.Since(start)
	if err != nil && exitCode == 0 {
		return &adapter.Result{
			Success:          false,
			ExitCode:         exitCode,
			Message:          fmt.Sprintf("Failed to install %s: %s", plan.SourceIdentifier, err.Error()),
			TechnicalDetails: out,
			ErrorCategory:    adapter.ErrExecutionFailed,
			Duration:         duration,
		}, nil
	}
	if exitCode != 0 {
		cat := categorizeFlatpakError(out)
		return &adapter.Result{
			Success:          false,
			ExitCode:         exitCode,
			Message:          fmt.Sprintf("flatpak install failed for %s (exit %d)", plan.SourceIdentifier, exitCode),
			TechnicalDetails: out,
			ErrorCategory:    cat,
			Duration:         duration,
		}, nil
	}
	return &adapter.Result{
		Success:          true,
		ExitCode:         0,
		Message:          fmt.Sprintf("Installed %s via Flatpak", plan.SourceIdentifier),
		TechnicalDetails: out,
		ChangedSystem:    true,
		Duration:         duration,
	}, nil
}

// Remove removes the Flatpak application identified by sourceIdentifier.
func (a *Adapter) Remove(sourceIdentifier string) (*adapter.Result, error) {
	start := time.Now()
	out, exitCode, err := a.exec.Run("flatpak", "remove", "--noninteractive", "-y", sourceIdentifier)
	duration := time.Since(start)
	if err != nil && exitCode == 0 {
		return &adapter.Result{
			Success:          false,
			ExitCode:         exitCode,
			Message:          fmt.Sprintf("Failed to remove %s: %s", sourceIdentifier, err.Error()),
			TechnicalDetails: out,
			ErrorCategory:    adapter.ErrExecutionFailed,
			Duration:         duration,
		}, nil
	}
	if exitCode != 0 {
		cat := categorizeFlatpakError(out)
		return &adapter.Result{
			Success:          false,
			ExitCode:         exitCode,
			Message:          fmt.Sprintf("flatpak remove failed for %s (exit %d)", sourceIdentifier, exitCode),
			TechnicalDetails: out,
			ErrorCategory:    cat,
			Duration:         duration,
		}, nil
	}
	return &adapter.Result{
		Success:          true,
		ExitCode:         0,
		Message:          fmt.Sprintf("Removed %s via Flatpak", sourceIdentifier),
		TechnicalDetails: out,
		ChangedSystem:    true,
		Duration:         duration,
	}, nil
}

// Verify checks that the Flatpak application is listed as installed after
// installation.
func (a *Adapter) Verify(plan *install.Plan) (*adapter.VerificationResult, error) {
	state, err := a.CheckInstalled(plan.SourceIdentifier)
	if err != nil {
		return nil, err
	}
	if state != adapter.StateInstalled {
		return &adapter.VerificationResult{
			Verified: false,
			Details:  fmt.Sprintf("flatpak application %q not found after installation", plan.SourceIdentifier),
		}, nil
	}
	return &adapter.VerificationResult{
		Verified: true,
		Details:  fmt.Sprintf("flatpak application %q is installed", plan.SourceIdentifier),
	}, nil
}

// categorizeFlatpakError maps Flatpak command output to an error category.
func categorizeFlatpakError(output string) string {
	lower := strings.ToLower(output)
	if strings.Contains(lower, "not found") || strings.Contains(lower, "no ref found") {
		return adapter.ErrPackageNotFound
	}
	if strings.Contains(lower, "permission denied") {
		return adapter.ErrPermissionDenied
	}
	if strings.Contains(lower, "network") || strings.Contains(lower, "connection") {
		return adapter.ErrNetworkError
	}
	return adapter.ErrExecutionFailed
}

