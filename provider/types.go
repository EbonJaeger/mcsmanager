package provider

const (
	// FileProvider is an update provider that simply downloads a given file.
	FileProvider = "FILE"

	// PaperProvider is an update provider that downloads the server jar from PaperMC.
	PaperProvider = "PAPER"
)

// Provider is an interface for a Minecraft server jar provider, such as PaperMC.
type Provider interface {
	Update(string) error
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
	Versions []string `json:"versions"`
}

// PaperBuilds is the representation of the Paper API response for a version.
type PaperBuilds struct {
	Builds struct {
		Latest string `json:"latest"`
	}
}