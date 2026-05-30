package engine

import (
	"errors"
	"fmt"

	adapter "github.com/mitslab-de/omniinstall/internal/adapters"
	"github.com/mitslab-de/omniinstall/internal/install"
)

// DefaultEngine is the standard implementation of the Engine interface.
// It selects the first available adapter that can handle the requested source
// type, validates and executes the plan, then hands off to verification.
type DefaultEngine struct {
	adapters        []adapter.Adapter
	progressHandler ProgressHandler
}

// NewDefaultEngine creates a DefaultEngine with the given adapters and an
// optional progress handler. A nil handler discards all progress events.
func NewDefaultEngine(adapters []adapter.Adapter, handler ProgressHandler) *DefaultEngine {
	return &DefaultEngine{
		adapters:        adapters,
		progressHandler: handler,
	}
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
	if plan == nil {
		return nil, errors.New("plan must not be nil")
	}

	appID := plan.ApplicationID
	e.emit(EventStarted, appID, fmt.Sprintf("Starting installation of %s", appID))

	// 1. Validate the plan.
	e.emit(EventValidating, appID, "Validating install plan")
	if err := plan.Validate(); err != nil {
		e.emit(EventFailed, appID, "Plan validation failed: "+err.Error())
		return &Result{
			Success:       false,
			ApplicationID: appID,
			Message:       "Plan validation failed: " + err.Error(),
			ErrorCategory: "invalid_plan",
		}, nil
	}

	// 2. Select the adapter.
	e.emit(EventSelectingAdapter, appID, fmt.Sprintf("Selecting adapter for source type %s", plan.SourceType))
	a, err := e.selectAdapter(plan)
	if err != nil {
		e.emit(EventFailed, appID, err.Error())
		return &Result{
			Success:       false,
			ApplicationID: appID,
			Message:       err.Error(),
			ErrorCategory: adapter.ErrBackendUnavailable,
		}, nil
	}

	// 3. Execute installation.
	e.emit(EventExecuting, appID, fmt.Sprintf("Installing %s via %s", appID, a.Name()))
	result, err := a.Install(plan)
	if err != nil {
		e.emit(EventFailed, appID, "Adapter error: "+err.Error())
		return &Result{
			Success:       false,
			ApplicationID: appID,
			Message:       "Adapter error: " + err.Error(),
			ErrorCategory: adapter.ErrExecutionFailed,
		}, nil
	}
	if !result.Success {
		e.emit(EventFailed, appID, result.Message)
		return &Result{
			Success:       false,
			ApplicationID: appID,
			Message:       result.Message,
			ErrorCategory: result.ErrorCategory,
		}, nil
	}

	// 4. Verify installation.
	e.emit(EventVerifying, appID, fmt.Sprintf("Verifying installation of %s", appID))
	vr, err := a.Verify(plan)
	if err != nil {
		// Verification failure is non-fatal but surfaced in the result.
		e.emit(EventCompleted, appID, fmt.Sprintf("Installed %s (verification error: %s)", appID, err.Error()))
		return &Result{
			Success:       true,
			ApplicationID: appID,
			Message:       fmt.Sprintf("Installed %s (verification error: %s)", appID, err.Error()),
		}, nil
	}

	if !vr.Verified {
		e.emit(EventCompleted, appID, fmt.Sprintf("Installed %s (verification failed: %s)", appID, vr.Details))
		return &Result{
			Success:       true,
			ApplicationID: appID,
			Message:       fmt.Sprintf("Installed %s (verification failed: %s)", appID, vr.Details),
		}, nil
	}

	e.emit(EventCompleted, appID, fmt.Sprintf("Successfully installed %s", appID))
	return &Result{
		Success:       true,
		ApplicationID: appID,
		Message:       fmt.Sprintf("Successfully installed %s", appID),
	}, nil
}

// Remove removes the application identified by applicationID. It attempts
// removal through the first available adapter.
func (e *DefaultEngine) Remove(applicationID string) (*Result, error) {
	if applicationID == "" {
		return nil, errors.New("applicationID must not be empty")
	}

	e.emit(EventStarted, applicationID, fmt.Sprintf("Starting removal of %s", applicationID))

	a := e.firstAvailableAdapter()
	if a == nil {
		msg := "no available adapter for removal"
		e.emit(EventFailed, applicationID, msg)
		return &Result{
			Success:       false,
			ApplicationID: applicationID,
			Message:       msg,
			ErrorCategory: adapter.ErrBackendUnavailable,
		}, nil
	}

	e.emit(EventExecuting, applicationID, fmt.Sprintf("Removing %s via %s", applicationID, a.Name()))
	result, err := a.Remove(applicationID)
	if err != nil {
		e.emit(EventFailed, applicationID, "Adapter error: "+err.Error())
		return &Result{
			Success:       false,
			ApplicationID: applicationID,
			Message:       "Adapter error: " + err.Error(),
			ErrorCategory: adapter.ErrExecutionFailed,
		}, nil
	}

	if result.Success {
		e.emit(EventCompleted, applicationID, fmt.Sprintf("Successfully removed %s", applicationID))
	} else {
		e.emit(EventFailed, applicationID, result.Message)
	}

	return &Result{
		Success:       result.Success,
		ApplicationID: applicationID,
		Message:       result.Message,
		ErrorCategory: result.ErrorCategory,
	}, nil
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

// firstAvailableAdapter returns the first adapter that reports IsAvailable(),
// or nil if none are available.
func (e *DefaultEngine) firstAvailableAdapter() adapter.Adapter {
	for _, a := range e.adapters {
		if a.IsAvailable() {
			return a
		}
	}
	return nil
}
