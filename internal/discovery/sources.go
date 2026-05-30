package discovery

import "github.com/mitslab-de/omniinstall/internal/source"

// MVPSources returns the hardcoded source mappings for the 7 MVP applications.
// The map key is the canonical application ID.
func MVPSources() map[string][]source.Source {
	return map[string][]source.Source{
		"obs-studio": {
			{
				ApplicationID:    "obs-studio",
				SourceType:       source.TypeAPT,
				SourceIdentifier: "obs-studio",
				TrustLevel:       source.TrustOfficial,
				RiskLevel:        source.RiskLow,
			},
			{
				ApplicationID:    "obs-studio",
				SourceType:       source.TypeFlatpak,
				SourceIdentifier: "com.obsproject.Studio",
				TrustLevel:       source.TrustVerified,
				RiskLevel:        source.RiskLow,
			},
		},
		"vlc": {
			{
				ApplicationID:    "vlc",
				SourceType:       source.TypeAPT,
				SourceIdentifier: "vlc",
				TrustLevel:       source.TrustOfficial,
				RiskLevel:        source.RiskLow,
			},
			{
				ApplicationID:    "vlc",
				SourceType:       source.TypeFlatpak,
				SourceIdentifier: "org.videolan.VLC",
				TrustLevel:       source.TrustVerified,
				RiskLevel:        source.RiskLow,
			},
		},
		"firefox": {
			{
				ApplicationID:    "firefox",
				SourceType:       source.TypeAPT,
				SourceIdentifier: "firefox",
				TrustLevel:       source.TrustOfficial,
				RiskLevel:        source.RiskLow,
			},
			{
				ApplicationID:    "firefox",
				SourceType:       source.TypeFlatpak,
				SourceIdentifier: "org.mozilla.firefox",
				TrustLevel:       source.TrustVerified,
				RiskLevel:        source.RiskLow,
			},
		},
		"bitwarden": {
			{
				ApplicationID:    "bitwarden",
				SourceType:       source.TypeFlatpak,
				SourceIdentifier: "com.bitwarden.desktop",
				TrustLevel:       source.TrustVerified,
				RiskLevel:        source.RiskLow,
			},
		},
		"git": {
			{
				ApplicationID:    "git",
				SourceType:       source.TypeAPT,
				SourceIdentifier: "git",
				TrustLevel:       source.TrustOfficial,
				RiskLevel:        source.RiskLow,
			},
		},
		"docker": {
			{
				ApplicationID:    "docker",
				SourceType:       source.TypeVendor,
				SourceIdentifier: "docker-ce",
				TrustLevel:       source.TrustVerified,
				RiskLevel:        source.RiskMedium,
			},
		},
		"visual-studio-code": {
			{
				ApplicationID:    "visual-studio-code",
				SourceType:       source.TypeVendor,
				SourceIdentifier: "code",
				TrustLevel:       source.TrustVerified,
				RiskLevel:        source.RiskLow,
			},
			{
				ApplicationID:    "visual-studio-code",
				SourceType:       source.TypeFlatpak,
				SourceIdentifier: "com.visualstudio.code",
				TrustLevel:       source.TrustVerified,
				RiskLevel:        source.RiskLow,
			},
		},
	}
}
