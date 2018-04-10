# Clair Run in Minikube

Before we go forward with scanner, we need to ensure that clair is already running. Here, we need two configuration files:

- config.yaml
- clair-kubernetes.yaml

These should be available in `docs/examples/clair` directory.

Then we need to run the following commands from the repository root,

```console
kubectl create secret generic clairsecret --from-file=docs/examples/clair/clair-config.yaml

kubectl create -f docs/examples/clair/clair-kubernetes.yaml
```

Finally, clair should be run in `https://192.168.99.100:30060`.
