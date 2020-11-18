package cmd

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/DataDrake/cli-ng/cmd"
	"github.com/EbonJaeger/mcsmanager/provider"
	"github.com/EbonJaeger/mcsmanager/config"
)

// Init sets up everything required to start a Minecraft server
var Init = cmd.Sub{
	Name:  "init",
	Alias: "i",
	Short: "Initialize the setup for a Minecraft server",
	Args:  &DownloaderArgs{},
	Run:   InitServer,
}

// InitServer sets up the Minecraft server directory
func InitServer(root *cmd.Root, c *cmd.Sub) {
	if !c.Args.(*DownloaderArgs).IsValid() {
		PrintDownloaderUsage(c)
		return
	}
	args := c.Args.(*DownloaderArgs).Args
	prefix, err := root.Flags.(*GlobalFlags).GetPathPrefix()
	if err != nil {
		Log.Fatalf("Error getting the working directory: %s\n", err)
	}

	// Check if our dependencies are installed
	Log.Infoln("Checking for installed dependencies...")
	missingDeps := make([]string, 0)
	if !isCommandAvailable("tmux") {
		missingDeps = append(missingDeps, "tmux")
	}
	if len(missingDeps) > 0 {
		Log.Fatalln("Some dependencies are missing! Please install the following, and try again:", strings.Join(missingDeps, ", "))
	}
	Log.Goodln("All dependencies are installed!")

	// Create the server config
	Log.Infof("Creating server config at '%s'\n", filepath.Join(prefix, "config.toml"))
	if err := config.CreateFile(prefix); err != nil {
		Log.Fatalf("Error creating server config: %s\n", err)
	}

	conf := config.Default()
	if err := conf.Save(prefix); err != nil {
		Log.Fatalf("Error saving default config: %s\n", err)
	}

	// Download the server jar
	fileName := conf.MainSettings.ServerFile
	outFile := filepath.Join(prefix, fileName)

	// Figure out our upgrade provider
	prov := provider.MatchProvider(args)
	if prov == nil {
		Log.Fatalf("Unable to get a download provider")
	}

	Log.Infoln("Downloading new server jar...")
	if err := prov.Update(outFile); err != nil {
		Log.Fatalln("Error downloading file:", err)
	} else {
		Log.Goodln("Server jar updated!")
	}
}

func isCommandAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}
