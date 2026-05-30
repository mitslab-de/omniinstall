// Package logging defines structured log entry types for OmniInstall.
//
// All core components should emit structured log entries using this package.
// Logs must support debugging, diagnostics, and future auditability.
package logging

import (
	"time"

	"github.com/mitslab-de/omniinstall/internal/source"
)

// Action identifies the high-level operation being logged.
type Action string

const (
	ActionSearch  Action = "search"
	ActionResolve Action = "resolve"
	ActionInstall Action = "install"
	ActionRemove  Action = "remove"
	ActionVerify  Action = "verify"
)

// Entry is a single structured log record.
type Entry struct {
	// Timestamp is when the action occurred.
	Timestamp time.Time

	// Action identifies the high-level operation.
	Action Action

	// ApplicationID is the canonical application identifier, if applicable.
	ApplicationID string

	// SourceType identifies the backend used, if applicable.
	SourceType source.Type

	// Result is "success" or "failure".
	Result string

	// Duration is how long the operation took.
	Duration time.Duration

	// RiskLevel is the assessed risk level at the time of the action.
	RiskLevel source.RiskLevel

	// VerificationStatus is the post-install verification outcome.
	VerificationStatus string

	// ErrorCategory is a stable machine-readable error category on failure.
	ErrorCategory string
}
