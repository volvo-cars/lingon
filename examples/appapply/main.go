// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package main

//go:generate rm -rf out
//go:generate go run github.com/volvo-cars/lingon/cmd/explode -in ../../pkg/kube/testdata/grafana.yaml -out out/explode
//go:generate -command kygo go run github.com/volvo-cars/lingon/cmd/kygo
//go:generate kygo -in ../../pkg/kube/testdata/grafana.yaml -out out/grafana -app grafana -group=false -clean-name=false
//go:generate kygo -in ../../pkg/kube/testdata/grafana.yaml -out out/grafanagrouped -app grafana -group -clean-name
//go:generate kygo -in ../../pkg/kube/testdata/external-secrets.yaml -out out/external-secrets -app external-secrets -pkg externalsecrets -group -clean-name
//go:generate kygo -in ../../pkg/kube/testdata/argocd.yaml -out out/argocd -app argocd -group -clean-name
//go:generate kygo -in ../../pkg/kube/testdata/cilium.yaml -out out/cilium -app cilium -group -clean-name
//go:generate kygo -in ../../pkg/kube/testdata/karpenter.yaml -out out/karpenter -app karpenter -group -clean-name
//go:generate kygo -in ../../pkg/kube/testdata/tekton.yaml -out out/tekton -app tekton -group -clean-name

// Need extra CRDs for the following:
// go:generate kygo -in ../../pkg/kube/testdata/istio.yaml -out out/istio.yaml -app istio -group -clean-name
// go:generate kygo -in ../../pkg/kube/testdata/spark.yaml -out out/spark.yaml -app spark -group -clean-name

import (
	"fmt"

	"appapply/out/argocd"
	"appapply/out/cilium"
	extsecrets "appapply/out/external-secrets"
	"appapply/out/grafana"
	"appapply/out/karpenter"
	"appapply/out/tekton"
)

func main() {
	if err := grafana.New().Export("out/export/grafana"); err != nil {
		fmt.Println(err)
		return
	}
	if err := extsecrets.New().Export("out/export/extsecrets"); err != nil {
		fmt.Println(err)
		return
	}
	if err := argocd.New().Export("out/export/argocd"); err != nil {
		fmt.Println(err)
		return
	}
	if err := cilium.New().Export("out/export/cilium"); err != nil {
		fmt.Println(err)
		return
	}
	if err := karpenter.New().Export("out/export/karpenter"); err != nil {
		fmt.Println(err)
		return
	}
	if err := tekton.New().Export("out/export/tekton"); err != nil {
		fmt.Println(err)
		return
	}
}
