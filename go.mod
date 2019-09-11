module github.com/Ankr-network/ankrctl

go 1.12

replace github.com/tendermint/tendermint => github.com/Ankr-network/tendermint v0.31.5-0.20190719093344-1f8077fcd482

require (
	github.com/Ankr-network/dccn-common v0.0.0-20190729064917-c6a667db8f77
	github.com/Masterminds/semver v1.4.2 // indirect
	github.com/btcsuite/btcd v0.0.0-20190824003749-130ea5bddde3 // indirect
	github.com/cyphar/filepath-securejoin v0.2.2 // indirect
	github.com/fatih/color v1.7.0
	github.com/gobwas/glob v0.2.3
	github.com/golang/protobuf v1.3.2
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/mattn/go-isatty v0.0.9 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20190826022208-cac0b30c2563 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644
	github.com/spf13/cobra v0.0.5
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/viper v1.4.0
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/tendermint/go-amino v0.15.0 // indirect
	github.com/tendermint/tendermint v0.32.3 // indirect
	golang.org/x/crypto v0.0.0-20190909091759-094676da4a83
	google.golang.org/genproto v0.0.0-20180831171423-11092d34479b // indirect
	google.golang.org/grpc v1.23.0
	gopkg.in/yaml.v2 v2.2.2
	k8s.io/apimachinery v0.0.0-20190831074630-461753078381 // indirect
	k8s.io/helm v2.14.3+incompatible
)
