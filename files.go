package mcsmanager

import (
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/stretchr/stew/slice"
)

// RemoveOldFiles will delete files in the given directory if they were last modified
// after a certain period of time.
func RemoveOldFiles(path string, maxAge int, exemptFiles ...string) (total int, err error) {
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

// RemoveTooManyFiles will remove the oldest files in a directory until
// the number of files in the directory is one under the limit.
func RemoveTooManyFiles(path string, maxFiles int, exemptFiles ...string) (total int, err error) {
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
