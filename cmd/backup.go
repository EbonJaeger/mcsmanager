package cmd

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/DataDrake/cli-ng/cmd"
	"github.com/EbonJaeger/mcsmanager/config"
	"github.com/EbonJaeger/mcsmanager/tmux"
	"gopkg.in/djherbis/times.v1"
)

// Backup archives the Minecraft server files.
var Backup = cmd.CMD{
	Name:  "backup",
	Alias: "b",
	Short: "Backup all server files into a .tar.gz archive",
	Args:  &BackupArgs{},
	Run:   ArchiveServer,
}

// BackupArgs contains the command arguments for the backup command.
type BackupArgs struct{}

// ArchiveServer adds all directories and files of the server into a Gzip'd tar archive.
func ArchiveServer(root *cmd.RootCMD, c *cmd.CMD) {
	// Check if the server is currently running
	if tmux.IsSessionRunning() {
		log.Warnln("Please stop the server before trying to archive it!")
		return
	}

	log.Infoln("Archiving server files...")

	// Get our backup directory path
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error trying to initialize config: %s\n", err.Error())
	}
	backupDir := filepath.Join(cwd, "backups")

	// Check if the backup directory exists
	if _, err := os.Stat(backupDir); os.IsNotExist(err) { // Dir does not exist
		log.Infoln("Backup directory does not exist! Creating it...")
		err = os.Mkdir(backupDir, 0755)
		if err != nil {
			log.Fatalln("Unable to create backups directory:", err.Error())
		}
		log.Goodln("Backup directory created!")
	}

	// Check for backups that are too old
	pruned, err := checkOldBackups(backupDir)
	if err != nil {
		log.Fatalf("Unable to remove old backups: %s\n", err.Error())
	}
	if pruned > 0 {
		log.Infof("Removed %d archive(s) due to age.\n", pruned)
	}

	// Check for too many backups
	pruned, err = checkNumBackups(backupDir)
	if err != nil {
		log.Fatalf("Unable to remove old backups: %s\n", err.Error())
	}
	if pruned > 0 {
		log.Infof("Removed %d archive(s) because over backup limit.\n", pruned)
	}

	// Create archive file
	tarFile, err := createArchiveFile(backupDir)
	defer tarFile.Close()
	if err != nil {
		log.Fatalf("Error while adding files to archive: %s\n", err.Error())
	}

	// Create our file writers
	fileWriter := gzip.NewWriter(tarFile)
	defer fileWriter.Close()
	tarFileWriter := tar.NewWriter(fileWriter)
	defer tarFileWriter.Close()

	src, err := os.Open(".")
	defer src.Close()
	if err != nil {
		log.Fatalf("Error while adding files to archive: %s\n", err.Error())
	}

	// Archive the server directory recursively
	archiveDir(src, tarFileWriter)
	log.Goodln("Server file archive created!")
}

func createArchiveFile(dir string) (*os.File, error) {
	currentTime := time.Now()
	timeStr := currentTime.Format("2006-01-02T15:04:05-0700") // ISO-8601 format
	tarPath := filepath.Join(dir, timeStr+".tar.gz")

	return os.Create(tarPath)
}

func checkOldBackups(path string) (int, error) {
	maxAge := config.Conf.BackupSettings.MaxAge

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
		t, err := times.Stat(filepath.Join(path, fi.Name()))
		if err != nil {
			return 0, err
		}

		cur := time.Now()
		difference := cur.Sub(t.ModTime())
		if difference.Hours() > float64(maxAge) { // Archive is older than max age, delete it
			err = os.Remove(filepath.Join(path, fi.Name()))
			if err != nil {
				return prunedCount, err
			}
			prunedCount++
		}
	}

	return prunedCount, nil
}

func checkNumBackups(path string) (int, error) {
	dir, err := os.Open(path)
	defer dir.Close()
	if err != nil {
		return 0, err
	}

	files, err := dir.Readdir(-1)
	if err != nil {
		return 0, err
	}

	numArchives := len(files)
	maxArchives := config.Conf.BackupSettings.MaxBackups
	pruned := 0
	if numArchives > maxArchives {
		// Find oldest backups to remove
		toRemove := make([]os.FileInfo, 0)
		sort.Slice(files, func(i, j int) bool {
			return files[i].ModTime().Before(files[j].ModTime())
		})
		numToRemove := numArchives - maxArchives + 1 // We should be at the limit after the current backup
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

func archiveDir(dir *os.File, w *tar.Writer) {
	files, err := dir.Readdir(-1)
	if err != nil {
		log.Errorf("Error read files in directory: %s\n", err.Error())
	}

	// Iterate through all files
	for _, fileInfo := range files {
		// Exclude backups directory
		if fileInfo.Name() == "backups" {
			continue
		}

		if fileInfo.IsDir() { // File is actually a directory
			nestedDir, err := os.Open(dir.Name() + string(filepath.Separator) + fileInfo.Name())
			defer nestedDir.Close()
			if err != nil {
				log.Errorf("Error opening directory for archiving: %s\n", err.Error())
				continue
			}

			// Write directory header to archive
			header, _ := tar.FileInfoHeader(fileInfo, "")
			header.Name = dir.Name() + string(filepath.Separator) + fileInfo.Name()
			err = w.WriteHeader(header)
			if err != nil {
				log.Errorf("Error adding directory to archive: %s\n", err.Error())
			}

			// Recurse and archive everything else in this directory
			archiveDir(nestedDir, w)
		} else { // File is a file, archive it normally
			archiveFile(dir, fileInfo, w)
		}
	}
}

func archiveFile(dir *os.File, fi os.FileInfo, w *tar.Writer) {
	file, err := os.Open(dir.Name() + string(filepath.Separator) + fi.Name())
	defer file.Close()
	if err != nil {
		log.Errorf("Error writing file to archive: %s\n", err.Error())
	}

	// Create tar header
	header := new(tar.Header)
	header.Name = file.Name()
	header.Size = fi.Size()
	header.Mode = int64(fi.Mode())
	header.ModTime = fi.ModTime()

	err = w.WriteHeader(header)
	if err != nil {
		log.Errorf("Error writing file to archive: %s\n", err.Error())
	}

	_, err = io.Copy(w, io.Reader(file))
	if err != nil {
		log.Errorf("Error writing file to archive: %s\n", err.Error())
	}
}
