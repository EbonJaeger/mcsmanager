package config

import (
	"bufio"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type config struct {
	MainSettings   mainSettings   `toml:"main_settings"`
	JavaSettings   javaSettings   `toml:"java_settings"`
	ServerSettings serverSettings `toml:"server_settings"`
}

type mainSettings struct {
	ServerFile string `toml:"server_file_name"`
}

type javaSettings struct {
	StartingMemory int32    `toml:"starting_memory"`
	MaxMemory      int32    `toml:"maximum_memory"`
	Flags          []string `toml:"java_flags"`
}

type serverSettings struct {
	Flags []string `toml:"jar_flags"`
}

// Conf holds all of the configuration settings
var Conf config

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error trying to initialize config: %s", err.Error())
	}

	configPath := filepath.Join(cwd, "config.toml")
	err = createConfigFile(cwd)
	if err != nil {
		log.Fatalln("Error while creating config file:", err)
	}

	Conf = config{}
	_, err = toml.DecodeFile(configPath, &Conf)
	if err != nil {
		log.Fatalf("Error trying to decode config: %s", err.Error())
	}
}

func createConfigFile(cwd string) error {
	configPath := filepath.Join(cwd, "config.toml")

	_, err := os.Stat(configPath)
	if !os.IsNotExist(err) { // Config file already exists
		return nil
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

var configString = `
[main_settings]
server_file_name = "minecraft_server.jar"

[java_settings]
starting_memory = "2G"
maximum_memory = "2G"
java_flags = []

[server_settings]
jar_flags = [
    "nogui"
]
`
