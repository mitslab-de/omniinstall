package logging

// Emitter emits structured log entries.
type Emitter interface {
	Emit(entry Entry)
}

// EmitFunc adapts a function to the Emitter interface.
type EmitFunc func(entry Entry)

// Emit emits the entry using f.
func (f EmitFunc) Emit(entry Entry) {
	f(entry)
}
