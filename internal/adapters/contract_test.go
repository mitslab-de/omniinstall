package adapter_test

import (
	"testing"
	"time"

	adapter "github.com/mitslab-de/omniinstall/internal/adapters"
	"github.com/mitslab-de/omniinstall/internal/install"
	"github.com/mitslab-de/omniinstall/internal/source"
)

// RunAdapterContract runs the behavioral contract test suite against any
// Adapter implementation. Call this from each adapter's test package to
// ensure it satisfies the full adapter contract.
func RunAdapterContract(t *testing.T, a adapter.Adapter) {
	t.Helper()

	t.Run("NameIsNonEmpty", func(t *testing.T) {
		if a.Name() == "" {
			t.Error("adapter Name() must return a non-empty string")
		}
	})

	t.Run("CanHandleUnknownReturnsFalse", func(t *testing.T) {
		if a.CanHandle("unknown-source-type-xyz") {
			t.Error("CanHandle must return false for unknown source type")
		}
	})
}

// contractMockAdapter is a minimal adapter implementation used to exercise
// the contract test helper itself.
type contractMockAdapter struct {
	name      string
	available bool
}

func (m *contractMockAdapter) Name() string       { return m.name }
func (m *contractMockAdapter) IsAvailable() bool  { return m.available }
func (m *contractMockAdapter) CanHandle(t source.Type) bool { return t == source.TypeAPT }
func (m *contractMockAdapter) CheckInstalled(id string) (adapter.InstalledState, error) {
	return adapter.StateUnknown, nil
}
func (m *contractMockAdapter) Install(plan *install.Plan) (*adapter.Result, error) {
	return &adapter.Result{Success: true, Message: "ok", ChangedSystem: true, Duration: time.Millisecond}, nil
}
func (m *contractMockAdapter) Remove(id string) (*adapter.Result, error) {
	return &adapter.Result{Success: true, Message: "ok", ChangedSystem: true, Duration: time.Millisecond}, nil
}
func (m *contractMockAdapter) Verify(plan *install.Plan) (*adapter.VerificationResult, error) {
	return &adapter.VerificationResult{Verified: true, Details: "ok"}, nil
}

func TestAdapterContractWithMock(t *testing.T) {
	m := &contractMockAdapter{name: "mock", available: true}
	RunAdapterContract(t, m)
}

func TestResultDurationField(t *testing.T) {
	r := &adapter.Result{
		Success:  true,
		Message:  "done",
		Duration: 500 * time.Millisecond,
	}
	if r.Duration != 500*time.Millisecond {
		t.Errorf("expected 500ms, got %v", r.Duration)
	}
}
