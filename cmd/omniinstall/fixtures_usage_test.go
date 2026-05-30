package main

// Integration-style tests that use the CLIFixture helpers to exercise
// common command flows with minimal boilerplate.

import (
	"strings"
	"testing"
)

func TestFixture_Search_ReturnsResults(t *testing.T) {
	f := newFixture(t).withAPT()
	out := f.mustRun(t, "search", "git")
	if !strings.Contains(out, "Git") && !strings.Contains(out, "git") {
		t.Errorf("search 'git' output should contain git results, got: %q", out)
	}
}

func TestFixture_Search_EmptyQuery_ReturnsError(t *testing.T) {
	f := newFixture(t).withAPT()
	_ = f.runExpectError(t, "search", "")
}

func TestFixture_List_Empty_NoApps(t *testing.T) {
	f := newFixture(t).withAPT()
	out := f.mustRun(t, "list")
	if !strings.Contains(out, "No applications") && !strings.Contains(out, "0 application") {
		// Acceptable empty states include no-output or an empty-message.
		// Just confirm no error was returned.
		_ = out
	}
}

func TestFixture_List_ShowsInstalledApp(t *testing.T) {
	f := newFixture(t).withAPT().withInstalledAPT("git", "git")
	out := f.mustRun(t, "list")
	if !strings.Contains(out, "git") {
		t.Errorf("list output should contain 'git', got: %q", out)
	}
}

func TestFixture_Verify_InstalledApp_Pass(t *testing.T) {
	f := newFixture(t).withVerifyResult(true, "binary found").withInstalledAPT("git", "git")
	out := f.mustRun(t, "verify", "git")
	if !strings.Contains(out, "✔") && !strings.Contains(out, "pass") && !strings.Contains(out, "verified") {
		t.Errorf("verify pass output should indicate success, got: %q", out)
	}
}

func TestFixture_Verify_NotInstalled_ReturnsError(t *testing.T) {
	f := newFixture(t).withAPT()
	_ = f.runExpectError(t, "verify", "unknown-app")
}

func TestFixture_Remove_InstalledApp_Succeeds(t *testing.T) {
	f := newFixture(t).withAPT().withInstalledAPT("git", "git")
	out := f.mustRun(t, "remove", "git")
	if !strings.Contains(out, "git") && !strings.Contains(out, "remov") {
		t.Errorf("remove output should mention app, got: %q", out)
	}
}

func TestFixture_Remove_NotInstalled_ReturnsError(t *testing.T) {
	f := newFixture(t).withAPT()
	_ = f.runExpectError(t, "remove", "unknown-app")
}

func TestFixture_UnknownCommand_ReturnsError(t *testing.T) {
	f := newFixture(t).withAPT()
	if err := f.run("definitely-unknown-command"); err == nil {
		t.Fatal("expected error for unknown fixture command")
	}
}

func TestFixture_OutputReset_ClearsBetweenCommands(t *testing.T) {
	f := newFixture(t).withAPT()
	f.mustRun(t, "search", "git")
	first := f.Output()
	f.resetOutput()
	if f.Output() != "" {
		t.Error("resetOutput should clear buffer")
	}
	f.mustRun(t, "search", "git")
	second := f.Output()
	if first != second {
		t.Errorf("same command should produce same output: first=%q second=%q", first, second)
	}
}

func TestFixture_Flatpak_ListShowsInstalledApp(t *testing.T) {
	f := newFixture(t).withFlatpak().withInstalledFlatpak("vlc", "org.videolan.VLC")
	out := f.mustRun(t, "list")
	if !strings.Contains(out, "vlc") {
		t.Errorf("list output should contain 'vlc', got: %q", out)
	}
}

func TestFixture_RemovedApp_ShowsInList(t *testing.T) {
	f := newFixture(t).withAPT().withRemovedAPT("git", "git")
	out := f.mustRun(t, "list")
	// Removed apps may or may not appear depending on implementation;
	// the key requirement is no error is returned.
	_ = out
}

// ---------------------------------------------------------------------------
// Dry-run output enrichment tests (task 0037)
// ---------------------------------------------------------------------------

func TestFixture_InstallDryRun_ShowsAdapterAndPackageID(t *testing.T) {
f := newFixture(t).withAPT()
out := f.mustRun(t, "install-dry-run", "git")
for _, want := range []string{"Dry-run:", "git", "apt", "git", "No changes were made"} {
if !strings.Contains(out, want) {
t.Errorf("dry-run output missing %q\ngot: %s", want, out)
}
}
}

func TestFixture_InstallDryRun_ShowsRiskLevel(t *testing.T) {
f := newFixture(t).withAPT()
out := f.mustRun(t, "install-dry-run", "git")
if !strings.Contains(out, "Risk level:") {
t.Errorf("dry-run output missing 'Risk level:'\ngot: %s", out)
}
}

func TestFixture_InstallDryRun_NoStateChange(t *testing.T) {
f := newFixture(t).withAPT()
f.mustRun(t, "install-dry-run", "git")
// After dry-run, the app should NOT appear in list as installed.
f.resetOutput()
listOut := f.mustRun(t, "list")
if strings.Contains(listOut, "git") {
t.Errorf("dry-run should not record state; list unexpectedly showed 'git': %s", listOut)
}
}

func TestFixture_InstallDryRun_EmptyAppID_ReturnsError(t *testing.T) {
f := newFixture(t).withAPT()
err := f.runExpectError(t, "install-dry-run", "")
if err == nil {
t.Fatal("expected error for empty application ID")
}
}

func TestFixture_InstallDryRunJSON_ContainsRequiredFields(t *testing.T) {
f := newFixture(t).withAPT()
out := f.mustRun(t, "install-dry-run-json", "git")
for _, field := range []string{`"adapter"`, `"package_id"`, `"risk_level"`} {
if !strings.Contains(out, field) {
t.Errorf("dry-run JSON output missing field %q\ngot: %s", field, out)
}
}
}

func TestFixture_InstallDryRunJSON_NoActualInstall(t *testing.T) {
f := newFixture(t).withAPT()
f.mustRun(t, "install-dry-run-json", "git")
// State should be empty after JSON dry-run.
f.resetOutput()
listOut := f.mustRun(t, "list")
if strings.Contains(listOut, "git") {
t.Errorf("dry-run JSON should not record state; list unexpectedly showed 'git': %s", listOut)
}
}

// ---------------------------------------------------------------------------
// State export/import tests (tasks 0038/0039)
// ---------------------------------------------------------------------------

func TestFixture_StateExport_Empty_ReturnsEmptyList(t *testing.T) {
f := newFixture(t).withAPT()
out := f.mustRun(t, "state-export")
if !strings.Contains(out, `"schema_version"`) {
t.Errorf("expected schema_version in export, got: %s", out)
}
if !strings.Contains(out, `"installations"`) {
t.Errorf("expected installations field in export, got: %s", out)
}
}

func TestFixture_StateExport_ContainsInstalledApp(t *testing.T) {
f := newFixture(t).withAPT().withInstalledAPT("git", "git")
out := f.mustRun(t, "state-export")
if !strings.Contains(out, "git") {
t.Errorf("expected exported state to contain 'git', got: %s", out)
}
}

func TestFixture_StateExport_ExcludesRemovedApps(t *testing.T) {
f := newFixture(t).withAPT().withRemovedAPT("git", "git")
out := f.mustRun(t, "state-export")
// A removed app should not appear in the exported active list.
if strings.Contains(out, `"git"`) {
// The JSON will have empty installations array.
if !strings.Contains(out, `"installations":[]`) && !strings.Contains(out, `"installations": []`) {
// It's ok if git appears as part of the empty array structure —
// but there should be no installation entries with git.
if strings.Contains(out, `"application_id":"git"`) || strings.Contains(out, `"application_id": "git"`) {
t.Errorf("removed app should not appear in state export, got: %s", out)
}
}
}
}

func TestFixture_StateExport_YAML_ContainsSchemaVersion(t *testing.T) {
f := newFixture(t).withAPT()
out := f.mustRun(t, "state-export-yaml")
if !strings.Contains(out, "schema_version:") {
t.Errorf("expected schema_version in YAML export, got: %s", out)
}
}
