package provider

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cheggaaa/pb/v3"
)

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(url string, filepath string) error {
	out, err := os.OpenFile(filepath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
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

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code not ok: %d", resp.StatusCode)
	}

	// Create our progress bar to report the download progress
	bar := pb.New64(resp.ContentLength)
	bar.Set(pb.SIBytesPrefix, true)
	bar.SetWriter(os.Stdout)
	bar.Start()

	reader := bar.NewProxyReader(resp.Body)

	// Copy the downloaded bytes to our out file
	if _, err = io.Copy(out, reader); err != nil {
		return err
	}

	bar.Finish()

	return nil
}

// Verify makes sure that the downloaded file's hash matches what the expected hash is.
// The hasing function used is `sha256`.
func Verify(path string, expected string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	sha256 := sha256.New()
	if _, err := io.Copy(sha256, file); err != nil {
		return err
	}

	if hash := hex.EncodeToString(sha256.Sum(nil)); hash != expected {
		return fmt.Errorf("hash mismatch: got %s, but expected %s", hash, expected)
	}

	return nil
}
