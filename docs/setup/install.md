---
title: Install
description: Scanner Install
menu:
  product_scanner_0.1.0:
    identifier: install-scanner
    name: Install
    parent: setup
    weight: 10
product_name: scanner
menu_name: product_scanner_0.1.0
section_menu_id: setup
---

# Installation Guide

Scanner can be installed via a script or as a Helm chart. Installer will deploy Clair with its PostgreSQL database and Scanner as a Kubernetes {validating webhook admission controller](https://kubernetes.io/docs/admin/admission-controllers/#validatingadmissionwebhook-alpha-in-18-beta-in-19) for workloads.

## Using Script

To install Scanner in your Kubernetes cluster, run the following command:

```console
$ curl -fsSL https://raw.githubusercontent.com/soter/scanner/0.1.0/hack/deploy/scanner.sh | bash
```

After successful installation, you should have a `scanner-***` pod running in the `kube-system` namespace.

```console
$ kubectl get pods -n kube-system
NAME                                    READY     STATUS    RESTARTS   AGE
clair-65855cdd5c-bdljz                  1/1       Running   0          44m
clair-postgresql-6dc5fdcbc-6sh62        1/1       Running   0          45m
scanner-7845944d7f-ftpdr                1/1       Running   0          43m
```

#### Customizing Installer

The installer script and associated yaml files can be found in the [/hack/deploy](https://github.com/soter/scanner/tree/0.1.0/hack/deploy) folder. You can see the full list of flags available to installer using `-h` flag.

```console
$ curl -fsSL https://raw.githubusercontent.com/soter/scanner/0.1.0/hack/deploy/scanner.sh | bash -s -- -h
scanner.sh - install Docker image scanner

scanner.sh [options]

options:
-h, --help                         show brief help
-n, --namespace=NAMESPACE          specify namespace (default: kube-system)
    --rbac                         create RBAC roles and bindings (default: true)
    --postgres-storage-class       name of storage class used to store Clair PostgreSQL data (default: standard)
    --docker-registry              docker registry used to pull scanner images (default: appscode)
    --image-pull-secret            name of secret used to pull scanner operator images
    --run-on-master                run scanner operator on master
    --enable-validating-webhook    enable/disable validating webhooks for Scanner
    --enable-analytics             send usage events to Google Analytics (default: true)
    --uninstall                    uninstall scanner
    --purge                        purges Clair installation
```

If you would like to run Scanner pod in `master` instances, pass the `--run-on-master` flag:

```console
$ curl -fsSL https://raw.githubusercontent.com/soter/scanner/0.1.0/hack/deploy/scanner.sh \
    | bash -s -- --run-on-master [--rbac]
```

Scanner will be installed in a `kube-system` namespace by default. If you would like to run Scanner pod in `scanner` namespace, pass the `--namespace=scanner` flag:

```console
$ kubectl create namespace scanner
$ curl -fsSL https://raw.githubusercontent.com/soter/scanner/0.1.0/hack/deploy/scanner.sh \
    | bash -s -- --namespace=scanner [--run-on-master]
```

If you are using a private Docker registry, you need to pull the following image:

 - [soter/scanner](https://hub.docker.com/r/soter/scanner)
 - [soter/clair](https://hub.docker.com/r/soter/clair)

To pass the address of your private registry and optionally a image pull secret use flags `--docker-registry` and `--image-pull-secret` respectively.

```console
$ kubectl create namespace scanner
$ curl -fsSL https://raw.githubusercontent.com/soter/scanner/0.1.0/hack/deploy/scanner.sh \
    | bash -s -- --docker-registry=MY_REGISTRY [--image-pull-secret=SECRET_NAME]
```

Scanner implements [validating admission webhooks](https://kubernetes.io/docs/admin/admission-controllers/#validatingadmissionwebhook-alpha-in-18-beta-in-19) to scan Kubernetes workload types. This is enabled by default for Kubernetes 1.9.0 or later releases. To disable this feature, pass the `--enable-validating-webhook=false` flag.

```console
$ curl -fsSL https://raw.githubusercontent.com/soter/scanner/0.1.0/hack/deploy/scanner.sh \
    | bash -s -- --enable-validating-webhook=false
```

## Using Helm
Scanner can be installed via [Helm](https://helm.sh/) using the [chart](https://github.com/soter/scanner/tree/0.1.0/chart/scanner) from [AppsCode Charts Repository](https://github.com/appscode/charts). To install the chart with the release name `my-release`:

```console
# Mac OSX amd64:
curl -fsSL -o onessl https://github.com/kubepack/onessl/releases/download/0.1.0/onessl-darwin-amd64 \
  && chmod +x onessl \
  && sudo mv onessl /usr/local/bin/

# Linux amd64:
curl -fsSL -o onessl https://github.com/kubepack/onessl/releases/download/0.1.0/onessl-linux-amd64 \
  && chmod +x onessl \
  && sudo mv onessl /usr/local/bin/

# Linux arm64:
curl -fsSL -o onessl https://github.com/kubepack/onessl/releases/download/0.1.0/onessl-linux-arm64 \
  && chmod +x onessl \
  && sudo mv onessl /usr/local/bin/

# Kubernetes 1.9.0 or later
$ helm repo add appscode https://charts.appscode.com/stable/
$ helm repo update
$ helm dependency up appscode/scanner
$ helm install appscode/scanner --name my-release \
  --set apiserver.ca="$(onessl get kube-ca)" \
  --set apiserver.enableValidatingWebhook=true
```

To see the detailed configuration options, visit [here](https://github.com/soter/scanner/tree/master/chart/scanner).

### Installing in GKE Cluster

If you are installing Scanner on a GKE cluster, you will need cluster admin permissions to install Scanner. Run the following command to grant admin permision to the cluster.

```console
# get current google identity
$ gcloud info | grep Account
Account: [user@example.org]

$ kubectl create clusterrolebinding cluster-admin-binding --clusterrole=cluster-admin --user=user@example.org
```


## Verify installation
To check if Scanner pods have started, run the following command:
```console
$ kubectl get pods --all-namespaces -l 'app in (clair,scanner)' -w
NAMESPACE     NAME                       READY     STATUS    RESTARTS   AGE
kube-system   clair-65855cdd5c-bdljz     1/1       Running   0          1h
kube-system   scanner-7845944d7f-ftpdr   1/1       Running   0          1h
```

Once the scanner pod is running, you can cancel the above command by typing `Ctrl+C`.

Now, to confirm apiservice have been registered by the scanner, run the following command:
```console
$ kubectl get apiservice | grep scanner
v1alpha1.admission.scanner.soter.ac    1h
v1alpha1.scanner.soter.ac              1h
```

Now, you are ready to [scan your first image](/docs/guides/README.md) using Scanner.


## Configuring RBAC
Scanner installer will create 1 user facing cluster roles:

| ClusterRole           | Aggregates To     | Desription                            |
|-----------------------|-------------------|---------------------------------------|
| appscode:scanner:view | admin, edit, view | Allows read-only access to Scanner api services, intended to be granted within a namespace using a RoleBinding. |

These user facing roles supports [ClusterRole Aggregation](https://kubernetes.io/docs/admin/authorization/rbac/#aggregated-clusterroles) feature in Kubernetes 1.9 or later clusters.


## Using kubectl for Restic
```console
# Get Restic YAML
$ kubectl get deloyments.scanner.soter.ac -n <namespace> <name> -o yaml
```


## Detect Scanner version
To detect Scanner version, exec into the scanner pod and run `scanner version` command.

```console
$ POD_NAMESPACE=kube-system
$ POD_NAME=$(kubectl get pods -n $POD_NAMESPACE -l app=scanner -o jsonpath={.items[0].metadata.name})
$ kubectl exec -it $POD_NAME -c scanner -n $POD_NAMESPACE scanner version

Version = 0.1.0
VersionStrategy = tag
Os = alpine
Arch = amd64
CommitHash = 85b0f16ab1b915633e968aac0ee23f877808ef49
GitBranch = release-0.5
GitTag = 0.1.0
CommitTimestamp = 2017-10-10T05:24:23
```
