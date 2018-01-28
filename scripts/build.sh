#!/usr/bin/env sh

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o hardhat .
