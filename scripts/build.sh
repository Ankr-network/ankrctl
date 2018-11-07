#!/bin/bash

set -eo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
OUT_DIR="$DIR/../out"
mkdir -p $OUT_DIR

go build \
  -o $OUT_DIR/doctl \
  -ldflags "-X github.com/Ankr-network/dccn-cli/Build=`git rev-parse --short HEAD`" \
  github.com/Ankr-network/dccn-cli/cmd/doctl

chmod +x $OUT_DIR/doctl
