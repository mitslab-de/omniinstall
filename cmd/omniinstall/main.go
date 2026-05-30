// Command omniinstall is the OmniInstall CLI entry point.
//
// OmniInstall is a universal Linux application discovery and installation platform.
// It allows users to find and install software without needing to understand
// package managers, package formats, or distribution-specific commands.
//
// Usage:
//
//	omniinstall <command> [arguments]
//
// Available commands:
//
//	search <query>       Search for an application
//	install <app>        Install an application
//	remove <app>         Remove an application
//	explain <app>        Show source selection details for an application
//	list                 List installed applications
//	version              Show version information
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// version is the current OmniInstall version.
const version = "0.1.0-dev"

const (
	ExitCodeSuccess            = 0
	ExitCodeUsage              = 2
	ExitCodeNotFound           = 3
	ExitCodeBackendUnavailable = 4
	ExitCodeExecutionFailed    = 5
)

func main() {
	os.Exit(runCLI(os.Args[1:]))
}

func runCLI(args []string) int {
	if err := run(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return exitCodeForError(err)
	}
	return ExitCodeSuccess
}

func exitCodeForError(err error) int {
	if err == nil {
		return ExitCodeSuccess
	}
	msg := strings.ToLower(err.Error())
	switch {
	case strings.HasPrefix(msg, "usage:"),
		strings.Contains(msg, "unknown command"):
		return ExitCodeUsage
	case strings.Contains(msg, " not found"),
		strings.Contains(msg, "no sources available for"),
		strings.Contains(msg, "not recorded as installed"):
		return ExitCodeNotFound
	case strings.Contains(msg, "no compatible source found"),
		strings.Contains(msg, "no available adapter for source type"),
		strings.Contains(msg, "backend_unavailable"):
		return ExitCodeBackendUnavailable
	default:
		return ExitCodeExecutionFailed
	}
}

func run(args []string) error {
	if len(args) == 0 {
		printUsage()
		return nil
	}

	switch args[0] {
	case "version", "--version", "-v":
		fmt.Printf("OmniInstall %s\n", version)
		return nil
	case "help", "--help", "-h":
		printUsage()
		return nil
	case "search":
		jsonOutput, positional, err := parseJSONFlag(args[1:])
		if err != nil {
			return err
		}
		if len(positional) > 1 {
			return errors.New("usage: omniinstall search <query> [--json]")
		}
		query := ""
		if len(positional) == 1 {
			query = positional[0]
		}

		app := newApp(os.Stdout)
		if jsonOutput {
			return app.searchJSON(query)
		}
		return app.search(query)
	case "install":
		jsonOutput, positional, err := parseInstallFlags(args[1:])
		if err != nil {
			return err
		}
		dryRun := false
		appID := ""
		remaining := make([]string, 0, len(positional))
		for _, a := range positional {
			if a == "--dry-run" {
				dryRun = true
				continue
			}
			remaining = append(remaining, a)
		}
		if len(remaining) > 1 {
			return errors.New("usage: omniinstall install <application> [--dry-run] [--json]")
		}
		if len(remaining) == 1 {
			appID = remaining[0]
		}
		app := newApp(os.Stdout)
		if dryRun {
			if jsonOutput {
				return app.installDryRunJSON(appID)
			}
			return app.installDryRun(appID)
		}
		return app.install(appID)
	case "remove":
		appID := ""
		if len(args) > 1 {
			appID = args[1]
		}
		app := newApp(os.Stdout)
		return app.remove(appID)
	case "verify":
		appID := ""
		if len(args) > 1 {
			appID = args[1]
		}
		app := newApp(os.Stdout)
		return app.verify(appID)
	case "explain":
		jsonOutput, positional, err := parseJSONFlag(args[1:])
		if err != nil {
			return err
		}
		if len(positional) > 1 {
			return errors.New("usage: omniinstall explain <application> [--json]")
		}
		appID := ""
		if len(positional) == 1 {
			appID = positional[0]
		}

		app := newApp(os.Stdout)
		if jsonOutput {
			return app.explainJSON(appID)
		}
		return app.explain(appID)
	case "list":
		jsonOutput, positional, err := parseJSONFlag(args[1:])
		if err != nil {
			return err
		}
		if len(positional) > 0 {
			return errors.New("usage: omniinstall list [--json]")
		}

		app := newApp(os.Stdout)
		if jsonOutput {
			return app.listJSON()
		}
		return app.list()
	case "state":
		if len(args) < 2 {
			return errors.New("usage: omniinstall state <export|import> [options]")
		}
		switch args[1] {
		case "export":
			yamlOutput := false
			for _, a := range args[2:] {
				if a == "--yaml" {
					yamlOutput = true
				}
			}
			app := newApp(os.Stdout)
			return app.stateExport(yamlOutput)
		case "import":
			filePath := ""
			force := false
			dryRun := false
			for _, a := range args[2:] {
				switch a {
				case "--force":
					force = true
				case "--dry-run":
					dryRun = true
				default:
					if !strings.HasPrefix(a, "--") && filePath == "" {
						filePath = a
					}
				}
			}
			app := newApp(os.Stdout)
			return app.stateImport(filePath, force, dryRun)
		default:
			return fmt.Errorf("unknown state subcommand %q — use 'export' or 'import'", args[1])
		}
	default:
		return fmt.Errorf("unknown command %q — run 'omniinstall help' for usage", args[0])
	}
}

func printUsage() {
	fmt.Print(`OmniInstall — universal Linux application installation

Usage:
  omniinstall <command> [arguments]

Commands:
  search <query>    Search for an application (use --json for machine output)
  install <app>     Install an application (use --dry-run to preview, --json for machine output)
  remove <app>      Remove an application
  verify <app>      Re-run verification for an installed application
  explain <app>     Show why a source was selected for an application (use --json for machine output)
  list              List installed applications (use --json for machine output)
  version           Show version information

Run 'omniinstall help' for more information.
`)
}

func parseJSONFlag(args []string) (bool, []string, error) {
	jsonOutput := false
	positional := make([]string, 0, len(args))

	for _, arg := range args {
		if arg == "--json" {
			if jsonOutput {
				return false, nil, errors.New("usage: duplicate --json flag")
			}
			jsonOutput = true
			continue
		}
		positional = append(positional, arg)
	}

	return jsonOutput, positional, nil
}

// parseInstallFlags parses the --json flag from install command arguments,
// leaving all other arguments (including --dry-run) as positional for further
// processing by the caller.
func parseInstallFlags(args []string) (jsonOutput bool, positional []string, err error) {
	return parseJSONFlag(args)
}
