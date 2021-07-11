# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v1.2.0] - 2021-07-11

### Added

- Config option to set backup directory
  - May be a relative or absolute path
- Config option to exclude files/paths from being backed up
  - Backup directory is always excluded

### Fixed

- Excluded files while pruning not working
- Unable to start a server because the logs directory doesn't exist

### Changed

- Backup command now prints the progress of archiving server files
- Backup command now prints how long it took to archive the server files
- Compression of backup archives is now optional by specifying a compression level
- Init command arguments are now optional
  - This allows the tool to just save the config, and nothing else
- Use a progress bar to show the download progress

## [v1.1.1] - 2021-02-01

### Fixed

- Return wrong error when attaching to a tmux session

## [v1.1.0] - 2021-02-01

### Added

- Provider system to download server jars
  - Users can specify a direct URL, or a provider and version, e.g. `paper 1.16.5` to fetch the server from
- Global flag to set the path to the local server, e.g. `mcsmanager start -p /home/minecraft/server1`
- Show download speed when downloading a file

### Changed

- Increase wait time for a server to shut down to 20 seconds
- Make error messages more helpful

## [v1.0.0] - 2019-12-29

Initial release.

[unreleased]: https://github.com/EbonJaeger/mcsmanager/compare/v1.2.0...master
[v1.2.0]: https://github.com/EbonJaeger/mcsmanager/compare/v1.1.1...v1.2.0
[v1.1.1]: https://github.com/EbonJaeger/mcsmanager/compare/v1.1.0...v1.1.1
[v1.1.0]: https://github.com/EbonJaeger/mcsmanager/compare/v1.0.0...v1.1.0
[v1.0.0]: https://github.com/EbonJaeger/mcsmanager/compare/3d043fd...v1.0.0
