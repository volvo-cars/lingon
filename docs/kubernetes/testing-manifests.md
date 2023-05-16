# Testing Manifests

## KWOK

or Kubernetes without Kubelet

```shell
# go to https://kwok.sigs.k8s.io/docs/user/install/
# or just 
brew install kwok

```

### Build kubernetes binaries

> If you want to use docker, you can skip this step

The kubernetes SIG release team issues only linux binaries.
Using Docker solve that issue for non-Linux system (OSX and Windows).
Here, we outline how to compile kubernetes binaries to avoid running docker.

```shell
TMP=~/tmp
KUBE_VERSION="v1.27.1"

mkdir -p "$TMP" && cd "$TMP"

wget https://dl.k8s.io/"$KUBE_VERSION"/kubernetes-src.tar.gz -O - | tar xz
#
# wget https://dl.k8s.io/"$KUBE_VERSION"/kubernetes-src.tar.gz && \
# tar xzf kubernetes-src.tar.gz

make WHAT=cmd/kube-apiserver
make WHAT=cmd/kube-controller-manager
make WHAT=cmd/kube-scheduler


export KUBE_BIN="$TMP/_output/local/bin/$(go env GOOS)/$(go env GOARCH)"
```

## Start a cluster

### With binaries

```shell
# get the path of the binaries (from previous step)
# KUBE_BIN="$TMP/_output/local/bin/$(go env GOOS)/$(go env GOARCH)"
  
kwokctl create cluster \
   --name fake \
   --runtime binary \
   --kube-admission \
   --kube-authorization \
   --kubeconfig "$TMP"/kubeconfig \
   --kube-controller-manager-binary "$KUBE_BIN"/kube-controller-manager \
   --kube-apiserver-binary "$KUBE_BIN"/kube-apiserver \
   --kube-scheduler-binary "$KUBE_BIN"/kube-scheduler

```

### With docker

```shell
kwokctl create cluster \
   --name fake \
   --kube-admission \
   --kube-authorization \
   --kubeconfig "$TMP"/kubeconfig
```

### Set kubeconfig and context

In order to avoid setting `--kubeconfig` and `--context` at every kubectl command

```shell
export KUBECONFIG=$TMP/kubeconfig"
```

### Create fake nodes

Feel free to specify taints on nodes.
For more information, see [kwok](https://github.com/kubernetes-sigs/kwok/blob/main/test/kwok/fake-node.yaml)

```shell
for i in $(seq 1 10);
do
  export NODE="node-$i"
  cat << EOH >> node.yaml
---
apiVersion: v1
kind: Node
metadata:
  annotations:
    kwok.x-k8s.io/node: fake
    node.alpha.kubernetes.io/ttl: "0"
  labels:
    beta.kubernetes.io/arch: amd64
    beta.kubernetes.io/os: linux
    kubernetes.io/arch: arm64
    kubernetes.io/hostname: ${NODE}
    kubernetes.io/os: linux
    kubernetes.io/role: agent
    node-role.kubernetes.io/agent: ""
    type: kwok-controller
  name: ${NODE}
spec:
EOH

done

kubectl  --context kwok-fake --kubeconfig "$TMP"/kubeconfig apply -f "$TMP"/node.yaml
```

### Deploy manifests to fake cluster

> NOTE: the CRDs must be deployed to the cluster before testing custom resources.

#### Create a test

```go
package app_test

import (
    "testing"
    "example.com/example/app"
    "github.com/volvo-cars/lingon/pkg/kube"
)
func TestAppExport(t *testing.T) {
    if err := kube.Export(app.New() /*, kube.WithExportOptions(...) */); err != nil {
        t.Error(err)
    }
}
```

#### Run test

This command will generate the manifests as a file and, with the help of `kubectl`,
deploys it to the `kwok-fake` cluster.

```shell
go test ./app && \
kubectl --context kwok-fake --kubeconfig "$TMP"/kubeconfig \
  apply -f ./app/out
```

### Cleanup

```shell
kwokctl delete cluster --name fake
```

## Kubescore

[Kubescore](https://github.com/zegl/kube-score)

Export the manifest to YAML and run kubescore on them.

To get a score programmatically, there is an [example here](../../pkg/testutil/example_score_test.go).
