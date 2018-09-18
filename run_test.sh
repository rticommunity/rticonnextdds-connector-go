#!/bin/sh

export LD_LIBRARY_PATH=$PWD/rticonnextdds-connector/lib/x64Linux2.6gcc4.4.5:$LD_LIBRARY_PATH
go test -v ./test
