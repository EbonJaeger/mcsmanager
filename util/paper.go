package util

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stretchr/stew/slice"
)

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

const paperVersionsURL = "https://papermc.io/api/v1/paper"
const paperBuildsURL = "https://papermc.io/api/v1/paper/%s"
const paperDownloadURL = "https://papermc.io/api/v1/paper/%s/%s/download"

// UpdatePaper gets the latest build of Paper from their website
// for the given Minecraft version.
func UpdatePaper(version string, filepath string) error {
	// See if we actially have a valid version
	resp, err := http.Get(paperVersionsURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	versions := &PaperVersions{}
	err = dec.Decode(versions)
	if err != nil {
		return err
	}

	if !slice.Contains(versions.Versions, version) {
		return fmt.Errorf("server version not found: %s", version)
	}

	// Figure out the latest build
	url := fmt.Sprintf(paperBuildsURL, version)
	resp2, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp2.Body.Close()

	dec = json.NewDecoder(resp2.Body)
	builds := &PaperBuilds{}
	err = dec.Decode(builds)
	if err != nil {
		return err
	}

	// Make the download URL and get the file
	build := builds.Builds.Latest
	url = fmt.Sprintf(paperDownloadURL, version, build)
	return DownloadFile(url, filepath)
}
