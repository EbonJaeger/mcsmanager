package provider

const (
	FILE_PROVIDER  = "FILE"
	PAPER_PROVIDER = "PAPER"
)

type Provider interface {
	Update(string) error
}

// File is an update provider that downloads a new server version from a given URL.
type File struct {
	Url string
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
