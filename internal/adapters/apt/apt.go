// Package apt provides the APT adapter for OmniInstall.
//
// The APT adapter handles Debian/Ubuntu APT package manager operations.
// It implements the adapter.Adapter interface.
package apt

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	adapter "github.com/mitslab-de/omniinstall/internal/adapters"
	"github.com/mitslab-de/omniinstall/internal/install"
	"github.com/mitslab-de/omniinstall/internal/source"
)

// AdapterName is the human-readable name for the APT adapter.
const AdapterName = "apt"

// Name is an alias kept for backwards compatibility.
const Name = AdapterName

// CanHandle reports whether the APT adapter supports the given source type.
func CanHandle(sourceType source.Type) bool {
	return sourceType == source.TypeAPT
}

// executor abstracts command execution for testability.
type executor interface {
	// LookPath reports whether a binary exists in PATH.
	LookPath(file string) (string, error)
	// Run runs the named command with the given arguments and returns stdout,
	// the exit code, and any execution error.
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

// Adapter is the APT implementation of the adapter.Adapter interface.
type Adapter struct {
	exec executor
}

// New returns a production-ready APT Adapter.
func New() *Adapter {
	return &Adapter{exec: &defaultExecutor{}}
}

// newWithExecutor returns an APT Adapter with a custom executor for testing.
func newWithExecutor(e executor) *Adapter {
	return &Adapter{exec: e}
}

// NewWithExecutorForTest exposes newWithExecutor for package-level tests.
// It accepts an interface value and wraps it.
func NewWithExecutorForTest(e interface {
	LookPath(string) (string, error)
	Run(string, ...string) (string, int, error)
}) *Adapter {
	return &Adapter{exec: e}
}

// Name returns the adapter's human-readable name.
func (a *Adapter) Name() string { return AdapterName }

// IsAvailable returns true if apt-get is present on the current system.
func (a *Adapter) IsAvailable() bool {
	_, err := a.exec.LookPath("apt-get")
	return err == nil
}

// CanHandle returns true for the APT source type.
func (a *Adapter) CanHandle(sourceType source.Type) bool {
	return CanHandle(sourceType)
}

// CheckInstalled returns the installation state for the given package name.
func (a *Adapter) CheckInstalled(sourceIdentifier string) (adapter.InstalledState, error) {
	out, exitCode, _ := a.exec.Run("dpkg-query", "-W", "-f", "${Status}", sourceIdentifier)
	if exitCode != 0 {
		return adapter.StateNotInstalled, nil
	}
	if strings.Contains(out, "install ok installed") {
		return adapter.StateInstalled, nil
	}
	return adapter.StateNotInstalled, nil
}

// Install installs the package described in the plan using apt-get.
func (a *Adapter) Install(plan *install.Plan) (*adapter.Result, error) {
	start := time.Now()
	out, exitCode, err := a.exec.Run("apt-get", "install", "-y", plan.SourceIdentifier)
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
			Message:          fmt.Sprintf("apt-get install failed for %s (exit %d)", plan.SourceIdentifier, exitCode),
			TechnicalDetails: out,
			ErrorCategory:    cat,
			Duration:         duration,
		}, nil
	}
	return &adapter.Result{
		Success:          true,
		ExitCode:         0,
		Message:          fmt.Sprintf("Installed %s via APT", plan.SourceIdentifier),
		TechnicalDetails: out,
		ChangedSystem:    true,
		Duration:         duration,
	}, nil
}

// Remove removes the package identified by sourceIdentifier using apt-get.
func (a *Adapter) Remove(sourceIdentifier string) (*adapter.Result, error) {
	start := time.Now()
	out, exitCode, err := a.exec.Run("apt-get", "remove", "-y", sourceIdentifier)
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
			Message:          fmt.Sprintf("apt-get remove failed for %s (exit %d)", sourceIdentifier, exitCode),
			TechnicalDetails: out,
			ErrorCategory:    cat,
			Duration:         duration,
		}, nil
	}
	return &adapter.Result{
		Success:          true,
		ExitCode:         0,
		Message:          fmt.Sprintf("Removed %s via APT", sourceIdentifier),
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

// categorizeExitCode maps an apt-get exit code and output to an error category.
func categorizeExitCode(exitCode int, output string) string {
	switch exitCode {
	case 100:
		return adapter.ErrPackageNotFound
	}
	lower := strings.ToLower(output)
	if strings.Contains(lower, "permission denied") || strings.Contains(lower, "operation not permitted") {
		return adapter.ErrPermissionDenied
	}
	if strings.Contains(lower, "unable to connect") || strings.Contains(lower, "network") {
		return adapter.ErrNetworkError
	}
	if strings.Contains(lower, "dependency") || strings.Contains(lower, "conflict") {
		return adapter.ErrDependencyConflict
	}
	return adapter.ErrExecutionFailed
}
