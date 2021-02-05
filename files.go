package mcsmanager

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/stretchr/stew/slice"
)

// Archive builds a tar archive of the given path.
//
// We use the  `filepath.Walk()` function to go through the entire
// file tree starting from the path that is passed in. This means we
// don't have to do a bunch of extra recursive logic for nested directories
// and having different code paths for files and directories.
func Archive(path string, w *tar.Writer, total int, excludes ...string) error {
	dir, err := os.Open(path)
	if err != nil {
		return err
	}
	defer dir.Close()

	current := 0

	// Walk the file tree from the server's root path
	return filepath.Walk(path, func(child string, info os.FileInfo, err error) error {
		// Don't try to add the root dir to the archive
		if path == child {
			return nil
		}

		// Don't archive files or directories that should be excluded
		for _, exclude := range excludes {
			if strings.Contains(child, exclude) {
				return nil
			}
		}

		// Write file header to archive
		name := strings.TrimPrefix(child, path+"/")
		header, err := tar.FileInfoHeader(info, name)
		if err != nil {
			return err
		}
		header.Name = name

		if err = w.WriteHeader(header); err != nil {
			return err
		}

		// Copy the item to the archive
		if !info.IsDir() {
			file, err := os.Open(child)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err = io.Copy(w, io.Reader(file)); err != nil {
				return err
			}

			// Print our progress
			current++
			fmt.Printf("\r%s", strings.Repeat(" ", 80))
			fmt.Printf("\rArchiving files... %d / %d", current, total)
		}

		return nil
	})
}

// CountFiles walks a directory tree and counts all files present.
func CountFiles(path string, excludes ...string) (count int, err error) {
	err = filepath.Walk(path, func(child string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return err
		}
		for _, e := range excludes {
			if strings.Contains(child, e) {
				return nil
			}
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})

	return
}

// PruneOld will delete files in the given directory if they were last modified
// after a certain period of time.
func PruneOld(path string, maxAge int, exemptFiles ...string) (total int, err error) {
	if maxAge == -1 { // -1 to disable age pruning
		return
	}

	// Max age is in days, convert it to hours
	maxAge = maxAge * 24

	// Check if we can access the path
	if _, err = os.Stat(path); err != nil {
		return
	}

	// Open the directory
	dir, err := os.Open(path)
	if err != nil {
		return
	}
	defer dir.Close()

	// Get all the files in the directory
	files, err := dir.Readdir(-1)
	if err != nil {
		return
	}

	// Iterate over the files
	for _, file := range files {
		// Check for exempt files
		if exemptFiles != nil {
			if slice.Contains(exemptFiles, file.Name()) {
				continue
			}
		}

		// Calculate time difference
		cur := time.Now()
		difference := cur.Sub(file.ModTime())
		// Remove file if needed
		if difference.Hours() > float64(maxAge) {
			if err = os.Remove(filepath.Join(path, file.Name())); err != nil {
				return
			}
			total++
		}
	}

	return
}

// Prune will remove the oldest files in a directory until
// the number of files in the directory is one under the limit.
func Prune(path string, maxFiles int, exemptFiles ...string) (total int, err error) {
	if maxFiles == -1 { // -1 to disable pruning
		return
	}

	// Check of the logs dir exists
	if _, err = os.Stat(path); err != nil {
		return
	}

	// Open the directory
	dir, err := os.Open(path)
	if err != nil {
		return
	}
	defer dir.Close()

	// Get the files in the directory
	files, err := dir.Readdir(-1)
	if err != nil {
		return
	}

	// Check for exempt files
	if exemptFiles != nil {
		for i, file := range files {
			if slice.Contains(exemptFiles, file.Name()) {
				// Remove exempt files from list of files that can be deleted
				files = append(files[:i], files[i+1:]...)
			}
		}
	}

	// Check if there are too many files
	numFiles := len(files)
	if numFiles < maxFiles {
		return
	}

	// Sort files so oldest is first
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})

	// Get files to remove
	toRemove := make([]os.FileInfo, 0)
	numToRemove := numFiles - maxFiles + 1 // We should be at the limit after pruning
	for i := 0; i < numToRemove; i++ {
		toRemove = append(toRemove, files[i])
	}

	// Remove each file
	for _, file := range toRemove {
		if err = os.Remove(filepath.Join(path, file.Name())); err != nil {
			return
		}
		total++
	}

	return
}
