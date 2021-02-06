package mcsmanager

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCountFiles(t *testing.T) {
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

	result, err := CountFiles(dir)
	if err != nil {
		t.Fatalf("error counting files: %s\n", err)
	}

	// Check if the result is correct
	if result != len(files) {
		t.Fatalf("counted the wrong number of files: expected %d, actual: %d", len(files), result)
	}
}

func TestCountFilesNested(t *testing.T) {
	// Create temp dir to test in
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "nested"), 0755); err != nil {
		t.Fatalf("error creating nested test dir: %s\n", err)
	}

	files := []string{
		"file1.txt",
		"file2.png",
	}

	nestedFiles := []string{
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
		defer f.Close()
	}

	for _, file := range nestedFiles {
		path := filepath.Join(dir, "nested", file)
		f, err := os.Create(path)
		if err != nil {
			t.Fatalf("error creating test file: %s\n", err)
		}
		f.Close()
	}

	result, err := CountFiles(dir)
	if err != nil {
		t.Fatalf("error counting files: %s\n", err)
	}

	// Check if the result is correct
	if result != (len(files) + len(nestedFiles)) {
		t.Fatalf("counted the wrong number of files: expected %d, actual: %d", len(files)+len(nestedFiles), result)
	}
}

func TestCountFilesWithExclusions(t *testing.T) {
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

	result, err := CountFiles(dir, ".png")
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
