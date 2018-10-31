#!/bin/sh

if [[ "$OSTYPE" == "linux-gnu" ]]; then
	export LD_LIBRARY_PATH=$PWD/rticonnextdds-connector/lib/x64Linux2.6gcc4.4.5:$LD_LIBRARY_PATH
	go test -v -race -coverprofile=coverage.txt -covermode=atomic 
elif [[ "$OSTYPE" == "darwin"* ]]; then
	export DYLD_LIBRARY_PATH=rticonnextdds-connector/lib/x64Darwin16clang8.0:$DYLD_LIBRARY_PATH
	go test -v -race -coverprofile=coverage.txt -covermode=atomic 
fi
