.PHONY: build_win build_mac build_lin win_env mac_env lin_env build

CLIENT_URL=client.dccn.ankr.com
TEST_DEV_URL=client-dev.dccn.ankr.com
TEST_STAGE_URL=client-stage.dccn.ankr.com
ANKR_CHAIN_URL = https://chain-01.dccn.ankr.com;https://chain-02.dccn.ankr.com;https://chain-03.dccn.ankr.com
ANKR_CHAIN_URL_STAGE = https://chain-stage-01.dccn.ankr.com;https://chain-stage-02.dccn.ankr.com
ANKR_CHAIN_URL_DEV = https://chain-stage-01.dccn.ankr.com;https://chain-stage-02.dccn.ankr.com

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


build_lin_dev: GOOS=linux
build_lin_dev: GOARCH=amd64
build_lin_dev: GOEXE=ankrctl_$(GOOS)_$(GOARCH)
build_lin_dev:
	@echo "Building linux executable dev"
	CGO_ENABLED=0 \
        GOOS=$(GOOS) \
        GOARCH=$(GOARCH) \
        go build -a \
        -installsuffix cgo \
        -ldflags="-w -s -X github.com/Ankr-network/dccn-cli/commands.clientURL=$(TEST_DEV_URL) -X github.com/Ankr-network/dccn-cli/commands.tendermintURL=$(ANKR_CHAIN_URL_DEV)" \
        -o build/$(GOEXE) \
        cmd/ankrctl/main.go


build_lin_stage: GOOS=linux
build_lin_stage: GOARCH=amd64
build_lin_stage: GOEXE=ankrctl_$(GOOS)_$(GOARCH)
build_lin_stage:
	@echo "Building linux executable stage"
	CGO_ENABLED=0 \
        GOOS=$(GOOS) \
        GOARCH=$(GOARCH) \
        go build -a \
        -installsuffix cgo \
        -ldflags="-w -s -X github.com/Ankr-network/dccn-cli/commands.clientURL=$(TEST_STAGE_URL) -X github.com/Ankr-network/dccn-cli/commands.tendermintURL=$(ANKR_CHAIN_URL_STAGE)" \
        -o build/$(GOEXE) \
        cmd/ankrctl/main.go


build_lin_prod: GOOS=linux
build_lin_prod: GOARCH=amd64
build_lin_prod: GOEXE=ankrctl_$(GOOS)_$(GOARCH)
build_lin_prod:
	@echo "Building linux executable prod"
	CGO_ENABLED=0 \
        GOOS=$(GOOS) \
        GOARCH=$(GOARCH) \
        go build -a \
        -installsuffix cgo \
        -ldflags="-w -s -X github.com/Ankr-network/dccn-cli/commands.clientURL=$(CLIENT_URL) -X github.com/Ankr-network/dccn-cli/commands.tendermintURL=$(ANKR_CHAIN_URL)" \
        -o build/$(GOEXE) \
        cmd/ankrctl/main.go

clean:
	@echo "Cleaning up all the builds"
	@rm -f build/*

build_all: clean build_win build_mac build_lin