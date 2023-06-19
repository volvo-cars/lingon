// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package nats

import (
	"context"
	"errors"
	"os"
	"os/exec"

	promoperatorv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/volvo-cars/lingon/pkg/kube"
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	"github.com/volvo-cars/lingoneks/nats/jetstream"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// validate the struct implements the interface
var _ kube.Exporter = (*Nats)(nil)

// Nats contains kubernetes manifests
type Nats struct {
	kube.App

	// embedded struct
	CRD
	Surveyor

	NS *corev1.Namespace

	BoxDeploy          *appsv1.Deployment
	ConfigCM           *corev1.ConfigMap
	PDB                *policyv1.PodDisruptionBudget
	SA                 *corev1.ServiceAccount
	STS                *appsv1.StatefulSet
	SVC                *corev1.Service
	ServiceMonitor     *promoperatorv1.ServiceMonitor
	TestRequestReplyPO *corev1.Pod
	// DashboardCM        *corev1.ConfigMap
}

// New creates a new Nats
func New() *Nats {
	return &Nats{
		CRD: CRD{
			AccountsNatsIoCRD:        jetstream.AccountsJetstreamNatsIoCRD,
			ConsumersNatsIoCRD:       jetstream.ConsumersJetstreamNatsIoCRD,
			StreamsNatsIoCRD:         jetstream.StreamsJetstreamNatsIoCRD,
			StreamtemplatesNatsIoCRD: jetstream.StreamTemplatesJetstreamNatsIoCRD,
		},
		Surveyor:           *NewSurveyor(),
		NS:                 NS,
		BoxDeploy:          BoxDeploy,
		ConfigCM:           cm.ConfigMap(),
		PDB:                PDB,
		SA:                 SA,
		STS:                STS,
		SVC:                SVC,
		ServiceMonitor:     ServiceMonitor,
		TestRequestReplyPO: TestRequestReplyPO,
		// DashboardCM:        DashboardNatsCM,
	}
}

type CRD struct {
	AccountsNatsIoCRD        *apiextensionsv1.CustomResourceDefinition
	ConsumersNatsIoCRD       *apiextensionsv1.CustomResourceDefinition
	StreamsNatsIoCRD         *apiextensionsv1.CustomResourceDefinition
	StreamtemplatesNatsIoCRD *apiextensionsv1.CustomResourceDefinition
}

const (
	appName           = "nats"
	defaultContainer  = appName
	namespace         = "nats"
	version           = "2.9.16"
	replicas          = 3
	storageClass      = "gp2"
	ImgNats           = "nats:" + version + "-alpine"
	sidecarVersion    = "0.10.1"
	ImgConfigReloader = "natsio/nats-server-config-reloader:" + sidecarVersion
	ImgPromExporter   = "natsio/prometheus-nats-exporter:" + sidecarVersion
)

var (
	NS = ku.Namespace(namespace, BaseLabels(), nil)
	SA = ku.ServiceAccount(appName, namespace, BaseLabels(), nil)
)

var matchLabels = map[string]string{
	ku.AppLabelName:     appName,
	ku.AppLabelInstance: appName,
}

func BaseLabels() map[string]string {
	return ku.MergeLabels(
		matchLabels, map[string]string{
			"app":                appName,
			ku.AppLabelComponent: appName,
			ku.AppLabelPartOf:    appName,
			ku.AppLabelVersion:   version,
			ku.AppLabelManagedBy: "lingon",
		},
	)
}

// Apply applies the kubernetes objects to the cluster
func (a *Nats) Apply(ctx context.Context) error {
	return Apply(ctx, a)
}

// Export exports the kubernetes objects to YAML files in the given directory
func (a *Nats) Export(dir string) error {
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
