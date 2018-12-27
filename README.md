```
akrctl is a command line interface for the Ankr Network distributed cloud computing network.

Usage:
  akrctl [command]

Available Commands:
  compute     compute commands

Flags:
  -u, --hub-url string        Override default Ankr Hub endpoint
  -h, --help                  help for akrctl

Use `akrctl [command] --help` for more information about a command.
```

## Installing `akrctl`

There are two ways to install `akrctl`:  
* Building a development version from source.
* Building it with Docker.

### Option 1 — Building the Development Version from Source

If you have a Go environment configured, you can install the development version of `akrctl` from the source.(below procedure tested in go version `go1.11.2 darwin/amd64`)

```
git clone -b feat/swdev-79-dccncli https://github.com/Ankr-network/dccn-cli.git $GOPATH/src/github.com/Ankr-network/dccncli
cd $GOPATH/src/github.com/Ankr-network/dccncli
dep ensure
go build -o akrcli cmd/dccncli/main.go
./akrctl any_akrctl_command
```

### Option 2 — Building with Docker

If you have Docker configured, you can build a Docker image using `akrcli`'s and run `akrctl` within a container. 
First, using `git clone` as in the `Option 1` to get the source code and then build docker image using the `Dockerfile.dep` file. (below procedure tested on Docker version 18.09.0): 

```
docker build -f Dockerfile.dep -t akrctl .
```

Then you can run it within a container: 

```
docker run --rm -p 50051:50051 akrctl any_akrctl_command
```

## Run with Docker Image on the ECR repository
If you are able to login AWS ECR you can run it within docker environment. (below procedure tested on Docker version 18.09.0)

* Command for login AWS ECR: 
```
eval $(aws ecr get-login --no-include-email --region us-west-2)
```
* Run akrctl with docker as following example:
```
docker run --rm -p 50051:50051 815280425737.dkr.ecr.us-west-2.amazonaws.com/dccn_ecr:akrctl any_akrctl_command
```

## Examples

`akrctl` is able to interact with all of your Ankr Network distributed cloud computing network resources. 
Below are a few common usage examples: 

* List all Tasks:
```
akrctl compute task list -u <addr_of_hub>
```
* Create a Task:
```
akrctl compute task create <taskname> --region <region> --zone <zone> -u <addr_of_hub>
```
* Delete a Task:
```
akrctl compute task delete <taskid> -f -u <addr_of_hub>
```

## Building and dependencies

`akrcli`'s dependencies are managed with [`dep`](https://github.com/golang/dep). 
To add dependencies, use [`dep ensure -add github.com/foo/bar`](https://github.com/golang/dep#adding-a-dependency)

* Initialize the dependency in vendor folder and create `Gopkg.toml` and `Gopkg.lock`:
```
dep init
```

* If any dependency like branch and version changed in `Gopkg.toml`, update the `Gopkg.lock` file related section and update the vendor folder:
```
dep ensure -update
```

* Checking the dependency:
```
dep status
```