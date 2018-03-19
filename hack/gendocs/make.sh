#!/usr/bin/env bash

pushd $GOPATH/src/github.com/soter/scanner/hack/gendocs
go run main.go
popd
