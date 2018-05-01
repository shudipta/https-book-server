# Scanner
[Scanner by AppsCode](https://github.com/soter/scanner) - Docker Image Scanner.

## TL;DR;

```console
$ helm repo add appscode https://charts.appscode.com/stable/
$ helm repo update
$ helm install soter/scanner
```

## Introduction

This chart bootstraps a [Scanner server](https://github.com/soter/scanner) deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.9+

## Installing the Chart
To install the chart with the release name `my-release`:

```console
$ helm install soter/scanner --name my-release
```

The command deploys Scanner server on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `my-release`:

```console
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the Scanner chart and their default values.


| Parameter                           | Description                                                       | Default            |
| ----------------------------------- | ----------------------------------------------------------------- | ------------------ |
| `replicaCount`                      | Number of Scanner server replicas to create (only 1 is supported) | `1`                |
| `scanner.registry`                  | Docker registry used to pull Scanner image                        | `soter`            |
| `scanner.repository`                | Scanner container image                                           | `scanner`          |
| `scanner.tag`                       | Scanner container image tag                                       | `0.1.0`            |
| `imagePullPolicy`                   | container image pull policy                                       | `IfNotPresent`     |
| `criticalAddon`                     | If true, installs Scanner server as critical addon                | `false`            |
| `rbac.create`                       | If `true`, create and use RBAC resources                          | `true`             |
| `serviceAccount.create`             | If `true`, create a new service account                           | `true`             |
| `serviceAccount.name`               | Service account to be used. If not set and `serviceAccount.create` is `true`, a name is generated using the fullname template | `` |
| `apiserver.groupPriorityMinimum`    | The minimum priority the group should have.                       | 10000              |
| `apiserver.versionPriority`         | The ordering of this API inside of the group.                     | 15                 |
| `apiserver.enableValidatingWebhook` | Enable validating webhooks for Kubernetes workloads               | false              |
| `apiserver.enableMutatingWebhook`   | Enable mutating webhooks for Kubernetes workloads                 | false              |
| `apiserver.ca`                      | CA certificate used by main Kubernetes api server                 | ``                 |
| `enableAnalytics`                   | Send usage events to Google Analytics                             | `true`             |


Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example:

```console
$ helm install --name my-release --set image.tag=v0.2.1 soter/scanner
```

Alternatively, a YAML file that specifies the values for the parameters can be provided while
installing the chart. For example:

```console
$ helm install --name my-release --values values.yaml soter/scanner
```

## RBAC
By default the chart will not install the recommended RBAC roles and rolebindings.

You need to have the flag `--authorization-mode=RBAC` on the api server. See the following document for how to enable [RBAC](https://kubernetes.io/docs/admin/authorization/rbac/).

To determine if your cluster supports RBAC, run the following command:

```console
$ kubectl api-versions | grep rbac
```

If the output contains "beta", you may install the chart with RBAC enabled (see below).

### Enable RBAC role/rolebinding creation

To enable the creation of RBAC resources (On clusters with RBAC). Do the following:

```console
$ helm install --name my-release soter/scanner --set rbac.create=true
```
