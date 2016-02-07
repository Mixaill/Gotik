#!/bin/bash
killall Kotik
export GODIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"/../GODIR
export GOPATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"/../GOPATH
mkdir $GODIR
wget -O $GODIR/go.tar.gz https://storage.googleapis.com/golang/go1.6rc2.linux-amd64.tar.gz
tar -xzf $GODIR/go.tar.gz -C $GODIR
export GOROOT=$GODIR/go
git pull origin master
$GOROOT/bin/go get
$GOROOT/bin/go build .
./Kotik -dev=false
