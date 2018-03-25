#!/usr/bin/env sh

dep ensure
go get -u golang.org/x/tools/cmd/cover
go get -u github.com/mitchellh/gox
