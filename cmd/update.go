package cmd

import (
	"os"
	"path/filepath"
	"strings"

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
	Args:  &UpdateArgs{},
	Run:   UpdateServer,
}

// UpdateArgs contains the command arguments for the update command
type UpdateArgs struct {
	Args []string `desc:"URL to the server jar to download, or a provider and version, e.g. \"paper 1.16.4\""`
}

// UpdateServer downloads the specified server file.
func UpdateServer(root *cmd.RootCMD, c *cmd.CMD) {
	args := c.Args.(*UpdateArgs).Args
	if len(args) != 1 || len(args) != 2 {
		printUsage()
		return
	}

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
	var prov provider.Provider
	if len(args) == 1 {
		prov = provider.File{URL: args[0]}
	} else if len(args) == 2 {
		providerType := strings.ToUpper(args[0])
		switch providerType {
		case provider.PaperProvider:
			prov = provider.Paper{Version: args[1]}
		default:
			log.Fatalf("Unknown provider type: %s\n", providerType)
		}
	}

	log.Infoln("Downloading new server jar...")
	if err = prov.Update(outFile); err != nil {
		log.Fatalln("Error downloading file:", err)
	} else {
		log.Goodln("Server jar updated!")
	}
}

func printUsage() {
	log.Errorln("Incorrect number of args!")
	log.Errorln("")
	log.Errorln("USAGE:")
	log.Errorln("\tmcsmanager update <url>")
	log.Errorln("OR")
	log.Errorln("\tmcsmanager update <provider> <version>")
	log.Errorln("")
	log.Errorln("PROVIDERS:")
	log.Errorln("\tpaper")
}
