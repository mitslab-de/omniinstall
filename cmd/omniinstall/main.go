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
		query := ""
		if len(args) > 1 {
			query = args[1]
		}
		app := newApp(os.Stdout)
		return app.search(query)
	case "install":
		appID := ""
		if len(args) > 1 {
			appID = args[1]
		}
		app := newApp(os.Stdout)
		return app.install(appID)
	case "remove":
		appID := ""
		if len(args) > 1 {
			appID = args[1]
		}
		app := newApp(os.Stdout)
		return app.remove(appID)
	case "explain":
		appID := ""
		if len(args) > 1 {
			appID = args[1]
		}
		app := newApp(os.Stdout)
		return app.explain(appID)
	case "list":
		app := newApp(os.Stdout)
		return app.list()
	default:
		return fmt.Errorf("unknown command %q — run 'omniinstall help' for usage", args[0])
	}
}

func printUsage() {
	fmt.Print(`OmniInstall — universal Linux application installation

Usage:
  omniinstall <command> [arguments]

Commands:
  search <query>    Search for an application
  install <app>     Install an application
  remove <app>      Remove an application
  explain <app>     Show why a source was selected for an application
  list              List installed applications
  version           Show version information

Run 'omniinstall help' for more information.
`)
}
