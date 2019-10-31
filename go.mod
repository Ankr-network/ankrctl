module github.com/Ankr-network/ankrctl

go 1.12

replace github.com/tendermint/tendermint => github.com/Ankr-network/tendermint v0.31.5-0.20191016091852-c60735e225bb

replace github.com/go-interpreter/wagon => github.com/Ankr-network/wagon v0.9.0-0.20191015152132-a57bd86fecb0

require (
	github.com/Ankr-network/ankr-chain v0.0.0-20191031013519-221574d16373
	github.com/Ankr-network/dccn-common v0.0.0-20191014090437-9fa44d3777fe
	github.com/Masterminds/semver v1.4.2 // indirect
	github.com/cyphar/filepath-securejoin v0.2.2 // indirect
	github.com/fatih/color v1.7.0
	github.com/gobwas/glob v0.2.3
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/grpc-gateway v1.11.4-0.20191029091745-69669120b0e0
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/mattn/go-isatty v0.0.9 // indirect
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644
	github.com/spf13/cobra v0.0.5
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	github.com/tendermint/tendermint v0.32.6
	golang.org/x/crypto v0.0.0-20190909091759-094676da4a83
	google.golang.org/appengine v1.4.0 // indirect
	google.golang.org/grpc v1.24.0
	gopkg.in/yaml.v2 v2.2.3
	k8s.io/apimachinery v0.0.0-20190831074630-461753078381 // indirect
	k8s.io/helm v2.14.3+incompatible
)
