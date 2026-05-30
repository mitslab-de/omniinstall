package logging_test

import (
	"testing"
	"time"

	"github.com/mitslab-de/omniinstall/internal/logging"
	"github.com/mitslab-de/omniinstall/internal/source"
)

func TestEntry_Fields(t *testing.T) {
	e := &logging.Entry{
		Timestamp:     time.Now(),
		Action:        logging.ActionInstall,
		ApplicationID: "obs-studio",
		SourceType:    source.TypeAPT,
		Result:        "success",
		Duration:      2 * time.Second,
		RiskLevel:     source.RiskLow,
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
	if e.Action != logging.ActionInstall {
		t.Errorf("expected action %q, got %q", logging.ActionInstall, e.Action)
	}
	if e.ApplicationID != "obs-studio" {
		t.Errorf("expected application_id 'obs-studio', got %q", e.ApplicationID)
	}
	if e.Result != "success" {
		t.Errorf("expected result 'success', got %q", e.Result)
	}
}

func TestActionConstants(t *testing.T) {
	actions := []logging.Action{
		logging.ActionSearch,
		logging.ActionResolve,
		logging.ActionInstall,
		logging.ActionRemove,
		logging.ActionVerify,
	}
	for _, a := range actions {
		if a == "" {
			t.Error("expected non-empty action constant")
		}
	}
}
