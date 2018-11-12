#!/usr/bin/env bash

set -euo pipefail

TEMP_DIR=$(mktemp -d)
trap "{ rm -rf $TEMP_DIR; }" EXIT

mkdir -p "$TEMP_DIR/uli/linux"
mkdir -p "$TEMP_DIR/uli/osx"

echo "Building Linux binary"
env GOOS=linux GOARCH=amd64 go build -o "$TEMP_DIR/uli/linux/uli" uli.go
echo "Done."

echo "Building OS X binary"
env GOOS=darwin GOARCH=amd64 go build -o "$TEMP_DIR/uli/osx/uli" uli.go
echo "Done."

OUT_DIR="$(cd "$(dirname "$BASH_SOURCE")"; pwd -P)"
cd "$TEMP_DIR"
zip -r "$OUT_DIR/uli.zip" uli

