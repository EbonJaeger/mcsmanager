package provider

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/EbonJaeger/mcsmanager/util"
	"github.com/stretchr/stew/slice"
)

const paperVersionsURL = "https://papermc.io/api/v1/paper"
const paperBuildsURL = "https://papermc.io/api/v1/paper/%s"
const paperDownloadURL = "https://papermc.io/api/v1/paper/%s/%s/download"

// getLatestBuild queries the Paper API to get the latest build number for the
// version of Minecraft we were given.
func (p Paper) getLatestBuild() (string, error) {
	url := fmt.Sprintf(paperBuildsURL, p.Version)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	builds := &PaperBuilds{}
	err = dec.Decode(builds)
	if err != nil {
		return "", err
	}

	return builds.Builds.Latest, nil
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
func (p Paper) Update(filepath string) error {
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

	url := fmt.Sprintf(paperDownloadURL, p.Version, build)
	return util.DownloadFile(url, filepath)
}
