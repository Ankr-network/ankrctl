```
akrctl is a command line interface for the Ankr Network distributed cloud computing network.

Usage:
  akrctl [command]

Available Commands:
  compute     compute commands

Flags:
  -u, --hub-url string        Override default Ankr Hub endpoint
  -h, --help                  help for dccncli

Use "akrctl [command] --help" for more information about a command.
```

## Installing `akrctl`

There are two ways to install `akrctl`:  building a development version from source, or building it with Docker.

### Option 1 — Building the Development Version from Source

If you have a Go environment configured, you can install the development version of `akrctl` from the command line.

```
git clone -b feat/swdev-79-dccncli https://github.com/Ankr-network/dccn-cli.git $GOPATH/src/github.com/Ankr-network/dccncli
cd $GOPATH/src/github.com/Ankr-network/dccncli
dep ensure
go build -o akrcli cmd/dccncli/main.go
```

### Option 2 — Building with Docker

If you have Docker configured, you can build a Docker image using `akrcli`'s and run `akrctl` within a container. first use the same step to get the source, then instead of go build the executable do the following docker build:

```
docker build -t akrctl .
```

Then you can run it within a container.

```
docker run --rm -p 50051:50051 akrctl -u ankr_hub_address any_akrctl_command
```

## Examples

`akrctl` is able to interact with all of your AnkrNetwork distributed cloud computing network resources. Below are a few common usage examples. 

* List all Tasks on your account:
```
akrctl compute task list
```
* Create a Task:
```
akrctl compute task create <name> --region <region-slug> --zone <zone-slug>
```

## Building and dependencies

`dccncli`'s dependencies are managed with [`dep`](https://github.com/golang/dep). To add dependencies, use [`dep ensure -add github.com/foo/bar`](https://github.com/golang/dep#adding-a-dependency)

* Initialize the dependency in vendor folder and create "Gopkg.toml" and "Gopkg.lock":
```
dep init
```

* If any dependency like branch and version changed in "Gopkg.toml", update the "Gopkg.lock" and vendor:
```
dep ensure -update
```

* Checking the dependency:
```
dep status
```