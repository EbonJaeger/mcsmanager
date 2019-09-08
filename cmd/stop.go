package cmd

import (
	"github.com/DataDrake/cli-ng/cmd"
	"github.com/EbonJaeger/mcsmanager/tmux"
)

// Stop attempts to stop a Minecraft server.
var Stop = cmd.CMD{
	Name:  "stop",
	Alias: "st",
	Short: "Stop the Minecraft server",
	Args:  &StopArgs{},
	Run:   StopServer,
}

// StopArgs contains the command arguments for the stop command.
type StopArgs struct{}

// StopServer stops the Minecraft server
func StopServer(root *cmd.RootCMD, c *cmd.CMD) {
	log.Infoln("Attempting to stop the server...")

	// Stop the server gracefully
	stopCmd := tmux.Exec("stop")
	// TODO: out doesn't work as expected
	_, err := stopCmd.Output()
	if err != nil {
		log.Errorln("Could not stop the server normally! Attempting to force close...", err)
		/*killCmd := tmux.KillSession()
		killCmd.Run()
		log.Warnln("Server session force-killed!")*/
		return
	}

	log.Goodln("Server stopped successfully!")
}
