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
)

// version is the current OmniInstall version.
const version = "0.1.0-dev"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
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
