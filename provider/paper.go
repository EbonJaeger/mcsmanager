package provider

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
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

// getLatestBuild queries the Paper API to get the latest build number for the
// version of Minecraft we were given.
func (p Paper) getLatestBuild() (int, error) {
	url := fmt.Sprintf(paperVersionsEndpoint, p.Version)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("failed to get builds for version '%s': %d", p.Version, resp.StatusCode)
	}

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

	url := fmt.Sprintf(paperBuildEndpoint, p.Version, build)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to get build info: %d", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	var b PaperBuild
	if err = dec.Decode(&b); err != nil {
		return err
	}

	url = fmt.Sprintf(paperDownloadEndpoint, p.Version, build, b.Download.Application.Name)
	if err = DownloadFile(url, filepath); err != nil {
		return err
	}

	return verifyDownload(filepath, b.Download.Application.Hash)
}

// verifyDownload makes sure that the downloaded file's hash matches what the API says its hash should be.
func verifyDownload(path string, expectedHash string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	sha256 := sha256.New()
	if _, err := io.Copy(sha256, file); err != nil {
		return err
	}

	hash := hex.EncodeToString(sha256.Sum(nil))
	if hash != expectedHash {
		return fmt.Errorf("hash mismatch: got %s, but expected %s", hash, expectedHash)
	}

	return nil
}
