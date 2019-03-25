# Installing `ankrctl`

There are three ways to install and run `ankrctl`:  
* Building it from Source.
* Building it with Docker.
* Run with Docker Image on the docker hub repository

## Option 1 — Building it from Source

If you have a Go environment configured, you can install the development version of `ankrctl` from the source.(below procedure tested in go version `go1.11.2 darwin/amd64`)

```
git clone https://github.com/Ankr-network/dccn-cli.git $GOPATH/src/github.com/Ankr-network/dccncli
cd $GOPATH/src/github.com/Ankr-network/dccncli
dep ensure --vendor-only
go build -o akrcli cmd/dccncli/main.go
./ankrctl <any_ankrctl_command>
```

## Option 2 — Building with Docker

If you have Docker configured, you can build a Docker image and run `ankrctl` within a container. 
First, get the source as in Option 1 and then build docker image: 

```
docker build -t ankrctl .
```

Then you can run it within a container: 

```
docker run -it ankrctl:latest
/ # ankrctl <any_ankrctl_command>
```

## Option 3 - Run with Docker Image on the docker hub repository
Below procedure tested on Docker version 18.09.0:

* Run ankrctl with docker as following example:
```
docker run -it ankrnetwork/ankrctl
/ # ankrctl <any_ankrctl_command>
```


## Building and dependencies

`akrcli`'s dependencies are managed with [`dep`](https://github.com/golang/dep). 
To add dependencies, use [`dep ensure -add github.com/foo/bar`](https://github.com/golang/dep#adding-a-dependency)

* Initialize the dependency in vendor folder and create `Gopkg.toml` and `Gopkg.lock`:
```
dep init
```

* If any dependency like branch and version changed in `Gopkg.toml`, update the `Gopkg.lock` and vendor:
```
dep ensure -update
```

* Checking the dependency:
```
dep status
```