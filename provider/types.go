package provider

import (
	"strings"
)

const (
	// FileProvider is an update provider that simply downloads a given file.
	FileProvider = "FILE"

	// PaperProvider is an update provider that downloads the server jar from PaperMC.
	PaperProvider = "PAPER"
)

// Provider is an interface for a Minecraft server jar provider, such as PaperMC.
type Provider interface {
	Download(string) error
}

// File is an update provider that downloads a new server version from a given URL.
type File struct {
	URL string
}

// Paper is an update provider that downloads a new Paper server version.
type Paper struct {
	Version string
}

// PaperVersions is the representation of all Paper versions returned by the API.
type PaperVersions struct {
	VersionGroups []string `json:"version_groups"`
	Versions      []string `json:"versions"`
}

// PaperBuilds is the representation of the Paper API response for a version.
type PaperBuilds struct {
	Builds []int `json:"builds"`
}

// MatchProvider creates and returns a provider for the given command arguments.
func MatchProvider(args []string) (prov Provider) {
	if len(args) == 1 {
		prov = File{URL: args[0]}
	} else if len(args) == 2 {
		providerType := strings.ToUpper(args[0])

		switch providerType {
		case PaperProvider:
			prov = Paper{Version: args[1]}
		default:
			prov = nil
		}
	}

	return
}
