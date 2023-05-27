// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	aws "github.com/golingon/terraproviders/aws/5.0.1"
	"github.com/golingon/terraproviders/aws/5.0.1/provider"
	"github.com/volvo-cars/lingon/pkg/terra"
	"github.com/volvo-cars/lingoneks/pkg/infra"
	"github.com/volvo-cars/lingoneks/pkg/terraclient"
	"golang.org/x/exp/slog"
)

const (
	name   = "platypusbootstrap"
	region = "eu-north-1"
)

var S = terra.String

type runParams struct {
	AWSParams AWSParams
	TFLabels  map[string]string
	KLabels   map[string]string
	Apply     bool
	Plan      bool
}
type AWSParams struct {
	Region  string
	Profile string
}

type StackConfig struct {
	terraclient.Stack
	Provider *aws.Provider
}

type s3Stack struct {
	StackConfig
	infra.Bucket
}

func main() {
	var apply, plan bool
	var profile string

	flag.BoolVar(
		&apply,
		"apply",
		false,
		"Apply the terraform changes (default: false)",
	)
	flag.BoolVar(
		&plan,
		"plan",
		false,
		"Plan the terraform changes (default: false)",
	)
	flag.StringVar(
		&profile,
		"profile",
		"",
		"name of the aws profile in ~/.aws/config (default: none)",
	)
	flag.Parse()

	if profile == "" {
		slog.Error("no profile defined")
		return
	}
	ap := AWSParams{
		// BackendS3Key: "lingon-tf-experiment",
		Region:  region,
		Profile: profile,
	}
	p := runParams{
		Apply:     apply,
		Plan:      plan,
		AWSParams: ap,
		TFLabels: map[string]string{
			infra.TagEnv: "dev",
			"terraform":  "true",
		},
		KLabels: map[string]string{
			infra.TagEnv: "dev",
		},
	}

	if err := run(p); err != nil {
		slog.Error("run", "err", err)
		os.Exit(1)
	}
	slog.Info("done")
}

func StepSep(name string) {
	fmt.Printf("\n\n> %s  \n =====================\n\n", name)
}

func run(p runParams) error {
	slog.Info("run", "params", p)
	ctx := context.Background()

	tf := terraclient.NewClient(
		terraclient.WithDefaultPlan(p.Plan),
		terraclient.WithDefaultApply(p.Apply),
	)

	StepSep("bucket")

	bucketName := name + "-lingon"
	if err := infra.ValidateName(bucketName); err != nil {
		return err
	}
	slog.Info(
		"bucket stack",
		slog.String("name", bucketName),
	)
	s3 := s3Stack{
		StackConfig: StackConfig{
			Stack: terraclient.Stack{Name: bucketName},
			Provider: aws.NewProvider(
				aws.ProviderArgs{
					Profile: S(p.AWSParams.Profile),
					Region:  S(p.AWSParams.Region),
					DefaultTags: []provider.DefaultTags{
						{
							Tags: infra.Ttags(p.TFLabels),
						},
					},
				},
			),
		},
		Bucket: *infra.NewBucket(bucketName),
	}

	if err := tf.Run(ctx, &s3); err != nil {
		return fmt.Errorf("tfrun: bucket %w", err)
	}

	if !s3.IsStateComplete() {
		slog.Info("VPC state not in sync, finishing here. Is it Applied ?")
		return fmt.Errorf("incomplete")
	}

	return nil
}
