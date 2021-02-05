package config

// Root is the root-level of our server configuration structure.
type Root struct {
	MainSettings   mainSettings   `toml:"main_settings"`
	JavaSettings   javaSettings   `toml:"java_settings"`
	ServerSettings serverSettings `toml:"server_settings"`
	BackupSettings backupSettings `toml:"backup_settings"`
}

type mainSettings struct {
	ServerFile string `toml:"server_file_name"`
	ServerName string `toml:"server_name"`
	MaxLogs    int    `toml:"max_log_count"`
	MaxAge     int    `toml:"max_log_age"`
}

type javaSettings struct {
	StartingMemory string    `toml:"starting_memory"`
	MaxMemory      string    `toml:"maximum_memory"`
	Flags          *[]string `toml:"java_flags"`
}

type serverSettings struct {
	Flags *[]string `toml:"jar_flags"`
}

type backupSettings struct {
	BackupDir     string    `toml:"backup_dir" comment:"Path can be an absolute or relative path"`
	ExcludedPaths *[]string `toml:"excluded_paths" comment:"Files that have any of these in their path will not be archived. The backup directory is always excluded"`
	MaxBackups    int       `toml:"max_number_backups"`
	MaxAge        int       `toml:"days_to_keep"`
}
