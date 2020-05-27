#!/bin/bash

set -e
set -x

CONTENT_FOLDER=content
OUT_FILE=content.go

# cleaning up
rm -rf $CONTENT_FOLDER/$OUT_FILE

# Install go-bindata
go get -u github.com/go-bindata/go-bindata/...

# convert icons into go code
cd $CONTENT_FOLDER
$GOPATH/bin/go-bindata -o $OUT_FILE -pkg $CONTENT_FOLDER images/...
cd ..

