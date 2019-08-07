.PHONY: build_win build_mac build_lin win_env mac_env lin_env build

CLIENT_URL=client.dccn.ankr.com

build = CGO_ENABLED=0 \
    GOOS=$(GOOS) \
    GOARCH=$(GOARCH) \
    go build -a \
    -installsuffix cgo \
    -ldflags="-w -s -X github.com/Ankr-network/dccn-cli/commands.clientURL=$(CLIENT_URL) -X github.com/Ankr-network/dccn-cli/commands.tendermintURL=$(ANKR_CHAIN_URL)" \
    -o build/$(GOEXE) \
    cmd/ankrctl/main.go

build_win: GOOS=windows
build_win: GOARCH=amd64
build_win: GOEXE=ankrctl_$(GOOS)_$(GOARCH).exe
build_win:
	@echo "Building win executable"
	@$(build)

build_mac: GOOS=darwin
build_mac: GOARCH=amd64
build_mac: GOEXE=ankrctl_$(GOOS)_$(GOARCH)
build_mac:
	@echo "Building mac executable"
	@$(build)

build_lin: GOOS=linux
build_lin: GOARCH=amd64
build_lin: GOEXE=ankrctl_$(GOOS)_$(GOARCH)
build_lin:
	@echo "Building linux executable"
	@$(build)

	"make build_lin"
build_lin: GOOS=linux
build_lin: GOARCH=amd64
build_lin: GOEXE=ankrctl_$(GOOS)_$(GOARCH)
build_lin:
	@echo "Building linux executable"
	CGO_ENABLED=0 \
        GOOS=$(GOOS) \
        GOARCH=$(GOARCH) \
        go build -a \
        -installsuffix cgo \
        -ldflags="-w -s -X github.com/Ankr-network/dccn-cli/commands.clientURL=$(URL_BRANCH) -X github.com/Ankr-network/dccn-cli/commands.tendermintURL=$(TENDERMINT_URL)" \
        -o build/$(GOEXE) \
        cmd/ankrctl/main.go


clean:
	@echo "Cleaning up all the builds"
	@rm -f build/*

build_all: clean build_win build_mac build_lin