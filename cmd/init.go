package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/DataDrake/cli-ng/cmd"
	"github.com/EbonJaeger/mcsmanager/config"
	"github.com/EbonJaeger/mcsmanager/provider"
)

// Init sets up everything required to start a Minecraft server
var Init = cmd.CMD{
	Name:  "init",
	Alias: "i",
	Short: "Initialize the setup for a Minecraft server",
	Args:  &DownloaderArgs{},
	Run:   InitServer,
}

// InitServer sets up the Minecraft server directory
func InitServer(root *cmd.RootCMD, c *cmd.CMD) {
	if !c.Args.(*DownloaderArgs).IsValid() {
		PrintDownloaderUsage(c)
		return
	}
	args := c.Args.(*DownloaderArgs).Args

	// Check if our dependencies are installed
	log.Infoln("Checking for installed dependencies...")
	missingDeps := make([]string, 0)
	if !isCommandAvailable("tmux") {
		missingDeps = append(missingDeps, "tmux")
	}
	// Notify the user
	if len(missingDeps) > 0 {
		log.Fatalln("Some dependencies are missing! Please install the following, and try again:", strings.Join(missingDeps, ", "))
	}
	log.Goodln("All dependencies are installed!")

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	fileName := config.Conf.MainSettings.ServerFile
	outFile := filepath.Join(cwd, fileName)

	// Figure out our upgrade provider
	prov := provider.MatchProvider(args)
	if prov == nil {
		log.Fatalf("Unable to get a download provider")
	}

	log.Infoln("Downloading new server jar...")
	if err = prov.Update(outFile); err != nil {
		log.Fatalln("Error downloading file:", err)
	} else {
		log.Goodln("Server jar updated!")
	}
}

func isCommandAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}
