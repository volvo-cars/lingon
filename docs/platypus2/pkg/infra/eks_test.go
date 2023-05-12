// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package infra

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/volvo-cars/lingon/pkg/terra"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
	aws "github.com/volvo-cars/lingoneks/providers/aws/4.66.1"
)

func testExportValidateStack(
	t *testing.T, ctx context.Context,
	stack terra.Exporter,
) {
	workdir := filepath.Join(os.TempDir(), t.Name())
	if err := os.RemoveAll(workdir); err != nil {
		t.Errorf(
			"failed removing temporary work dir %s: %s",
			workdir,
			err.Error(),
		)
		return
	}
	if err := os.MkdirAll(workdir, os.ModePerm); err != nil {
		t.Errorf(
			"failed creating temporary work dir %s: %s",
			workdir, err.Error(),
		)
	}
	file, err := os.CreateTemp(workdir, "main_*.tf")
	if err != nil {
		t.Errorf(
			"failed creating temporary file in %s: %s",
			workdir,
			err.Error(),
		)
		return
	}
	defer os.Remove(file.Name()) // clean up

	if err := terra.Export(stack, terra.WithExportWriter(file)); err != nil {
		t.Errorf(
			"failed exporting stack to %s: %s",
			file.Name(),
			err.Error(),
		)
		return
	}
	tf, err := tfexec.NewTerraform(workdir, "terraform")
	tu.AssertNoError(t, err, "creating terraform runtime")
	if err := tf.Init(ctx); err != nil {
		tu.AssertNoError(t, err, "initialising terraform config")
	}
	tfValidate, err := tf.Validate(ctx)
	tu.AssertNoError(t, err, "validating terraform config")
	tu.AssertEqual(t, 0, len(tfValidate.Diagnostics))
	for _, diag := range tfValidate.Diagnostics {
		t.Log(diag.Summary)
	}
}

func TestEKS(t *testing.T) {
	type awsStack struct {
		terra.Stack

		Provider *aws.Provider
		Cluster  `validate:"required"`
	}
	eks := NewCluster(
		ClusterOpts{
			Name:    "test",
			Version: "1.24",
			VPCID:   "123456",
			PrivateSubnetIDs: [3]string{
				"a", "b", "c",
			},
		},
	)
	stack := awsStack{
		Provider: aws.NewProvider(aws.ProviderArgs{}),
		Cluster:  *eks,
	}
	ctx := context.Background()
	testExportValidateStack(t, ctx, &stack)
}
