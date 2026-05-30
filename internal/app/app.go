// Package app defines the Application domain model for OmniInstall.
//
// An application is the user-facing identity of software.
// It is independent of any specific package manager or format.
package app

import (
	"errors"
	"regexp"
	"strings"
)

// validIDPattern checks that an application ID is lowercase kebab-case ASCII.
var validIDPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// Application represents a user-facing software application in OmniInstall.
// It is separate from any package-manager-specific representation.
type Application struct {
	// ID is the canonical OmniInstall identifier (lowercase, kebab-case).
	ID string

	// DisplayName is the human-readable application name.
	DisplayName string

	// Summary is a short one-line description.
	Summary string

	// Categories are user-facing category labels.
	Categories []string

	// Aliases are alternative search names.
	Aliases []string

	// Homepage is the upstream project URL.
	Homepage string

	// License is the SPDX license identifier.
	License string

	// Publisher is the name of the publisher or maintainer.
	Publisher string

	// Description is the full user-facing description.
	Description string
}

// Validate checks that the Application has all required fields and valid format.
func (a *Application) Validate() error {
	if a.ID == "" {
		return errors.New("application id is required")
	}
	if !validIDPattern.MatchString(a.ID) {
		return errors.New("application id must be lowercase kebab-case ASCII")
	}
	if strings.TrimSpace(a.DisplayName) == "" {
		return errors.New("display_name is required")
	}
	if strings.TrimSpace(a.Summary) == "" {
		return errors.New("summary is required")
	}
	if len(a.Categories) == 0 {
		return errors.New("at least one category is required")
	}
	return nil
}
