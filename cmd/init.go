package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/DataDrake/cli-ng/cmd"
	"github.com/EbonJaeger/mcsmanager/config"
	"github.com/EbonJaeger/mcsmanager/util"
)

// Init sets up everything required to start a Minecraft server
var Init = cmd.CMD{
	Name:  "init",
	Alias: "i",
	Short: "Initialize the setup for a Minecraft server",
	Args:  &InitArgs{},
	Run:   InitServer,
}

// InitArgs contains the command arguments for the init command
type InitArgs struct {
	URL string `desc:"Location of a server jar to download"`
}

// InitServer sets up the Minecraft server directory
func InitServer(root *cmd.RootCMD, c *cmd.CMD) {
	// Get the command args
	args := c.Args.(*InitArgs)
	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

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

	// Download the specified server jar
	log.Infoln("Downloading server jar...")
	fileName := config.Conf.MainSettings.ServerFile
	outFile := filepath.Join(cwd, fileName)
	if err = util.DownloadFile(args.URL, outFile); err != nil {
		log.Fatalln("Error downloading file:", err)
	}
	log.Goodln("Server jar downloaded!")
}

func isCommandAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}
