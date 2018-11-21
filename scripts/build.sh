#!/bin/bash

set -eo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
OUT_DIR="$DIR/../out"
mkdir -p $OUT_DIR

go build \
  -o $OUT_DIR/dccncli \
  -ldflags "-X github.com/Ankr-network/dccn-cli/Build=`git rev-parse --short HEAD`" \
  github.com/Ankr-network/dccn-cli/cmd/dccncli

chmod +x $OUT_DIR/dccncli
