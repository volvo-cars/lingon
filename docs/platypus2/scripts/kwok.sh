#!/usr/bin/env bash
# Copyright (c) 2023 Volvo Car Corporation
# SPDX-License-Identifier: Apache-2.0


## KWOK is kubernetes without kubelet
## a fake cluster, perfect to test manifest

echo
echo '██  ██ ██     ██  ██████  ██   ██ '
echo '█  ██  ██     ██ ██    ██ ██  ██   '
echo '████   ██  █  ██ ██    ██ █████    '
echo '█  ██  ██ ███ ██ ██    ██ ██  ██   '
echo '█   ██  ███ ███   ██████  ██   ██  '
echo

set -exuo pipefail

command -v kwok > /dev/null || { echo "need to install kwok -> brew install kwok "; exit 1 ; }
command -v go > /dev/null
command -v git > /dev/null

ROOT_DIR=$(git rev-parse --show-toplevel)
TMP="$ROOT_DIR"/out/kwok
KUBE_VERSION="v1.27.2"
KUBE_BIN="$TMP/bin"
DEBUG=0

pushd "$ROOT_DIR"


# tool downloads and builds kubernetes binaries
# cache at every stage to speed up the process
function tool {
    pushd "$TMP" > /dev/null
    mkdir -p "$TMP/kubernetes" && pushd "$TMP/kubernetes" > /dev/null
    PREFIX_KUBEBIN="$TMP/kubernetes/_output/local/bin/$(go env GOOS)/$(go env GOARCH)/"

    # download if not there
    [ ! -f "$TMP/kubernetes/README.md" ] && \
      wget https://dl.k8s.io/"$KUBE_VERSION"/kubernetes-src.tar.gz -O - | tar xz

    # build if not there
    [ ! -f "$PREFIX_KUBEBIN/kube-apiserver" ] && make WHAT=cmd/kube-apiserver
    [ ! -f "$PREFIX_KUBEBIN/kube-controller-manager" ] && make WHAT=cmd/kube-controller-manager
    [ ! -f "$PREFIX_KUBEBIN/kube-scheduler" ] && make WHAT=cmd/kube-scheduler

    # copy binaries
    mkdir -p "$TMP/bin"
    cp "$PREFIX_KUBEBIN/kube-apiserver" "$TMP/bin/"
    cp "$PREFIX_KUBEBIN/kube-controller-manager" "$TMP/bin/"
    cp "$PREFIX_KUBEBIN/kube-scheduler" "$TMP/bin/"
    [ $DEBUG ] && rm -rf "$TMP/kubernetes"

    popd > /dev/null
    popd > /dev/null
}

function fake_node {
    NODE_FILE=node.yaml

    rm -rf "$NODE_FILE"

    for i in $(seq 1 10);
    do
        export NODE="node-$i"
        cat << EOH >> $NODE_FILE
---
apiVersion: v1
kind: Node
metadata:
  annotations:
    kwok.x-k8s.io/node: fake
    node.alpha.kubernetes.io/ttl: "0"
  labels:
    beta.kubernetes.io/arch: arm64
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

}

function main {
  mkdir -p "$TMP" && pushd "$TMP" > /dev/null

  [ ! -f "$KUBE_BIN"/kube-apiserver ] && tool

  kwokctl create cluster \
          --name fake \
          --runtime binary \
          --kube-admission \
          --kube-authorization \
          --kubeconfig "$TMP"/kubeconfig \
          --kube-controller-manager-binary "$KUBE_BIN"/kube-controller-manager \
          --kube-apiserver-binary "$KUBE_BIN"/kube-apiserver \
          --kube-scheduler-binary "$KUBE_BIN"/kube-scheduler

  pushd "$TMP" > /dev/null

  set +x
  echo " create fake nodes"
  fake_node
  set -x

  kubectl  --context kwok-fake --kubeconfig "$TMP"/kubeconfig apply -f "$TMP"/node.yaml

  set +x
  echo
  echo "To access the fake cluster:"
  echo
  echo "    export KUBECONFIG=$TMP/kubeconfig"
  echo
  echo

}

main
