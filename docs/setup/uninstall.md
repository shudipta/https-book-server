---
title: Uninstall
description: Scanner Uninstall
menu:
  product_scanner_0.1.0:
    identifier: uninstall-scanner
    name: Uninstall
    parent: setup
    weight: 20
product_name: scanner
menu_name: product_scanner_0.1.0
section_menu_id: setup
---
# Uninstall Scanner

To uninstall Scanner, run the following command:

```console
$ curl -fsSL https://raw.githubusercontent.com/soter/scanner/0.1.0/hack/deploy/scanner.sh \
    | bash -s -- --uninstall [--namespace=NAMESPACE]

+ kubectl delete deployment -l app=scanner -n kube-system
deployment "scanner-operator" deleted
+ kubectl delete service -l app=scanner -n kube-system
service "scanner-operator" deleted
+ kubectl delete secret -l app=scanner -n kube-system
No resources found
+ kubectl delete serviceaccount -l app=scanner -n kube-system
No resources found
+ kubectl delete clusterrolebindings -l app=scanner -n kube-system
No resources found
+ kubectl delete clusterrole -l app=scanner -n kube-system
No resources found
```

The above command will leave the Scanner crd objects as-is. If you wish to **nuke** Clair installation, also pass the `--purge` flag.
