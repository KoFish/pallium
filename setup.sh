#!/bin/sh

CWD=`realpath $(dirname $0)`
BINDATA_PKG="github.com/jteeuwen/go-bindata"

if [ ! -n "${GOPATH}" ]; then
    echo "GOPATH not set"
    exit 255
fi

echo " # Get go-bindata"
go get "$BINDATA_PKG"
echo " # Install go-bindata"
go install "$BINDATA_PKG"

echo " # Generate files"
exec $GOPATH/bin/go-bindata -pkg="storage" -o="${CWD}/storage/schemas.go" -prefix="${CWD}/storage" "${CWD}/storage/schemas/..."
