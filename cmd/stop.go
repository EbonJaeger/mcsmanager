package cmd

import (
	"time"

	"github.com/DataDrake/cli-ng/cmd"
	"github.com/EbonJaeger/mcsmanager/config"
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
	// Get the server name
	name := config.Conf.MainSettings.ServerName

	// Check if the server is already stopped
	if !tmux.IsServerRunning(name) {
		log.Warnln("The Minecraft server is already stopped!")
		return
	}

	log.Infoln("Attempting to stop the server...")

	// Stop the server gracefully
	err := tmux.Exec("stop", name)

	// Wait 10 seconds for server to stop
	done := make(chan bool)
	go pollSessions(done, name)
	stopped := <-done

	if !stopped || err != nil {
		log.Errorln("Could not stop the server normally! Attempting to force close...")
		tmux.KillWindow(name)
		log.Warnln("Server window force-killed!")
		return
	}

	log.Println("")
	log.Goodln("Server stopped successfully!")
}

func pollSessions(done chan bool, name string) {
	ticker := time.NewTicker(1 * time.Second)
	tickCount := 0
	for {
		select {
		case <-ticker.C: // Tick received
			tickCount++
			if !tmux.IsServerRunning(name) { // Session no longer running
				done <- true
			} else { // Session still running
				if tickCount == 10 { // Stop polling after 10 seconds
					done <- false
					ticker.Stop()
					log.Println("")
					return
				}

				log.Printf("\rWaiting up to 10 seconds: %d", tickCount)
			}
		}
	}
}
