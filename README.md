# mcsmanager

Fast and simple Minecraft server manager tool written in Go.

[![Report](https://goreportcard.com/badge/github.com/EbonJaeger/mcsmanager)](https://goreportcard.com/report/github.com/EbonJaeger/mcsmanager) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
--------

`mcsmanager` allows server administrators to easily start and manage Minecraft servers using simple commands and an easy-to-understand config.

## Dependencies

Here are the dependencies you will need to run this tool:
- `tmux`

If you want to build it for yourself, you will need:
- Go

## Installing

To download and install this using the Go tools, run `go get github.com/EbonJaeger/mcsmanager/cmd/mcsmanager`.

If you wish to build from source, run `go install ./cmd/mcsmanager` from the root of the project.

## Usage

You can view the full help by running `mcsmanager help`.

`mcsmanager CMD [args]`, where `CMD` is any one of:

- `attach|a` : Open the server console
- `backup|b` : Backup all server files into a .tar.gz archive
- `exec|e <args>` : Executes a command in the Minecraft server, e.g. `mcsmanager exec "say Hello there!"`. This can be used for automated messages before server restarts. :)
- `init|i <URL>` : Initialize the setup for a Minecraft server. The tool will download the server jar for you, so you don't have to.
- `start|s` : Start the Minecraft server
- `stop|t`  : Stop the Minecraft server
- `update|u <URL>` OR `<provider> <version>` : Update the jar file for the Minecraft server. Currently, Paper is the only provider supported.

## License

Copyright © 2019-2021 Evan Maddock (EbonJaeger)  

`mcsmanager` is available under the terms of the Apache-2.0 license
