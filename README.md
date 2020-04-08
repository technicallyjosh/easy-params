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
  help        Help about any command
  ls          List parameters by path
  migrate     Migrate parameters by path
  put         Put parameter by path
  rm          Remove parameter by path

Flags:
      --config string   config file (default is $HOME/.ezparams.yaml)
  -h, --help            help for ezparams
      --region string   AWS Region to use
  -l, --useLocalTime    convert UTC to local time (default true)
      --version         Show version

Use "ezparams [command] --help" for more information about a command.
```

## Basic Usage

Right now it uses your default profile found in ~/.aws/credentials. You can change the region to use
on any command but for now it pulls from your [default] defined.
