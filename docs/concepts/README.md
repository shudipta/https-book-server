---
title: Concepts | Scanner
menu:
  product_scanner_0.1.0:
    identifier: concepts-readme
    name: Readme
    parent: concepts
    weight: -1
product_name: scanner
menu_name: product_scanner_0.1.0
section_menu_id: concepts
url: /products/scanner/0.1.0/concepts/
aliases:
  - /products/scanner/0.1.0/concepts/README/
---

# Concepts

Concepts help you learn about the different parts of the Scanner and the abstractions it uses.

- What is Scanner?
  - [Overview](/docs/concepts/what-is-scanner/overview.md). Provides a conceptual introduction to Scanner, including the problems it solves and its high-level architecture.
- Custom Resource Definitions
  - [Restic](/docs/concepts/crds/restic.md). Introduces the concept of `Restic` for configuring [restic](https://restic.net) in a Kubernetes native way.
  - [Recovery](/docs/concepts/crds/recovery.md). Introduces the concept of `Recovery` to restore a backup taken using Scanner.
  - [Repository](/docs/concepts/crds/repository.md) Introduce concept of `Repository` that represents restic repository in a Kubernetes native way.
  - [Snapshot](/docs/concepts/crds/snapshot.md) Introduce concept of `Snapshot` that represents backed up snapshots in a Kubernetes native way.