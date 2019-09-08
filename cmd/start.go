package cmd

import (
	"fmt"
	"strings"

	"github.com/DataDrake/cli-ng/cmd"
	"github.com/EbonJaeger/mcsmanager/config"
	"github.com/EbonJaeger/mcsmanager/tmux"
)

// Start attempts to start a Minecraft server.
var Start = cmd.CMD{
	Name:  "start",
	Alias: "s",
	Short: "Start the Minecraft server",
	Args:  &StartArgs{},
	Run:   StartServer,
}

// StartArgs contains the command arguments for the start command.
type StartArgs struct{}

// StartServer starts a Minecraft server.
func StartServer(root *cmd.RootCMD, c *cmd.CMD) {
	log.Infoln("Starting Minecraft server...")

	// Check for already running server
	sessions, _ := tmux.ListSessions()
	if strings.Contains(sessions, tmux.SessionName) {
		log.Warnln("A server session is already running!")
		return
	}

	// Build the Java command to start the server
	javaCmd := buildJavaCmd()

	// Create tmux session stuff
	// TODO: out doesn't work as expected
	_, err := tmux.CreateSession(javaCmd)
	if err != nil {
		log.Fatalln("Error creating tmux session!", err)
	}

	log.Goodln("Server started!")
}

func buildJavaCmd() string {
	// Set the memory flags
	javaCmd := fmt.Sprintf("java -Xms%s -Xmx%s", config.Conf.JavaSettings.StartingMemory, config.Conf.JavaSettings.MaxMemory)

	// Add any JVM flags
	if len(config.Conf.JavaSettings.Flags) > 0 {
		javaCmd = javaCmd + " " + strings.Join(config.Conf.JavaSettings.Flags, " ")
	}

	// Set the jar file
	javaCmd = javaCmd + fmt.Sprintf(" -jar %s", config.Conf.MainSettings.ServerFile)

	// Add any jar flags
	if len(config.Conf.ServerSettings.Flags) > 0 {
		javaCmd = javaCmd + " " + strings.Join(config.Conf.ServerSettings.Flags, " ")
	}

	return javaCmd
}
