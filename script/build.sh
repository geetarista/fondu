#!/usr/bin/env bash

PLATFORMS="darwin/386 darwin/amd64 freebsd/386 freebsd/amd64 freebsd/arm linux/386 linux/amd64 linux/arm windows/386 windows/amd64"

set -e

for platform in $PLATFORMS; do
  goos=${platform%/*}
  goarch=${platform#*/}
  output="${goos}-${goarch}"
  echo "Building ${output}..."

  CGO_ENABLED=0 GOOS=$goos GOARCH=$goarch go build -o bin/$output
done
