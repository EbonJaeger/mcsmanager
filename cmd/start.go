package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/EbonJaeger/mcsmanager"
	"github.com/EbonJaeger/mcsmanager/config"
	"github.com/EbonJaeger/mcsmanager/tmux"
)

// Start attempts to start a Minecraft server.
var Start = cmd.Sub{
	Name:  "start",
	Alias: "s",
	Short: "Start the Minecraft server",
	Run:   StartServer,
}

// StartServer starts a Minecraft server.
func StartServer(root *cmd.Root, c *cmd.Sub) {
	Log.Infoln("Starting Minecraft server...")

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
	if tmux.IsServerRunning(name) {
		Log.Warnln("A server session is already running!")
		return
	}

	// Check if the Minecraft EULA has been accepted
	if !isEulaAccepted(prefix) {
		Log.Warnln("The Minecraft EULA has not been accepted!")
		Log.Warnln("The server will not start until the EULA has been accepted.")
		Log.Warnln("Open 'eula.txt' in a text editor, and change the line 'eula=false' to 'eula=true'.")
		return
	}

	// Remove old Logs
	LogsDir := filepath.Join(prefix, "logs")
	pruned, err := mcsmanager.RemoveOldFiles(LogsDir, conf.MainSettings.MaxAge, "latest.log")
	if err != nil {
		Log.Fatalf("Unable to remove old backups: %s\n", err.Error())
	}
	if pruned > 0 {
		Log.Infof("Removed %d log(s) due to age.\n", pruned)
	}

	// Remove too many Logs
	pruned, err = mcsmanager.RemoveTooManyFiles(LogsDir, conf.MainSettings.MaxLogs, "latest.log")
	if err != nil {
		Log.Fatalf("Unable to remove old backups: %s\n", err.Error())
	}
	if pruned > 0 {
		Log.Infof("Removed %d log(s) because over log limit.\n", pruned)
	}

	// Build the Java command to start the server
	javaCmd := buildJavaCmd(conf, prefix)

	// Create tmux window
	// TODO: out doesn't work as expected
	_, err = tmux.CreateSession(javaCmd, name)
	if err != nil {
		Log.Fatalln("Error creating tmux session!", err)
	}

	Log.Goodln("Server started!")
}

func buildJavaCmd(conf config.Root, prefix string) string {
	// Set the memory flags
	javaCmd := fmt.Sprintf("java -Xms%s -Xmx%s", conf.JavaSettings.StartingMemory, conf.JavaSettings.MaxMemory)

	// Add any JVM flags
	if len(*conf.JavaSettings.Flags) > 0 {
		javaCmd = javaCmd + " " + strings.Join(*conf.JavaSettings.Flags, " ")
	}

	// Set the jar file
	jarPath := filepath.Join(prefix, conf.MainSettings.ServerFile)
	javaCmd = javaCmd + fmt.Sprintf(" -jar %s", jarPath)

	// Add any jar flags
	if len(*conf.ServerSettings.Flags) > 0 {
		javaCmd = javaCmd + " " + strings.Join(*conf.ServerSettings.Flags, " ")
	}

	return javaCmd
}

func isEulaAccepted(prefix string) bool {
	// Check if the file exists
	eulaPath := filepath.Join(prefix, "eula.txt")
	_, err := os.Stat(eulaPath)
	if os.IsNotExist(err) { // EULA file doesn't exist yet
		// Let the server run to generate the eula.txt file
		Log.Warnln("EULA file doesn't exist yet! Starting server to generate it.")
		Log.Warnln("The server will not start until the EULA has been accepted.")
		Log.Warnln("Open 'eula.txt' in a text editor, and change the line 'eula=false' to 'eula=true'.")
		return true
	}

	// Open the file
	file, err := os.Open(eulaPath)
	if err != nil {
		Log.Fatalln("Unable to open EULA file:", err)
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
