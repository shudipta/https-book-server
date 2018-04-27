#!/bin/bash
set -eou pipefail

apiServices=(v1alpha1.admission)

echo "checking kubeconfig context"
kubectl config current-context || { echo "Set a context (kubectl use-context <context>) out of the following:"; echo; kubectl config get-contexts; exit 1; }
echo ""

# http://redsymbol.net/articles/bash-exit-traps/
function cleanup {
    rm -rf $ONESSL ca.crt ca.key server.crt server.key
}
trap cleanup EXIT

# ref: https://github.com/appscodelabs/libbuild/blob/master/common/lib.sh#L55
inside_git_repo() {
    git rev-parse --is-inside-work-tree > /dev/null 2>&1
    inside_git=$?
    if [ "$inside_git" -ne 0 ]; then
        echo "Not inside a git repository"
        exit 1
    fi
}

detect_tag() {
    inside_git_repo

    # http://stackoverflow.com/a/1404862/3476121
    git_tag=$(git describe --exact-match --abbrev=0 2>/dev/null || echo '')

    commit_hash=$(git rev-parse --verify HEAD)
    git_branch=$(git rev-parse --abbrev-ref HEAD)
    commit_timestamp=$(git show -s --format=%ct)

    if [ "$git_tag" != '' ]; then
        TAG=$git_tag
        TAG_STRATEGY='git_tag'
    elif [ "$git_branch" != 'master' ] && [ "$git_branch" != 'HEAD' ] && [[ "$git_branch" != release-* ]]; then
        TAG=$git_branch
        TAG_STRATEGY='git_branch'
    else
        hash_ver=$(git describe --tags --always --dirty)
        TAG="${hash_ver}"
        TAG_STRATEGY='commit_hash'
    fi

    export TAG
    export TAG_STRATEGY
    export git_tag
    export git_branch
    export commit_hash
    export commit_timestamp
}

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

# ref: https://stackoverflow.com/a/7069755/244009
# ref: https://jonalmeida.com/posts/2013/05/26/different-ways-to-implement-flags-in-bash/
# ref: http://tldp.org/LDP/abs/html/comparison-ops.html

export SCANNER_NAMESPACE=kube-system
export SCANNER_SERVICE_ACCOUNT=scanner
export SCANNER_ENABLE_RBAC=true
export SCANNER_RUN_ON_MASTER=0
export SCANNER_ENABLE_VALIDATING_WEBHOOK=false
export SCANNER_DOCKER_REGISTRY=soter
export SCANNER_SERVER_TAG=canary
export SCANNER_IMAGE_PULL_SECRET=
export SCANNER_IMAGE_PULL_POLICY=IfNotPresent
export SCANNER_ENABLE_ANALYTICS=true
export SCANNER_UNINSTALL=0

export APPSCODE_ENV=${APPSCODE_ENV:-prod}
export SCRIPT_LOCATION="curl -fsSL https://raw.githubusercontent.com/soter/scanner/master/"
if [ "$APPSCODE_ENV" = "dev" ]; then
    detect_tag
    export SCRIPT_LOCATION="cat "
    export SCANNER_SERVER_TAG=$TAG
    export SCANNER_IMAGE_PULL_POLICY=Always
fi

KUBE_APISERVER_VERSION=$(kubectl version -o=json | $ONESSL jsonpath '{.serverVersion.gitVersion}')
$ONESSL semver --check='>=1.9.0' $KUBE_APISERVER_VERSION
if [ $? -eq 0 ]; then
    export SCANNER_ENABLE_VALIDATING_WEBHOOK=true
fi

show_help() {
    echo "scanner.sh - install scanner"
    echo " "
    echo "scanner.sh [options]"
    echo " "
    echo "options:"
    echo "-h, --help                         show brief help"
    echo "-n, --namespace=NAMESPACE          specify namespace (default: kube-system)"
    echo "    --rbac                         create RBAC roles and bindings (default: true)"
    echo "    --docker-registry              docker registry used to pull scanner images (default: appscode)"
    echo "    --image-pull-secret            name of secret used to pull scanner operator images"
    echo "    --run-on-master                run scanner operator on master"
    echo "    --enable-validating-webhook    enable/disable validating webhooks for Scanner"
    echo "    --enable-analytics             send usage events to Google Analytics (default: true)"
    echo "    --uninstall                    uninstall scanner"
}

while test $# -gt 0; do
    case "$1" in
        -h|--help)
            show_help
            exit 0
            ;;
        -n)
            shift
            if test $# -gt 0; then
                export SCANNER_NAMESPACE=$1
            else
                echo "no namespace specified"
                exit 1
            fi
            shift
            ;;
        --namespace*)
            export SCANNER_NAMESPACE=`echo $1 | sed -e 's/^[^=]*=//g'`
            shift
            ;;
        --docker-registry*)
            export SCANNER_DOCKER_REGISTRY=`echo $1 | sed -e 's/^[^=]*=//g'`
            shift
            ;;
        --image-pull-secret*)
            secret=`echo $1 | sed -e 's/^[^=]*=//g'`
            export SCANNER_IMAGE_PULL_SECRET="name: '$secret'"
            shift
            ;;
        --enable-validating-webhook*)
            val=`echo $1 | sed -e 's/^[^=]*=//g'`
            if [ "$val" = "false" ]; then
                export SCANNER_ENABLE_VALIDATING_WEBHOOK=false
            fi
            shift
            ;;
        --enable-analytics*)
            val=`echo $1 | sed -e 's/^[^=]*=//g'`
            if [ "$val" = "false" ]; then
                export SCANNER_ENABLE_ANALYTICS=false
            fi
            shift
            ;;
        --rbac*)
            val=`echo $1 | sed -e 's/^[^=]*=//g'`
            if [ "$val" = "false" ]; then
                export SCANNER_SERVICE_ACCOUNT=default
                export SCANNER_ENABLE_RBAC=false
            fi
            shift
            ;;
        --run-on-master)
            export SCANNER_RUN_ON_MASTER=1
            shift
            ;;
        --uninstall)
            export SCANNER_UNINSTALL=1
            shift
            ;;
        *)
            show_help
            exit 1
            ;;
    esac
done

if [ "$SCANNER_UNINSTALL" -eq 1 ]; then
    # delete webhooks and apiservices
    kubectl delete validatingwebhookconfiguration -l app=scanner
    kubectl delete mutatingwebhookconfiguration -l app=scanner
    kubectl delete apiservice -l app=scanner
    # delete scanner operator
    kubectl delete deployment -l app=scanner --namespace $SCANNER_NAMESPACE
    kubectl delete service -l app=scanner --namespace $SCANNER_NAMESPACE
    kubectl delete secret -l app=scanner --namespace $SCANNER_NAMESPACE
    # delete RBAC objects, if --rbac flag was used.
    kubectl delete serviceaccount -l app=scanner --namespace $SCANNER_NAMESPACE
    kubectl delete clusterrolebindings -l app=scanner
    kubectl delete clusterrole -l app=scanner
    kubectl delete rolebindings -l app=scanner --namespace $SCANNER_NAMESPACE
    kubectl delete role -l app=scanner --namespace $SCANNER_NAMESPACE

    echo "waiting for scanner operator pod to stop running"
    for (( ; ; )); do
       pods=($(kubectl get pods --all-namespaces -l app=scanner -o jsonpath='{range .items[*]}{.metadata.name} {end}'))
       total=${#pods[*]}
        if [ $total -eq 0 ] ; then
            break
        fi
       sleep 2
    done

    echo
    echo "Successfully uninstalled Scanner!"
    exit 0
fi

echo "checking whether extended apiserver feature is enabled"
$ONESSL has-keys configmap --namespace=kube-system --keys=requestheader-client-ca-file extension-apiserver-authentication || { echo "Set --requestheader-client-ca-file flag on Kubernetes apiserver"; exit 1; }
echo ""

env | sort | grep SCANNER*
echo ""

export SERVICE_SERVING_CERT_CA=$(cat ca.crt | $ONESSL base64)
export TLS_SERVING_CERT=$(cat server.crt | $ONESSL base64)
export TLS_SERVING_KEY=$(cat server.key | $ONESSL base64)
export KUBE_CA=$($ONESSL get kube-ca | $ONESSL base64)

export CLAIR_API_SERVING_CERT_CA=$(cat clair-cert/ca.crt | $ONESSL base64)
export CLAIR_API_CLIENT_CERT=$(cat clair-cert/client@soter.ac.crt | $ONESSL base64)
export CLAIR_API_CLIENT_KEY=$(cat clair-cert/client@soter.ac.key | $ONESSL base64)

${SCRIPT_LOCATION}hack/deploy/deployment.yaml | $ONESSL envsubst | kubectl apply -f -

if [ "$SCANNER_ENABLE_RBAC" = true ]; then
    kubectl create serviceaccount $SCANNER_SERVICE_ACCOUNT --namespace $SCANNER_NAMESPACE
    kubectl label serviceaccount $SCANNER_SERVICE_ACCOUNT app=scanner --namespace $SCANNER_NAMESPACE
    ${SCRIPT_LOCATION}hack/deploy/rbac-list.yaml | $ONESSL envsubst | kubectl auth reconcile -f -
fi

if [ "$SCANNER_RUN_ON_MASTER" -eq 1 ]; then
    kubectl patch deploy scanner -n $SCANNER_NAMESPACE \
      --patch="$(${SCRIPT_LOCATION}hack/deploy/run-on-master.yaml)"
fi

if [ "$SCANNER_ENABLE_VALIDATING_WEBHOOK" = true ]; then
    ${SCRIPT_LOCATION}hack/deploy/validating-webhook.yaml | $ONESSL envsubst | kubectl apply -f -
fi

echo
echo "waiting until scanner deployment is ready"
$ONESSL wait-until-ready deployment scanner --namespace $SCANNER_NAMESPACE || { echo "Scanner deployment failed to be ready"; exit 1; }

echo "waiting until scanner apiservice is available"
for api in "${apiServices[@]}"; do
    $ONESSL wait-until-ready apiservice ${api}.scanner.soter.ac || { echo "Scanner apiservice $api failed to be ready"; exit 1; }
done

echo
echo "Successfully installed Scanner!"
