#!/bin/bash

go get -v github.com/stretchr/testify
curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.17.1

golangci-lint run

if [[ "$OSTYPE" == "linux-gnu" ]]; then
	export LD_LIBRARY_PATH=$PWD/rticonnextdds-connector/lib/x64Linux2.6gcc4.4.5:$LD_LIBRARY_PATH
	go test -v -race -coverprofile=coverage.txt -covermode=atomic 
elif [[ "$OSTYPE" == "darwin"* ]]; then
	export DYLD_LIBRARY_PATH=rticonnextdds-connector/lib/x64Darwin16clang8.0:$DYLD_LIBRARY_PATH
	go test -v -race -coverprofile=coverage.txt -covermode=atomic 
fi
