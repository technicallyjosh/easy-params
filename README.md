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
$ easy-params
An easy AWS Parameter Store CLI

Usage:
  easy-params [flags]
  easy-params [command]

Available Commands:
  diff        Shows the difference recursively between 2 paths.
  get         A brief description of your command
  help        Help about any command
  ls          List parameters by path
  migrate     Migrate parameters by path
  pull
  put         Put parameter by path
  rm          Remove parameter(s) by path

Flags:
      --config string   config file (default is $HOME/.easy-params.yaml)
  -h, --help            help for easy-params
      --load-config     load aws config from ~/.aws/config (default true)
  -l, --local-time      convert UTC to local time (default true)
      --region string   AWS region to use
      --version         show version

Use "easy-params [command] --help" for more information about a command.
```

## Basic Usage

Right now it uses your default profile found in ~/.aws/credentials. You can change the region to use
on any command but for now it pulls from your [default] defined. This will also load ~/.aws/config
by default.

## Commands

### `ls`

Lists parameters in specified path.

```console
$ easy-params ls --help
List parameters by path

Usage:
  easy-params ls <path> [flags]

Flags:
  -d, --decrypt     decrypt "SecureString" values (default true)
  -e, --env         output plain .env format
  -h, --help        help for ls
  -p, --plain       plain text instead of table
  -r, --recursive   recursively get values based on path (default true)
  -v, --values      display values

Global Flags:
      --config string   config file (default is $HOME/.easy-params.yaml)
      --load-config     load aws config from ~/.aws/config (default true)
  -l, --local-time      convert UTC to local time (default true)
      --region string   AWS region to use
      --version         show version
```

### `get`

Get a single parameter by path.

```console
Get parameter value by path

Usage:
  easy-params get <path> [flags]

Flags:
  -d, --decrypt   decrypt "SecureString" value (default true)
  -h, --help      help for get

Global Flags:
      --config string   config file (default is $HOME/.easy-params.yaml)
      --load-config     load aws config from ~/.aws/config (default true)
  -l, --local-time      convert UTC to local time (default true)
      --region string   AWS region to use
      --version         show version
```

### `put`

Put a parameter to the specified path.

```console
$ easy-params put --help
Put parameter by path

Usage:
  easy-params put <path> <value> [flags]

Flags:
  -c, --context string   context mode for setting many values.
  -h, --help             help for put
  -o, --overwrite        overwrite param if exists.
  -t, --type string      type of parameter. (default "SecureString")

Global Flags:
      --config string   config file (default is $HOME/.easy-params.yaml)
      --load-config     load aws config from ~/.aws/config (default true)
  -l, --local-time      convert UTC to local time (default true)
      --region string   AWS region to use
      --version         show version
```

### `rm`

Remove a parameter by path. You can also remove recursively by path with the --recursive flag.

```console
$ easy-params rm --help
Remove parameter(s) by path

Usage:
  easy-params rm <path(s)> [flags]

Flags:
  -h, --help        help for rm
      --recursive   remove all children on path recursively

Global Flags:
      --config string   config file (default is $HOME/.easy-params.yaml)
      --load-config     load aws config from ~/.aws/config (default true)
  -l, --local-time      convert UTC to local time (default true)
      --region string   AWS region to use
      --version         show version
```

### `migrate`

Migrate parameters from one path to another. Supports region to region. _This command will use the
same region if `region-to` is not specified._

```console
$ easy-params migrate --help
Migrate parameters by path

Usage:
  easy-params migrate <source path> [destination path] [flags]

Flags:
  -h, --help                 help for migrate
      --overwrite            overwrite destination params
  -f, --region-from string   the region to migrate from
  -t, --region-to string     the region to migrate to

Global Flags:
      --config string   config file (default is $HOME/.easy-params.yaml)
      --load-config     load aws config from ~/.aws/config (default true)
  -l, --local-time      convert UTC to local time (default true)
      --region string   AWS region to use
      --version         show version
```

### `diff`

Simple diff between 2 paths. Can also diff values.

```console
$ easy-params diff --help
Shows the difference recursively between 2 paths.

Usage:
  easy-params diff <path 1> <path 2> [flags]

Flags:
  -d, --decrypt           decrypt "SecureString" values (default true)
  -h, --help              help for diff
  -v, --values            show value diffs
  -w, --width-limit int   width limit of value output

Global Flags:
      --config string   config file (default is $HOME/.easy-params.yaml)
      --load-config     load aws config from ~/.aws/config (default true)
  -l, --local-time      convert UTC to local time (default true)
      --region string   AWS region to use
      --version         show version
```
