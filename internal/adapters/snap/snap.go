// Package snap provides the Snap adapter for OmniInstall.
//
// The Snap adapter handles Canonical Snap package manager operations.
// It implements the adapter.Adapter interface.
package snap

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	adapter "github.com/mitslab-de/omniinstall/internal/adapters"
	"github.com/mitslab-de/omniinstall/internal/install"
	"github.com/mitslab-de/omniinstall/internal/source"
)

// AdapterName is the human-readable name for the Snap adapter.
const AdapterName = "snap"

// CanHandle reports whether the Snap adapter supports the given source type.
func CanHandle(sourceType source.Type) bool {
	return sourceType == source.TypeSnap
}

// executor abstracts command execution for testability.
type executor interface {
	LookPath(file string) (string, error)
	Run(name string, args ...string) (string, int, error)
}

// defaultExecutor executes real OS commands.
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

// Adapter is the Snap implementation of the adapter.Adapter interface.
type Adapter struct {
	exec executor
}

// New returns a production-ready Snap Adapter.
func New() *Adapter {
	return &Adapter{exec: &defaultExecutor{}}
}

// NewWithExecutorForTest returns a Snap Adapter with a custom executor for testing.
func NewWithExecutorForTest(e interface {
	LookPath(string) (string, error)
	Run(string, ...string) (string, int, error)
}) *Adapter {
	return &Adapter{exec: e}
}

// Name returns the adapter's human-readable name.
func (a *Adapter) Name() string { return AdapterName }

// IsAvailable returns true if snap is present on the current system.
func (a *Adapter) IsAvailable() bool {
	_, err := a.exec.LookPath("snap")
	return err == nil
}

// CanHandle returns true for the Snap source type.
func (a *Adapter) CanHandle(sourceType source.Type) bool {
	return CanHandle(sourceType)
}

// CheckInstalled returns the installation state for the given snap name.
//
// snap list <name> exits 0 and prints a table row if the snap is installed,
// or exits non-zero with an error message if it is not present.
func (a *Adapter) CheckInstalled(sourceIdentifier string) (adapter.InstalledState, error) {
	out, exitCode, err := a.exec.Run("snap", "list", sourceIdentifier)
	if err != nil && exitCode == 0 {
		return adapter.StateUnknown, fmt.Errorf("snap list failed: %w", err)
	}
	if exitCode != 0 {
		lower := strings.ToLower(out)
		if strings.Contains(lower, "no snaps are installed") ||
			strings.Contains(lower, "is not installed") ||
			strings.Contains(lower, "not found") {
			return adapter.StateNotInstalled, nil
		}
		return adapter.StateUnknown, fmt.Errorf("snap list returned unexpected output (exit %d): %q", exitCode, out)
	}
	// A zero-exit with output containing the snap name means it's installed.
	// An empty output (edge case) is treated as not-installed.
	trimmed := strings.TrimSpace(out)
	if trimmed == "" {
		return adapter.StateNotInstalled, nil
	}
	return adapter.StateInstalled, nil
}

// Install installs the snap described in the plan using snap install.
func (a *Adapter) Install(plan *install.Plan) (*adapter.Result, error) {
	start := time.Now()
	out, exitCode, err := a.exec.Run("snap", "install", plan.SourceIdentifier)
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
		cat := categorizeExitCode(exitCode, out)
		return &adapter.Result{
			Success:          false,
			ExitCode:         exitCode,
			Message:          fmt.Sprintf("snap install failed for %s (exit %d)", plan.SourceIdentifier, exitCode),
			TechnicalDetails: out,
			ErrorCategory:    cat,
			Duration:         duration,
		}, nil
	}
	return &adapter.Result{
		Success:          true,
		ExitCode:         0,
		Message:          fmt.Sprintf("Installed %s via Snap", plan.SourceIdentifier),
		TechnicalDetails: out,
		ChangedSystem:    true,
		Duration:         duration,
	}, nil
}

// Remove removes the snap identified by sourceIdentifier using snap remove.
func (a *Adapter) Remove(sourceIdentifier string) (*adapter.Result, error) {
	start := time.Now()
	out, exitCode, err := a.exec.Run("snap", "remove", sourceIdentifier)
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
		cat := categorizeExitCode(exitCode, out)
		return &adapter.Result{
			Success:          false,
			ExitCode:         exitCode,
			Message:          fmt.Sprintf("snap remove failed for %s (exit %d)", sourceIdentifier, exitCode),
			TechnicalDetails: out,
			ErrorCategory:    cat,
			Duration:         duration,
		}, nil
	}
	return &adapter.Result{
		Success:          true,
		ExitCode:         0,
		Message:          fmt.Sprintf("Removed %s via Snap", sourceIdentifier),
		TechnicalDetails: out,
		ChangedSystem:    true,
		Duration:         duration,
	}, nil
}

// Verify checks that the commands listed in the plan's verification rules
// exist in PATH after installation.
func (a *Adapter) Verify(plan *install.Plan) (*adapter.VerificationResult, error) {
	for _, rule := range plan.Verification {
		if rule.Command == "" {
			continue
		}
		_, err := a.exec.LookPath(rule.Command)
		if err != nil {
			return &adapter.VerificationResult{
				Verified: false,
				Details:  fmt.Sprintf("command %q not found in PATH after installation", rule.Command),
			}, nil
		}
	}
	return &adapter.VerificationResult{
		Verified: true,
		Details:  "all verification rules passed",
	}, nil
}

// categorizeExitCode maps a snap exit code and output to an error category.
func categorizeExitCode(exitCode int, output string) string {
	lower := strings.ToLower(output)
	if strings.Contains(lower, "snap not found") || strings.Contains(lower, "not found in store") {
		return adapter.ErrPackageNotFound
	}
	if strings.Contains(lower, "access denied") || strings.Contains(lower, "permission denied") {
		return adapter.ErrPermissionDenied
	}
	if strings.Contains(lower, "unable to contact snap store") || strings.Contains(lower, "post https") {
		return adapter.ErrNetworkError
	}
	if strings.Contains(lower, "conflict") {
		return adapter.ErrDependencyConflict
	}
	_ = exitCode
	return adapter.ErrExecutionFailed
}
