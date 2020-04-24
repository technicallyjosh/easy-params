[![Release](https://img.shields.io/github/release/technicallyjosh/easy-params.svg)](https://github.com/technicallyjosh/easy-params/releases/latest)
[![Build](https://github.com/technicallyjosh/easy-params/workflows/Build/badge.svg)](https://github.com/technicallyjosh/easy-params/actions?query=workflow%3ABuild)
[![Release](https://github.com/technicallyjosh/easy-params/workflows/Release/badge.svg)](https://github.com/technicallyjosh/easy-params/actions?query=workflow%3ARelease)

# Easy SSM Parameters

Simple SSM Parameter Store interactions via a CLI.

## Why?

I've been using AWS SSM Parameter store for a while now for all my applications
and the management is a nightmare. This CLI aims to simplify the management of
parameters in AWS SSM.

## Installation

You can check out the [releases page](https://github.com/technicallyjosh/easy-params/releases).

### OSX

#### Homebrew

```console
$ brew tap technicallyjosh/easy-params
$ brew install easy-params
```

## CLI

```console
$ ezparams
An easy AWS Parameter Store CLI

Usage:
  ezparams [flags]
  ezparams [command]

Available Commands:
  diff        Shows the difference recursively between 2 paths.
  help        Help about any command
  ls          List parameters by path
  migrate     Migrate parameters by path
  put         Put parameter by path
  rm          Remove parameter(s) by path

Flags:
      --config string   config file (default is $HOME/.ezparams.yaml)
  -h, --help            help for ezparams
  -l, --local-time      convert UTC to local time (default true)
      --region string   AWS region to use
      --version         show version

Use "ezparams [command] --help" for more information about a command.
```

## Basic Usage

Right now it uses your default profile found in ~/.aws/credentials. You can change the region to use
on any command but for now it pulls from your [default] defined.

## Commands

### `ls`

Lists parameters in specified path.

```console
$ ezparams ls --help
List parameters by path

Usage:
  ezparams ls <path> [flags]

Flags:
  -d, --decrypt     decrypt "SecureString" values (default true)
  -h, --help        help for ls
  -p, --plain       plain text instead of table
  -r, --recursive   recursively get values based on path (default true)
  -v, --values      display values

Global Flags:
      --config string   config file (default is $HOME/.ezparams.yaml)
  -l, --local-time      convert UTC to local time (default true)
      --region string   AWS region to use
      --version         show version
```

### `put`

Put a parameter to the specified path.

```console
$ ezparams put --help
Put parameter by path

Usage:
  ezparams put <path> <value> [flags]

Flags:
  -h, --help          help for put
  -o, --overwrite     Overwrite param if exists.
  -t, --type string   Type of parameter. (default "SecureString")

Global Flags:
      --config string   config file (default is $HOME/.ezparams.yaml)
  -l, --local-time      convert UTC to local time (default true)
      --region string   AWS region to use
      --version         show version
```

### `rm`

Remove a parameter by path. You can also remove recursively by path with the --recursive flag.

```console
$ ezparams rm --help
Remove parameter(s) by path

Usage:
  ezparams rm <path(s)> [flags]

Flags:
  -h, --help        help for rm
      --recursive   remove all children on path recursively

Global Flags:
      --config string   config file (default is $HOME/.ezparams.yaml)
  -l, --local-time      convert UTC to local time (default true)
      --region string   AWS region to use
      --version         show version
```

### `migrate`

Migrate parameters from one path to another. Supports region to region. _This command will use the
same region if `region-to` is not specified._

```console
$ ezparams migrate --help
Migrate parameters by path

Usage:
  ezparams migrate <source path> [destination path] [flags]

Flags:
  -h, --help                 help for migrate
      --overwrite            overwrite destination params
  -f, --region-from string   the region to migrate from
  -t, --region-to string     the region to migrate to

Global Flags:
      --config string   config file (default is $HOME/.ezparams.yaml)
  -l, --local-time      convert UTC to local time (default true)
      --region string   AWS region to use
      --version         show version
```

### `diff`

```console
$ ezparams diff --help
Shows the difference recursively between 2 paths.

Usage:
  ezparams diff <path 1> <path 2> [flags]

Flags:
  -h, --help     help for diff
  -v, --values   show value diffs

Global Flags:
      --config string   config file (default is $HOME/.ezparams.yaml)
  -l, --local-time      convert UTC to local time (default true)
      --region string   AWS region to use
      --version         show version
```
