package cmd

import (
	"bufio"
	"os"
	"strings"

	"github.com/DataDrake/cli-ng/cmd"
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
	// Check for already running server
	sessions, _ := tmux.ListSessions()
	if !strings.Contains(sessions, tmux.SessionName) {
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
		if err := tmux.Attach(); err != nil {
			log.Fatalln("Unable to attach to session:", err)
		}

		log.Goodln("Closed server console!")
	} else {
		log.Goodln("Exiting!")
	}
}
