package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/EbonJaeger/mcsmanager/config"
	"github.com/EbonJaeger/mcsmanager/properties"
	"github.com/EbonJaeger/mcsmanager/tmux"
	"github.com/dustin/go-humanize"
)

// Status prints information about a Minecraft server.
var Status = cmd.Sub{
	Name:  "status",
	Alias: "n",
	Short: "View status info about a Minecraft server",
	Flags: &StatusFlags{},
	Run:   ServerStatus,
}

// StatusFlags holds the various flags that can be passed
// to the status command.
type StatusFlags struct {
	ShowAll      bool `short:"a" desc:"Print all extra server info"`
	ShowGameplay bool `short:"g" desc:"Print gameplay properties"`
	ShowRcon     bool `short:"r" desc:"Print Rcon configuration information"`
	ShowWorld    bool `short:"w" desc:"Print main world informations"`
}

// Color codes for terminal colors, matching colors from Waterlog.
const (
	red   = "\033[49;38;5;160m"
	green = "\033[49;38;5;040m"
	blue  = "\033[49;38;5;045m"
	reset = "\033[0m"
)

// ServerStatus handles the `Status` command and prints out various
// information about a Minecraft server.
func ServerStatus(root *cmd.Root, c *cmd.Sub) {
	prefix, err := root.Flags.(*GlobalFlags).GetPathPrefix()
	if err != nil {
		Log.Fatalf("Error getting the working directory: %s\n", err)
	}

	conf, err := config.Load(prefix)
	if err != nil {
		Log.Fatalf("Error loading server config: %s\n", err)
	}

	// Open and read the server.properties file
	path := filepath.Join(prefix, "server.properties")
	f, err := os.Open(path)
	if err != nil {
		Log.Fatalf("Error opening server.properties file: %s\n", err)
	}
	defer f.Close()

	raw, err := io.ReadAll(f)
	if err != nil {
		Log.Fatalf("Error reading server.properties file: %s\n", err)
	}

	// Read the server properties from the file
	props, err := properties.Read(raw)
	if err != nil {
		Log.Fatalf("Error reading props: %s\n", err)
	}

	print(conf.MainSettings.ServerName, conf.JavaSettings.MaxMemory, c.Flags.(*StatusFlags), props)
}

// print will write various server settings in a nice and readable
// format to stdout.
func print(name string, maxMemory string, flags *StatusFlags, props properties.Map) {
	var running string
	if tmux.IsServerRunning(name) {
		running = fmt.Sprintf("%sYES", green)
	} else {
		running = fmt.Sprintf("%sNO", red)
	}

	// Make the memory printout a bit nicer
	bytesDisplay := maxMemory
	if m, err := humanize.ParseBytes(bytesDisplay); err == nil {
		bytesDisplay = humanize.Bytes(m)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Print status header
	fmt.Fprintf(tw, "%s========== Status of '%s' ==========\n", blue, name)
	fmt.Fprintf(tw, "%sServer Address:\t%s%s\t%sServer Port:\t%s%s\n", blue, reset, props["server-ip"], blue, reset, props["server-port"])
	fmt.Fprintf(tw, "%sAllocated Memory:\t%s%s\t%sMax Players:\t%s%s\n", blue, reset, bytesDisplay, blue, reset, props["max-players"])
	fmt.Fprintf(tw, "%sRunning: %s\n", blue, running)

	// Print general gameplay settings
	if flags.ShowAll || flags.ShowGameplay {
		fmt.Fprintln(tw, "")

		fmt.Fprintf(tw, "%sGameplay Options:\n", blue)
		fmt.Fprintf(tw, "\t%sGamemode: \t%s%s \t%sDifficulty: \t%s%s\n", blue, reset, props["gamemode"], blue, reset, props["difficulty"])
		fmt.Fprintf(tw, "\t%sPVP Enabled: \t%s%s \t%sWhitelist Enabled: \t%s%s\n", blue, reset, props["pvp"], blue, reset, props["white-list"])

		fmt.Fprintln(tw, "")

		fmt.Fprintf(tw, "\t%sSpawning:\n", blue)
		fmt.Fprintf(tw, "\t\t%sAnimals: \t%s%s\n", blue, reset, props["spawn-animals"])
		fmt.Fprintf(tw, "\t\t%sNPCs: \t%s%s\n", blue, reset, props["spawn-npcs"])
		fmt.Fprintf(tw, "\t\t%sMonsters: \t%s%s\n", blue, reset, props["spawn-monsters"])
	}

	// Print Rcon settings
	if flags.ShowAll || flags.ShowRcon {
		fmt.Fprintln(tw, "")

		rconEnabled, _ := props["enable-rcon"].(bool)
		fmt.Fprintf(tw, "%sRcon:\n", blue)
		fmt.Fprintf(tw, "\t%sEnabled: \t%s%t\n", blue, reset, rconEnabled)
		if rconEnabled {
			fmt.Fprintf(tw, "\t%sRcon address: \t%s%s\n", blue, reset, props["server-ip"])
			fmt.Fprintf(tw, "\t%sRcon port: \t%s%s\n", blue, reset, props["rcon.port"])

			if props["rcon.password"] == "" {
				fmt.Fprintf(tw, "\t%sRcon is enabled, but no password is set!\n", red)
			}
		}
	}

	// Print main world settings
	if flags.ShowAll || flags.ShowWorld {
		fmt.Fprintln(tw, "")

		fmt.Fprintf(tw, "%sWorld:\n", blue)
		fmt.Fprintf(tw, "\t%sWorld Name: \t%s%s\n", blue, reset, props["level-name"])
		fmt.Fprintf(tw, "\t%sSeed: \t%s%s\n", blue, reset, props["level-seed"])
		fmt.Fprintf(tw, "\t%sType: \t%s%s\n", blue, reset, props["level-type"])
	}

	fmt.Fprintf(tw, "%s=============================================\n", blue)
	tw.Flush()
}
