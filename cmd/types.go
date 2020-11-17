package cmd

import (
	"github.com/DataDrake/cli-ng/cmd"
)

// DownloaderArgs contains the command arguments for commands that download
// a server jar file.
type DownloaderArgs struct {
	Args []string `desc:"URL to the server jar to download, or a provider and version, e.g. \"paper 1.16.4\""`
}

func (a DownloaderArgs) IsValid() bool {
	return len(a.Args) == 1 || len(a.Args) == 2
}

// PrintDownloaderUsage prints the proper usage to the user.
func PrintDownloaderUsage(cmd *cmd.CMD) {
	log.Errorln("Incorrect number of args!")
	log.Errorln("")
	log.Errorln("USAGE:")
	log.Errorf("\tmcsmanager %s <url>\n", cmd.Name)
	log.Errorln("OR")
	log.Errorf("\tmcsmanager %s <provider> <version>\n", cmd.Name)
	log.Errorln("")
	log.Errorln("PROVIDERS:")
	log.Errorln("\tpaper")
}
