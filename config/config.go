package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type config struct {
	MainSettings   mainSettings   `mapstructure:"main_settings"`
	JavaSettings   javaSettings   `mapstructure:"java_settings"`
	ServerSettings serverSettings `mapstructure:"server_settings"`
}

type mainSettings struct {
	ServerFile string `mapstructure:"server_file_name"`
}

type javaSettings struct {
	StartingMemory int32    `mapstructure:"starting_memory"`
	MaxMemory      int32    `mapstructure:"maximum_memory"`
	Flags          []string `mapstructure:"java_flags"`
}

type serverSettings struct {
	Flags []string `mapstructure:"jar_flags"`
}

// Conf holds all of the configuration settings
var Conf config

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error trying to initialize config: %s", err.Error())
	}

	configPath := filepath.Join(cwd, "config.toml")
	Conf = config{}
	_, err = toml.DecodeFile(configPath, &Conf)
	if err != nil {
		log.Fatalf("Error trying to decode config: %s", err.Error())
	}
}
