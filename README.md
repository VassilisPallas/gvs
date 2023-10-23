<p align="center">
  <img alt="GoReleaser Logo" src="./logo.png" height="140" />
  <h2 align="center">Golang Version Switcher</h2>
  <p align="center">
    <a href="https://github.com/VassilisPallas/gvs/actions?workflow=build"><img alt="Github Actions status" src="https://github.com/VassilisPallas/gvs/actions/workflows/actions.yml/badge.svg?branch=main&event=push"/></a>
    <a hre="https://opensource.org/licenses/MIT"><img alt="MIT License" src="https://img.shields.io/badge/License-MIT-blue.svg"/></a>
    <a hre="https://app.fossa.io/projects/git%2Bgithub.com%2FVassilisPallas%2Fgvs?ref=badge_shield"><img alt="FOSSA status" src="https://app.fossa.io/api/projects/git%2Bgithub.com%2FVassilisPallas%2Fgvs.svg?type=shield"/></a>
  </p>
</p>

## Table of Contents

- [Intro](#intro)
- [About](#about)
- [Installation](#installation)
- [Usage](#usage)
    - [Use the dropdown to select a version](#use-the-dropdown-to-select-a-version)
    - [See all versions including release candidates (rc)](#see-all-versions-including-release-candidates-rc)
    - [Install latest version](#install-latest-version)
    - [Refresh version list](#refresh-version-list)
    - [Help](#help)
- [Contributions](#contributions)
    - [Getting started with the project](#getting-started-with-the-project)
- [Licence](#licence)

## Intro

`gvs` allows you to quickly install and use different versions of Go via the command line. The installation is easy. Once installed, simply select the version you desire from the dropdown.

**Example:**
```sh
$ gvs
Use the arrow keys to navigate: ↓ ↑ → ←
? Select go version: 
  ▸ 1.21.3
    1.21.2
    1.21.1
    1.21.0
    1.20.10

✔ 1.21.3
Downloading...
Compare Checksums...
Unzipping...
Installing version...
1.21.3 version is installed!

$ go version
go version go1.21.3 darwin/arm64
```

## About

gvs is a version manager for go, designed to be installed per-user, and invoked per-shell. gvs works on any POSIX-compliant shell (sh, dash, ksh, zsh, bash), in particular on these platforms: unix and macOS.

> Windows will be supported in a later version.

## Installation

TBD

## Usage

**Before start using gvs, read the below:**

> gvs installs the `go` and `gofmt` binaries in `$HOME/bin/`. Make sure to append to your profile file: `export PATH=$PATH:$HOME/bin`, otherwise the terminal will not be able to find them.

### Use the dropdown to select a version

```sh
$ gvs
Use the arrow keys to navigate: ↓ ↑ → ←
? Select go version: 
    1.21.3
  ▸ 1.21.2
    1.21.1
    1.21.0
    1.20.10

✔ 1.21.2
Downloading...
Compare Checksums...
Unzipping...
Installing version...
1.21.2 version is installed!

$ go version
go version go1.21.2 darwin/arm64
```

1. Select the version you want to be installed by using the up and down arrows.
2. Hit **Enter** to select the desired version.

### See all versions including release candidates (rc)

To see a list with all versions, stable and unstable (release candidates) ones, just use the `--show-all` flag.

```sh
$ gvs --show-all
Use the arrow keys to navigate: ↓ ↑ → ←
? Select go version: 
  ▸ 1.21.3 (stable)
    1.21.2 (stable)
    1.21.1 (stable)
    1.21.0 (stable)
    1.21rc4 (unstable)
    1.21rc3 (unstable)
    1.21rc2 (unstable)
```

### Install latest version

In order to install the latest stable version, use the `--install-latest`.

```sh
$ gvs --install-latest
Downloading...
Compare Checksums...
Unzipping...
Installing version...
1.21.3 version is installed!
```

### Delete unused versions

Every time you install a new versions, gvs keeps the previous installed versions, so you can easily chnage between them. If you want to delete all the unused versions and keep only the current one, use the `--delete-unused` flag.

In the below example, the versions `1.20` and `1.19` are previously installed, and since they are not used (neither of them is the current version you use), they will be deleted.

```sh
$ gvs --delete-unused
Deleting go1.20.
go1.20 is deleted.
Deleting go1.19.
go1.19 is deleted.
All the unused version are deleted!
```

### Refresh version list

gvs caches the versions that are fetched from `https://go.dev/dl` in order to avoid overloading the server with requests.

The cache expires after a week, but if for any reason you'd like to force the fetch, you can use the `--refresh-versions` flag.

```sh
$ gvs --refresh-versions
Use the arrow keys to navigate: ↓ ↑ → ←
? Select go version: 
  ▸ 1.21.3 (stable)
    1.21.2 (stable)
    1.21.1 (stable)
    1.21.0 (stable)
    1.21rc4 (unstable)
    1.21rc3 (unstable)
    1.21rc2 (unstable)
```

> You can combine the flags `--refresh-versions` and `-show-all` to refresh the list and see all the versions.

### Help

For more help you can use the `--help` flag.

## Contributions

Contributions are very much welcomed! In order to contribute, the required version is `1.21.0`.
For more information regarding contributions, feel free to read the [contributing.md](./.github/contributing.md) file.

### Getting started with the project

After cloning the repository on you local machine, you can run the below commands:

```sh
# install the dependencies
make install-deps

# run the tests
make test

# build the project
make build
```

There a many Makefile recipes, feel free to run `make`. The default `make` command will print a short documentation for each recipe.

> Recipes like `format`, `lint`, `vet` and `test` also run in Github Actions, but it should be nice to run them yourself as well before committing the changes. If you think this needs to be a git hook, I'll be more than happy to include it.

## Licence

See [LICENSE.md](./LICENSE.md)

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FVassilisPallas%2Fgvs.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2FVassilisPallas%2Fgvs?ref=badge_large)