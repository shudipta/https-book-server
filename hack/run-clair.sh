#!/usr/bin/env bash
set -xe

GOPATH=$(go env GOPATH)
PACKAGE_NAME=https-book-server
REPO_ROOT="$GOPATH/src/github.com/shudipta/$PACKAGE_NAME"

pushd $REPO_ROOT
go run cert-generator/certgen.go

if [ -x "$(command -v onessl)" ]; then
    export ONESSL=onessl
else
    curl -fsSL -o onessl https://github.com/kubepack/onessl/releases/download/0.1.0/onessl-linux-amd64
    chmod +x onessl
    export ONESSL=./onessl
fi

export CA_CERT=$(cat cert-generator/ca.crt | $ONESSL base64)
export CLIENT_CERT=$(cat cert-generator/cl.crt | $ONESSL base64)
export CLIENT_KEY=$(cat cert-generator/cl.key | $ONESSL base64)

kubectl delete secret clairsecret
kubectl delete secret -l app=clair
kubectl delete svc -l app=clair
kubectl delete rc -l app=clair
kubectl delete svc -l app=postgres
kubectl delete rc -l app=postgres
kubectl create secret generic clairsecret --from-file=./hack/deploy/config.yaml
cat ./hack/deploy/clair-kubernetes.yaml | $ONESSL envsubst | kubectl apply -f -

pushd $REPO_ROOT/cert-generator
rm -rf ca.crt ca.key srv.crt srv.key cl.crt cl.key
popd
popd