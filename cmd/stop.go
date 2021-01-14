package cmd

import (
	"time"

	"github.com/DataDrake/cli-ng/cmd"
	"github.com/EbonJaeger/mcsmanager/config"
	"github.com/EbonJaeger/mcsmanager/tmux"
)

// Stop attempts to stop a Minecraft server.
var Stop = cmd.Sub{
	Name:  "stop",
	Alias: "t",
	Short: "Stop the Minecraft server",
	Run:   StopServer,
}

// StopServer stops the Minecraft server
func StopServer(root *cmd.Root, c *cmd.Sub) {
	prefix, err := root.Flags.(*GlobalFlags).GetPathPrefix()
	if err != nil {
		Log.Fatalf("Error getting the working directory: %s\n", err)
	}

	conf, err := config.Load(prefix)
	if err != nil {
		Log.Fatalf("Error loading server config: %s\n", err)
	}

	name := conf.MainSettings.ServerName

	// Check if the server is already stopped
	if !tmux.IsServerRunning(name) {
		Log.Warnln("The Minecraft server is already stopped!")
		return
	}

	Log.Infoln("Attempting to stop the server...")

	// Stop the server gracefully
	err = tmux.Exec("stop", name)

	// Wait 20 seconds for server to stop
	done := make(chan bool)
	go pollSessions(done, name)
	stopped := <-done

	if !stopped || err != nil {
		Log.Errorln("Could not stop the server normally! Attempting to force close...")
		tmux.KillWindow(name)
		Log.Warnln("Server window force-killed!")
		return
	}

	Log.Println("")
	Log.Goodln("Server stopped successfully!")
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
				if tickCount == 20 { // Stop polling after 20 seconds
					done <- false
					ticker.Stop()
					Log.Println("")
					return
				}

				Log.Printf("\rWaiting up to 20 seconds: %d", tickCount)
			}
		}
	}
}
