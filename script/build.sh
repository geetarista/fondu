#!/usr/bin/env bash

PLATFORMS="darwin/386 darwin/amd64 freebsd/386 freebsd/amd64 freebsd/arm linux/386 linux/amd64 linux/arm windows/386 windows/amd64"

set -e

mkdir -p bin

for platform in $PLATFORMS; do
  CGO_ENABLED=0 GOOS=${platform%/*} GOARCH=${platform#*/}
  go build -o bin/$GOOS-$GOARCH
done
