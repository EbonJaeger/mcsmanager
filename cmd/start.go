package cmd

import (
	"github.com/DataDrake/cli-ng/cmd"
)

// Start attempts to start a Minecraft server
var Start = cmd.CMD{
	Name:  "start",
	Alias: "s",
	Short: "Start the Minecraft server",
	Args:  &StartArgs{},
	Run:   StartServer,
}

// StartArgs contains the command arguments for the start command
type StartArgs struct{}

// StartServer starts a Minecraft server
func StartServer(root *cmd.RootCMD, c *cmd.CMD) {
	log.Goodln("Server started!")
}
