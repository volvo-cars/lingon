// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package promcrd

import (
	"context"
	"errors"
	"os"
	"os/exec"

	"github.com/volvo-cars/lingon/pkg/kube"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// validate the struct implements the interface
var _ kube.Exporter = (*PromCrd)(nil)

const namespace = "monitoring"

// PromCrd contains kubernetes manifests
type PromCrd struct {
	kube.App

	Namespace                                 *corev1.Namespace
	AlertmanagerconfigsMonitoringCoreosComCRD *apiextensionsv1.CustomResourceDefinition
	AlertmanagersMonitoringCoreosComCRD       *apiextensionsv1.CustomResourceDefinition
	PodmonitorsMonitoringCoreosComCRD         *apiextensionsv1.CustomResourceDefinition
	ProbesMonitoringCoreosComCRD              *apiextensionsv1.CustomResourceDefinition
	PrometheusagentsMonitoringCoreosComCRD    *apiextensionsv1.CustomResourceDefinition
	PrometheusesMonitoringCoreosComCRD        *apiextensionsv1.CustomResourceDefinition
	PrometheusrulesMonitoringCoreosComCRD     *apiextensionsv1.CustomResourceDefinition
	ServicemonitorsMonitoringCoreosComCRD     *apiextensionsv1.CustomResourceDefinition
	ThanosrulersMonitoringCoreosComCRD        *apiextensionsv1.CustomResourceDefinition
}

// New creates a new PromCrd
func New() *PromCrd {
	return &PromCrd{
		Namespace: &corev1.Namespace{
			TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Namespace"},
			ObjectMeta: metav1.ObjectMeta{Name: namespace},
		},
		AlertmanagerconfigsMonitoringCoreosComCRD: AlertmanagerconfigsMonitoringCoreosComCRD,
		AlertmanagersMonitoringCoreosComCRD:       AlertmanagersMonitoringCoreosComCRD,
		PodmonitorsMonitoringCoreosComCRD:         PodmonitorsMonitoringCoreosComCRD,
		ProbesMonitoringCoreosComCRD:              ProbesMonitoringCoreosComCRD,
		PrometheusagentsMonitoringCoreosComCRD:    PrometheusagentsMonitoringCoreosComCRD,
		PrometheusesMonitoringCoreosComCRD:        PrometheusesMonitoringCoreosComCRD,
		PrometheusrulesMonitoringCoreosComCRD:     PrometheusrulesMonitoringCoreosComCRD,
		ServicemonitorsMonitoringCoreosComCRD:     ServicemonitorsMonitoringCoreosComCRD,
		ThanosrulersMonitoringCoreosComCRD:        ThanosrulersMonitoringCoreosComCRD,
	}
}

// Apply applies the kubernetes objects to the cluster
func (a *PromCrd) Apply(ctx context.Context) error {
	return Apply(ctx, a)
}

// Export exports the kubernetes objects to YAML files in the given directory
func (a *PromCrd) Export(dir string) error {
	return kube.Export(a, kube.WithExportOutputDirectory(dir))
}

// Apply applies the kubernetes objects contained in Exporter to the cluster
func Apply(ctx context.Context, km kube.Exporter) error {
	cmd := exec.CommandContext(
		ctx,
		"kubectl",
		"apply",
		"--server-side=true",
		"-f",
		"-",
	)
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
