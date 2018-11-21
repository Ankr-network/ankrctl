#!/bin/bash

set -o pipefail

ver=$1

if [[ -z "$ver" ]]; then
	echo "usage: $0 <version>"
	exit 1
fi

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
OUTPUT_DIR="${DIR}/../builds/${ver}/release"

for r in $(ls ${OUTPUT_DIR}/dccncli-${ver}-*); do
	name=$(basename $r)
	echo "uploading $name"
	github-release upload \
		--user ankrnetwork \
		--repo dccncli \
		--tag v${ver} \
		--name $name \
		--file $r
done

