#!/usr/bin/env sh

glide up
go get -u golang.org/x/tools/cmd/cover
go get -u github.com/mitchellh/gox
