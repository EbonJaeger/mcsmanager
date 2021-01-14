package provider

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer
// interface and we can pass this into io.TeeReader() which will report progress on each
// write cycle.
type WriteCounter struct {
	Current uint64
	Start   time.Time
	Total   uint64
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

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code not ok: %d", resp.StatusCode)
	}

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{
		Start: time.Now(),
		Total: uint64(resp.ContentLength),
	}
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
	wc.Current += uint64(n)
	wc.PrintProgress()
	return n, nil
}

// PrintProgress prints the progress of the current download to Stdout.
func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 80))

	// Calculate the current transfer rate
	rate := uint64(float64(wc.Current) / time.Since(wc.Start).Seconds())

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s / %s complete (%s/s)", humanize.Bytes(wc.Current), humanize.Bytes(wc.Total), humanize.Bytes(rate))
}
