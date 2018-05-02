---
title: Weclome | Scanner
description: Welcome to Scanner
menu:
  product_scanner_0.1.0:
    identifier: readme-scanner
    name: Readme
    parent: welcome
    weight: -1
product_name: scanner
menu_name: product_scanner_0.1.0
section_menu_id: welcome
url: /products/scanner/0.1.0/welcome/
aliases:
  - /products/scanner/0.1.0/
  - /products/scanner/0.1.0/README/
---

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

From here you can learn all about Scanner's architecture and how to deploy and use Scanner.

- [Concepts](/docs/concepts/). Concepts explain some significant aspect of Scanner. This is where you can learn about what Scanner does and how it does it.

- [Setup](/docs/setup/). Setup contains instructions for installing
  the Scanner in various cloud providers.

- [Guides](/docs/guides/). Guides show you how to perform tasks with Scanner.

- [Reference](/docs/reference/). Detailed exhaustive lists of
command-line options, configuration options, API definitions, and procedures.

We're always looking for help improving our documentation, so please don't hesitate to [file an issue](https://github.com/soter/scanner/issues/new) if you see some problem. Or better yet, submit your own [contributions](/docs/CONTRIBUTING.md) to help
make our docs better.

---

**Scanner binaries collects anonymous usage statistics to help us learn how the software is being used and how we can improve it. To disable stats collection, run the caner with the flag** `--enable-analytics=false`.

---
