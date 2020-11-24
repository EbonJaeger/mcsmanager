package provider

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stretchr/stew/slice"
)

const paperVersionsURL = "https://papermc.io/api/v2/projects/paper"
const paperBuildsURL = "https://papermc.io/api/v2/projects/paper/versions/%s"
const paperDownloadURL = "https://papermc.io/api/v2/projects/paper/versions/%s/builds/%d/downloads/%s"

// getLatestBuild queries the Paper API to get the latest build number for the
// version of Minecraft we were given.
func (p Paper) getLatestBuild() (int, error) {
	url := fmt.Sprintf(paperBuildsURL, p.Version)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	builds := &PaperBuilds{}
	err = dec.Decode(builds)
	if err != nil {
		return 0, err
	}

	return builds.Builds[len(builds.Builds)-1], nil
}

// validateVersion queries the Paper API to see if we have a valid version string.
func (p Paper) validateVersion() (bool, error) {
	resp, err := http.Get(paperVersionsURL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	versions := &PaperVersions{}
	err = dec.Decode(versions)
	if err != nil {
		return false, err
	}

	return slice.Contains(versions.Versions, p.Version), nil
}

// Update gets the latest build of Paper from their website
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
	build, err := p.getLatestBuild()
	if err != nil {
		return err
	}

	download := fmt.Sprintf("paper-%s-%d.jar", p.Version, build)

	url := fmt.Sprintf(paperDownloadURL, p.Version, build, download)
	return DownloadFile(url, filepath)
}
