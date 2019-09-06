package cmd

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/DataDrake/cli-ng/cmd"
)

// Init sets up everything required to start a Minecraft server
var Init = cmd.CMD{
	Name:  "init",
	Alias: "i",
	Short: "Initialize the setup for a Minecraft server",
	Args:  &InitArgs{},
	Run:   InitServer,
}

// InitArgs contains the command arguments for the init command
type InitArgs struct{}

// InitServer sets up the Minecraft server directory
func InitServer(root *cmd.RootCMD, c *cmd.CMD) {
	log.Infoln("Setting up Minecraft server")

	// Check if our dependencies are installed
	log.Infoln("Checking for installed dependencies...")
	missingDeps := make([]string, 0)
	if !isCommandAvailable("tmux") {
		missingDeps = append(missingDeps, "tmux")
	}

	if len(missingDeps) > 0 {
		log.Fatalln("Some dependencies are missing! Please install the following, and try again:",
			strings.Join(missingDeps, ", "))
	}

	log.Goodln("All dependencies are installed!")

	// Create config file
	log.Infoln("Generating config file...")
	err := createConfigFile()
	if err != nil {
		log.Fatalln("Unable to generate config:", err)
	}
	log.Goodln("Configuration file generated")

	// Create server directory if it does not exist
	log.Infoln("Creating server files directory...")
	createServerDir()
}

func createConfigFile() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	configPath := filepath.Join(cwd, "config.toml")

	_, err = os.Stat(configPath)
	if !os.IsNotExist(err) { // Config file already exists
		return errors.New("Config file already exists")
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(configString)
	err = writer.Flush()
	return err
}

func createServerDir() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Error trying to set up server:", err.Error())
	}

	serverPath := filepath.Join(cwd, "serverfiles")
	_, err = os.Stat(serverPath)
	if os.IsNotExist(err) { // Check if the server directory exists. If not, attempt to create it
		err2 := os.MkdirAll(serverPath, 0755)
		if err2 != nil { // Error creating server directory
			log.Fatalln("Unable to create server directory:", err2.Error())
		}

		log.Goodln("Server directory created")
	} else { // A server directory already exists
		log.Fatalln("Directory already exists! Aborting.")
	}
}

func isCommandAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}

var configString = `
[main_settings]
server_file_name = "paperclip.jar"

[java_settings]
starting_memory = 2
maximum_memory = 2
java_flags = []

[server_settings]
jar_flats = [
    "nogui"
]
`
