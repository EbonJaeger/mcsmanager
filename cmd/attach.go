package cmd

import (
	"bufio"
	"os"

	"github.com/DataDrake/cli-ng/cmd"
	"github.com/EbonJaeger/mcsmanager/config"
	"github.com/EbonJaeger/mcsmanager/tmux"
)

// Attach opens the server console.
var Attach = cmd.Sub{
	Name:  "attach",
	Alias: "a",
	Short: "Open the server console",
	Args:  &AttachArgs{},
	Run:   AttachToSession,
}

// AttachArgs contains the command arguments for the stop command.
type AttachArgs struct{}

// AttachToSession attaches to the server console if the server is running.
func AttachToSession(root *cmd.Root, c *cmd.Sub) {
	prefix, err := root.Flags.(*GlobalFlags).GetPathPrefix()
	if err != nil {
		Log.Fatalf("Error getting the working directory: %s\n", err)
	}

	conf, err := config.Load(prefix)
	if err != nil {
		Log.Fatalf("Error loading server config: %s\n", err)
	}

	name := conf.MainSettings.ServerName

	// Check for already running server
	if !tmux.IsServerRunning(name) {
		Log.Warnln("Server is not currently running!")
		return
	}

	// Inform the user how the console works
	Log.Infoln("Attention!")
	Log.Infoln("To leave the console, press Ctrl+B then 'd'")
	Log.Warnln("Warning! Do not press Ctrl+C to exit! You will force-close your server!")
	Log.Println("")
	Log.Print("     Continue? [y/N] ")

	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		Log.Fatalln("Error while reading input:", err)
	}

	if char == 'y' || char == 'Y' {
		Log.Infoln("Opening server console...")
		if err := tmux.Attach(name); err != nil {
			Log.Fatalln("Unable to attach to session:", err)
		}

		Log.Goodln("Closed server console!")
	} else {
		Log.Goodln("Exiting!")
	}
}
