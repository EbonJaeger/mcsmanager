package cmd

import (
	"github.com/DataDrake/cli-ng/cmd"
	"github.com/EbonJaeger/mcsmanager/config"
	"github.com/EbonJaeger/mcsmanager/tmux"
)

// Exec sends a command to the Minecraft server
var Exec = cmd.CMD{
	Name:  "exec",
	Alias: "e",
	Short: "Executes a command in the Minecraft server",
	Args:  &ExecArgs{},
	Run:   Execute,
}

// ExecArgs contains the command arguments for the execute command
type ExecArgs struct {
	Command string `desc:"The full command to send to the server. Use quotes for multiple words."`
}

// Execute will run a command on a running Minecraft server
func Execute(root *cmd.RootCMD, c *cmd.CMD) {
	// Get the server name
	name := config.Conf.MainSettings.ServerName

	// Check if the server is running
	if !tmux.IsServerRunning(name) {
		log.Warnln("The Minecraft server is not running!")
		return
	}

	// Get the args
	args := c.Args.(*ExecArgs)

	// Send the command to the server
	err := tmux.Exec(args.Command, name)

	// Show any errors to the user
	if err != nil {
		log.Fatalf("Error while sending command: %s", err.Error())
	} else {
		log.Goodln("Command sent successfully!")
	}
}
