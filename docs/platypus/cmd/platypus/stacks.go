// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"os"
	"os/exec"

	aws "github.com/golingon/terraproviders/aws/4.60.0"
	"github.com/golingon/terraproviders/aws/4.60.0/provider"
	"github.com/volvo-cars/lingon/docs/platypus/pkg/infra/awsvpc"
	"github.com/volvo-cars/lingon/docs/platypus/pkg/infra/cluster_eks"
	"github.com/volvo-cars/lingon/docs/platypus/pkg/platform/grafana"
	"github.com/volvo-cars/lingon/docs/platypus/pkg/platform/karpenter"
	"github.com/volvo-cars/lingon/docs/platypus/pkg/terraclient"

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
	awsvpc.AWSVPC
}

type eksStack struct {
	AWSStackConfig
	cluster_eks.Cluster
}

type rdsSnapshotStack struct {
	AWSStackConfig
	Snapshot *aws.DbSnapshot `validate:"required"`
}

type grafanaStack struct {
	AWSStackConfig
	grafana.RDSPostgres
}

type karpenterStack struct {
	AWSStackConfig
	karpenter.Infra
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

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}
