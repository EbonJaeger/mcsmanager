package cmd

import (
	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/EbonJaeger/mcsmanager/config"
	"github.com/EbonJaeger/mcsmanager/tmux"
)

// Exec sends a command to the Minecraft server
var Exec = cmd.Sub{
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
func Execute(root *cmd.Root, c *cmd.Sub) {
	prefix, err := root.Flags.(*GlobalFlags).GetPathPrefix()
	if err != nil {
		Log.Fatalf("Error getting the working directory: %s\n", err)
	}

	conf, err := config.Load(prefix)
	if err != nil {
		Log.Fatalf("Error loading server config: %s\n", err)
	}

	name := conf.MainSettings.ServerName // Get the server name

	// Check if the server is running
	if !tmux.IsServerRunning(name) {
		Log.Warnln("The Minecraft server is not running!")
		return
	}

	// Get the args
	args := c.Args.(*ExecArgs)

	// Send the command to the server
	err = tmux.Exec(args.Command, name)
	if err != nil {
		Log.Fatalf("Error while sending command: %s", err.Error())
	} else {
		Log.Goodln("Command sent successfully!")
	}
}
