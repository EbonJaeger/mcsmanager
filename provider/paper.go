package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/stretchr/stew/slice"
)

const (
	paperProjectEndpoint  = "https://papermc.io/api/v2/projects/paper"
	paperVersionsEndpoint = "https://papermc.io/api/v2/projects/paper/versions/%s"
	paperBuildEndpoint    = "https://papermc.io/api/v2/projects/paper/versions/%s/builds/%d"
	paperDownloadEndpoint = "https://papermc.io/api/v2/projects/paper/versions/%s/builds/%d/downloads/%s"
)

// ErrAErrAlreadyUpToDate is an error returned when the server is already
// at the latest build for the given version.
var ErrAlreadyUpToDate = errors.New("paper version is already at the latest build")

// Paper is an update provider that downloads a new Paper server version.
type Paper struct {
	Version string
}

// Download gets the latest build of Paper from their website
// for the given Minecraft version.
func (p Paper) Download(filepath string) error {
	// See if we actially have a valid version
	valid, err := p.validateVersion()
	if err != nil {
		return err
	}
	if !valid {
		return fmt.Errorf("server version not found: %s", p.Version)
	}

	// Get the latest build number
	b, err := getLatestBuild(p.Version)
	if err != nil {
		return err
	}

	// Check if we have the version we're currently running saved
	saved, err := Load(".paper_build.json")
	if err != nil {
		return fmt.Errorf("unable to read old version: %s", err.Error())
	}

	// Check if the current version and build matches the latest
	if saved.Version == b.Version {
		if saved.Build == b.Build {
			return ErrAlreadyUpToDate
		}
	}

	// Download the actual jar file
	url := fmt.Sprintf(paperDownloadEndpoint, p.Version, b.Build, b.Download.Application.Name)
	if err = DownloadFile(url, filepath); err != nil {
		return err
	}

	// Save the new version and build to disk
	if err = b.Save(".paper_build.json"); err != nil {
		return fmt.Errorf("unable to save version file: %s", err.Error())
	}

	// Verify the downloaded file
	return Verify(filepath, b.Download.Application.Hash)
}

// validateVersion queries the Paper API to see if we have a valid version string.
func (p Paper) validateVersion() (bool, error) {
	resp, err := http.Get(paperProjectEndpoint)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("failed to get version list: %d", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	versions := &PaperVersions{}
	err = dec.Decode(versions)
	if err != nil {
		return false, err
	}

	return slice.Contains(versions.Versions, p.Version), nil
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

// PaperBuild holds the API response data for a particular Paper build.
type PaperBuild struct {
	Build    int           `json:"build"`
	Download PaperDownload `json:"downloads"`
	Version  string        `json:"version"`
}

// PaperDownload contains information about a build's file.
type PaperDownload struct {
	Application PaperApplication `json:"application"`
}

// PaperApplication holds the name and hash of a file for a Paper build.
type PaperApplication struct {
	Name string `json:"name"`
	Hash string `json:"sha256"`
}

// Load reads saved version information from a file.
// If the file does not exist, this function returns
// an empty struct and no error.
func Load(path string) (*PaperBuild, error) {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &PaperBuild{}, nil
		}
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b := PaperBuild{}
	dec := json.NewDecoder(file)
	if err = dec.Decode(&b); err != nil {
		return nil, err
	}

	return &b, nil
}

// Save write a Paper build to a file on disk.
func (p PaperBuild) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(p)
}

// getLatestBuild queries the Paper API to get the latest build for the
// version of Minecraft we were given.
func getLatestBuild(version string) (*PaperBuild, error) {
	// Get the list of builds for the given version
	url := fmt.Sprintf(paperVersionsEndpoint, version)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get builds for version '%s': %d", version, resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	builds := &PaperBuilds{}
	err = dec.Decode(builds)
	if err != nil {
		return nil, err
	}

	build := builds.Builds[len(builds.Builds)-1]

	// Get the latest build info for this version
	url = fmt.Sprintf(paperBuildEndpoint, version, build)
	buildResp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer buildResp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get build info: %d", buildResp.StatusCode)
	}

	dec = json.NewDecoder(buildResp.Body)
	var b PaperBuild
	if err = dec.Decode(&b); err != nil {
		return nil, err
	}

	return &b, nil
}
