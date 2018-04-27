# Clair Setup Procedure

Before we go forward with scanner, we need to ensure that clair is already running.

## Run Clair in Minikube

 Here, we need two configuration files and we have provided two sets of these files in `docs/examples/clair` directory.

For using without tls,

- config.yaml
- clair-kubernetes.yaml

For using with tls,

- config-tls.yaml
- clair-kubernetes-tls.yaml


### Without TLS Configuration

If we want to run clair without tls configuration, then we just have to run the following commands from the repository root

```console
$ kubectl create secret generic clairsecret --from-file=docs/examples/clair/config.yaml
secret "clairsecret" created

$ kubectl create -f docs/examples/clair/clair-kubernetes.yaml
service "clairsvc" created
replicationcontroller "clair" created
replicationcontroller "clair-postgres" created
service "postgres" created
```

### With TLS Configuration

We can easily generate necessary certificates for tls configuration using `onessl` tool.

First, we'll install onessl by,

```console
$ curl -fsSL -o onessl https://github.com/kubepack/onessl/releases/download/0.1.0/onessl-linux-amd64
$ chmod +x onessl
$ export ONESSL=./onessl
```

#### For Api Part

If we want to use clair api with tls configuration, then we have to use one certificate for server and another for client to verify the identity of themselves. Usually, another certificate called CA certificate is used to verify identity. Here, clair is server and scanner is client. So, here we need certificates for CA, server and client.

Now, we will generate necessary certificates for tls secure clair api,

```console
$ $ONESSL create ca-cert --cert-dir=clair-cert/
Wrote ca certificates in  clair-cert/

$ $ONESSL create server-cert server --cert-dir=clair-cert/ --domains=clairsvc.default.svc --ips=192.168.99.100
Wrote server certificates in clair-cert/

$ $ONESSL create client-cert client --cert-dir=clair-cert/ --organization=soter.ac
Wrote client certificates in clair-cert/
```

First one will create `ca.crt` and `ca.key`. Second one will create `server.crt` and `server.key`. And the third one will create `client@soter.ac.crt` and `client@soter.ac.key`.

Now, under api field in `config-tls.yaml` file, we have to provide the path for `ca.crt`, `server.crt` and `server.key` files like followings:

```yaml
  api:
    # API server port
    port: 6060

    # Health server port
    # This is an unencrypted endpoint useful for load balancers to check to healthiness of the clair server.
    healthport: 6061

    # Deadline before an API request will respond with a 503
    timeout: 900s

    # 32-bit URL-safe base64 key used to encrypt pagination tokens
    # If one is not provided, it will be generated.
    # Multiple clair instances in the same cluster need the same value.
    paginationkey:

    # Optional PKI configuration
    # If you want to easily generate client certificates and CAs, try the following projects:
    # https://github.com/coreos/etcd-ca
    # https://github.com/cloudflare/cfssl
    servername:
    cafile: /var/clair-api-serving-cert/ca.crt
    keyfile: /var/clair-api-serving-cert/srv.key
    certfile: /var/clair-api-serving-cert/srv.crt
```

The client certificate will be used for scanner.

#### For Notifier Part

If we want to use clair notifier with tls configuration, we have to do the same as for clair api. In this case, clair is client and scanner is server. Here, we need certificates for CA and client.

Now, we will generate necessary certificates for tls secure clair notifier,

```console
$ export SCANNER_NAMESPACE=kube-system

$ $ONESSL create ca-cert
Wrote ca certificates in path/to/gopath/src/github.com/soter/scanner

$ $ONESSL create server-cert server --domains=scanner.$SCANNER_NAMESPACE.svc
Wrote server certificates in path/to/gopath/src/github.com/soter/scanner

$ $ONESSL create client-cert client --organization=clair
Wrote client certificates in path/to/gopath/src/github.com/soter/scanner
```

First one will create `ca.crt` and `ca.key`. Second one will create `server.crt` and `server.key`. And the third one will create `client@clair.crt` and `client@clair.key`.

Now, under notifier field in `config-tls.yaml` file, we have to provide the path for `ca.crt`, `client@clair.crt` and `client@clair.key` files like followings:

```yaml
  notifier:
    # Number of attempts before the notifier is marked as failed to be sent
    attempts: 3

    # Duration before a failed notification is retried
    renotifyinterval: 5m

    http:
      # Optional endpoint that will receive notifications via POST requests
      endpoint: https://scanner.kube-system.svc:443/audit-log

      # Optional PKI configuration
      # If you want to easily generate client certificates and CAs, try the following projects:
      # https://github.com/cloudflare/cfssl
      # https://github.com/coreos/etcd-ca
      servername:
      cafile: /var/clair-notifier-serving-cert/ca.crt
      keyfile: /var/clair-notifier-serving-cert/cl.key
      certfile: /var/clair-notifier-serving-cert/cl.crt

      # Optional HTTP Proxy: must be a valid URL (including the scheme).
      proxy:
```

In `config-tls.yaml` file, we have added for both api and notifier. We have used the path where the certificate files are mounted as volume in the pod running for clair.

Then, we need to run the following commands.

```console
# Exporting certificates into variables for clair notifier.
export CLAIR_NOTIFIER_SERVING_CERT_CA=$(cat ca.crt | $ONESSL base64)
export CLAIR_NOTIFIER_CLIENT_CERT=$(cat client@clair.crt | $ONESSL base64)
export CLAIR_NOTIFIER_CLIENT_KEY=$(cat client@clair.key | $ONESSL base64)

# Exporting certificates into variables for clair notifier.
$ export CLAIR_API_SERVING_CERT_CA=$(cat clair-cert/ca.crt | $ONESSL base64)
$ export CLAIR_API_SERVER_CERT=$(cat clair-cert/server.crt | $ONESSL base64)
$ export CLAIR_API_SERVER_KEY=$(cat clair-cert/server.key | $ONESSL base64)

# Running clair
$ kubectl create secret generic clairsecret --from-file=docs/examples/clair/config.yaml
secret "clairsecret" created

$ kubectl label secret clairsecret app=clair
secret "clairsecret" labeled

$ cat docs/examples/clair/clair-kubernetes.yaml | $ONESSL envsubst | kubectl apply -f -
service "clairsvc" created
replicationcontroller "clair" created
replicationcontroller "clair-postgres" created
service "postgres" created
```

Finally, clair should be run in `https://192.168.99.100:30060`.
