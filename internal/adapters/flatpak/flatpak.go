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
//
// It checks the system installation first, then the user installation.
// A non-zero exit code from `flatpak info` means the application is not found
// in that scope; any stderr that looks like a backend error (permissions,
// no installations configured) is surfaced as StateUnknown.
func (a *Adapter) CheckInstalled(sourceIdentifier string) (adapter.InstalledState, error) {
	// Try system scope first, then user scope.
	for _, scope := range []string{"--system", "--user"} {
		out, exitCode, err := a.exec.Run("flatpak", "info", scope, sourceIdentifier)
		if err != nil && exitCode == 0 {
			// Command could not be launched at all.
			return adapter.StateUnknown, fmt.Errorf("flatpak info failed: %w", err)
		}
		if exitCode != 0 {
			// Surface backend errors so callers can distinguish "not found"
			// from "flatpak is broken".
			if errState, errMsg := classifyFlatpakInfoError(out); errState != "" {
				return errState, fmt.Errorf("%s", errMsg)
			}
			// Normal "not installed in this scope" — continue to next scope.
			continue
		}
		// Output may include an "ID:" or "Application:" line; confirm it.
		if parseFlatpakInfoInstalled(out) {
			return adapter.StateInstalled, nil
		}
	}
	return adapter.StateNotInstalled, nil
}

// parseFlatpakInfoInstalled returns true when `flatpak info` output contains
// a non-empty ID or Application field, confirming the app is installed.
func parseFlatpakInfoInstalled(out string) bool {
	out = strings.TrimSpace(out)
	if out == "" {
		return false
	}
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		// flatpak info prints "ID:         <id>" or "Application: <id>"
		if strings.HasPrefix(line, "ID:") || strings.HasPrefix(line, "Application:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
				return true
			}
		}
	}
	// If we got output but no recognised fields, assume installed
	// (older flatpak versions may omit the ID header).
	return true
}

// classifyFlatpakInfoError inspects error output from a failed `flatpak info`
// invocation and returns (StateUnknown, message) when the failure indicates a
// backend problem rather than a simple "not installed". Returns ("", "") when
// the error is a normal "package not found" situation.
func classifyFlatpakInfoError(out string) (adapter.InstalledState, string) {
	lower := strings.ToLower(out)
	if strings.Contains(lower, "no installations") {
		return adapter.StateUnknown, "flatpak has no installations configured"
	}
	if strings.Contains(lower, "permission denied") {
		return adapter.StateUnknown, "flatpak permission denied"
	}
	return "", ""
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
