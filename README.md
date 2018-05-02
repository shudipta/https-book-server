[![Go Report Card](https://goreportcard.com/badge/github.com/soter/scanner)](https://goreportcard.com/report/github.com/soter/scanner)
[![Build Status](https://travis-ci.org/soter/scanner.svg?branch=master)](https://travis-ci.org/soter/scanner)
[![codecov](https://codecov.io/gh/soter/scanner/branch/master/graph/badge.svg)](https://codecov.io/gh/soter/scanner)
[![Docker Pulls](https://img.shields.io/docker/pulls/soter/scanner.svg)](https://hub.docker.com/r/soter/scanner/)
[![Slack](https://slack.appscode.com/badge.svg)](https://slack.appscode.com)
[![Twitter](https://img.shields.io/twitter/follow/appscodehq.svg?style=social&logo=twitter&label=Follow)](https://twitter.com/intent/follow?screen_name=AppsCodeHQ)

# Scanner
Scanner by AppsCode is a Docker image scanner. It uses [Clair](https://github.com/coreos/clair) for the static analysis of vulnerabilities in Docker containers. Using Scanner, you can backup Kubernetes volumes mounted in following types of workloads:

- Deployment
- DaemonSet
- ReplicaSet
- ReplicationController
- StatefulSet
- Pod
- Job
- CronJob
- Openshift DeploymentConfig

## Supported Versions
Please pick a version of Scanner that matches your Kubernetes installation.

| Scanner Version                                                            | Docs                                                      | Kubernetes Version |
|----------------------------------------------------------------------------|-----------------------------------------------------------|--------------------|
| [0.1.0](https://github.com/appscode/scanner/releases/tag/0.1.0) (uses CRD) | [User Guide](https://appscode.com/products/scanner/0.1.0) | 1.9.x +            |

## Installation

To install Scanner, please follow the guide [here](https://appscode.com/products/scanner/0.1.0/setup/install).

## Using Scanner
Want to learn how to use Scanner? Please start [here](https://appscode.com/products/scanner/0.1.0).

## Contribution guidelines
Want to help improve Scanner? Please start [here](https://appscode.com/products/scanner/0.1.0/welcome/contributing).

---

**Scanner binaries collects anonymous usage statistics to help us learn how the software is being used and how we can improve it. To disable stats collection, run the scanner with the flag** `--enable-analytics=false`.

---

## Acknowledgement
 - Many thanks to CoreOS for [Clair](https://github.com/coreos/clair) project.

## Support
We use Slack for public discussions. To chit chat with us or the rest of the community, join us in the [AppsCode Slack team](https://appscode.slack.com/messages/CAER85GPK/details/) channel `#scanner`. To sign up, use our [Slack inviter](https://slack.appscode.com/).

If you have found a bug with Scanner or want to request for new features, please [file an issue](https://github.com/appscode/scanner/issues/new).
