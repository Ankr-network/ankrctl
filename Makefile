.PHONE: build
export CGO=0

my_d=$(shell pwd)
OUT_D = $(shell echo $${OUT_D:-$(my_d)/builds})

GOOS = linux
UNAME_S := $(shell uname -s)
UNAME_P := $(shell uname -p)
ifeq ($(UNAME_S),Darwin)
  NATIVE = darwin
  GOARCH = 386
endif

GOARCH = amd64
ifneq ($(UNAME_P), x86_64)
  GOARCH = 386
endif

native: _build
native:
	@mv $(OUT_D)/dccncli_$(GOOS)_$(GOARCH) $(OUT_D)/dccncli
	@echo "built $(OUT_D)/dccncli"

_build:
	@mkdir -p builds
	@echo "building dccncli"
	@cd cmd/dccncli && env GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(OUT_D)/dccncli_$(GOOS)_$(GOARCH)
	@echo "built $(OUT_D)/dccncli_$(GOOS)_$(GOARCH)"


build_mac: GOOS = darwin
build_mac: GOARCH = 386
build_mac: _build

build_linux_386: GOARCH = 386
build_linux_386: _build

build_linux_amd64: GOARCH = amd64
build_linux_amd64: _build

clean:
	@rm -rf builds

_base_docker_cntr:
	docker build -f Dockerfile.build . -t dccncli_builder

docker_build: _base_docker_cntr
docker_build:
	@mkdir -p $(OUT_D)
	@docker build -f Dockerfile.cntr . -t dccncli_local
	@docker run --rm \
		-v $(OUT_D):/copy \
		-it --entrypoint /usr/bin/rsync \
		dccncli_local -av /app/ /copy/
	@docker run --rm \
		-v $(OUT_D):/copy \
		-it --entrypoint /bin/chown \
		alpine -R $(shell whoami | id -u): /copy
	@echo "Built binaries to $(OUT_D)"
	@echo "Created a local Docker container. To use, run: docker run --rm -it dccncli_local"
