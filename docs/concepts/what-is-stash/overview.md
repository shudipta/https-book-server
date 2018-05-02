---
title: Scanner Overview
description: Scanner Overview
menu:
  product_scanner_0.1.0:
    identifier: overview-concepts
    name: Overview
    parent: what-is-scanner
    weight: 10
product_name: scanner
menu_name: product_scanner_0.1.0
section_menu_id: concepts
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
