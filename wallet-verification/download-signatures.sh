#!/usr/bin/env bash

set -e -o pipefail

BASE_URL="https://downloads.skycoin.net/wallet/"

if [ -z $VERSION ]; then
    echo "VERSION must be set"
    exit 1
fi

mkdir -p "$VERSION"
pushd "$VERSION"

curl -o skycoin-${VERSION}-bin-linux-arm.tar.gz.asc ${BASE_URL}skycoin-${VERSION}-bin-linux-arm.tar.gz.asc
curl -o skycoin-${VERSION}-bin-linux-x64.tar.gz.asc ${BASE_URL}skycoin-${VERSION}-bin-linux-x64.tar.gz.asc
curl -o skycoin-${VERSION}-gui-linux-x64.AppImage.asc ${BASE_URL}skycoin-${VERSION}-gui-linux-x64.AppImage.asc
curl -o skycoin-${VERSION}-bin-win-x64.zip.asc ${BASE_URL}skycoin-${VERSION}-bin-win-x64.zip.asc
curl -o skycoin-${VERSION}-bin-win-x86.zip.asc ${BASE_URL}skycoin-${VERSION}-bin-win-x86.zip.asc
curl -o skycoin-${VERSION}-gui-win-setup.exe.asc ${BASE_URL}skycoin-${VERSION}-gui-win-setup.exe.asc
curl -o skycoin-${VERSION}-bin-osx-darwin-x64.zip.asc ${BASE_URL}skycoin-${VERSION}-bin-osx-darwin-x64.zip.asc
curl -o skycoin-${VERSION}-gui-osx-x64.zip.asc ${BASE_URL}skycoin-${VERSION}-gui-osx-x64.zip.asc
curl -o skycoin-${VERSION}-gui-osx.dmg.asc ${BASE_URL}skycoin-${VERSION}-gui-osx.dmg.asc
