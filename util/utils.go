package util

import (
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/stretchr/stew/slice"
)

// RemoveOldFiles will delete files in the given directory if they were last modified
// after a certain period of time.
func RemoveOldFiles(path string, maxAge int, exemptFiles ...string) (int, error) {
	if maxAge == -1 { // -1 to disable age pruning
		return 0, nil
	}

	maxAge = maxAge * 24 // Max age is in days, convert it to hours

	dir, err := os.Open(path)
	defer dir.Close()
	if err != nil {
		return 0, err
	}

	files, err := dir.Readdir(-1)
	if err != nil {
		return 0, err
	}

	prunedCount := 0
	for _, fi := range files {
		if exemptFiles != nil { // Check for exempt files
			if slice.Contains(exemptFiles, fi.Name()) {
				continue
			}
		}

		cur := time.Now()
		difference := cur.Sub(fi.ModTime())
		if difference.Hours() > float64(maxAge) { // Log is older than max age, delete it
			err = os.Remove(filepath.Join(path, fi.Name()))
			if err != nil {
				return prunedCount, err
			}
			prunedCount++
		}
	}

	return prunedCount, nil
}

// RemoveTooManyFiles will remove the oldest files in a directory until
// the number of files in the directory is one under the limit.
func RemoveTooManyFiles(path string, maxFiles int, exemptFiles ...string) (int, error) {
	if maxFiles == -1 { // -1 to disable pruning
		return 0, nil
	}

	dir, err := os.Open(path)
	defer dir.Close()
	if err != nil {
		return 0, err
	}

	files, err := dir.Readdir(-1)
	if err != nil {
		return 0, err
	}

	// Check for exempt files
	if exemptFiles != nil {
		for i, a := range files {
			if slice.Contains(exemptFiles, a.Name()) {
				// Remove exempt files from list of files that can be deleted
				files = append(files[:i], files[i+1:]...)
			}
		}
	}

	numFiles := len(files)
	pruned := 0
	if numFiles >= maxFiles {
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

		for _, fi := range toRemove {
			err = os.Remove(filepath.Join(path, fi.Name()))
			if err != nil {
				return pruned, err
			}
			pruned++
		}
	}

	return pruned, nil
}
