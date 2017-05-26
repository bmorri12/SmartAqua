#!/bin/sh

# restore dependency packages. 
cd $GOPATH/src/github.com/bmorri12/SmartAqua
cp -r Godeps/_workspace/src/* $GOPATH/src

# install binaries
go install -v github.com/bmorri12/SmartAqua/services/...
