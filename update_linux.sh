#!/bin/bash
killall Kotik
export GOPATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"/../GOPATH
git pull origin master
go get
go get -u
go build .
./Kotik -dev=false
