// Package dnf provides the DNF adapter for OmniInstall.
//
// The DNF adapter handles Fedora/RHEL/CentOS DNF package manager operations.
// It implements the adapter.Adapter interface.
package dnf

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	adapter "github.com/mitslab-de/omniinstall/internal/adapters"
	"github.com/mitslab-de/omniinstall/internal/install"
	"github.com/mitslab-de/omniinstall/internal/source"
)

// AdapterName is the human-readable name for the DNF adapter.
const AdapterName = "dnf"

// CanHandle reports whether the DNF adapter supports the given source type.
func CanHandle(sourceType source.Type) bool {
	return sourceType == source.TypeDNF
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

// Adapter is the DNF implementation of the adapter.Adapter interface.
type Adapter struct {
	exec executor
}

// New returns a production-ready DNF Adapter.
func New() *Adapter {
	return &Adapter{exec: &defaultExecutor{}}
}

// NewWithExecutorForTest returns a DNF Adapter with a custom executor for testing.
func NewWithExecutorForTest(e interface {
	LookPath(string) (string, error)
	Run(string, ...string) (string, int, error)
}) *Adapter {
	return &Adapter{exec: e}
}

// Name returns the adapter's human-readable name.
func (a *Adapter) Name() string { return AdapterName }

// IsAvailable returns true if dnf is present on the current system.
func (a *Adapter) IsAvailable() bool {
	_, err := a.exec.LookPath("dnf")
	return err == nil
}

// CanHandle returns true for the DNF source type.
func (a *Adapter) CanHandle(sourceType source.Type) bool {
	return CanHandle(sourceType)
}

// CheckInstalled returns the installation state for the given package name.
//
// rpm -q <package> returns exit 0 and the package NEVRA if installed, or
// exit 1 with "package <name> is not installed" if not present.
func (a *Adapter) CheckInstalled(sourceIdentifier string) (adapter.InstalledState, error) {
	out, exitCode, err := a.exec.Run("rpm", "-q", sourceIdentifier)
	if err != nil && exitCode == 0 {
		return adapter.StateUnknown, fmt.Errorf("rpm -q failed: %w", err)
	}
	if exitCode != 0 {
		lower := strings.ToLower(out)
		if strings.Contains(lower, "not installed") {
			return adapter.StateNotInstalled, nil
		}
		return adapter.StateUnknown, fmt.Errorf("rpm -q returned unexpected output (exit %d): %q", exitCode, out)
	}
	trimmed := strings.TrimSpace(out)
	if trimmed == "" {
		return adapter.StateNotInstalled, nil
	}
	return adapter.StateInstalled, nil
}

// Install installs the package described in the plan using dnf.
func (a *Adapter) Install(plan *install.Plan) (*adapter.Result, error) {
	start := time.Now()
	out, exitCode, err := a.exec.Run("dnf", "install", "-y", plan.SourceIdentifier)
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
			Message:          fmt.Sprintf("dnf install failed for %s (exit %d)", plan.SourceIdentifier, exitCode),
			TechnicalDetails: out,
			ErrorCategory:    cat,
			Duration:         duration,
		}, nil
	}
	return &adapter.Result{
		Success:          true,
		ExitCode:         0,
		Message:          fmt.Sprintf("Installed %s via DNF", plan.SourceIdentifier),
		TechnicalDetails: out,
		ChangedSystem:    true,
		Duration:         duration,
	}, nil
}

// Remove removes the package identified by sourceIdentifier using dnf.
func (a *Adapter) Remove(sourceIdentifier string) (*adapter.Result, error) {
	start := time.Now()
	out, exitCode, err := a.exec.Run("dnf", "remove", "-y", sourceIdentifier)
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
			Message:          fmt.Sprintf("dnf remove failed for %s (exit %d)", sourceIdentifier, exitCode),
			TechnicalDetails: out,
			ErrorCategory:    cat,
			Duration:         duration,
		}, nil
	}
	return &adapter.Result{
		Success:          true,
		ExitCode:         0,
		Message:          fmt.Sprintf("Removed %s via DNF", sourceIdentifier),
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

// categorizeExitCode maps a dnf exit code and output to an error category.
func categorizeExitCode(exitCode int, output string) string {
	lower := strings.ToLower(output)
	if strings.Contains(lower, "no package") || strings.Contains(lower, "no match") {
		return adapter.ErrPackageNotFound
	}
	if strings.Contains(lower, "permission denied") || strings.Contains(lower, "not authorized") {
		return adapter.ErrPermissionDenied
	}
	if strings.Contains(lower, "cannot download") || strings.Contains(lower, "curl error") {
		return adapter.ErrNetworkError
	}
	if strings.Contains(lower, "conflict") || strings.Contains(lower, "dependency") {
		return adapter.ErrDependencyConflict
	}
	_ = exitCode
	return adapter.ErrExecutionFailed
}
