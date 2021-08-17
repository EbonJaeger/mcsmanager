package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/EbonJaeger/mcsmanager"
	"github.com/EbonJaeger/mcsmanager/config"
	"github.com/EbonJaeger/mcsmanager/tmux"
)

// BackupFlags holds the flags for the backup command.
type BackupFlags struct {
	Level int `short:"l" long:"level" desc:"Set the compression format to use; 0: no compression; 1: gzip"`
}

// Backup archives the Minecraft server files.
var Backup = cmd.Sub{
	Name:  "backup",
	Alias: "b",
	Short: "Backup all server files into a tar archive with optional compression",
	Flags: &BackupFlags{},
	Run:   ArchiveServer,
}

// ArchiveServer adds all directories and files of the server into a tar archive with optional compression.
func ArchiveServer(root *cmd.Root, c *cmd.Sub) {
	level := c.Flags.(*BackupFlags).Level
	if level < 0 || level > 1 {
		Log.Fatalln("Compression level must be between 0 and 1")
	}

	prefix, err := root.Flags.(*GlobalFlags).GetPathPrefix()
	if err != nil {
		Log.Fatalf("Error getting the working directory: %s\n", err)
	}

	conf, err := config.Load(prefix)
	if err != nil {
		Log.Fatalf("Error loading server config: %s\n", err)
	}

	// Check if the server is currently running
	if tmux.IsServerRunning(conf.MainSettings.ServerName) {
		Log.Warnln("Please stop the server before trying to archive it!")
		return
	}

	Log.Infoln("Archiving server files...")

	// Get our backup directory path
	var backupDir string
	if filepath.IsAbs(conf.BackupSettings.BackupDir) {
		backupDir = conf.BackupSettings.BackupDir
	} else {
		backupDir = filepath.Join(prefix, conf.BackupSettings.BackupDir)
	}

	// Check if the backup directory exists
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		Log.Infoln("Backup directory does not exist! Creating it...")
		if err = os.Mkdir(backupDir, 0755); err != nil {
			Log.Fatalln("Unable to create backups directory: %s\n", err)
		}
		Log.Goodln("Backup directory created!")
	}

	// Check for backups that are too old
	if pruned, err := mcsmanager.PruneOld(backupDir, conf.BackupSettings.MaxAge); err == nil {
		if pruned > 0 {
			Log.Infof("Removed %d archive(s) due to age.\n", pruned)
		}
	} else {
		Log.Fatalf("Unable to remove old backups: %s\n", err)
	}

	// Check for too many backups
	if pruned, err := mcsmanager.Prune(backupDir, conf.BackupSettings.MaxBackups); err == nil {
		if pruned > 0 {
			Log.Infof("Removed %d archive(s) because over backup limit.\n", pruned)
		}
	} else {
		Log.Fatalf("Unable to remove old backups: %s\n", err)
	}

	exclusions := append(*conf.BackupSettings.ExcludedPaths, conf.BackupSettings.BackupDir)

	// Create archive file
	tarFile, err := createArchive(backupDir, level)
	if err != nil {
		Log.Fatalf("Error while adding files to archive: %s\n", err)
	}
	defer tarFile.Close()

	// Create our file writers
	var w *tar.Writer
	switch level {
	case 0:
		w = tar.NewWriter(tarFile)
	case 1:
		compressor, err := gzip.NewWriterLevel(tarFile, gzip.BestCompression)
		if err != nil {
			Log.Fatalln("Failed to create archive compressor: %s\n", err)
		}
		defer compressor.Close()
		w = tar.NewWriter(compressor)
	}
	defer w.Close()

	// Add all of the server files to the archive
	start := time.Now()
	err = mcsmanager.Archive(prefix, w, exclusions...)
	diff := time.Since(start)

	if err != nil {
		Log.Errorf("Error adding files to archive: %s\n", err)
	}
	Log.Goodf("Server backup archive created in %v\n", diff)
}

func createArchive(dir string, level int) (*os.File, error) {
	currentTime := time.Now()
	timeStr := currentTime.Format("2006-01-02T15:04:05-0700") // ISO-8601 format

	var path string
	switch level {
	case 0: // No compression
		path = filepath.Join(dir, timeStr+".tar")
	case 1: // gzip compression
		path = filepath.Join(dir, timeStr+".tar.gz")
	default:
		return nil, fmt.Errorf("compression level not supported: %d", level)
	}

	return os.Create(path)
}
