package cmd

import (
	log2 "log"
	"os"

	"github.com/DataDrake/cli-ng/cmd"
	"github.com/DataDrake/waterlog"
	"github.com/DataDrake/waterlog/format"
	"github.com/DataDrake/waterlog/level"
)

var log *waterlog.WaterLog

type GlobalFlags struct{}

// Root is the main command for the application
var Root *cmd.RootCMD

func init() {
	Root = &cmd.RootCMD{
		Name:  "mcsmanager",
		Short: "Minecraft Server Manager",
		Flags: &GlobalFlags{},
	}

	// Initialize subcommands
	Root.RegisterCMD(&Init)
	Root.RegisterCMD(&Start)

	// Initialize logging
	log = waterlog.New(os.Stdout, "", log2.Ltime)
	log.SetLevel(level.Info)
	log.SetFormat(format.Min)
}
