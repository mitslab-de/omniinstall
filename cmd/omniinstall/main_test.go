package main

import (
	"errors"
	"strings"
	"testing"
)

func TestRun_NoArgs(t *testing.T) {
	if err := run([]string{}); err != nil {
		t.Errorf("expected no error for no args, got: %v", err)
	}
}

func TestRun_Version(t *testing.T) {
	if err := run([]string{"version"}); err != nil {
		t.Errorf("expected no error for 'version', got: %v", err)
	}
	if err := run([]string{"--version"}); err != nil {
		t.Errorf("expected no error for '--version', got: %v", err)
	}
	if err := run([]string{"-v"}); err != nil {
		t.Errorf("expected no error for '-v', got: %v", err)
	}
}

func TestRun_Help(t *testing.T) {
	if err := run([]string{"help"}); err != nil {
		t.Errorf("expected no error for 'help', got: %v", err)
	}
	if err := run([]string{"--help"}); err != nil {
		t.Errorf("expected no error for '--help', got: %v", err)
	}
}

func TestRun_UnknownCommand(t *testing.T) {
	err := run([]string{"notacommand"})
	if err == nil {
		t.Error("expected error for unknown command, got nil")
	}
	if !strings.Contains(err.Error(), "notacommand") {
		t.Errorf("expected error message to contain command name, got: %v", err)
	}
}

func TestRunCLI_SuccessExitCode(t *testing.T) {
	if code := runCLI([]string{"version"}); code != ExitCodeSuccess {
		t.Fatalf("expected success exit code %d, got %d", ExitCodeSuccess, code)
	}
}

func TestRunCLI_UsageExitCode(t *testing.T) {
	if code := runCLI([]string{"search"}); code != ExitCodeUsage {
		t.Fatalf("expected usage exit code %d, got %d", ExitCodeUsage, code)
	}
}

func TestRunCLI_NotFoundExitCode(t *testing.T) {
	if code := runCLI([]string{"install", "definitely-not-a-real-app"}); code != ExitCodeNotFound {
		t.Fatalf("expected not-found exit code %d, got %d", ExitCodeNotFound, code)
	}
}

func TestExitCodeForError_BackendUnavailable(t *testing.T) {
	err := errors.New("install engine error: no available adapter for source type apt")
	if code := exitCodeForError(err); code != ExitCodeBackendUnavailable {
		t.Fatalf("expected backend-unavailable exit code %d, got %d", ExitCodeBackendUnavailable, code)
	}
}

func TestExitCodeForError_ExecutionFailedDefault(t *testing.T) {
	err := errors.New("installation failed: dependency conflict")
	if code := exitCodeForError(err); code != ExitCodeExecutionFailed {
		t.Fatalf("expected execution-failed exit code %d, got %d", ExitCodeExecutionFailed, code)
	}
}
