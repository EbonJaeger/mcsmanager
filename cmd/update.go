package cmd

import (
	"os"
	"path/filepath"

	"github.com/DataDrake/cli-ng/cmd"
	"github.com/EbonJaeger/mcsmanager/config"
	"github.com/EbonJaeger/mcsmanager/tmux"
	"github.com/EbonJaeger/mcsmanager/util"
)

// Update downloads a server jar from the given URL.
var Update = cmd.CMD{
	Name:  "update",
	Alias: "u",
	Short: "Update the jar file for the Minecraft server",
	Args:  &UpdateArgs{},
	Run:   UpdateServer,
}

// UpdateArgs contains the command arguments for the update command
type UpdateArgs struct {
	Provider string `desc:"The server software we're trying to update"`
	Version  string `desc:"The Minecraft version to download"`
}

// UpdateServer downloads the specified server file.
func UpdateServer(root *cmd.RootCMD, c *cmd.CMD) {
	args := c.Args.(*UpdateArgs)

	// Check if the provider is supported
	if args.Provider != "paper" {
		log.Fatalln("Only Paper is supported by this command at this time. :(")
	}

	// Get the server name
	name := config.Conf.MainSettings.ServerName

	// Check if the server is running
	if tmux.IsServerRunning(name) {
		log.Warnln("The server is currently running! Please close it before updating.")
		return
	}

	// Get the current working directory
	log.Infoln("Downloading new server jar...")
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	// Download the specified server jar
	fileName := config.Conf.MainSettings.ServerFile
	outFile := filepath.Join(cwd, fileName)
	err = util.UpdatePaper(args.Version, outFile)
	if err != nil {
		log.Fatalln("Error downloading file:", err)
	}
	log.Goodln("Server jar updated!")
}
