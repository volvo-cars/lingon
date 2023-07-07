// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package karpenter

import (
	"os"
	"testing"

	karpentercore "github.com/aws/karpenter-core/pkg/apis"
	karpenter "github.com/aws/karpenter/pkg/apis"
	promoperatorv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/volvo-cars/lingon/pkg/kube"
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
	karpentercrd "github.com/volvo-cars/lingoneks/karpenter/crd"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestExport(t *testing.T) {
	_ = os.RemoveAll("out")

	app := New(
		Opts{
			ClusterName:            "REPLACE_ME_CLUSTER_NAME",
			ClusterEndpoint:        "REPLACE_ME_CLUSTER_ENDPOINT",
			IAMRoleArn:             "REPLACE_ME_ROLE_ARN",
			DefaultInstanceProfile: "REPLACE_ME_DEFAULT_INSTANCE_PROFILE",
			InterruptQueue:         "REPLACE_ME_INTERRUPT_QUEUE",
		},
	)

	tu.AssertNoError(t, kube.Export(app, kube.WithExportOutputDirectory("out")))

	tu.AssertNoError(
		t,
		kube.Export(karpentercrd.New(), kube.WithExportOutputDirectory("out")),
		"karpenter crd export",
	)

	tu.AssertNoError(
		t, kube.Export(
			NewProvisioners(
				ProvisionersOpts{
					ClusterName:       "myclustername",
					AvailabilityZones: [3]string{"AZ1", "AZ2", "AZ3"},
				},
			), kube.WithExportOutputDirectory("out"),
		),
	)

	ly, err := ku.ListYAMLFiles("out")
	tu.AssertNoError(t, err, "list yaml files")

	defaultSerz := func() runtime.Decoder {
		utilruntime.Must(promoperatorv1.AddToScheme(scheme.Scheme))
		utilruntime.Must(karpentercore.AddToScheme(scheme.Scheme))
		utilruntime.Must(karpenter.AddToScheme(scheme.Scheme))
		utilruntime.Must(apiextensions.AddToScheme(scheme.Scheme))
		return scheme.Codecs.UniversalDeserializer()
	}

	tu.AssertNoError(
		t, kube.Import(
			kube.WithImportAppName("karpenter"),
			kube.WithImportManifestFiles(ly),
			kube.WithImportPackageName("karpenter"),
			kube.WithImportSerializer(defaultSerz()),
			kube.WithImportRemoveAppName(true),
			kube.WithImportOutputDirectory("out/go"),
		), "import",
	)
}
