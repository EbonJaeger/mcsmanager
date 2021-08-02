package config

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// CreateFile creates a blank config file in the given path.
func CreateFile(prefix string) error {
	path := filepath.Join(prefix, "config.toml")

	if _, err := os.Stat(path); err != nil {
		// Server path doesn't exist. Create it first
		if os.IsNotExist(err) {
			if mkDirErr := os.Mkdir(prefix, 0755); err != nil {
				// Return if the mkdir failed for any other reason
				if !os.IsExist(mkDirErr) {
					return mkDirErr
				}
			}

			// Create the config file
			if _, createErr := os.Create(path); createErr != nil {
				return createErr
			}
		} else {
			return err
		}
	}

	return nil
}

// Default returns a new config with default settings.
func Default() Root {
	return Root{
		MainSettings: mainSettings{
			ServerFile: "minecraft_server.jar",
			ServerName: "Server 1",
			MaxLogs:    10,
			MaxAge:     7,
		},

		JavaSettings: javaSettings{
			StartingMemory: "2G",
			MaxMemory:      "2G",
			Flags:          &[]string{},
		},

		ServerSettings: serverSettings{
			Flags: &[]string{"nogui"},
		},

		BackupSettings: backupSettings{
			BackupDir:     "backups",
			ExcludedPaths: &[]string{},
			MaxBackups:    10,
			MaxAge:        7,
		},
	}
}

// Load reads a config from a config file in the given path.
func Load(prefix string) (conf Root, err error) {
	path := filepath.Join(prefix, "config.toml")

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	decoder := toml.NewDecoder(file)
	if _, err = decoder.Decode(&conf); err != nil {
		return
	}

	return
}

// Save writes the given config to the disk at the given path.
func (c Root) Save(prefix string) error {
	path := filepath.Join(prefix, "config.toml")

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	encoder := toml.NewEncoder(writer)

	if err := encoder.Encode(c); err != nil {
		return err
	}

	var buf bytes.Buffer
	if _, err := writer.Write(buf.Bytes()); err != nil {
		return err
	}

	return writer.Flush()
}
