package cmd

import (
	"path/filepath"

	"github.com/DataDrake/cli-ng/cmd"
	"github.com/EbonJaeger/mcsmanager/config"
	"github.com/EbonJaeger/mcsmanager/provider"
	"github.com/EbonJaeger/mcsmanager/tmux"
)

// Update downloads a server jar from the given URL.
var Update = cmd.Sub{
	Name:  "update",
	Alias: "u",
	Short: "Update the jar file for the Minecraft server",
	Args:  &DownloaderArgs{},
	Run:   UpdateServer,
}

// UpdateServer downloads the specified server file.
func UpdateServer(root *cmd.Root, c *cmd.Sub) {
	if !c.Args.(*DownloaderArgs).IsValid() {
		PrintDownloaderUsage(c)
		return
	}
	args := c.Args.(*DownloaderArgs).Args

	prefix, err := root.Flags.(*GlobalFlags).GetPathPrefix()
	if err != nil {
		Log.Fatalf("Error getting the working directory: %s\n", err)
	}

	conf, err := config.Load(prefix)
	if err != nil {
		Log.Fatalf("Error loading server config: %s\n", err)
	}

	name := conf.MainSettings.ServerName

	// Check if the server is running
	if tmux.IsServerRunning(name) {
		Log.Warnln("The server is currently running! Please close it before updating.")
		return
	}

	fileName := conf.MainSettings.ServerFile
	outFile := filepath.Join(prefix, fileName)

	// Figure out our upgrade provider
	prov := provider.MatchProvider(args)
	if prov == nil {
		Log.Fatalf("Unable to get a download provider")
	}

	Log.Infoln("Downloading new server jar...")
	if err := prov.Download(outFile); err != nil {
		Log.Fatalln("Error downloading file:", err)
	} else {
		Log.Goodln("Server jar updated!")
	}
}
