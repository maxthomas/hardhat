#!/usr/bin/env sh

mkdir -p $GOPATH/bin
go get -u github.com/golang/dep/cmd/dep
which dep
