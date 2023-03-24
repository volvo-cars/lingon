// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package kube_test

import (
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/volvo-cars/lingon/pkg/kube"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestExplode(t *testing.T) {
	f := "./testdata/karpenter.yaml"
	fp, err := os.Open(f)
	tu.AssertNoError(t, err, "failed to open file")
	out := "./out/explode"
	tu.AssertNoError(t, os.RemoveAll(out), "failed to remove out dir")
	defer os.RemoveAll(out)

	// explode them into individual files
	err = kube.Explode(fp, out)
	tu.AssertNoError(t, err, "failed to explode")
	tu.AssertNoError(t, fp.Close(), "failed to close")

	got, err := kube.ListYAMLFiles("./out/explode")
	tu.AssertNoError(t, err, "failed to list go files")
	sort.Strings(got)

	want := []string{
		"out/explode/_cluster/rbac/1_karpenter-admin_cr.yaml",
		"out/explode/_cluster/rbac/1_karpenter-core_cr.yaml",
		"out/explode/_cluster/rbac/1_karpenter_cr.yaml",
		"out/explode/_cluster/rbac/2_karpenter-core_crb.yaml",
		"out/explode/_cluster/rbac/2_karpenter_crb.yaml",
		"out/explode/_cluster/webhook/4_defaulting.webhook.karpenter.k8s.aws_mutatingwebhookconfigurations.yaml",
		"out/explode/_cluster/webhook/4_defaulting.webhook.karpenter.sh_mutatingwebhookconfigurations.yaml",
		"out/explode/_cluster/webhook/4_validation.webhook.config.karpenter.sh_validatingwebhookconfigurations.yaml",
		"out/explode/_cluster/webhook/4_validation.webhook.karpenter.k8s.aws_validatingwebhookconfigurations.yaml",
		"out/explode/_cluster/webhook/4_validation.webhook.karpenter.sh_validatingwebhookconfigurations.yaml",
		"out/explode/karpenter/1_karpenter_role.yaml",
		"out/explode/karpenter/1_karpenter_sa.yaml",
		"out/explode/karpenter/1_karpenter_svc.yaml",
		"out/explode/karpenter/2_config-logging_cm.yaml",
		"out/explode/karpenter/2_karpenter-cert_secrets.yaml",
		"out/explode/karpenter/2_karpenter-global-settings_cm.yaml",
		"out/explode/karpenter/2_karpenter_rb.yaml",
		"out/explode/karpenter/3_karpenter_deploy.yaml",
		"out/explode/karpenter/4_karpenter_pdb.yaml",
		"out/explode/kube-system/1_karpenter-dns_role.yaml",
		"out/explode/kube-system/2_karpenter-dns_rb.yaml",
	}
	if diff := tu.Diff(got, want); diff != "" {
		t.Error(tu.Callers(), diff)
	}
}

func TestExplodeReader(t *testing.T) {
	f := "./testdata/reader.yaml"
	out := "./out/explodereader"
	err := os.RemoveAll(out)
	tu.AssertNoError(t, err, "failed to remove out dir")
	defer os.RemoveAll(out)

	manifest, err := kube.ReadManifest(f)
	tu.AssertNoError(t, err, "failed to read manifest")
	fakeManifest := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: fake	
`
	manifest = append(manifest, fakeManifest)
	joined := strings.Join(manifest, "\n---\n")
	r := strings.NewReader(joined)

	// explode them into individual files
	err = kube.Explode(r, out)
	tu.NotNil(t, err)

	got, err := kube.ListYAMLFiles(out)
	tu.AssertNoError(t, err, "failed to list go files")

	want := []string{
		"out/explodereader/default/3_webapp_deploy.yaml",
	}
	if diff := tu.Diff(got, want); diff != "" {
		t.Error(tu.Callers(), diff)
	}
}
