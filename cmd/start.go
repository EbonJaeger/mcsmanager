package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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
	if tmux.IsSessionRunning() {
		log.Warnln("A server session is already running!")
		return
	}

	// Check if the Minecraft EULA has been accepted
	if !isEulaAccepted() {
		log.Warnln("The Minecraft EULA has not been accepted!")
		log.Warnln("The server will not start until the EULA has been accepted.")
		log.Warnln("Open 'eula.txt' in a text editor, and change the line 'eula=false' to 'eula=true'.")
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

func isEulaAccepted() bool {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Error trying to read EULA:", err)
	}

	// Check if the file exists
	eulaPath := filepath.Join(cwd, "eula.txt")
	_, err = os.Stat(eulaPath)
	if os.IsNotExist(err) { // EULA file doesn't exist yet
		// Let the server run to generate the eula.txt file
		log.Warnln("EULA file doesn't exist yet! Starting server to generate it.")
		log.Warnln("The server will not start until the EULA has been accepted.")
		log.Warnln("Open 'eula.txt' in a text editor, and change the line 'eula=false' to 'eula=true'.")
		return true
	}

	// Open the file
	file, err := os.Open(eulaPath)
	if err != nil {
		log.Fatalln("Unable to open EULA file:", err)
	}
	defer file.Close()

	// Read the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() { // For each line
		line := scanner.Text()
		if strings.HasPrefix(line, "eula") { // Line starts with eula
			value := strings.Split(line, "=")[1]
			if strings.ToLower(value) == "true" {
				return true
			}
			return false
		}
	}

	return false
}
