mcsmanager
--------

Fast and simple Minecraft server manager tool written in Go.

[![Report](https://goreportcard.com/badge/github.com/EbonJaeger/mcsmanager)](https://goreportcard.com/report/github.com/EbonJaeger/mcsmanager) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
--------

`mcsmanager` allows server administrators to easily start and manage Minecraft servers using simple commands and an easy-to-understand config.

## Dependencies
Here are the dependencies you will need to run this tool:
- `tmux`

If you want to build it for yourself, you will need:
- Golang
- make

## Build
To build the project, all you have to do is run `make`. Easy!

## Installing
If you're using a pre-compiled binary, simply place it where ever you want to use it, or put it into `/usr/bin`.

If you're building the project, you can use `sudo make install` to put it in `/usr/bin` for you.

## Usage
`mcsmanager CMD [args]`, where `CMD` is any one of:

- `attach|a` : Open the server console
- `backup|b` : Backup all server files into a .tar.gz archive
- `exec|e <args>` : Executes a command in the Minecraft server, e.g. `mcsmanager exec "say Hello there!"`. This can be used for automated messages before server restarts. :)
- `init|i <URL>` : Initialize the setup for a Minecraft server. The tool will download the server jar for you, so you don't have to.
- `start|s` : Start the Minecraft server
- `stop|st` : Stop the Minecraft server
- `update|u <provider> <version>` : Update the jar file for the Minecraft server. Currently only Paper is supported, since Mojang does not provide an API to easily get a particular version.

## License
Copyright © 2019-2020 Evan Maddock (EbonJaeger)  
Makefile.waterlog © Bryan Meyers (DataDrake) under the Apache 2.0 license, Makefile adapted from Bryan Meyers (DataDrake)

`mcsmanager` is available under the terms of the Apache-2.0 license
