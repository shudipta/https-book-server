#!/usr/bin/env bash
set -eou pipefail

while test $# -gt 0; do
    case "$1" in
        --uninstall)
            kubectl delete secret -l app=clair
            kubectl delete service -l app=clair
            kubectl delete service -l app=postgres
            kubectl delete rc -l app=clair
            kubectl delete rc -l app=postgres
            exit 1
            ;;
    esac
done

# https://stackoverflow.com/a/677212/244009
if [ -x "$(command -v onessl)" ]; then
    export ONESSL=onessl
else
    # ref: https://stackoverflow.com/a/27776822/244009
    case "$(uname -s)" in
        Darwin)
            curl -fsSL -o onessl https://github.com/kubepack/onessl/releases/download/0.1.0/onessl-darwin-amd64
            chmod +x onessl
            export ONESSL=./onessl
            ;;

        Linux)
            curl -fsSL -o onessl https://github.com/kubepack/onessl/releases/download/0.1.0/onessl-linux-amd64
            chmod +x onessl
            export ONESSL=./onessl
            ;;

        CYGWIN*|MINGW32*|MSYS*)
            curl -fsSL -o onessl.exe https://github.com/kubepack/onessl/releases/download/0.1.0/onessl-windows-amd64.exe
            chmod +x onessl.exe
            export ONESSL=./onessl.exe
            ;;
        *)
            echo 'other OS'
            ;;
    esac
fi

export SCANNER_NAMESPACE=kube-system

export APPSCODE_ENV=${APPSCODE_ENV:-prod}
export SCRIPT_LOCATION="curl -fsSL https://raw.githubusercontent.com/soter/scanner/master/"
if [ "$APPSCODE_ENV" = "dev" ]; then
    export SCRIPT_LOCATION="cat "
fi

echo "running clair"

# Exporting certificates for clair notifier.
export CLAIR_NOTIFIER_SERVING_CERT_CA=$(cat pki/scanner/ca.crt | $ONESSL base64)
export NOTIFIER_CLIENT_CERT=$(cat pki/scanner/client@clair.crt | $ONESSL base64)
export NOTIFIER_CLIENT_KEY=$(cat pki/scanner/client@clair.key | $ONESSL base64)

# Exporting certificates for clair api.
export CLAIR_API_SERVING_CERT_CA=$(cat pki/clair/ca.crt | $ONESSL base64)
export CLAIR_API_SERVER_CERT=$(cat pki/clair/server.crt | $ONESSL base64)
export CLAIR_API_SERVER_KEY=$(cat pki/clair/server.key | $ONESSL base64)

# Running clair
kubectl create secret generic clairsecret --from-file=docs/examples/clair/config.yaml
kubectl label secret clairsecret app=clair
${SCRIPT_LOCATION}hack/deploy/clair/clair.yaml | $ONESSL envsubst | kubectl apply -f -
