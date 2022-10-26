#! /bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

build() {
    EXT=""
    [[ $GOOS = "windows" ]] && EXT=".exe"
    echo "Building ${GOOS} ${GOARCH}"
    go build \
        -trimpath \
        -o ./bin/pg-${GOOS}-${GOARCH}${EXT} .
}

### multi arch binary build
GOOS=linux GOARCH=arm64 build
GOOS=linux GOARCH=amd64 build

GOOS=darwin GOARCH=arm64 build

GOOS=windows GOARCH=amd64 build