#!/usr/bin/env bash
set -xe

GOPATH=$(go env GOPATH)
PACKAGE_NAME=https-book-server
REPO_ROOT="$GOPATH/src/github.com/shudipta/$PACKAGE_NAME"

pushd $REPO_ROOT
go build -o hack/docker/https-book-server book-server/book_server.go

docker build -t shudipta/https-book-server:v1 hack/docker/
docker save shudipta/https-book-server:v1 | pv | (eval $(minikube docker-env) && docker load)
docker push shudipta/https-book-server:v1

rm -rf hack/docker/https-book-server
popd