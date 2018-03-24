#!/usr/bin/env sh

set -e

cd cmd/hardhat

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o hardhat .
