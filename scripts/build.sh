#!/usr/bin/env bash
set -Eeuo pipefail

BUILD_GOOS=$(go env GOOS)
if [ -n "${GOOS+x}" ] && [ -n "$GOOS" ]; then
  BUILD_GOOS=$GOOS
fi

BUILD_GOARCH=$(go env GOARCH)
if [ -n "${GOARCH+x}" ] && [ -n "$GOARCH" ]; then
  BUILD_GOARCH=$GOARCH
fi

SERVICE_NAME="bookie-api";
GIT_REF=$(git describe --always)
VERSION=commit-$GIT_REF

for directory in ./src/cmd/* ; do
  component=$(basename $directory)
  out=./build/$component

  echo "building $component"
  echo "  → service: $SERVICE_NAME"
  echo "  → version: $VERSION"
  echo "  → output: $out"

  CGO_ENABLED=0 go build -o $out -v \
    -ldflags "-X main.version=$VERSION -X main.serviceName=$SERVICE_NAME -X main.componentName=$component" \
    $directory
done