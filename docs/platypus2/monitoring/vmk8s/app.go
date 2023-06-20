// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package vmk8s

import (
	"context"
	"errors"
	"os"
	"os/exec"

	"github.com/VictoriaMetrics/operator/api/victoriametrics/v1beta1"
	"github.com/volvo-cars/lingon/pkg/kube"
	corev1 "k8s.io/api/core/v1"
)

// validate the struct implements the interface
var _ kube.Exporter = (*Vmk8S)(nil)

// Vmk8S contains kubernetes manifests
type Vmk8S struct {
	kube.App

	Operator         *Operator
	Grafana          *Grafana
	KubeStateMetrics *KubeStateMetrics
	NodeExporter     *NodeExporter
	VMAlertManager   *VMAlertManager

	MonKubeScheduler *MonKubeScheduler
	MonAPIServer     *MonAPIServer
	MonCoreDNS       *MonCoreDNS
	MonETCD          *MonETCD

	K8SRules *K8SRules

	DashboardK8SGlobalCM     *corev1.ConfigMap
	DashboardK8SNamespacesCM *corev1.ConfigMap
	DashboardK8SPodsCM       *corev1.ConfigMap

	CadvisorNodeScrape *v1beta1.VMNodeScrape
	ProbesNodeScrape   *v1beta1.VMNodeScrape
	KubeletNodeScrape  *v1beta1.VMNodeScrape

	VMSingle                   *v1beta1.VMSingle
	DashboardBackupManagerCM   *corev1.ConfigMap
	VMHealthAlertRules         *v1beta1.VMRule
	VMAgent                    *v1beta1.VMAgent
	VMAgentCM                  *corev1.ConfigMap
	VMAgentAlertRules          *v1beta1.VMRule
	VMK8sSA                    *corev1.ServiceAccount
	VMSingleAlertRules         *v1beta1.VMRule
	DashboardVictoriaMetricsCM *corev1.ConfigMap
}

// New creates a new Vmk8S
func New() *Vmk8S {
	return &Vmk8S{
		Operator:         NewOperator(),
		Grafana:          NewGrafana(),
		KubeStateMetrics: NewKubeStateMetrics(),
		NodeExporter:     NewNodeExporter(),
		VMAlertManager:   NewVMAlertManager(),

		MonAPIServer:     NewMonAPIServer(),
		MonKubeScheduler: NewMonKubeScheduler(),
		MonCoreDNS:       NewMonCoreDNS(),
		MonETCD:          NewMonETCD(),

		K8SRules: NewK8SRules(),

		CadvisorNodeScrape: CadvisorNodeScrape,
		ProbesNodeScrape:   ProbesNodeScrape,
		KubeletNodeScrape:  KubeletNodeScrape,

		DashboardK8SGlobalCM:     DashboardK8SGlobalCM,
		DashboardK8SNamespacesCM: DashboardK8SNamespacesCM,
		DashboardK8SPodsCM:       DashboardK8SPodsCM,

		VMK8sSA:                    VMSA,
		DashboardVictoriaMetricsCM: DashboardVictoriaMetricsCM,
		DashboardBackupManagerCM:   DashboardBackupManagerCM,
		VMAgent:                    VMAgent,
		VMAgentCM:                  VMAgentCM,
		VMAgentAlertRules:          VMAgentAlertRules,
		VMSingle:                   VMSingle,
		VMHealthAlertRules:         VMHealthAlertRules,
		VMSingleAlertRules:         VMSingleAlertRules,
	}
}

// Apply applies the kubernetes objects to the cluster
func (a *Vmk8S) Apply(ctx context.Context) error {
	return Apply(ctx, a)
}

// Export exports the kubernetes objects to YAML files in the given directory
func (a *Vmk8S) Export(dir string) error {
	return kube.Export(a, kube.WithExportOutputDirectory(dir))
}

// Apply applies the kubernetes objects contained in Exporter to the cluster
func Apply(ctx context.Context, km kube.Exporter) error {
	cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", "-")
	cmd.Env = os.Environ()        // inherit environment in case we need to use kubectl from a container
	stdin, err := cmd.StdinPipe() // pipe to pass data to kubectl
	if err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	go func() {
		defer func() {
			err = errors.Join(err, stdin.Close())
		}()
		if errEW := kube.Export(
			km,
			kube.WithExportWriter(stdin),
			kube.WithExportAsSingleFile("stdin"),
		); errEW != nil {
			err = errors.Join(err, errEW)
		}
	}()

	if errS := cmd.Start(); errS != nil {
		return errors.Join(err, errS)
	}

	// waits for the command to exit and waits for any copying
	// to stdin or copying from stdout or stderr to complete
	return errors.Join(err, cmd.Wait())
}

// P converts T to *T, useful for basic types
func P[T any](t T) *T {
	return &t
}
