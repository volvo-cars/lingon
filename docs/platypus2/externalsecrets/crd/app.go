// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package crd

import (
	"context"
	"errors"
	"os"
	"os/exec"

	"github.com/volvo-cars/lingon/pkg/kube"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func New() CRD {
	return CRD{
		CrdACRAccessTokensGenerators:        AcraccesstokensGeneratorsCrd,
		CrdClusterExternalSecrets:           ClusterexternalsecretsCrd,
		CrdClusterSecretStores:              ClustersecretstoresCrd,
		CrdECRAuthorizationTokensGenerators: EcrauthorizationtokensGeneratorsCrd,
		CrdExternalSecrets:                  ExternalsecretsCrd,
		CrdFakesGenerators:                  FakesGeneratorsCrd,
		CrdGCRAccessTokensGenerators:        GcraccesstokensGeneratorsCrd,
		CrdPasswordsGenerators:              PasswordsGeneratorsCrd,
		CrdPushSecrets:                      PushsecretsCrd,
		CrdSecretStores:                     SecretstoresCrd,
	}
}

type CRD struct {
	kube.App

	CrdACRAccessTokensGenerators        *apiextv1.CustomResourceDefinition
	CrdClusterExternalSecrets           *apiextv1.CustomResourceDefinition
	CrdClusterSecretStores              *apiextv1.CustomResourceDefinition
	CrdECRAuthorizationTokensGenerators *apiextv1.CustomResourceDefinition
	CrdExternalSecrets                  *apiextv1.CustomResourceDefinition
	CrdFakesGenerators                  *apiextv1.CustomResourceDefinition
	CrdGCRAccessTokensGenerators        *apiextv1.CustomResourceDefinition
	CrdPasswordsGenerators              *apiextv1.CustomResourceDefinition
	CrdPushSecrets                      *apiextv1.CustomResourceDefinition
	CrdSecretStores                     *apiextv1.CustomResourceDefinition
}

// Apply applies the kubernetes objects to the cluster
func (a *CRD) Apply(ctx context.Context) error {
	return Apply(ctx, a)
}

// Export exports the kubernetes objects to YAML files in the given directory
func (a *CRD) Export(dir string) error {
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
