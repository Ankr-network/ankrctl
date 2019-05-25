# Introduction

```
ankrctl is a command line interface for the Ankr's distributed cloud computing network to provide access to most account and app commands.

Usage:
  ankrctl [command]

Available Commands:
  app         app commands
  chart       chart commands
  cluster     cluster commands
  help        Help about any command
  namespace   namespace commands
  user        user commands
  wallet      wallet commands

Flags:
  -u, --hub-url string        Override default Ankr Hub endpoint
  -h, --help                  help for ankrctl

Use `ankrctl [command] --help` for more information about a command.
```

# Prerequisites

You will need a local computer with ankrctl installed by following the project's [installation instructions](doc/install.md).

This reference is for the typical ankrctl's operations. 

# Invoking Commands

In ankrctl individual features are invoked by giving the utility a command, one or more sub-commands, and sometimes one or more options specifying particular values. Commands are grouped under three main categories:

* [user](doc/user.md) for user account operation and authentication
* [wallet](doc/wallet.md) for managing user's keys and tokens
* [app](doc/app.md) for managing user's application
* [cluster](doc/cluster.md) for managing user's cluster
* [chart](doc/chart.md) for managing user's chart
* [namespace](doc/namespace.md) for managing user's chart

To see an overview of all commands, you can invoke ankrctl by itself. To see all available commands under one of the three main categories, you can use ankrctl category, like ankrctl app. For a usage guide on a specific command, enter the command with the --help flag, i.e. ankrctl app --help.