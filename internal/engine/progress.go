package engine

// EventKind classifies a progress event emitted by the Install Engine.
type EventKind string

const (
	// EventStarted is emitted when an install or remove operation begins.
	EventStarted EventKind = "started"

	// EventValidating is emitted while the engine validates the install plan.
	EventValidating EventKind = "validating"

	// EventSelectingAdapter is emitted when the engine is choosing an adapter.
	EventSelectingAdapter EventKind = "selecting_adapter"

	// EventExecuting is emitted when the adapter is executing the operation.
	EventExecuting EventKind = "executing"

	// EventVerifying is emitted when the post-install verification is running.
	EventVerifying EventKind = "verifying"

	// EventCompleted is emitted when the operation finishes successfully.
	EventCompleted EventKind = "completed"

	// EventFailed is emitted when the operation fails.
	EventFailed EventKind = "failed"
)

// ProgressEvent is a structured event emitted by the Install Engine as it
// works through an install or remove operation.
type ProgressEvent struct {
	// Kind classifies the event.
	Kind EventKind

	// ApplicationID is the application being acted on.
	ApplicationID string

	// Message is a user-facing description of what is happening.
	Message string
}

// ProgressHandler is a function that receives progress events from the engine.
// Implementations must not block. A nil handler is valid (events are discarded).
type ProgressHandler func(ProgressEvent)
