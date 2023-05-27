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

command -v helm > /dev/null
command -v go > /dev/null
command -v git > /dev/null

ROOT_DIR=$(git rev-parse --show-toplevel)
VALUES_DIR="$ROOT_DIR"/docs/platypus2/scripts
TEMPD="$ROOT_DIR"/out
KYGO="$TEMPD"/kygo

DEBUG=0
pushd "$ROOT_DIR"



# build a version of kygo with all possible CRDs
function tool() {
  pushd $TEMPD > /dev/null
  git clone --depth 1 "https://github.com/veggiemonk/lingonweb"
  popd > /dev/null

  pushd "$TEMPD"/lingonweb > /dev/null
  [ $DEBUG ] && printf  "\n replace github.com/volvo-cars/lingon => ../../ \n" >> go.mod
  go build -o kygo ./cmd/kygo && mv kygo "$TEMPD"
  popd > /dev/null
  [ $DEBUG ] && rm -rf "$TEMPD"/lingonweb

}

function install_repo() {
  helm repo add external-secrets https://charts.external-secrets.io
  helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
  helm repo add aws-ebs-csi-driver https://kubernetes-sigs.github.io/aws-ebs-csi-driver
  helm repo add kube-state-metrics https://kubernetes.github.io/kube-state-metrics
  helm repo add vm https://victoriametrics.github.io/helm-charts/
  helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server
  helm repo add nats https://nats-io.github.io/k8s/helm/charts/
  helm repo add benthos https://benthosdev.github.io/benthos-helm-chart/
  helm repo add autoscaler https://kubernetes.github.io/autoscaler
  helm repo update
}

function manifests() {
  rm -rf "$TEMPD"/manifests
  mkdir -p $TEMPD/manifests && pushd $TEMPD/manifests > /dev/null

  helm template external-secrets external-secrets/external-secrets \
    --namespace=external-secrets --create-namespace --set installCRDs=true | \
    $KYGO -out "externalsecrets" -app external-secrets -pkg externalsecrets

  #
  # MONITORING
  #
  helm template metrics-server metrics-server/metrics-server --namespace=monitoring --values="$VALUES_DIR"/metrics-server.values.yaml | \
    $KYGO -out "monitoring/metrics-server" -app metrics-server -pkg metricsserver

  helm template promcrd prometheus-community/prometheus-operator-crds | \
    $KYGO -out "monitoring/promcrd" -app prometheus -pkg promcrd -group=false -clean-name=false

  helm template kube-promtheus-stack prometheus-community/kube-prometheus-stack --namespace=monitoring | \
    $KYGO -out "monitoring/promstack" -app kube-prometheus-stack -pkg promstack

  helm template vm vm/victoria-metrics-single --namespace=monitoring --values "$VALUES_DIR"/victoriametrics-single.values.yaml | \
    $KYGO -out "monitoring/victoriametrics" -app victoria-metrics -pkg victoriametrics

  #
  # NATS
  #
  helm template nats nats/nats --namespace=nats --values "$VALUES_DIR"/nats.values.yaml | \
    $KYGO -out "nats" -app nats -pkg nats

  helm template surveyor nats/surveyor --namespace=surveyor --values "$VALUES_DIR"/surveyor.values.yaml | \
    $KYGO -out "nats/surveyor" -app surveyor -pkg surveyor

  helm template benthos benthos/benthos --namespace=benthos --values "$VALUES_DIR"/benthos.values.yaml | \
    $KYGO -out "nats/benthos" -app benthos -pkg benthos

  wget https://github.com/nats-io/nack/releases/latest/download/crds.yml -O - | \
    $KYGO -out "nats/jetstream" -app jetstream -pkg jetstream -group=false -clean-name=false

  #
  # Karpenter
  #
  helm template karpenter oci://public.ecr.aws/karpenter/karpenter --namespace=karpenter \
    --create-namespace \
    --version "v0.27.5" \
    --set serviceAccount.annotations."eks\.amazonaws\.com/role-arn"=REPLACE_ME_KARPENTER_IAM_ROLE_ARN \
    --set settings.aws.clusterName=REPLACE_ME_CLUSTER_NAME \
    --set settings.aws.defaultInstanceProfile=KarpenterNodeInstanceProfile-REPLACE_ME_CLUSTER_NAME \
    --set settings.aws.interruptionQueueName=REPLACE_ME_CLUSTER_NAME \
    --set controller.resources.requests.cpu=1 \
    --set controller.resources.requests.memory=1Gi \
    --set controller.resources.limits.cpu=1 \
    --set controller.resources.limits.memory=1Gi \
  | $KYGO -out "karpenter" -app karpenter -pkg karpenter



  wget https://raw.githubusercontent.com/aws/karpenter/main/pkg/apis/crds/karpenter.sh_provisioners.yaml -O - >> karpenter/mani.yaml
  wget https://raw.githubusercontent.com/aws/karpenter/main/pkg/apis/crds/karpenter.sh_machines.yaml -O - >> karpenter/mani.yaml
  wget https://raw.githubusercontent.com/aws/karpenter/main/pkg/apis/crds/karpenter.k8s.aws_awsnodetemplates.yaml -O -  >> karpenter/mani.yaml

  $KYGO -in karpenter/mani.yaml -out "karpenter/crd" -pkg karpentercrd -app karpenter -group=false -clean-name=false

  helm template autoscaler autoscaler/cluster-autoscaler  --set 'autoDiscovery.clusterName'="REPLACE_ME"  | \
    $KYGO -out "autoscaler" -pkg autoscaler -app autoscaler

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