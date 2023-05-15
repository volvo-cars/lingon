// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package testutil_test

import (
	"sort"
	"testing"

	"github.com/rogpeppe/go-internal/txtar"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestFolder2Txtar(t *testing.T) {
	ar, err := tu.Folder2Txtar("../kube/testdata")
	tu.AssertNoError(t, err)
	tu.IsNotEqual(t, 0, len(ar.Files))
	want := []string{
		"../kube/testdata/argocd.yaml",
		"../kube/testdata/cilium.yaml",
		"../kube/testdata/external-secrets.yaml",
		"../kube/testdata/go/tekton/app.go",
		"../kube/testdata/go/tekton/app_test.go",
		"../kube/testdata/go/tekton/cluster-role-binding.go",
		"../kube/testdata/go/tekton/cluster-role.go",
		"../kube/testdata/go/tekton/config-map.go",
		"../kube/testdata/go/tekton/custom-resource-definition.go",
		"../kube/testdata/go/tekton/deployment.go",
		"../kube/testdata/go/tekton/horizontal-pod-autoscaler.go",
		"../kube/testdata/go/tekton/mutating-webhook-configuration.go",
		"../kube/testdata/go/tekton/namespace.go",
		"../kube/testdata/go/tekton/role-binding.go",
		"../kube/testdata/go/tekton/role.go",
		"../kube/testdata/go/tekton/secret.go",
		"../kube/testdata/go/tekton/service-account.go",
		"../kube/testdata/go/tekton/service.go",
		"../kube/testdata/go/tekton/validating-webhook-configuration.go",
		"../kube/testdata/golden/cm-comment.golden",
		"../kube/testdata/golden/cm-comment.yaml",
		"../kube/testdata/golden/configmap.golden",
		"../kube/testdata/golden/configmap.yaml",
		"../kube/testdata/golden/deployment.golden",
		"../kube/testdata/golden/deployment.yaml",
		"../kube/testdata/golden/dirty.txt",
		"../kube/testdata/golden/dirty.yaml",
		"../kube/testdata/golden/export_embedded_struct.txt",
		"../kube/testdata/golden/export_embedded_struct_explode.txt",
		"../kube/testdata/golden/export_embedded_struct_with_explode_and_name_file_func.txt",
		"../kube/testdata/golden/export_embedded_struct_with_explode_and_name_file_func_as_JSON.txt",
		"../kube/testdata/golden/export_embedded_struct_with_name_file_func.txt",
		"../kube/testdata/golden/export_remove_secrets.txt",
		"../kube/testdata/golden/export_tekton.txt",
		"../kube/testdata/golden/empty.golden",
		"../kube/testdata/golden/empty.yaml",
		"../kube/testdata/golden/encode.txt",
		"../kube/testdata/golden/log.golden",
		"../kube/testdata/golden/reader.yaml",
		"../kube/testdata/golden/secret.golden",
		"../kube/testdata/golden/secret.yaml",
		"../kube/testdata/golden/service.golden",
		"../kube/testdata/golden/service.yaml",
		"../kube/testdata/grafana.yaml",
		"../kube/testdata/istio.yaml",
		"../kube/testdata/karpenter.yaml",
		"../kube/testdata/spark.yaml",
		"../kube/testdata/tekton-updated.yaml",
		"../kube/testdata/tekton.yaml",
	}
	filenames := make([]string, 0, len(ar.Files))
	for _, f := range ar.Files {
		filenames = append(filenames, f.Name)
	}
	sort.Strings(want)
	tu.AssertEqualSlice(t, want, filenames)
}

func TestVerifyGo(t *testing.T) {
	ar := &txtar.Archive{
		Files: []txtar.File{
			{
				Name: "main.go",
				Data: []byte(`package main
func main() { fmt.Println("Hello, world!") }
`),
			},
			{
				Name: "bla.go",
				Data: []byte(`package main
oops I did it again
func main() { fmt.Println("Hello, world!") }
`),
			},
		},
	}

	want := "bla.go:2:1: expected declaration, found oops"
	tu.AssertError(t, tu.VerifyGo(ar), want)

	// test the generated tekton example
	ar, err := tu.Folder2Txtar("../kube/testdata/go/tekton")
	tu.AssertNoError(t, err)
	tu.IsNotEqual(t, 0, len(ar.Files))
	tu.AssertNoError(t, tu.VerifyGo(ar))
}
