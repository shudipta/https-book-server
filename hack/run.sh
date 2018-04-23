#!/usr/bin/env bash
set -xe

GOPATH=$(go env GOPATH)
PACKAGE_NAME=https-book-server
REPO_ROOT="$GOPATH/src/github.com/shudipta/$PACKAGE_NAME"

pushd $REPO_ROOT
go build -o hack/docker/https-book-server book-server/book_server.go

docker build -t shudipta/https-book-server:v1 hack/docker/
docker push shudipta/https-book-server:v1
# docker save shudipta/https-book-server:v1 | pv | (eval $(minikube docker-env) && docker load)

go run cert-generator/certgen.go
if [ -x "$(command -v onessl)" ]; then
    export ONESSL=onessl
else
    curl -fsSL -o onessl https://github.com/kubepack/onessl/releases/download/0.1.0/onessl-linux-amd64
    chmod +x onessl
    export ONESSL=./onessl
fi

export CA_CERT=$(cat cert-generator/ca.crt | $ONESSL base64)
export SERVER_CERT=$(cat cert-generator/cl.crt | $ONESSL base64)
export SERVER_KEY=$(cat cert-generator/cl.key | $ONESSL base64)

kubectl delete secret -l app=https
kubectl delete svc -l app=https
kubectl delete deploy -l app=https
cat ./hack/deploy/deployment.yaml | $ONESSL envsubst | kubectl apply -f -

pushd $REPO_ROOT/cert-generator
# rm -rf ca.crt ca.key srv.crt srv.key cl.crt cl.key
popd
rm -rf hack/docker/https-book-server
popd