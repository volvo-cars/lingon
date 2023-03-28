// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube_test

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/volvo-cars/lingon/pkg/kube"
)

func ExampleExplode() {
	out := "./out/explode"
	_ = os.RemoveAll(out)
	defer func() {
		_ = os.RemoveAll(out)
	}()
	fp, err := os.Open("./testdata/karpenter.yaml")
	if err != nil {
		panic(fmt.Errorf("open file: %w", err))
	}
	defer func() {
		_ = fp.Close()
	}()

	// explode them into individual files
	if err = kube.Explode(fp, out); err != nil {
		panic(fmt.Errorf("explode manifest files: %w", err))
	}

	got, err := kube.ListYAMLFiles("./out/explode")
	if err != nil {
		panic(fmt.Errorf("list yaml files: %w", err))
	}
	// sort the files to make the output deterministic
	sort.Strings(got)

	for _, f := range got {
		fmt.Println(f)
	}
	// Output:
	//
	// out/explode/_cluster/rbac/1_karpenter-admin_cr.yaml
	// out/explode/_cluster/rbac/1_karpenter-core_cr.yaml
	// out/explode/_cluster/rbac/1_karpenter_cr.yaml
	// out/explode/_cluster/rbac/2_karpenter-core_crb.yaml
	// out/explode/_cluster/rbac/2_karpenter_crb.yaml
	// out/explode/_cluster/webhook/4_defaulting.webhook.karpenter.k8s.aws_mutatingwebhookconfigurations.yaml
	// out/explode/_cluster/webhook/4_defaulting.webhook.karpenter.sh_mutatingwebhookconfigurations.yaml
	// out/explode/_cluster/webhook/4_validation.webhook.config.karpenter.sh_validatingwebhookconfigurations.yaml
	// out/explode/_cluster/webhook/4_validation.webhook.karpenter.k8s.aws_validatingwebhookconfigurations.yaml
	// out/explode/_cluster/webhook/4_validation.webhook.karpenter.sh_validatingwebhookconfigurations.yaml
	// out/explode/karpenter/1_karpenter_role.yaml
	// out/explode/karpenter/1_karpenter_sa.yaml
	// out/explode/karpenter/1_karpenter_svc.yaml
	// out/explode/karpenter/2_config-logging_cm.yaml
	// out/explode/karpenter/2_karpenter-cert_secrets.yaml
	// out/explode/karpenter/2_karpenter-global-settings_cm.yaml
	// out/explode/karpenter/2_karpenter_rb.yaml
	// out/explode/karpenter/3_karpenter_deploy.yaml
	// out/explode/karpenter/4_karpenter_pdb.yaml
	// out/explode/kube-system/1_karpenter-dns_role.yaml
	// out/explode/kube-system/2_karpenter-dns_rb.yaml
}

func ExampleExplode_reader() {
	out := "./out/explodereader"
	_ = os.RemoveAll(out)
	defer func() {
		_ = os.RemoveAll(out)
	}()
	manifest, err := kube.ReadManifest("./testdata/golden/reader.yaml")
	if err != nil {
		panic(fmt.Errorf("read manifest files: %w", err))
	}
	fakeManifest := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: fake	
`
	manifest = append(manifest, fakeManifest)
	// join the manifest into a single string
	joined := strings.Join(manifest, "\n---\n")
	r := strings.NewReader(joined)

	// explode them into individual files
	// this will fail because the reader.yaml manifest is not valid
	_ = kube.Explode(r, out)

	// but we can still list the files already processed
	got, err := kube.ListYAMLFiles(out)
	if err != nil {
		panic(fmt.Errorf("list yaml files: %w", err))
	}

	for _, f := range got {
		fmt.Println(f)
	}
	// Output:
	// out/explodereader/default/3_webapp_deploy.yaml
}
