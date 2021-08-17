package mcsmanager

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
)

// Archive builds a tar archive of the given path.
//
// We use the `filepath.Walk()` function to go through the entire
// file tree starting from the path that is passed in. This means we
// don't have to do a bunch of extra recursive logic for nested directories
// and having different code paths for files and directories.
func Archive(path string, w *tar.Writer, exclusions ...string) error {
	dir := os.DirFS(path)

	// Count the number of files to archive
	count, err := CountFiles(path, exclusions...)
	if err != nil {
		return fmt.Errorf("error counting files: %s", err)
	}

	bar := pb.New(count)
	bar.SetTemplate(pb.Simple)
	bar.Set(pb.CleanOnFinish, true)
	bar.SetWriter(os.Stdout)
	bar.SetMaxWidth(80)
	bar.Start()

	// Walk the file tree from the server's root path
	err = fs.WalkDir(dir, ".", func(child string, entry fs.DirEntry, err error) error {
		// Don't try to add the root dir to the archive
		if child == "." {
			return nil
		}

		// Don't archive files or directories that should be excluded
		if isExempt(child, exclusions...) {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return err
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
		if !entry.IsDir() {
			file, err := dir.Open(child)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err = io.Copy(w, file); err != nil {
				return err
			}

			// Print our progress
			bar.Increment()
		}

		return nil
	})

	bar.Finish()
	return err
}

// CountFiles walks a directory tree and counts all files present.
func CountFiles(path string, exclusions ...string) (count int, err error) {
	err = fs.WalkDir(os.DirFS(path), ".", func(child string, dir fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return err
		}
		if isExempt(child, exclusions...) {
			return nil
		}
		if !dir.IsDir() {
			count++
		}
		return nil
	})

	return
}

// PruneOld will delete files in the given directory if they were last modified
// after a certain period of time.
func PruneOld(path string, maxAge int, exemptions ...string) (total int, err error) {
	if maxAge == -1 { // -1 to disable age pruning
		return
	}

	// Max age is in days, convert it to hours
	maxAge = maxAge * 24

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

	cur := time.Now()

	// Iterate over the files
	for _, file := range files {
		if isExempt(file.Name(), exemptions...) {
			continue
		}

		// Calculate time difference
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
func Prune(path string, maxFiles int, exemptions ...string) (total int, err error) {
	if maxFiles == -1 { // -1 to disable pruning
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
	for i, file := range files {
		if isExempt(file.Name(), exemptions...) {
			files = append(files[:i], files[i+1:]...)
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

	toRemove := make([]os.FileInfo, 0)
	// We should be at the limit after pruning
	numToRemove := numFiles - maxFiles + 1
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

// isExempt checks if a given file is in the list
// of exempt files, returning `true` or `false`.
func isExempt(path string, exemptions ...string) bool {
	for _, exemption := range exemptions {
		if strings.Contains(path, exemption) {
			return true
		}
	}

	return false
}
