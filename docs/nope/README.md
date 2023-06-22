# Nope - Kubernetes Operator for managing NATS Accounts

Nope is a Kubernetes Operator for managing [decentralized JWT authentication](https://docs.nats.io/running-a-nats-service/configuration/securing_nats/auth_intro/jwt) for accounts in NATS.

// TODO: write a better sentence above.

## Description

// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started

You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster

1. Install Instances of Custom Resources:

    ```sh
    kubectl apply -f config/samples/
    ```

2. Build and push your image to the location specified by `IMG`:

    ```sh
    make docker-build docker-push IMG=<some-registry>/nope:tag
    ```

3. Deploy the controller to the cluster with the image specified by `IMG`:

    ```sh
    make deploy IMG=<some-registry>/nope:tag
    ```

### Uninstall CRDs

To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller

UnDeploy the controller from the cluster:

```sh
make undeploy
```

## Contributing

// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works

This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out

Nope requires a running NATS instance to communicate with, and the NATS instance has to be configured
with an operator, system account and system user.

In [hack/natsserver](./hack/natsserver/readme.md) there is a Go program that runs an embedded NATS server with the necessary bootstrapping done, and writing the necessary files. Cool!

```bash
# In a terminal where you want to run the nats server
cd hack/natsserver
go run main.go

# In a terminal where you want to run the operator
export NATS_MAIN_URL="nats://0.0.0.0:4222"
export OPERATOR_SEED=$(pwd)/hack/natsserver/operator.nk
export NATS_CREDS=$(pwd)/hack/natsserver/sys_user.creds

# Install CRDs and run the operator
make install run

# If you want to test out some samples, first generate the YAML in config/samples/out
make test
# Apply some samples, e.g.
kubectl apply -f config/samples/out/account_sample.yaml
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions

If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
