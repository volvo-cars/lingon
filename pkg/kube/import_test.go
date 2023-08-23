// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube_test

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/volvo-cars/lingon/pkg/kube"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
	"golang.org/x/tools/txtar"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsbeta "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

const defaultImportOutputDir = "out/import"

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
			Name:   "import with CRDs and remove app name and group by kind",
			OutDir: filepath.Join(defaultImportOutputDir, "argocd"),
			Opts: []kube.ImportOption{
				kube.WithImportAppName("argocd"),
				kube.WithImportManifestFiles([]string{"testdata/argocd.yaml"}),
				kube.WithImportSerializer(defaultSerializer()),
				kube.WithImportRemoveAppName(true),
				kube.WithImportGroupByKind(true),
			},
			OutFiles: []string{
				"out/import/argocd/app.go",
				"out/import/argocd/cluster-role-binding.go",
				"out/import/argocd/cluster-role.go",
				"out/import/argocd/config-map.go",
				"out/import/argocd/custom-resource-definition.go",
				"out/import/argocd/deployment.go",
				"out/import/argocd/network-policy.go",
				"out/import/argocd/role-binding.go",
				"out/import/argocd/role.go",
				"out/import/argocd/secret.go",
				"out/import/argocd/service-account.go",
				"out/import/argocd/service.go",
				"out/import/argocd/stateful-set.go",
			},
		},
		{
			Name: "import with CRDs and remove app name containing dash and group by kind",
			OutDir: filepath.Join(
				defaultImportOutputDir,
				"external-secrets",
			),
			Opts: []kube.ImportOption{
				kube.WithImportAppName("external-secrets"),
				kube.WithImportManifestFiles([]string{"testdata/external-secrets.yaml"}),
				kube.WithImportSerializer(defaultSerializer()),
				kube.WithImportRemoveAppName(true),
				kube.WithImportGroupByKind(true),
			},
			OutFiles: []string{
				"out/import/external-secrets/app.go",
				"out/import/external-secrets/cluster-role-binding.go",
				"out/import/external-secrets/cluster-role.go",
				"out/import/external-secrets/custom-resource-definition.go",
				"out/import/external-secrets/deployment.go",
				"out/import/external-secrets/role-binding.go",
				"out/import/external-secrets/role.go",
				"out/import/external-secrets/secret.go",
				"out/import/external-secrets/service-account.go",
				"out/import/external-secrets/service.go",
				"out/import/external-secrets/validating-webhook-configuration.go",
			},
		},
		{
			Name:   "import with CRDs and remove app name and split by name",
			OutDir: filepath.Join(defaultImportOutputDir, "karpenter"),
			Opts: []kube.ImportOption{
				kube.WithImportAppName("karpenter"),
				kube.WithImportPackageName("karpenter"),
				kube.WithImportManifestFiles([]string{"testdata/karpenter.yaml"}),
				kube.WithImportSerializer(defaultSerializer()),
				kube.WithImportRemoveAppName(true),
				kube.WithImportGroupByKind(false),
			},
			OutFiles: []string{
				"out/import/karpenter/admin_cr.go",
				"out/import/karpenter/app.go",
				"out/import/karpenter/cert_secrets.go",
				"out/import/karpenter/config-logging_cm.go",
				"out/import/karpenter/core_cr.go",
				"out/import/karpenter/core_crb.go",
				"out/import/karpenter/cr.go",
				"out/import/karpenter/crb.go",
				"out/import/karpenter/defaulting.webhook..k8s.aws_mutatingwebhookconfigurations.go",
				"out/import/karpenter/defaulting.webhook..sh_mutatingwebhookconfigurations.go",
				"out/import/karpenter/deploy.go",
				"out/import/karpenter/dns_rb.go",
				"out/import/karpenter/dns_role.go",
				"out/import/karpenter/global-settings_cm.go",
				"out/import/karpenter/pdb.go",
				"out/import/karpenter/rb.go",
				"out/import/karpenter/role.go",
				"out/import/karpenter/sa.go",
				"out/import/karpenter/svc.go",
				"out/import/karpenter/validation.webhook..k8s.aws_validatingwebhookconfigurations.go",
				"out/import/karpenter/validation.webhook..sh_validatingwebhookconfigurations.go",
				"out/import/karpenter/validation.webhook.config..sh_validatingwebhookconfigurations.go",
			},
		},
		{
			Name:   "import with vanilla serializer and remove app name and group by kind",
			OutDir: filepath.Join(defaultImportOutputDir, "grafana"),
			Opts: []kube.ImportOption{
				kube.WithImportAppName("grafana"),
				kube.WithImportPackageName("grafana"),
				kube.WithImportManifestFiles([]string{"testdata/grafana.yaml"}),
				kube.WithImportRemoveAppName(true),
				kube.WithImportGroupByKind(true),
			},
			OutFiles: []string{
				"out/import/grafana/app.go",
				"out/import/grafana/cluster-role-binding.go",
				"out/import/grafana/cluster-role.go",
				"out/import/grafana/config-map.go",
				"out/import/grafana/deployment.go",
				"out/import/grafana/pod.go",
				"out/import/grafana/role-binding.go",
				"out/import/grafana/role.go",
				"out/import/grafana/secret.go",
				"out/import/grafana/service-account.go",
				"out/import/grafana/service.go",
			},
		},
		{
			Name:   "import with vanilla serializer and add methods",
			OutDir: filepath.Join(defaultImportOutputDir, "manifester"),
			Opts: []kube.ImportOption{
				kube.WithImportAppName("grafana"),
				kube.WithImportPackageName("grafana"),
				kube.WithImportManifestFiles([]string{"testdata/grafana.yaml"}),
				kube.WithImportRemoveAppName(true),
				kube.WithImportGroupByKind(true),
				kube.WithImportAddMethods(true),
			},
			OutFiles: []string{
				"out/import/manifester/app.go",
				"out/import/manifester/cluster-role-binding.go",
				"out/import/manifester/cluster-role.go",
				"out/import/manifester/config-map.go",
				"out/import/manifester/deployment.go",
				"out/import/manifester/pod.go",
				"out/import/manifester/role-binding.go",
				"out/import/manifester/role.go",
				"out/import/manifester/secret.go",
				"out/import/manifester/service-account.go",
				"out/import/manifester/service.go",
			},
		},
		{
			Name:   "import with group kind and clean name",
			OutDir: filepath.Join(defaultImportOutputDir, "tekton"),
			Opts: []kube.ImportOption{
				kube.WithImportAppName("tekton"),
				kube.WithImportPackageName("tekton"),
				kube.WithImportManifestFiles([]string{"testdata/tekton.yaml"}),
				kube.WithImportRemoveAppName(true),
				kube.WithImportGroupByKind(true),
				kube.WithImportAddMethods(true),
				kube.WithImportCleanUp(true),
			},
			OutFiles: []string{
				"out/import/tekton/app.go",
				"out/import/tekton/cluster-role-binding.go",
				"out/import/tekton/cluster-role.go",
				"out/import/tekton/config-map.go",
				"out/import/tekton/custom-resource-definition.go",
				"out/import/tekton/deployment.go",
				"out/import/tekton/horizontal-pod-autoscaler.go",
				"out/import/tekton/mutating-webhook-configuration.go",
				"out/import/tekton/namespace.go",
				"out/import/tekton/role-binding.go",
				"out/import/tekton/role.go",
				"out/import/tekton/secret.go",
				"out/import/tekton/service-account.go",
				"out/import/tekton/service.go",
				"out/import/tekton/validating-webhook-configuration.go",
			},
		},
		{
			Name:   "import list object",
			OutDir: filepath.Join(defaultImportOutputDir, "velero"),
			Opts: []kube.ImportOption{
				kube.WithImportAppName("velero"),
				kube.WithImportPackageName("velero"),
				kube.WithImportManifestFiles([]string{"testdata/velero.yaml"}),
				kube.WithImportRemoveAppName(true),
				kube.WithImportGroupByKind(true),
				kube.WithImportAddMethods(true),
				kube.WithImportCleanUp(true),
			},
			OutFiles: []string{
				"out/import/velero/app.go",
				"out/import/velero/cluster-role-binding.go",
				"out/import/velero/custom-resource-definition.go",
				"out/import/velero/deployment.go",
				"out/import/velero/namespace.go",
				"out/import/velero/service-account.go",
			},
		},
	}

	for _, tt := range TT {
		tc := tt
		t.Run(
			tt.Name, func(t *testing.T) {
				t.Parallel()
				var buf bytes.Buffer
				tc.Opts = append(
					tc.Opts,
					kube.WithImportOutputDirectory(tc.OutDir),
					kube.WithImportWriter(&buf),
				)
				err := kube.Import(tc.Opts...)
				tu.AssertNoError(t, err, "failed to import")

				// compare filenames
				ar := txtar.Parse(buf.Bytes())
				got := tu.Filenames(ar)
				sort.Strings(got)
				want := tc.OutFiles
				tu.AssertEqualSlice(t, want, got)
				tu.AssertNoError(t, tu.VerifyGo(ar))

				// compare content
				golden, err := txtar.ParseFile(
					filepath.Join(
						"testdata", "golden",
						strings.ReplaceAll(tc.Name, " ", "_")+".txt",
					),
				)
				tu.AssertNoError(t, err, "reading golden file")
				if diff := tu.DiffTxtarSort(ar, golden); diff != "" {
					t.Fatal(tu.Callers(), diff)
				}
			},
		)
	}
}

func TestImport_Error(t *testing.T) {
	err := kube.Import(
		kube.WithImportAppName("foo-app"),
		kube.WithImportPackageName("foo-package"),
		kube.WithImportManifestFiles([]string{"does-not-exists.yaml"}),
	)
	tu.IsNotNil(t, err)
	output := `import options: incompatible options
package name cannot contain a dash
file does not exist: does-not-exists.yaml`
	tu.AssertErrorMsg(t, err, output)
}

func TestImport_ErrorEmptyManifest(t *testing.T) {
	pathg := "testdata/golden/empty.golden"
	pathy := "testdata/golden/empty.yaml"
	golden, err := os.ReadFile(pathg)
	tu.AssertNoError(t, err, "read golden file")
	var buf bytes.Buffer
	err = kube.Import(
		kube.WithImportAppName("foo-app"),
		kube.WithImportPackageName("foopackage"),
		kube.WithImportManifestFiles([]string{pathy}),
		kube.WithImportWriter(&buf),
		kube.WithImportGroupByKind(false),
	)
	tu.AssertNoError(t, err, "failed to import")
	tu.AssertEqual(t, string(golden), buf.String())
}

func TestImport_SaveFromReader(t *testing.T) {
	filename := "testdata/grafana.yaml"
	file, err := os.Open(filename)
	tu.AssertNoError(t, err, fmt.Sprintf("failed to open file: %s", filename))
	out := filepath.Join(defaultImportOutputDir, "reader")
	var buf bytes.Buffer
	err = kube.Import(
		kube.WithImportAppName("grafana"),
		kube.WithImportPackageName("grafana"),
		kube.WithImportOutputDirectory(out),
		kube.WithImportWriter(&buf),
		kube.WithImportReader(file),
		kube.WithImportNameFieldFunc(kube.NameFieldFunc),
		kube.WithImportNameVarFunc(kube.NameVarFunc),
		kube.WithImportGroupByKind(false),
		kube.WithImportNameFileFunc(
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
		"out/import/reader/app.go",
		"out/import/reader/clusterrole-grafana-clusterrole.go",
		"out/import/reader/clusterrolebinding-grafana-clusterrolebinding.go",
		"out/import/reader/configmap-grafana-dashboards-default.go",
		"out/import/reader/configmap-grafana-test.go",
		"out/import/reader/configmap-grafana.go",
		"out/import/reader/deployment-grafana.go",
		"out/import/reader/pod-grafana-test.go",
		"out/import/reader/role-grafana.go",
		"out/import/reader/rolebinding-grafana.go",
		"out/import/reader/secret-grafana.go",
		"out/import/reader/service-grafana.go",
		"out/import/reader/serviceaccount-grafana-test.go",
		"out/import/reader/serviceaccount-grafana.go",
	}

	// compare filenames
	ar := txtar.Parse(buf.Bytes())
	got := tu.Filenames(ar)
	sort.Strings(got)
	tu.AssertEqualSlice(t, want, got)

	// compare content
	golden, err := txtar.ParseFile(
		filepath.Join("testdata", "golden", "import_save_from_reader.txt"),
	)
	tu.AssertNoError(t, err, "reading golden file")
	if diff := tu.DiffTxtarSort(ar, golden); diff != "" {
		t.Fatal(tu.Callers(), diff)
	}
}

func TestImport_MissingCRDs(t *testing.T) {
	filename := "testdata/istio.yaml"
	file, err := os.Open(filename)
	tu.AssertNoError(t, err, fmt.Sprintf("failed to open file: %s", filename))

	var buf bytes.Buffer

	err = kube.Import(
		kube.WithImportAppName("istio"),
		kube.WithImportPackageName("istio"),
		kube.WithImportOutputDirectory("manifests/"),
		kube.WithImportReader(file),
		kube.WithImportWriter(&buf),
		kube.WithImportGroupByKind(false),
		kube.WithImportNameFileFunc(
			func(m kubeutil.Metadata) string {
				return fmt.Sprintf(
					"%s-%s.go",
					strings.ToLower(m.Kind),
					m.Meta.Name,
				)
			},
		),
	)
	errmsg := "generate go: stdin: " +
		"no kind \"EnvoyFilter\" is registered for version " +
		"\"networking.istio.io/v1alpha3\" in scheme \"pkg/runtime/scheme.go:100\""
	tu.AssertErrorMsg(t, err, errmsg)
}

func TestImport_ReaderWriter(t *testing.T) {
	filename := "testdata/grafana.yaml"
	file, err := os.Open(filename)
	tu.AssertNoError(t, err, fmt.Sprintf("failed to open file: %s", filename))

	var buf bytes.Buffer

	err = kube.Import(
		kube.WithImportAppName("grafana"),
		kube.WithImportPackageName("grafana"),
		kube.WithImportOutputDirectory("manifests/"),
		kube.WithImportReader(file),
		kube.WithImportWriter(&buf),
		kube.WithImportGroupByKind(false),
		kube.WithImportNameFileFunc(
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

	// compare filenames
	ar := txtar.Parse(buf.Bytes())
	got := tu.Filenames(ar)
	sort.Strings(got)
	tu.AssertEqualSlice(t, want, got)

	// compare content
	golden, err := txtar.ParseFile(
		filepath.Join(
			"testdata",
			"golden",
			"import_reader_writer.txt",
		),
	)
	tu.AssertNoError(t, err, "reading golden file")
	if diff := tu.DiffTxtarSort(ar, golden); diff != "" {
		t.Fatal(tu.Callers(), diff)
	}
}

func TestImport_ConfigMapComments(t *testing.T) {
	out := filepath.Join(defaultImportOutputDir, "tekton")
	tu.AssertNoError(t, os.RemoveAll(out), "rm out dir")
	t.Cleanup(
		func() {
			tu.AssertNoError(t, os.RemoveAll(out), "rm out dir")
		},
	)

	err := kube.Import(
		kube.WithImportAppName("tekton"),
		kube.WithImportPackageName("tekton"),
		kube.WithImportOutputDirectory(out),
		kube.WithImportManifestFiles([]string{"testdata/tekton.yaml"}),
		kube.WithImportSerializer(defaultSerializer()),
		kube.WithImportRemoveAppName(true),
		kube.WithImportGroupByKind(true),
		kube.WithImportAddMethods(true),
		kube.WithImportCleanUp(false),
	)
	tu.AssertNoError(t, err, "failed to import")

	got, err := kubeutil.ListGoFiles(out)
	tu.AssertNoError(t, err, "failed to list go files")
	sort.Strings(got)

	want := []string{
		"out/import/tekton/app.go",
		"out/import/tekton/cluster-role-binding.go",
		"out/import/tekton/cluster-role.go",
		"out/import/tekton/config-map.go",
		"out/import/tekton/custom-resource-definition.go",
		"out/import/tekton/deployment.go",
		"out/import/tekton/horizontal-pod-autoscaler.go",
		"out/import/tekton/mutating-webhook-configuration.go",
		"out/import/tekton/namespace.go",
		"out/import/tekton/role-binding.go",
		"out/import/tekton/role.go",
		"out/import/tekton/secret.go",
		"out/import/tekton/service-account.go",
		"out/import/tekton/service.go",
		"out/import/tekton/validating-webhook-configuration.go",
	}

	tu.IsEqual(t, len(want), len(got))
	tu.AssertEqualSlice(t, want, got)

	src, err := os.ReadFile("out/import/tekton/config-map.go")
	tu.AssertNoError(t, err, "reading config-map.go")

	comments := []string{
		"\t\t   Contains pipelines version which can be queried by external\n\t\t   tools such as CLI. Elevated permissions are already given to\n\t\t   this ConfigMap such that even if we don't have access to\n\t\t   other resources in the namespace we still can have access to\n\t\t   this ConfigMap.\n",
		"\t\t   Setting this flag to \"enforce\" will enforce verification of tasks/pipeline. Failing to verify\n\t\t   will fail the taskrun/pipelinerun. \"warn\" will only log the err message and \"skip\"\n\t\t   will skip the whole verification\n",
		"\t\t   Setting this flag to \"false\" will stop Tekton from waiting for a\n\t\t   TaskRun's sidecar containers to be running before starting the first\n\t\t   step. This will allow Tasks to be run in environments that don't\n\t\t   support the DownwardAPI volume type, but may lead to unintended\n\t\t   behaviour if sidecars are used.\n\t\t   #\n\t\t   See https://github.com/tektoncd/pipeline/issues/4937 for more info.\n",
		"\t\t   Setting this flag to \"true\" enables CloudEvents for CustomRuns and Runs, as long as a\n\t\t   CloudEvents sink is configured in the config-defaults config map\n",
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
	sort.SliceStable(cc, func(i, j int) bool { return cc[i] < cc[j] })
	tu.AssertEqualSlice(t, comments, cc[:len(comments)])
}

func TestImport_VerboseLogger(t *testing.T) {
	out := filepath.Join(defaultImportOutputDir, "tekton")
	golden := filepath.Join("testdata", "golden", "log.golden")
	var bufLog, bufOut bytes.Buffer

	log := func(w io.Writer) *slog.Logger {
		replace := func(groups []string, a slog.Attr) slog.Attr {
			// Remove time.
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.Attr{}
			}
			// remove file:line
			if a.Key == slog.SourceKey {
				return slog.Attr{}
			}
			return a
		}
		return slog.New(
			slog.NewTextHandler(
				w,
				&slog.HandlerOptions{
					AddSource:   true,
					ReplaceAttr: replace,
				},
			).WithAttrs(
				[]slog.Attr{slog.String("app", "lingon")},
			),
		)
	}

	err := kube.Import(
		kube.WithImportAppName("tekton"),
		kube.WithImportPackageName("tekton"),
		kube.WithImportOutputDirectory(out),
		kube.WithImportManifestFiles([]string{"testdata/tekton.yaml"}),
		kube.WithImportSerializer(defaultSerializer()),
		kube.WithImportRemoveAppName(true),
		kube.WithImportGroupByKind(true),
		kube.WithImportAddMethods(true),
		kube.WithImportWriter(&bufOut),
		kube.WithImportVerbose(true),
		kube.WithImportIgnoreErrors(true),
		kube.WithImportLogger(log(&bufLog)),
		kube.WithImportCleanUp(true),
	)
	tu.AssertNoError(t, err, "failed to import")
	got := bufLog.String()
	want, err := os.ReadFile(golden)
	tu.AssertNoError(t, err, "reading golden file")
	tu.AssertEqual(t, string(want), got)
}

func TestImport_CleanUp(t *testing.T) {
	var buf bytes.Buffer
	err := kube.Import(
		kube.WithImportManifestFiles([]string{"testdata/golden/dirty.yaml"}),
		kube.WithImportWriter(&buf),
		kube.WithImportGroupByKind(false),
		kube.WithImportCleanUp(true),
	)
	tu.AssertNoError(t, err, "failed to import")

	ar, err := txtar.ParseFile("testdata/golden/dirty.txt")
	tu.AssertNoError(t, err, "reading golden file")
	want := txtar.Format(ar)
	tu.AssertEqual(t, string(want), buf.String())
}
