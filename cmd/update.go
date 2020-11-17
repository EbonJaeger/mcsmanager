package cmd

import (
	"os"
	"path/filepath"

	"github.com/DataDrake/cli-ng/cmd"
	"github.com/EbonJaeger/mcsmanager/config"
	"github.com/EbonJaeger/mcsmanager/provider"
	"github.com/EbonJaeger/mcsmanager/tmux"
)

// Update downloads a server jar from the given URL.
var Update = cmd.CMD{
	Name:  "update",
	Alias: "u",
	Short: "Update the jar file for the Minecraft server",
	Args:  &DownloaderArgs{},
	Run:   UpdateServer,
}

// UpdateServer downloads the specified server file.
func UpdateServer(root *cmd.RootCMD, c *cmd.CMD) {
	if !c.Args.(*DownloaderArgs).IsValid() {
		PrintDownloaderUsage(c)
		return
	}
	args := c.Args.(*DownloaderArgs).Args

	// Get the server name
	name := config.Conf.MainSettings.ServerName

	// Check if the server is running
	if tmux.IsServerRunning(name) {
		log.Warnln("The server is currently running! Please close it before updating.")
		return
	}

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
