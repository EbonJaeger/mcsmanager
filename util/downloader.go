package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/stretchr/stew/slice"
)

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer
// interface and we can pass this into io.TeeReader() which will report progress on each
// write cycle.
type WriteCounter struct {
	Total uint64
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

const paperVersionsURL = "https://papermc.io/api/v1/paper"
const paperBuildsURL = "https://papermc.io/api/v1/paper/%s"
const paperDownloadURL = "https://papermc.io/api/v1/paper/%s/%s/download"

// DownloadFromProvider downloads the latest file release of the given
// Minecraft version for the server software being used.
//
// Currently only Paper is supported.
func DownloadFromProvider(provider string, version string, filepath string) error {
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

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory. We pass an io.TeeReader
// into Copy() to report progress on the download.
func DownloadFile(url string, filepath string) error {

	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print("\n")

	err = os.Rename(filepath+".tmp", filepath)
	if err != nil {
		return err
	}

	err = os.Chmod(filepath, 0755)
	if err != nil {
		return err
	}

	return nil
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

// PrintProgress prints the progress of the current download to Stdout.
func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}
