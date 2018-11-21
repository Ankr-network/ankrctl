# dccncli [![Build Status](https://travis-ci.org/ankrnetwork/dccncli.svg?branch=master)](https://travis-ci.org/ankrnetwork/dccncli) [![GoDoc](https://godoc.org/github.com/Ankr-network/dccn-cli?status.svg)](https://godoc.org/github.com/Ankr-network/dccn-cli) [![Go Report Card](https://goreportcard.com/badge/github.com/Ankr-network/dccn-cli)](https://goreportcard.com/report/github.com/Ankr-network/dccn-cli)

```
dccncli is a command line interface for the AnkrNetwork API.

Usage:
  dccncli [command]

Available Commands:
  account     account commands
  auth        auth commands
  completion  completion commands
  compute     compute commands
  version     show the current version

Flags:
  -t, --access-token string   API V2 Access Token
  -u, --api-url string        Override default API V2 endpoint
  -c, --config string         config file (default is $HOME/.config/dccncli/config.yaml)
      --context string        authentication context name
  -h, --help                  help for dccncli
  -o, --output string         output format [text|json] (default "text")
      --trace                 trace api access
  -v, --verbose               verbose output

Use "dccncli [command] --help" for more information about a command.
```

## Installing `dccncli`

There are four ways to install `dccncli`: using a package manager, downloading a GitHub release, building a development version from source, or building it with [Docker](https://www.ankrnetwork.com/community/tutorials/the-docker-ecosystem-an-introduction-to-common-components).

### Option 1 – Using a Package Manager (Preferred)

A package manager allows you to install and keep up with new `dccncli` versions using only a few commands. Currently, `dccncli` is available as part of [Homebrew](https://brew.sh/) for macOS users and [Snap](https://snapcraft.io/) for GNU/Linux users.

You can use [Homebrew](https://brew.sh/) to install `dccncli` on macOS with this command:

```
brew install dccncli
```

You can use [Snap](https://snapcraft.io/) on [Snap-supported](https://snapcraft.io/docs/core/install) systems to install `dccncli` with this command:

```
sudo snap install dccncli
```
  #### Arch Linux
  Arch users not using snaps can install from the [AUR](https://aur.archlinux.org/packages/dccncli-bin/).

Support for Windows package managers is on the way.

### Option 2 — Downloading a Release from GitHub

Visit the [Releases page](https://github.com/Ankr-network/dccn-cli/releases) for the [`dccncli` GitHub project](https://github.com/Ankr-network/dccn-cli), and find the appropriate archive for your operating system and architecture.  You can download the archive from from your browser, or copy its URL and retrieve it to your home directory with `wget` or `curl`.

For example, with `wget`:

```
cd ~
wget https://github.com/Ankr-network/dccn-cli/releases/download/v1.11.0/dccncli-1.11.0-linux-amd64.tar.gz
```

Or with `curl`:

```
cd ~
curl -OL https://github.com/Ankr-network/dccn-cli/releases/download/v1.11.0/dccncli-1.11.0-linux-amd64.tar.gz
```

Extract the binary. On GNU/Linux or OS X systems, you can use `tar`.

```
tar xf ~/dccncli-1.11.0-linux-amd64.tar.gz
```

Or download and extract with this oneliner:
```
curl -sL https://github.com/Ankr-network/dccn-cli/releases/download/v1.11.0/dccncli-1.11.0-linux-amd64.tar.gz | tar -xzv
```

On Windows systems, you should be able to double-click the zip archive to extract the `dccncli` executable.

Move the `dccncli` binary to somewhere in your path. For example, on GNU/Linux and OS X systems:

```
sudo mv ~/dccncli /usr/local/bin
```

Windows users can follow [How to: Add Tool Locations to the PATH Environment Variable](https://msdn.microsoft.com/en-us/library/office/ee537574(v=office.14).aspx) in order to add `dccncli` to their `PATH`.

### Option 3 — Building the Development Version from Source

If you have a [Go environment](https://www.ankrnetwork.com/community/tutorials/how-to-install-go-1-6-on-ubuntu-16-04) configured, you can install the development version of `dccncli` from the command line.

```
go get -u github.com/Ankr-network/dccn-cli/cmd/dccncli
```

While the development version is a good way to take a peek at `dccncli`'s latest features before they get released, be aware that it may have bugs. Officially released versions will generally be more stable.

### Option 4 — Building with Docker

If you have [Docker](https://www.ankrnetwork.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04) configured, you can build a Docker image using `dccncli`'s [Dockerfile](https://github.com/Ankr-network/dccn-cli/blob/master/Dockerfile) and run `dccncli` within a container.

```
docker build -t dccncli .
```

Then you can run it within a container.

```
docker run --rm -e ANKRNETWORK_ACCESS_TOKEN="your_DO_token" dccncli any_dccncli_command
```

## Authenticating with AnkrNetwork

In order to use `dccncli`, you need to authenticate with AnkrNetwork by providing an access token, which can be created from the [Applications & API](https://cloud.ankrnetwork.com/settings/api/tokens) section of the Control Panel. You can learn how to generate a token by following the [AnkrNetwork API guide](https://www.ankrnetwork.com/community/tutorials/how-to-use-the-ankrnetwork-api-v2).

Docker users will have to use the `ANKRNETWORK_ACCESS_TOKEN` environmental variable to authenticate, as explained in the Installation section of this document.

If you're not using Docker to run `dccncli`, authenticate with the `auth init` command.

```
dccncli auth init
```

You will be prompted to enter the AnkrNetwork access token that you generated in the AnkrNetwork control panel.

```
AnkrNetwork access token: your_DO_token
```

After entering your token, you will receive confirmation that the credentials were accepted. If the token doesn't validate, make sure you copied and pasted it correctly.

```
Validating token: OK
```

This will create the necessary directory structure and configuration file to store your credentials.

### Logging in to multiple AnkrNetwork accounts

`dccncli` allows you to log in to multiple AnkrNetwork accounts at the same time and easily switch between them with the use of authentication contexts.

By default, a context named `default` is used. To create a new context, run `dccncli auth init --context new-context-name`. You may also pass the new context's name using the `ANKRNETWORK_CONTEXT` variable. You will be prompted for your API access token which will be associated with the new context.

To use a non-default context, pass the context name as described above to any `dccncli` command. To set a new default context, run `dccncli auth switch`. This command will save the current context to the config file and use it for all commands by default if a context is not specified.

The `--access-token` flag or `ANKRNETWORK_ACCESS_TOKEN` variable are acknowledged only if the `default` context is used. Otherwise, they will have no effect on what API access token is used. To temporarily override the access token if a different context is set as default, use `dccncli --context default --access-token your_DO_token ...`.

## Configuring Default Values

The `dccncli` configuration file is used to store your API Access Token as well as the defaults for command flags. If you find yourself using certain flags frequently, you can change their default values to avoid typing them every time. This can be useful when, for example, you want to change the username or port used for SSH.

On OS X and Linux, `dccncli`'s configuration file can be found at `${XDG_CONFIG_HOME}/dccncli/config.yaml` if the `${XDG_CONFIG_HOME}` environmental variable is set. Otherwise, the config will be written to `~/.config/dccncli/config.yaml`. For Windows users, the config will be available at `%LOCALAPPDATA%/dccncli/config/config.yaml`.

The configuration file was automatically created and populated with default properties when you authenticated with `dccncli` for the first time. The typical format for a property is `category.command.sub-command.flag: value`. For example, the property for the `force` flag with tag deletion is `tag.delete.force`.

To change the default SSH user used when connecting to a Task with `dccncli`, look for the `compute.ssh.ssh-user` property and change the value after the colon. In this example, we changed it to the username **sammy**.

```
. . .
compute.ssh.ssh-user: sammy
. . .
```

Save and close the file. The next time you use `dccncli`, the new default values you set will be in effect. In this example, that means that it will SSH as the **sammy** user (instead of the default **root** user) next time you log into a Task.

## Enabling Shell Auto-Completion

`dccncli` also has auto-completion support. It can be set up so that if you partially type a command and then press `TAB`, the rest of the command is automatically filled in. For example, if you type `dccncli comp<TAB><TAB> drop<TAB><TAB>` with auto-completion enabled, you'll see `dccncli compute task` appear on your command prompt.

**Note:** Shell auto-completion is not available for Windows users.

How you enable auto-completion depends on which operating system you're using. If you installed `dccncli` via Homebrew or Snap, auto-completion is activated automatically, though you may need to configure your local environment to enable it.

`dccncli` can generate an auto-completion script with the `dccncli completion your_shell_here` command. Valid arguments for the shell are Bash (`bash`) and ZSH (`zsh`). By default, the script will be printed to the command line output.  For more usage examples for the `completion` command, use `dccncli completion --help`.

### Linux

The most common way to use the `completion` command is by adding a line to your local profile configuration. At the end of your `~/.profile` file, add this line:

```
source <(dccncli completion your_shell_here)
```

Then refresh your profile.

```
source ~/.profile
```

### macOS

macOS users will have to install the `bash-completion` framework to use the auto-completion feature.

```
brew install bash-completion
```

After it's installed, load `bash_completion` by adding following line to your `.profile` or `.bashrc`/`.zshrc` file.

```
source $(brew --prefix)/etc/bash_completion
```


## Examples

`dccncli` is able to interact with all of your AnkrNetwork resources. Below are a few common usage examples. To learn more about the features available, see [the full tutorial on the AnkrNetwork community site](https://www.ankrnetwork.com/community/tutorials/how-to-use-dccncli-the-official-ankrnetwork-command-line-client).

* List all Tasks on your account:
```
dccncli compute task list
```
* Create a Task:
```
dccncli compute task create <name> --region <region-slug> --image <image-slug> --size <size-slug>
```
* Assign a Floating IP to a Task:
```
dccncli compute floating-ip-action assign <ip-addr> <task-id>
```
* Create a new A record for an existing domain:
```
dccncli compute domain records create --record-type A --record-name www --record-data <ip-addr> <domain-name>
```

`dccncli` also simplifies actions without an API endpoint. For instance, it allows you to SSH to your Task by name:
```
dccncli compute ssh <task-name>
```

By default, it assumes you are using the `root` user. If you want to SSH as a specific user, you can do that as well:
```
dccncli compute ssh <user>@<task-name>
```

## Building and dependencies

`dccncli`'s dependencies are managed with [`dep`](https://github.com/golang/dep). To add dependencies, use [`dep ensure -add github.com/foo/bar`](https://github.com/golang/dep#adding-a-dependency)

## More info

* [How To Use Dccncli, the Official AnkrNetwork Command-Line Client](https://www.ankrnetwork.com/community/tutorials/how-to-use-dccncli-the-official-ankrnetwork-command-line-client)
* [How To Work with AnkrNetwork Load Balancers Using Dccncli](https://www.ankrnetwork.com/community/tutorials/how-to-work-with-ankrnetwork-load-balancers-using-dccncli)
* [How To Secure Web Server Infrastructure With AnkrNetwork Cloud Firewalls Using Dccncli](https://www.ankrnetwork.com/community/tutorials/how-to-secure-web-server-infrastructure-with-ankrnetwork-cloud-firewalls-using-dccncli)
* [How To Work with AnkrNetwork Block Storage Using Dccncli](https://www.ankrnetwork.com/community/tutorials/how-to-work-with-ankrnetwork-block-storage-using-dccncli)
* [dccncli Releases](https://github.com/Ankr-network/dccn-cli/releases)
