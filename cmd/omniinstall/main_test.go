package main

import (
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
