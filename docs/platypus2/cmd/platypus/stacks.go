// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"os"
	"os/exec"

	aws "github.com/golingon/terraproviders/aws/5.6.2"
	"github.com/golingon/terraproviders/aws/5.6.2/provider"
	"github.com/volvo-cars/lingoneks/infra"
	"github.com/volvo-cars/lingoneks/karpenter"
	"github.com/volvo-cars/lingoneks/terraclient"

	"github.com/volvo-cars/lingon/pkg/terra"
)

func newAWSStackConfig(name string, p runParams) AWSStackConfig {
	return AWSStackConfig{
		Stack: terraclient.Stack{
			Name: name,
		},
		Backend:  newBackend(p.AWSParams, name),
		Provider: newProv(p.AWSParams, p.TFLabels),
	}
}

type AWSStackConfig struct {
	terraclient.Stack
	Backend  *backendS3    `validate:"required"`
	Provider *aws.Provider `validate:"required"`
}

type vpcStack struct {
	AWSStackConfig
	infra.AWSVPC
}

type eksStack struct {
	AWSStackConfig
	infra.Cluster
}

type karpenterStack struct {
	AWSStackConfig
	karpenter.Infra
}
type csiEbsStack struct {
	AWSStackConfig
	infra.CSI
}

func newBackend(p AWSParams, stateFile string) *backendS3 {
	return &backendS3{
		Bucket:  p.BackendS3Key,
		Key:     stateFile,
		Profile: p.Profile,
		Region:  p.Region,
	}
}

var _ terra.Backend = (*backendS3)(nil)

type backendS3 struct {
	Bucket  string `hcl:"bucket"`
	Key     string `hcl:"key"`
	Profile string `hcl:"profile"`
	Region  string `hcl:"region"`
}

func (b *backendS3) BackendType() string {
	return "s3"
}

func newProv(p AWSParams, labels map[string]string) *aws.Provider {
	l := make(map[string]terra.StringValue, len(labels))
	for k, v := range labels {
		l[k] = S(v)
	}

	return aws.NewProvider(
		aws.ProviderArgs{
			Profile: S(p.Profile),
			Region:  S(p.Region),
			DefaultTags: []provider.DefaultTags{
				{
					Tags: terra.Map(l),
				},
			},
		},
	)
}

func kubeconfigFromAWSCmd(
	ctx context.Context,
	profile string,
	clusterName, region string,
	kubeconfigPath string,
) error {
	cmd := exec.CommandContext(
		ctx,
		"aws",
		"--profile",
		profile,
		"eks",
		"update-kubeconfig",
		"--name",
		clusterName,
		"--kubeconfig",
		kubeconfigPath,
		"--alias",
		clusterName,
		"--region",
		region,
	)

	cmd.Env = os.Environ()

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
