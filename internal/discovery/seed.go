package discovery

import (
	"github.com/mitslab-de/omniinstall/internal/app"
)

// MVPCatalog returns a Catalog pre-loaded with MVP seed applications.
// These are the recommended set of applications from specs/01-APPLICATION_MODEL.md.
func MVPCatalog() *Catalog {
	catalog := NewCatalog()
	for _, a := range mvpApplications {
		// Errors here indicate a programming mistake in the seed data.
		if err := catalog.Add(a); err != nil {
			panic("invalid MVP seed application: " + err.Error())
		}
	}
	return catalog
}

// mvpApplications is the seed set of known applications for the OmniInstall MVP.
var mvpApplications = []*app.Application{
	{
		ID:          "obs-studio",
		DisplayName: "OBS Studio",
		Summary:     "Video recording and live streaming software.",
		Categories:  []string{"video", "streaming"},
		Aliases:     []string{"obs", "open broadcaster software", "open broadcaster studio"},
		Homepage:    "https://obsproject.com",
		License:     "GPL-2.0-or-later",
		Publisher:   "OBS Project",
		Description: "Free and open source software for video recording and live streaming.",
	},
	{
		ID:          "vlc",
		DisplayName: "VLC",
		Summary:     "Versatile media player and streaming server.",
		Categories:  []string{"video", "audio"},
		Aliases:     []string{"vlc media player", "videolan"},
		Homepage:    "https://www.videolan.org/vlc/",
		License:     "GPL-2.0-or-later",
		Publisher:   "VideoLAN",
		Description: "VLC is a free and open source cross-platform multimedia player.",
	},
	{
		ID:          "firefox",
		DisplayName: "Firefox",
		Summary:     "Fast, private and secure web browser.",
		Categories:  []string{"browsers"},
		Aliases:     []string{"mozilla firefox", "mozilla"},
		Homepage:    "https://www.mozilla.org/firefox/",
		License:     "MPL-2.0",
		Publisher:   "Mozilla Foundation",
		Description: "Firefox is a free and open-source web browser.",
	},
	{
		ID:          "bitwarden",
		DisplayName: "Bitwarden",
		Summary:     "Open source password manager.",
		Categories:  []string{"security"},
		Aliases:     []string{"bitwarden password manager"},
		Homepage:    "https://bitwarden.com",
		License:     "GPL-3.0",
		Publisher:   "Bitwarden Inc.",
		Description: "Bitwarden is a free and open-source password management service.",
	},
	{
		ID:          "git",
		DisplayName: "Git",
		Summary:     "Distributed version control system.",
		Categories:  []string{"development"},
		Aliases:     []string{"git scm", "git version control"},
		Homepage:    "https://git-scm.com",
		License:     "GPL-2.0-only",
		Publisher:   "Software Freedom Conservancy",
		Description: "Git is a free and open source distributed version control system.",
	},
	{
		ID:          "docker",
		DisplayName: "Docker",
		Summary:     "Container platform for building and running applications.",
		Categories:  []string{"development", "system-tools"},
		Aliases:     []string{"docker engine", "docker ce"},
		Homepage:    "https://www.docker.com",
		License:     "Apache-2.0",
		Publisher:   "Docker, Inc.",
		Description: "Docker is an open platform for developing, shipping, and running applications.",
	},
	{
		ID:          "visual-studio-code",
		DisplayName: "Visual Studio Code",
		Summary:     "Code editor for building and debugging applications.",
		Categories:  []string{"development"},
		Aliases:     []string{"vscode", "vs code", "code"},
		Homepage:    "https://code.visualstudio.com",
		License:     "MIT",
		Publisher:   "Microsoft",
		Description: "Visual Studio Code is a lightweight but powerful source code editor.",
	},
	{
		ID:          "curl",
		DisplayName: "curl",
		Summary:     "Command-line tool for transferring data with URLs.",
		Categories:  []string{"network", "development"},
		Aliases:     []string{"curl transfer", "libcurl"},
		Homepage:    "https://curl.se",
		License:     "curl",
		Publisher:   "Daniel Stenberg",
		Description: "curl is a command line tool and library for transferring data with URLs.",
	},
	{
		ID:          "htop",
		DisplayName: "htop",
		Summary:     "Interactive process viewer and system monitor.",
		Categories:  []string{"system-tools"},
		Aliases:     []string{"htop process viewer", "htop monitor"},
		Homepage:    "https://htop.dev",
		License:     "GPL-2.0-or-later",
		Publisher:   "htop contributors",
		Description: "htop is a cross-platform interactive process viewer.",
	},
	{
		ID:          "neovim",
		DisplayName: "Neovim",
		Summary:     "Hyperextensible text editor based on Vim.",
		Categories:  []string{"development", "editors"},
		Aliases:     []string{"nvim", "neo-vim"},
		Homepage:    "https://neovim.io",
		License:     "Apache-2.0",
		Publisher:   "Neovim contributors",
		Description: "Neovim is a Vim-based text editor engineered for extensibility and usability.",
	},
}
