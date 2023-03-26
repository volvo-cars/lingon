// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package kube_test

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rogpeppe/go-internal/txtar"
	"github.com/volvo-cars/lingon/pkg/kube"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsbeta "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

const defOutDir = "out/jamel"

func defaultSerializer() runtime.Decoder {
	// NEEDED FOR CRDS
	//
	_ = apiextensions.AddToScheme(scheme.Scheme)
	_ = apiextensionsbeta.AddToScheme(scheme.Scheme)
	return scheme.Codecs.UniversalDeserializer()
}

func TestImport(t *testing.T) {
	type args struct {
		Name     string
		OutDir   string
		Opts     []kube.ImportOption
		OutFiles []string
	}
	TT := []args{
		{
			Name:   "convert with CRDs and remove app name and group by kind",
			OutDir: filepath.Join(defOutDir, "argocd"),
			Opts: []kube.ImportOption{
				kube.WithAppName("argocd"),
				kube.WithOutputDirectory(filepath.Join(defOutDir, "argocd")),
				kube.WithManifestFiles([]string{"testdata/argocd.yaml"}),
				kube.WithSerializer(defaultSerializer()),
				kube.WithRemoveAppName(true),
				kube.WithGroupByKind(true),
			},
			OutFiles: []string{
				"out/jamel/argocd/app.go",
				"out/jamel/argocd/cluster-role-binding.go",
				"out/jamel/argocd/cluster-role.go",
				"out/jamel/argocd/config-map.go",
				"out/jamel/argocd/custom-resource-definition.go",
				"out/jamel/argocd/deployment.go",
				"out/jamel/argocd/network-policy.go",
				"out/jamel/argocd/role-binding.go",
				"out/jamel/argocd/role.go",
				"out/jamel/argocd/secret.go",
				"out/jamel/argocd/service-account.go",
				"out/jamel/argocd/service.go",
				"out/jamel/argocd/stateful-set.go",
			},
		}, {
			Name: "convert with CRDs and remove app name containing dash and group by kind",
			OutDir: filepath.Join(
				defOutDir,
				"external-secrets",
			),
			Opts: []kube.ImportOption{
				kube.WithAppName("external-secrets"),
				kube.WithPackageName("externalsecrets"),
				kube.WithOutputDirectory(
					filepath.Join(
						defOutDir,
						"external-secrets",
					),
				),
				kube.WithManifestFiles([]string{"testdata/external-secrets.yaml"}),
				kube.WithSerializer(defaultSerializer()),
				kube.WithRemoveAppName(true),
				kube.WithGroupByKind(true),
			},
			OutFiles: []string{
				"out/jamel/external-secrets/app.go",
				"out/jamel/external-secrets/cluster-role-binding.go",
				"out/jamel/external-secrets/cluster-role.go",
				"out/jamel/external-secrets/custom-resource-definition.go",
				"out/jamel/external-secrets/deployment.go",
				"out/jamel/external-secrets/role-binding.go",
				"out/jamel/external-secrets/role.go",
				"out/jamel/external-secrets/secret.go",
				"out/jamel/external-secrets/service-account.go",
				"out/jamel/external-secrets/service.go",
				"out/jamel/external-secrets/validating-webhook-configuration.go",
			},
		}, {
			Name:   "convert with CRDs and remove app name and split by name",
			OutDir: filepath.Join(defOutDir, "karpenter"),
			Opts: []kube.ImportOption{
				kube.WithAppName("karpenter"),
				kube.WithPackageName("karpenter"),
				kube.WithOutputDirectory(filepath.Join(defOutDir, "karpenter")),
				kube.WithManifestFiles([]string{"testdata/karpenter.yaml"}),
				kube.WithSerializer(defaultSerializer()),
				kube.WithRemoveAppName(true),
			},
			OutFiles: []string{
				"out/jamel/karpenter/admin_cr.go",
				"out/jamel/karpenter/app.go",
				"out/jamel/karpenter/cert_secrets.go",
				"out/jamel/karpenter/config-logging_cm.go",
				"out/jamel/karpenter/core_cr.go",
				"out/jamel/karpenter/core_crb.go",
				"out/jamel/karpenter/cr.go",
				"out/jamel/karpenter/crb.go",
				"out/jamel/karpenter/defaulting.webhook..k8s.aws_mutatingwebhookconfigurations.go",
				"out/jamel/karpenter/defaulting.webhook..sh_mutatingwebhookconfigurations.go",
				"out/jamel/karpenter/deploy.go",
				"out/jamel/karpenter/dns_rb.go",
				"out/jamel/karpenter/dns_role.go",
				"out/jamel/karpenter/global-settings_cm.go",
				"out/jamel/karpenter/pdb.go",
				"out/jamel/karpenter/rb.go",
				"out/jamel/karpenter/role.go",
				"out/jamel/karpenter/sa.go",
				"out/jamel/karpenter/svc.go",
				"out/jamel/karpenter/validation.webhook..k8s.aws_validatingwebhookconfigurations.go",
				"out/jamel/karpenter/validation.webhook..sh_validatingwebhookconfigurations.go",
				"out/jamel/karpenter/validation.webhook.config..sh_validatingwebhookconfigurations.go",
			},
		}, {
			Name:   "convert with vanilla serializer and remove app name and group by kind",
			OutDir: filepath.Join(defOutDir, "grafana"),
			Opts: []kube.ImportOption{
				kube.WithAppName("grafana"),
				kube.WithPackageName("grafana"),
				kube.WithOutputDirectory(filepath.Join(defOutDir, "grafana")),
				kube.WithManifestFiles([]string{"testdata/grafana.yaml"}),
				kube.WithRemoveAppName(true),
				kube.WithGroupByKind(true),
			},
			OutFiles: []string{
				"out/jamel/grafana/app.go",
				"out/jamel/grafana/cluster-role-binding.go",
				"out/jamel/grafana/cluster-role.go",
				"out/jamel/grafana/config-map.go",
				"out/jamel/grafana/deployment.go",
				"out/jamel/grafana/pod.go",
				"out/jamel/grafana/role-binding.go",
				"out/jamel/grafana/role.go",
				"out/jamel/grafana/secret.go",
				"out/jamel/grafana/service-account.go",
				"out/jamel/grafana/service.go",
			},
		}, {
			Name:   "convert grafana with vanilla serializer and implement Exporter",
			OutDir: filepath.Join(defOutDir, "manifester"),
			Opts: []kube.ImportOption{
				kube.WithAppName("grafana"),
				kube.WithPackageName("grafana"),
				kube.WithOutputDirectory(
					filepath.Join(
						defOutDir,
						"manifester",
					),
				),
				kube.WithManifestFiles([]string{"testdata/grafana.yaml"}),
				kube.WithRemoveAppName(true),
				kube.WithGroupByKind(true),
				kube.WithMethods(true),
			},
			OutFiles: []string{
				"out/jamel/manifester/app.go",
				"out/jamel/manifester/cluster-role-binding.go",
				"out/jamel/manifester/cluster-role.go",
				"out/jamel/manifester/config-map.go",
				"out/jamel/manifester/deployment.go",
				"out/jamel/manifester/pod.go",
				"out/jamel/manifester/role-binding.go",
				"out/jamel/manifester/role.go",
				"out/jamel/manifester/secret.go",
				"out/jamel/manifester/service-account.go",
				"out/jamel/manifester/service.go",
			},
		},
	}

	for _, tt := range TT {
		t.Run(
			tt.Name, func(t *testing.T) {
				tu.AssertNoError(t, os.RemoveAll(tt.OutDir), "rm out dir")
				err := kube.Import(tt.Opts...)
				tu.AssertNoError(t, err, "failed to import")
				got, err := kube.ListGoFiles(tt.OutDir)
				tu.AssertNoError(t, err, "failed to list go files")
				sort.Strings(got)
				want := tt.OutFiles
				if !cmp.Equal(want, got) {
					t.Error(tu.Diff(want, got))
				}
			},
		)
		_ = os.RemoveAll(tt.OutDir)
	}
}

func TestJamel_SaveFromReader(t *testing.T) {
	filename := "testdata/grafana.yaml"
	file, err := os.Open(filename)
	tu.AssertNoError(t, err, fmt.Sprintf("failed to open file: %s", filename))
	out := filepath.Join(defOutDir, "reader")
	tu.AssertNoError(t, os.RemoveAll(out), "rm out dir")
	defer tu.AssertNoError(t, os.RemoveAll(out), "rm out dir")

	err = kube.Import(
		kube.WithAppName("grafana"),
		kube.WithPackageName("grafana"),
		kube.WithOutputDirectory(out),
		kube.WithReader(file),
		kube.WithNameFieldFunc(kube.NameFieldFunc),
		kube.WithNameVarFunc(kube.NameVarFunc),
		kube.WithNameFileFunc(
			func(m kubeutil.Metadata) string {
				return fmt.Sprintf(
					"%s-%s.go",
					strings.ToLower(m.Kind),
					m.Meta.Name,
				)
			},
		),
	)
	tu.AssertNoError(t, err, "failed to import")

	got, err := kube.ListGoFiles(out)
	tu.AssertNoError(t, err, "failed to list go files")
	sort.Strings(got)

	want := []string{
		"out/jamel/reader/app.go",
		"out/jamel/reader/clusterrole-grafana-clusterrole.go",
		"out/jamel/reader/clusterrolebinding-grafana-clusterrolebinding.go",
		"out/jamel/reader/configmap-grafana-dashboards-default.go",
		"out/jamel/reader/configmap-grafana-test.go",
		"out/jamel/reader/configmap-grafana.go",
		"out/jamel/reader/deployment-grafana.go",
		"out/jamel/reader/pod-grafana-test.go",
		"out/jamel/reader/role-grafana.go",
		"out/jamel/reader/rolebinding-grafana.go",
		"out/jamel/reader/secret-grafana.go",
		"out/jamel/reader/service-grafana.go",
		"out/jamel/reader/serviceaccount-grafana-test.go",
		"out/jamel/reader/serviceaccount-grafana.go",
	}

	if len(got) != len(want) {
		t.Errorf("expected %d files, got %d", len(want), len(got))
	}
	if !cmp.Equal(want, got) {
		t.Error(tu.Diff(want, got))
	}
}

func TestJamel_ReaderWriter(t *testing.T) {
	filename := "testdata/grafana.yaml"
	file, err := os.Open(filename)
	tu.AssertNoError(t, err, fmt.Sprintf("failed to open file: %s", filename))

	var bufout bytes.Buffer

	err = kube.Import(
		kube.WithAppName("grafana"),
		kube.WithPackageName("grafana"),
		kube.WithOutputDirectory("manifests/"),
		kube.WithReader(file),
		kube.WithWriter(&bufout),
		kube.WithNameFileFunc(
			func(m kubeutil.Metadata) string {
				return fmt.Sprintf(
					"%s-%s.go",
					strings.ToLower(m.Kind),
					m.Meta.Name,
				)
			},
		),
	)
	tu.AssertNoError(t, err, "failed to import")

	want := []string{
		"manifests/app.go",
		"manifests/clusterrole-grafana-clusterrole.go",
		"manifests/clusterrolebinding-grafana-clusterrolebinding.go",
		"manifests/configmap-grafana-dashboards-default.go",
		"manifests/configmap-grafana-test.go",
		"manifests/configmap-grafana.go",
		"manifests/deployment-grafana.go",
		"manifests/pod-grafana-test.go",
		"manifests/role-grafana.go",
		"manifests/rolebinding-grafana.go",
		"manifests/secret-grafana.go",
		"manifests/service-grafana.go",
		"manifests/serviceaccount-grafana-test.go",
		"manifests/serviceaccount-grafana.go",
	}

	ar := txtar.Parse(bufout.Bytes())
	got := make([]string, len(ar.Files))
	for i, f := range ar.Files {
		got[i] = f.Name
	}
	sort.Strings(got)
	if !cmp.Equal(want, got) {
		t.Error(tu.Diff(want, got))
	}
}

func TestJamel_ConfigMapComments(t *testing.T) {
	out := filepath.Join(defOutDir, "tekton")
	tu.AssertNoError(t, os.RemoveAll(out), "rm out dir")
	defer tu.AssertNoError(t, os.RemoveAll(out), "rm out dir")

	err := kube.Import(
		kube.WithAppName("tekton"),
		kube.WithPackageName("tekton"),
		kube.WithOutputDirectory(out),
		kube.WithManifestFiles([]string{"testdata/tekton.yaml"}),
		kube.WithSerializer(defaultSerializer()),
		kube.WithRemoveAppName(true),
		kube.WithGroupByKind(true),
		kube.WithMethods(true),
	)
	tu.AssertNoError(t, err, "failed to import")

	got, err := kube.ListGoFiles(out)
	tu.AssertNoError(t, err, "failed to list go files")
	sort.Strings(got)

	want := []string{
		"out/jamel/tekton/app.go",
		"out/jamel/tekton/cluster-role-binding.go",
		"out/jamel/tekton/cluster-role.go",
		"out/jamel/tekton/config-map.go",
		"out/jamel/tekton/custom-resource-definition.go",
		"out/jamel/tekton/deployment.go",
		"out/jamel/tekton/horizontal-pod-autoscaler.go",
		"out/jamel/tekton/mutating-webhook-configuration.go",
		"out/jamel/tekton/namespace.go",
		"out/jamel/tekton/role-binding.go",
		"out/jamel/tekton/role.go",
		"out/jamel/tekton/secret.go",
		"out/jamel/tekton/service-account.go",
		"out/jamel/tekton/service.go",
		"out/jamel/tekton/validating-webhook-configuration.go",
	}

	if len(got) != len(want) {
		t.Errorf("expected %d files, got %d", len(want), len(got))
	}
	if !cmp.Equal(want, got) {
		t.Error(tu.Diff(want, got))
	}

	src, err := os.ReadFile("out/jamel/tekton/config-map.go")
	tu.AssertNoError(t, err, "reading config-map.go")

	comments := []string{
		"Contains pipelines version which can be queried by external\n\t\t   tools such as CLI. Elevated permissions are already given to\n\t\t   this ConfigMap such that even if we don't have access to\n\t\t   other resources in the namespace we still can have access to\n\t\t   this ConfigMap.",
		"Setting this flag to \"enforce\" will enforce verification of tasks/pipeline. Failing to verify\n\t\t   will fail the taskrun/pipelinerun. \"warn\" will only log the err message and \"skip\"\n\t\t   will skip the whole verification",
		"Setting this flag to \"false\" will stop Tekton from waiting for a\n\t\t   TaskRun's sidecar containers to be running before starting the first\n\t\t   step. This will allow Tasks to be run in environments that don't\n\t\t   support the DownwardAPI volume type, but may lead to unintended\n\t\t   behaviour if sidecars are used.\n\t\t   #\n\t\t   See https://github.com/tektoncd/pipeline/issues/4937 for more info.",
		"Setting this flag to \"true\" enables CloudEvents for CustomRuns and Runs, as long as a\n\t\t   CloudEvents sink is configured in the config-defaults config map",
	}
	set := token.NewFileSet()
	astFile, err := parser.ParseFile(set, "", src, parser.ParseComments)
	tu.AssertNoError(t, err, fmt.Sprintf("parsing %s", src))
	if len(astFile.Comments) < len(comments) {
		t.Errorf("not enough comments: %d", len(astFile.Comments))
	}
	cc := []string{}
	for _, comment := range astFile.Comments {
		cc = append(cc, comment.Text())
	}
	sort.SliceStable(
		cc, func(i, j int) bool {
			return cc[i] < cc[j]
		},
	)
	for i, comment := range comments {
		if diff := tu.Diff(cleanStr(comment), cleanStr(cc[i])); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func cleanStr(s string) string {
	return strings.Trim(
		strings.ReplaceAll(
			strings.ReplaceAll(s, "\t", ""),
			"\n",
			"",
		), " ",
	)
}
