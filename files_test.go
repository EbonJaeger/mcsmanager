package mcsmanager

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

var files = []string{
	"file1.txt",
	"file2.txt",
	filepath.Join("nested", "file3.txt"),
	filepath.Join("nested", "file4.txt"),
}

// setupTestDir creates a simple file tree in the given directory.
func setupTestDir(path string) error {
	if err := os.Mkdir(filepath.Join(path, "nested"), 0755); err != nil {
		return err
	}

	// Create our test files
	for _, file := range files {
		path := filepath.Join(path, file)
		err := os.WriteFile(path, []byte("test file contents"), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestArchiveDir(t *testing.T) {
	// Create temp dir to test in
	dir := t.TempDir()
	if err := setupTestDir(dir); err != nil {
		t.Fatalf("error creating test dir: %s\n", err)
	}

	// Create our tar writer
	archivePath := filepath.Join(dir, "archive.tar")
	tarFile, err := os.Create(archivePath)
	if err != nil {
		t.Fatalf("error creating archive file: %s\n", err)
	}
	w := tar.NewWriter(tarFile)

	// Archive the directory tree
	err = Archive(dir, w, "archive.tar")
	if err != nil {
		t.Fatalf("error writing file tree to archive: %s\n", err)
	}
	w.Close()
	tarFile.Close()

	// Extract the archive files to inspect them
	extractPath := filepath.Join(dir, "extracted")
	if err = os.Mkdir(extractPath, 0755); err != nil {
		t.Fatalf("error creating extraction dir: %s\n", err)
	}

	if err = extract(archivePath, extractPath); err != nil {
		t.Fatalf("error extracting archive: %s\n", err)
	}

	if err = verifyFiles(dir, extractPath); err != nil {
		t.Fatal(err)
	}
}

// extract extracts the entire file tree inside a tar archive to the given path.
func extract(archive, dir string) error {
	tarFile, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	extractPath, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	reader := tar.NewReader(tarFile)

	// Extract the file tree in the archive
	for {
		header, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		fi := header.FileInfo()
		path := filepath.Join(extractPath, header.Name)

		// If it's a directory, create it
		if fi.Mode().IsDir() {
			if err := os.MkdirAll(path, fi.Mode().Perm()); err != nil {
				return err
			}
			continue
		}

		// Create new file from archived file
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fi.Mode().Perm())
		if err != nil {
			return err
		}

		// Copy the contents of the file
		n, err := io.Copy(file, reader)
		if err != nil {
			return err
		}

		file.Close()
		if n != fi.Size() {
			return fmt.Errorf("error copying file '%s': wrote %d, but should have written %d", header.Name, n, fi.Size())
		}
	}

	return nil
}

// verifyFiles opens the original and extracted files and checks to see that they match.
func verifyFiles(basePath, extractedPath string) error {
	for _, fName := range files {
		expectedPath := filepath.Join(basePath, fName)
		actualPath := filepath.Join(extractedPath, fName)

		expectedFi, err := os.Stat(expectedPath)
		if err != nil {
			return err
		}
		actualFi, err := os.Stat(actualPath)
		if err != nil {
			return err
		}

		// Test the attributes
		if expectedFi.Size() != actualFi.Size() {
			return fmt.Errorf("file size does not match for '%s': expected %d, got %d", fName, expectedFi.Size(), actualFi.Size())
		}
		if expectedFi.Mode().Perm() != actualFi.Mode().Perm() {
			return fmt.Errorf("file perms do not match for '%s': expected %#o, got %#o", fName, expectedFi.Mode().Perm(), actualFi.Mode().Perm())
		}

		// Test the file contents
		expectedContents, err := os.ReadFile(expectedPath)
		if err != nil {
			return err
		}
		actualContents, err := os.ReadFile(actualPath)
		if err != nil {
			return err
		}
		if string(expectedContents) != string(actualContents) {
			return fmt.Errorf("file contents differ for file '%s': expected '%s', got '%s'", fName, string(expectedContents), string(actualContents))
		}
	}

	return nil
}

func TestCountFiles(t *testing.T) {
	// Create temp dir to test in
	dir := t.TempDir()
	if err := setupTestDir(dir); err != nil {
		t.Fatalf("error creating test dir: %s\n", err)
	}

	result, err := CountFiles(dir)
	if err != nil {
		t.Fatalf("error counting files: %s\n", err)
	}

	// Check if the result is correct
	if result != len(files) {
		t.Fatalf("counted the wrong number of files: expected %d, actual: %d", len(files), result)
	}
}

func TestCountFilesWithExclusions(t *testing.T) {
	// Create temp dir to test in
	dir := t.TempDir()
	if err := setupTestDir(dir); err != nil {
		t.Fatalf("error creating test dir: %s\n", err)
	}

	result, err := CountFiles(dir, "file2.txt")
	if err != nil {
		t.Fatalf("error counting files: %s\n", err)
	}

	// Check if the result is correct
	if result != len(files)-1 {
		t.Fatalf("counted the wrong number of files: expected %d, actual: %d", len(files)-1, result)
	}
}

func TestPruneFiles(t *testing.T) {
	// Create temp dir to test in
	dir := t.TempDir()

	files := []string{
		"file1.txt",
		"file2.png",
	}

	// Create our test files
	for _, file := range files {
		path := filepath.Join(dir, file)
		f, err := os.Create(path)
		if err != nil {
			t.Fatalf("error creating test file: %s\n", err)
		}
		f.Close()
	}

	result, err := Prune(dir, 2)
	if err != nil {
		t.Fatalf("error counting files: %s\n", err)
	}

	// Check if the result is correct
	if result != 1 {
		t.Fatalf("pruned wrong number of files: expected 1, pruned %d", result)
	}
}

func TestPruneFilesWithExemptions(t *testing.T) {
	// Create temp dir to test in
	dir := t.TempDir()

	files := []string{
		"file1.txt",
		"file2.png",
	}

	// Create our test files
	for _, file := range files {
		path := filepath.Join(dir, file)
		f, err := os.Create(path)
		if err != nil {
			t.Fatalf("error creating test file: %s\n", err)
		}
		f.Close()
	}

	result, err := Prune(dir, 2, ".txt")
	if err != nil {
		t.Fatalf("error counting files: %s\n", err)
	}

	// Check if the result is correct
	if result != 0 {
		t.Fatalf("pruned wrong number of files: expected 0, actual %d", result)
	}
}

func TestPruneOldFiles(t *testing.T) {
	// Create temp dir to test in
	dir := t.TempDir()

	// Create our test files
	path := filepath.Join(dir, "old.txt")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("error creating test file: %s\n", err)
	}
	f.Close()

	// Prune the files
	result, err := PruneOld(dir, 0)
	if err != nil {
		t.Fatalf("error counting files: %s\n", err)
	}

	// Check if the result is correct
	if result != 1 {
		t.Fatalf("pruned wrong number of files: expected 1, pruned %d", result)
	}
}

func TestPruneOldFilesWithExemptions(t *testing.T) {
	// Create temp dir to test in
	dir := t.TempDir()

	// Create our test files
	path := filepath.Join(dir, "old.txt")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("error creating test file: %s\n", err)
	}
	f.Close()

	// Prune the files
	result, err := PruneOld(dir, 0, ".txt")
	if err != nil {
		t.Fatalf("error counting files: %s\n", err)
	}

	// Check if the result is correct
	if result != 0 {
		t.Fatalf("pruned wrong number of files: expected 0, pruned %d", result)
	}
}
