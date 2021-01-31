package main

import (
	log2 "log"
	"os"

	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/DataDrake/waterlog"
	"github.com/DataDrake/waterlog/format"
	"github.com/DataDrake/waterlog/level"
	commands "github.com/EbonJaeger/mcsmanager/cmd"
)

func main() {
	root := &cmd.Root{
		Name:  "mcsmanager",
		Short: "Minecraft Server Manager",
		Flags: &commands.GlobalFlags{},
	}

	// Initialize logging
	logger := waterlog.New(os.Stdout, "", log2.Ltime)
	logger.SetLevel(level.Info)
	logger.SetFormat(format.Min)
	commands.Log = logger

	// Initialize subcommands
	cmd.Register(&commands.Init)
	cmd.Register(&commands.Exec)
	cmd.Register(&commands.Start)
	cmd.Register(&commands.Stop)
	cmd.Register(&commands.Attach)
	cmd.Register(&commands.Backup)
	cmd.Register(&commands.Update)

	root.Run()
}
