#!/usr/bin/env bash
# Copyright (c) 2023 Volvo Car Corporation
# SPDX-License-Identifier: Apache-2.0


## HELM in reverse is MLEH

echo
echo '███    ███ ██      ███████ ██   ██ '
echo '████  ████ ██      ██      ██   ██ '
echo '██ ████ ██ ██      █████   ███████ '
echo '██  ██  ██ ██      ██      ██   ██ '
echo '██      ██ ███████ ███████ ██   ██ '
echo '                                   '


set -exuo pipefail

ROOT_DIR=$(git rev-parse --show-toplevel)
VALUES_DIR="$ROOT_DIR"/docs/platypus2/scripts
TEMPD="$ROOT_DIR"/out
KYGO="$TEMPD"/kygo

pushd "$ROOT_DIR"

command -v helm > /dev/null
command -v go > /dev/null
command -v git > /dev/null


# build a version of kygo with all possible CRDs
function tool() {
  pushd $TEMPD > /dev/null
  git clone --depth 1 "https://github.com/veggiemonk/lingonweb"
  popd > /dev/null

  pushd "$TEMPD"/lingonweb > /dev/null
  go build -o kygo ./cmd/kygo && mv kygo "$TEMPD"
  popd > /dev/null
  rm -rf "$TEMPD"/lingonweb

}

function install_repo() {
  helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
  helm repo add aws-ebs-csi-driver https://kubernetes-sigs.github.io/aws-ebs-csi-driver
  helm repo add kube-state-metrics https://kubernetes.github.io/kube-state-metrics
  helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server
  helm repo add nats https://nats-io.github.io/k8s/helm/charts/
  helm repo update
}

function manifests() {
  mkdir -p $TEMPD
  pushd $TEMPD > /dev/null

  rm -rf metrics-server
  helm template metrics-server metrics-server/metrics-server --values "$VALUES_DIR"/metrics-server.values.yaml | \
    $KYGO -out metrics-server -app matrics-server -pkg metricsserver

  rm -rf promcrd
  helm template promcrd prometheus-community/prometheus-operator-crds | \
    $KYGO -out promcrd -app prometheus -pkg promcrd -group=false

  rm -rf promstack
  helm template promstack prometheus-community/kube-prometheus-stack --namespace=monitoring | \
    $KYGO -out promstack -app kube-prometheus-stack -pkg promstack

  rm -rf nats
  helm template nats nats/nats --namespace=nats --values "$VALUES_DIR"/nats.values.yaml | \
    $KYGO -out nats -app nats -pkg nats

  popd
}

function step() {
  set +x
  local name="$1"
  echo
  echo '   #' "$name"
  echo '   ======================'
  echo
  set -x
}

function main {

  mkdir -p "$TEMPD"

  step "build kygo"
  [ ! -f "$KYGO" ] && tool

  step "install/update repo"
  install_repo || true

  step "generate manifests"
  manifests

}

main