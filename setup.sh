#!/bin/sh

CWD=`pwd`
BINDATA_PKG="github.com/jteeuwen/go-bindata/..."

if [ ! -n "${GOPATH}" ]; then
    echo "GOPATH not set"
    exit 255
fi

echo " # Get and install go-bindata"
go get -u "$BINDATA_PKG"

echo " # Generate files"
exec $GOPATH/bin/go-bindata -pkg="storage" -o="${CWD}/storage/schemas.go" -prefix="${CWD}/storage" "${CWD}/storage/schemas/..."
