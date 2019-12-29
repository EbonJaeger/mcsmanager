package cmd

import (
	"bufio"
	"os"

	"github.com/DataDrake/cli-ng/cmd"
	"github.com/EbonJaeger/mcsmanager/config"
	"github.com/EbonJaeger/mcsmanager/tmux"
)

// Attach opens the server console.
var Attach = cmd.CMD{
	Name:  "attach",
	Alias: "a",
	Short: "Open the server console",
	Args:  &AttachArgs{},
	Run:   AttachToSession,
}

// AttachArgs contains the command arguments for the stop command.
type AttachArgs struct{}

// AttachToSession attaches to the server console if the server is running.
func AttachToSession(root *cmd.RootCMD, c *cmd.CMD) {
	// Get the server name
	name := config.Conf.MainSettings.ServerName

	// Check for already running server
	if !tmux.IsServerRunning(name) {
		log.Warnln("Server is not currently running!")
		return
	}

	// Inform the user how the console works
	log.Infoln("Attention!")
	log.Infoln("To leave the console, press Ctrl+B then 'd'")
	log.Warnln("Warning! Do not press Ctrl+C to exit! You will force-close your server!")
	log.Println("")
	log.Print("     Continue? [y/N] ")

	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		log.Fatalln("Error while reading input:", err)
	}

	if char == 'y' || char == 'Y' {
		log.Infoln("Opening server console...")
		if err := tmux.Attach(name); err != nil {
			log.Fatalln("Unable to attach to session:", err)
		}

		log.Goodln("Closed server console!")
	} else {
		log.Goodln("Exiting!")
	}
}
