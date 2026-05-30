package engine

import (
	"errors"
	"fmt"
	"strings"
	"time"

	adapter "github.com/mitslab-de/omniinstall/internal/adapters"
	"github.com/mitslab-de/omniinstall/internal/install"
	"github.com/mitslab-de/omniinstall/internal/logging"
	"github.com/mitslab-de/omniinstall/internal/source"
)

// DefaultEngine is the standard implementation of the Engine interface.
// It selects the first available adapter that can handle the requested source
// type, validates and executes the plan, then hands off to verification.
type DefaultEngine struct {
	adapters        []adapter.Adapter
	progressHandler ProgressHandler
	logger          logging.Emitter
	now             func() time.Time
}

// NewDefaultEngine creates a DefaultEngine with the given adapters and an
// optional progress handler. A nil handler discards all progress events.
func NewDefaultEngine(adapters []adapter.Adapter, handler ProgressHandler) *DefaultEngine {
	return &DefaultEngine{
		adapters:        adapters,
		progressHandler: handler,
		now:             time.Now,
	}
}

// WithLogger configures structured action logging for engine operations.
func (e *DefaultEngine) WithLogger(logger logging.Emitter) *DefaultEngine {
	e.logger = logger
	return e
}

func (e *DefaultEngine) emitLog(entry logging.Entry) {
	if e.logger == nil {
		return
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = e.now()
	}
	e.logger.Emit(entry)
}

// emit sends a progress event if a handler is configured.
func (e *DefaultEngine) emit(kind EventKind, applicationID, message string) {
	if e.progressHandler != nil {
		e.progressHandler(ProgressEvent{
			Kind:          kind,
			ApplicationID: applicationID,
			Message:       message,
		})
	}
}

// Install validates the plan, selects the appropriate adapter, executes the
// installation, and hands off to post-install verification.
func (e *DefaultEngine) Install(plan *install.Plan) (*Result, error) {
	start := time.Now()
	if plan == nil {
		return nil, errors.New("plan must not be nil")
	}

	appID := plan.ApplicationID
	e.emit(EventStarted, appID, fmt.Sprintf("Starting installation of %s", appID))

	// 1. Validate the plan.
	e.emit(EventValidating, appID, "Validating install plan")
	if err := plan.Validate(); err != nil {
		e.emit(EventFailed, appID, "Plan validation failed: "+err.Error())
		res := &Result{
			Success:       false,
			ApplicationID: appID,
			Message:       "Plan validation failed: " + err.Error(),
			ErrorCategory: "invalid_plan",
		}
		e.emitLog(logging.Entry{
			Action:        logging.ActionInstall,
			ApplicationID: appID,
			SourceType:    plan.SourceType,
			Result:        "failure",
			Duration:      time.Since(start),
			RiskLevel:     plan.RiskLevel,
			ErrorCategory: res.ErrorCategory,
		})
		return res, nil
	}

	// 2. Select the adapter.
	e.emit(EventSelectingAdapter, appID, fmt.Sprintf("Selecting adapter for source type %s", plan.SourceType))
	a, err := e.selectAdapter(plan)
	if err != nil {
		e.emit(EventFailed, appID, err.Error())
		res := &Result{
			Success:       false,
			ApplicationID: appID,
			Message:       err.Error(),
			ErrorCategory: adapter.ErrBackendUnavailable,
		}
		e.emitLog(logging.Entry{
			Action:        logging.ActionInstall,
			ApplicationID: appID,
			SourceType:    plan.SourceType,
			Result:        "failure",
			Duration:      time.Since(start),
			RiskLevel:     plan.RiskLevel,
			ErrorCategory: res.ErrorCategory,
		})
		return res, nil
	}

	// 2b. Run adapter preflight checks.
	e.emit(EventValidating, appID, fmt.Sprintf("Running preflight checks for %s", a.Name()))
	if err := preflightCheckAdapter(a, plan.SourceIdentifier); err != nil {
		msg := "Preflight check failed: " + err.Error()
		e.emit(EventFailed, appID, msg)
		res := &Result{
			Success:       false,
			ApplicationID: appID,
			Message:       msg,
			ErrorCategory: categorizePreflightError(err),
		}
		e.emitLog(logging.Entry{
			Action:        logging.ActionInstall,
			ApplicationID: appID,
			SourceType:    plan.SourceType,
			Result:        "failure",
			Duration:      time.Since(start),
			RiskLevel:     plan.RiskLevel,
			ErrorCategory: res.ErrorCategory,
		})
		return res, nil
	}

	// 3. Execute installation.
	e.emit(EventExecuting, appID, fmt.Sprintf("Installing %s via %s", appID, a.Name()))
	result, err := a.Install(plan)
	if err != nil {
		e.emit(EventFailed, appID, "Adapter error: "+err.Error())
		res := &Result{
			Success:       false,
			ApplicationID: appID,
			Message:       "Adapter error: " + err.Error(),
			ErrorCategory: adapter.ErrExecutionFailed,
		}
		e.emitLog(logging.Entry{
			Action:        logging.ActionInstall,
			ApplicationID: appID,
			SourceType:    plan.SourceType,
			Result:        "failure",
			Duration:      time.Since(start),
			RiskLevel:     plan.RiskLevel,
			ErrorCategory: res.ErrorCategory,
		})
		return res, nil
	}
	if !result.Success {
		e.emit(EventFailed, appID, result.Message)
		res := &Result{
			Success:       false,
			ApplicationID: appID,
			Message:       result.Message,
			ErrorCategory: result.ErrorCategory,
		}
		e.emitLog(logging.Entry{
			Action:        logging.ActionInstall,
			ApplicationID: appID,
			SourceType:    plan.SourceType,
			Result:        "failure",
			Duration:      time.Since(start),
			RiskLevel:     plan.RiskLevel,
			ErrorCategory: res.ErrorCategory,
		})
		return res, nil
	}

	// 4. Verify installation.
	e.emit(EventVerifying, appID, fmt.Sprintf("Verifying installation of %s", appID))
	verifyStart := time.Now()
	vr, err := a.Verify(plan)
	if err != nil {
		// Verification failure is non-fatal but surfaced in the result.
		e.emit(EventCompleted, appID, fmt.Sprintf("Installed %s (verification error: %s)", appID, err.Error()))
		res := &Result{
			Success:       true,
			ApplicationID: appID,
			Message:       fmt.Sprintf("Installed %s (verification error: %s)", appID, err.Error()),
		}
		e.emitLog(logging.Entry{
			Action:             logging.ActionVerify,
			ApplicationID:      appID,
			SourceType:         plan.SourceType,
			Result:             "failure",
			Duration:           time.Since(verifyStart),
			RiskLevel:          plan.RiskLevel,
			VerificationStatus: "error",
			ErrorCategory:      adapter.ErrExecutionFailed,
		})
		e.emitLog(logging.Entry{
			Action:        logging.ActionInstall,
			ApplicationID: appID,
			SourceType:    plan.SourceType,
			Result:        "success",
			Duration:      time.Since(start),
			RiskLevel:     plan.RiskLevel,
		})
		return res, nil
	}

	if !vr.Verified {
		e.emit(EventCompleted, appID, fmt.Sprintf("Installed %s (verification failed: %s)", appID, vr.Details))
		res := &Result{
			Success:       true,
			ApplicationID: appID,
			Message:       fmt.Sprintf("Installed %s (verification failed: %s)", appID, vr.Details),
		}
		e.emitLog(logging.Entry{
			Action:             logging.ActionVerify,
			ApplicationID:      appID,
			SourceType:         plan.SourceType,
			Result:             "failure",
			Duration:           time.Since(verifyStart),
			RiskLevel:          plan.RiskLevel,
			VerificationStatus: "failed",
			ErrorCategory:      "verification_failed",
		})
		e.emitLog(logging.Entry{
			Action:        logging.ActionInstall,
			ApplicationID: appID,
			SourceType:    plan.SourceType,
			Result:        "success",
			Duration:      time.Since(start),
			RiskLevel:     plan.RiskLevel,
		})
		return res, nil
	}

	e.emit(EventCompleted, appID, fmt.Sprintf("Successfully installed %s", appID))
	e.emitLog(logging.Entry{
		Action:             logging.ActionVerify,
		ApplicationID:      appID,
		SourceType:         plan.SourceType,
		Result:             "success",
		Duration:           time.Since(verifyStart),
		RiskLevel:          plan.RiskLevel,
		VerificationStatus: "verified",
	})
	e.emitLog(logging.Entry{
		Action:        logging.ActionInstall,
		ApplicationID: appID,
		SourceType:    plan.SourceType,
		Result:        "success",
		Duration:      time.Since(start),
		RiskLevel:     plan.RiskLevel,
	})
	return &Result{
		Success:       true,
		ApplicationID: appID,
		Message:       fmt.Sprintf("Successfully installed %s", appID),
	}, nil
}

// Remove removes the application identified by applicationID using the source
// metadata recorded at install time.
func (e *DefaultEngine) Remove(applicationID string, sourceType source.Type, sourceIdentifier string) (*Result, error) {
	start := time.Now()
	if applicationID == "" {
		return nil, errors.New("applicationID must not be empty")
	}
	if sourceType == "" {
		return nil, errors.New("sourceType must not be empty")
	}
	if sourceIdentifier == "" {
		return nil, errors.New("sourceIdentifier must not be empty")
	}

	e.emit(EventStarted, applicationID, fmt.Sprintf("Starting removal of %s", applicationID))

	a := e.selectRemovalAdapter(sourceType)
	if a == nil {
		msg := fmt.Sprintf("no available adapter for source type %s", sourceType)
		e.emit(EventFailed, applicationID, msg)
		res := &Result{
			Success:       false,
			ApplicationID: applicationID,
			Message:       msg,
			ErrorCategory: adapter.ErrBackendUnavailable,
		}
		e.emitLog(logging.Entry{
			Action:        logging.ActionRemove,
			ApplicationID: applicationID,
			SourceType:    sourceType,
			Result:        "failure",
			Duration:      time.Since(start),
			ErrorCategory: res.ErrorCategory,
		})
		return res, nil
	}

	if err := preflightCheckAdapter(a, sourceIdentifier); err != nil {
		msg := "Preflight check failed: " + err.Error()
		e.emit(EventFailed, applicationID, msg)
		res := &Result{
			Success:       false,
			ApplicationID: applicationID,
			Message:       msg,
			ErrorCategory: categorizePreflightError(err),
		}
		e.emitLog(logging.Entry{
			Action:        logging.ActionRemove,
			ApplicationID: applicationID,
			SourceType:    sourceType,
			Result:        "failure",
			Duration:      time.Since(start),
			ErrorCategory: res.ErrorCategory,
		})
		return res, nil
	}

	e.emit(EventExecuting, applicationID, fmt.Sprintf("Removing %s via %s", sourceIdentifier, a.Name()))
	result, err := a.Remove(sourceIdentifier)
	if err != nil {
		e.emit(EventFailed, applicationID, "Adapter error: "+err.Error())
		res := &Result{
			Success:       false,
			ApplicationID: applicationID,
			Message:       "Adapter error: " + err.Error(),
			ErrorCategory: adapter.ErrExecutionFailed,
		}
		e.emitLog(logging.Entry{
			Action:        logging.ActionRemove,
			ApplicationID: applicationID,
			SourceType:    sourceType,
			Result:        "failure",
			Duration:      time.Since(start),
			ErrorCategory: res.ErrorCategory,
		})
		return res, nil
	}

	if result.Success {
		e.emit(EventCompleted, applicationID, fmt.Sprintf("Successfully removed %s", applicationID))
	} else {
		e.emit(EventFailed, applicationID, result.Message)
	}

	res := &Result{
		Success:       result.Success,
		ApplicationID: applicationID,
		Message:       result.Message,
		ErrorCategory: result.ErrorCategory,
	}
	entryResult := "success"
	if !res.Success {
		entryResult = "failure"
	}
	e.emitLog(logging.Entry{
		Action:        logging.ActionRemove,
		ApplicationID: applicationID,
		SourceType:    sourceType,
		Result:        entryResult,
		Duration:      time.Since(start),
		ErrorCategory: res.ErrorCategory,
	})
	return res, nil
}

// selectAdapter returns the first available adapter that can handle the source
// type specified in the plan.
func (e *DefaultEngine) selectAdapter(plan *install.Plan) (adapter.Adapter, error) {
	for _, a := range e.adapters {
		if a.IsAvailable() && a.CanHandle(plan.SourceType) {
			return a, nil
		}
	}
	return nil, fmt.Errorf("no available adapter for source type %s", plan.SourceType)
}

// selectRemovalAdapter returns the first available adapter that supports the
// given source type.
func (e *DefaultEngine) selectRemovalAdapter(sourceType source.Type) adapter.Adapter {
	for _, a := range e.adapters {
		if a.IsAvailable() && a.CanHandle(sourceType) {
			return a
		}
	}
	return nil
}

func preflightCheckAdapter(a adapter.Adapter, sourceIdentifier string) error {
	if !a.IsAvailable() {
		return fmt.Errorf("backend %s is not available", a.Name())
	}
	_, err := a.CheckInstalled(sourceIdentifier)
	if err != nil {
		return fmt.Errorf("backend readiness check failed: %w", err)
	}
	return nil
}

func categorizePreflightError(err error) string {
	lower := strings.ToLower(err.Error())
	if strings.Contains(lower, "permission denied") || strings.Contains(lower, "operation not permitted") {
		return adapter.ErrPermissionDenied
	}
	return adapter.ErrBackendUnavailable
}
