package cmd

import (
	"os"

	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/DataDrake/waterlog"
)

// DownloaderArgs contains the command arguments for commands that download
// a server jar file.
type DownloaderArgs struct {
	Args []string `desc:"URL to the server jar to download, or a provider and version, e.g. \"paper 1.16.4\""`
}

// IsValid checks if the correct number of args were passed.
// For downloading, args must contain either just one arg with a URL,
// or two args with the provider name and version.
func (a DownloaderArgs) IsValid() bool {
	return len(a.Args) == 1 || len(a.Args) == 2
}

// PrintDownloaderUsage prints the proper usage to the user.
func PrintDownloaderUsage(sub *cmd.Sub) {
	Log.Errorln("Incorrect number of args!")
	Log.Errorln("")
	Log.Errorln("USAGE:")
	Log.Errorf("\tmcsmanager %s <url>\n", sub.Name)
	Log.Errorln("OR")
	Log.Errorf("\tmcsmanager %s <provider> <version>\n", sub.Name)
	Log.Errorln("")
	Log.Errorln("PROVIDERS:")
	Log.Errorln("\tpaper")
}

// GlobalFlags holds the flags for the root command.
type GlobalFlags struct {
	Path string `short:"p" long:"path" arg:"true" desc:"Set the path of the Minecraft server"`
}

// GetPathPrefix gets the server path from the command line flags. If there isn't one,
// return the current working directory.
func (f GlobalFlags) GetPathPrefix() (string, error) {
	prefix := f.Path
	if prefix == "" {
		var err error
		prefix, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	return prefix, nil
}

// Log is our logger via WaterLog.
// TODO: Make this not suck
var Log *waterlog.WaterLog
