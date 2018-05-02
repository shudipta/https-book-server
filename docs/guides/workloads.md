---
title: Workloads | Scanner
description: workloads of Scanner
menu:
  product_scanner_0.1.0:
    identifier: workloads-scanner
    name: Workloads
    parent: guides
    weight: 20
product_name: scanner
menu_name: product_scanner_0.1.0
section_menu_id: guides
---

> New to Scanner? Please start [here](/docs/concepts/README.md).

# Supported Workloads

Scanner supports the following types of Kubernetes workloads.

## Deployments
To backup a Deployment, create a Restic with matching selectors. You can find a full working demo in [examples folder](/docs/examples/workloads/deployment.yaml).

## ReplicaSets
To backup a ReplicaSet, create a Restic with matching selectors. You can find a full working demo in [examples folder](/docs/examples/workloads/replicaset.yaml).

## ReplicationControllers
To backup a ReplicationController, create a Restic with matching selectors. You can find a full working demo in [examples folder](/docs/examples/workloads/rc.yaml).

## DaemonSets
To backup a DaemonSet, create a Restic with matching selectors. You can find a full working demo in [examples folder](/docs/examples/workloads/daemonset.yaml). This example shows how Scanner can be used to backup host paths on all nodes of a cluster. First run a DaemonSet without nodeSelectors. This DaemonSet acts as a vector for Restic sidecar and mounts host paths that are to be backed up. In this example, we use a `busybox` container for this. Now, create a Restic that has a matching selector. This Restic also `spec.volumeMounts` the said host path and points to the host path in `spec.fileGroups`.

## StatefulSets
Kubernetes does not support adding sidecar to a StatefulSet after it is created. It is recommended to enable **mutating webhook** by providing `--enable-mutating-webhook` flag while installing scanner. To know more about how to provide various flag while installing scanner see [here](/docs/setup/install.md#customizing-installer).

If you don't want to enable **mutating webhook** then you have to add Scanner sidecar container to your StatefulSet manually. You can see the relevant portions of a working example below:

```yaml
apiVersion: apps/v1beta1
kind: StatefulSet
metadata:
  labels:
    app: statefulset-demo
  name: workload
  namespace: default
spec:
  replicas: 1
  serviceName: headless
  template:
    metadata:
      labels:
        app: statefulset-demo
      name: busybox
    spec:
      serviceAccountName: statefulset-demo
      containers:
      - image: busybox
        name: busybox
        imagePullPolicy: IfNotPresent
        command:
        - sleep
        - "3600"
        resources: {}
        volumeMounts:
        - mountPath: /source/data
          name: source-data
      - image: soter/scanner:0.1.0
        name: scanner
        imagePullPolicy: IfNotPresent
        args:
        - backup
        - --restic-name=statefulset-restic
        - --workload-kind=Statefulset
        - --workload-name=workload
        - --run-via-cron=true
        - --v=3
        env:
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: spec.nodeName
        - name: POD_NAME
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.name
        resources: {}
        volumeMounts:
        - mountPath: /tmp
          name: scanner-scratchdir
        - mountPath: /etc/scanner
          name: scanner-podinfo
        - mountPath: /source/data
          name: source-data
          readOnly: true
        - mountPath: /safe/data
          name: scanner-local
      volumes:
      - gitRepo:
          repository: https://github.com/soter/scanner-data.git
        name: source-data
      - hostPath:
          path: /data/scanner-test/restic-repo
          type: ""
        name: scanner-local
      - emptyDir: {}
        name: scanner-scratchdir
      - downwardAPI:
          defaultMode: 420
          items:
          - fieldRef:
              apiVersion: v1
              fieldPath: metadata.labels
            path: labels
        name: scanner-podinfo
```

You can find the full working demo in [examples folder](/docs/examples/workloads/statefulset.yaml). The section you should change for your own StatefulSet are:

 - `--restic-name` flag should be set to the name of the Restic used as configuration.
 - `--workload-kind` flag specifies the kind of workload (Deployment/Replicaset/RepliationController/DaemonSet/StatefulSet).
 - `--workload-name` flag specifies the name of workload where sidecar pod is added.

To learn about the meaning of various flags, please visit [here](/docs/reference/scanner_backup.md).

## Next Steps

- Learn how to use Scanner to backup a Kubernetes deployment [here](/docs/guides/backup.md).
- Learn about the details of Restic CRD [here](/docs/concepts/crds/restic.md).
- To restore a backup see [here](/docs/guides/restore.md).
- Learn about the details of Recovery CRD [here](/docs/concepts/crds/recovery.md).
- To run backup in offline mode see [here](/docs/guides/offline_backup.md)
- See the list of supported backends and how to configure them [here](/docs/guides/backends.md).
- Thinking about monitoring your backup operations? Scanner works [out-of-the-box with Prometheus](/docs/guides/monitoring.md).
- Learn about how to configure [RBAC roles](/docs/guides/rbac.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on Scanner? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
