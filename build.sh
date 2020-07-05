#!/usr/bin/env bash

PACKAGE="github.com/jtogrul/comake"

env GOOS=linux GOARCH=amd64 go build -o dist/comake_linux ${PACKAGE}
env GOOS=darwin GOARCH=amd64 go build -o dist/comake_darwin ${PACKAGE}
env GOOS=windows GOARCH=amd64 go build -o dist/comake_win64.exe ${PACKAGE}