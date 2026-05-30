// Package adapter defines the Adapter interface for OmniInstall.
//
// Adapters isolate backend-specific package manager behavior from the core.
// The Install Engine calls adapters through this stable interface.
package adapter

import (
	"time"

	"github.com/mitslab-de/omniinstall/internal/install"
	"github.com/mitslab-de/omniinstall/internal/source"
)

// InstalledState describes whether a package is currently installed.
type InstalledState string

const (
	StateInstalled    InstalledState = "installed"
	StateNotInstalled InstalledState = "not-installed"
	StateUnknown      InstalledState = "unknown"
)

// Result is the structured outcome of an adapter operation.
type Result struct {
	// Success indicates whether the operation completed successfully.
	Success bool

	// ExitCode is the exit code from the backend process, if applicable.
	ExitCode int

	// Message is a user-facing description of what happened.
	Message string

	// TechnicalDetails contains backend-specific diagnostic information.
	TechnicalDetails string

	// ErrorCategory is a stable machine-readable error category on failure.
	ErrorCategory string

	// ChangedSystem indicates whether the system state was actually changed.
	ChangedSystem bool

	// Duration is the time taken to complete the operation.
	Duration time.Duration
}

// VerificationResult is the outcome of a post-install verification check.
type VerificationResult struct {
	// Verified indicates whether verification succeeded.
	Verified bool

	// Details describes what was checked and what was found.
	Details string
}

// Adapter is the interface that all package manager backends must implement.
type Adapter interface {
	// Name returns the adapter's human-readable name.
	Name() string

	// IsAvailable returns true if this backend is usable on the current system.
	IsAvailable() bool

	// CanHandle returns true if this adapter supports the given source type.
	CanHandle(sourceType source.Type) bool

	// CheckInstalled returns the installation state for the given identifier.
	CheckInstalled(sourceIdentifier string) (InstalledState, error)

	// Install executes the install plan and returns a structured result.
	Install(plan *install.Plan) (*Result, error)

	// Remove removes the application identified by sourceIdentifier.
	Remove(sourceIdentifier string) (*Result, error)

	// Verify checks whether the installation of sourceIdentifier succeeded.
	Verify(plan *install.Plan) (*VerificationResult, error)
}

// ErrorCategories defines stable error category strings for adapter results.
const (
	ErrBackendUnavailable   = "backend_unavailable"
	ErrPackageNotFound      = "package_not_found"
	ErrPermissionDenied     = "permission_denied"
	ErrNetworkError         = "network_error"
	ErrDependencyConflict   = "dependency_conflict"
	ErrExecutionFailed      = "execution_failed"
	ErrUnsupportedOperation = "unsupported_operation"
)
