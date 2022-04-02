#!/usr/bin/env bash

set -euo pipefail

. $(dirname $0)/commons.sh

SCRIPTS_DIR="$(cd -P "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_DIR="${SCRIPTS_DIR}/../dist/"
SOURCE_DIR="${SCRIPTS_DIR}/../"
NAME="azapi"
BUILD_ARTIFACT="terraform-provider-${NAME}_v${VERSION}"

OS_ARCH=("freebsd:amd64"
  "freebsd:386"
  "freebsd:arm"
  "freebsd:arm64"
  "windows:amd64"
  "windows:386"
  "linux:amd64"
  "linux:386"
  "linux:arm"
  "linux:arm64"
  "darwin:amd64"
  "darwin:arm64")


function clean() {
  info "Cleaning $BUILD_DIR"
  rm -rf "$BUILD_DIR"
  mkdir -p "$BUILD_DIR"
}

function release() {
  info "Clean build directory"
  clean

  info "Attempting to build ${BUILD_ARTIFACT}"

  cd "$SOURCE_DIR"
  go mod download
  for os_arch in "${OS_ARCH[@]}" ; do
    OS=${os_arch%%:*}
    ARCH=${os_arch#*:}
    EXT=$([ "$OS" == "windows" ] && echo ".exe" || echo "")
    info "GOOS: ${OS}, GOARCH: ${ARCH}"
    (
      env GOOS="${OS}" GOARCH="${ARCH}" CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X 'github.com/Azure/terraform-provider-azapi/version.ProviderVersion=v${VERSION}'" -o "${BUILD_ARTIFACT}_${OS}_${ARCH}${EXT}"
      mv "${BUILD_ARTIFACT}_${OS}_${ARCH}${EXT}" "${BUILD_DIR}"
    )
  done
  cd "${BUILD_DIR}"
  cp ../scripts/dearmor.sh ./
}

release
