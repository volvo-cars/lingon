// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package testutil_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/dave/jennifer/jen"
	"github.com/volvo-cars/lingon/pkg/kube"
	"github.com/volvo-cars/lingon/pkg/kube/testdata/go/tekton"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsbeta "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

const (
	kubeAppPkgPath = "github.com/volvo-cars/lingon/pkg/kube"
)

var (
	defaultImportOutputDir = "out"
	diffOutputDir          = filepath.Join(defaultImportOutputDir, "diff")
)

func defaultSerializer() runtime.Decoder {
	// NEEDED FOR CRDS
	//
	_ = apiextensions.AddToScheme(scheme.Scheme)
	_ = apiextensionsbeta.AddToScheme(scheme.Scheme)
	return scheme.Codecs.UniversalDeserializer()
}

func TestDiffUpdate(t *testing.T) {
	_, _ = diffLatest(
		"tekton",
		"tekton",
		diffOutputDir,
		nil,
		tekton.New(),
		[]string{"testdata/tekton-updated.yaml"},
	)
}

func TestDiff2YAML(t *testing.T) {
	t.Skip("skipping test Diff for now.")

	tu.AssertNoError(t, os.RemoveAll(diffOutputDir), "remove output dir")
	t.Cleanup(
		func() {
			tu.AssertNoError(t, os.RemoveAll(diffOutputDir), "rm out dir")
		},
	)

	appName := "tekton"
	currentPkgName := "old"
	updatePkgName := "update"
	currentPkgDir := filepath.Join(diffOutputDir, currentPkgName)
	updatePkgDir := filepath.Join(diffOutputDir, updatePkgName)

	// export the tekton app we modified from Go to YAML files
	// then import it back to Go
	importGo(
		t,
		currentPkgDir,
		appName,
		currentPkgName,
		[]string{"testdata/tekton.yaml"},
	)
	importGo(
		t,
		updatePkgDir,
		appName,
		updatePkgName,
		[]string{"testdata/tekton-updated.yaml"},
	)
	generateMain(t, currentPkgName, updatePkgName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// #nosec G204 - we are not passing any user input to the command
	err := exec.CommandContext(
		ctx,
		"go",
		"run",
		filepath.Join(diffOutputDir, "main.go"),
	).Run()
	if err != nil {
		m, err := os.ReadFile(filepath.Join(diffOutputDir, "main.go"))
		tu.AssertNoError(t, err, "read main.go")
		t.Fatalf("cannot run main.go: %s", string(m))
	}
}

func generateMain(t *testing.T, pkgName ...string) {
	diffPkgDir := kubeAppPkgPath + "/" + diffOutputDir
	f := jen.NewFile("main")
	stmts := []*jen.Statement{}
	for _, p := range pkgName {
		f.ImportName(diffPkgDir+"/"+p, p)
		stmts = append(
			stmts,
			jen.Id("err"+p).Op(":=").Qual(
				diffPkgDir+"/"+p,
				"New",
			).Call().Dot("Export").Call(
				jen.Lit(
					filepath.Join(
						diffOutputDir,
						"exported"+p,
					),
				),
			),
			jen.If(jen.Id("err"+p).Op("!=").Nil()).Block(
				jen.Panic(jen.Id("err"+p)),
			),
		)
	}

	f.Func().Id("main").Params().BlockFunc(
		func(g *jen.Group) {
			for _, s := range stmts {
				g.Add(s)
			}
		},
	)

	err := f.Save(filepath.Join(diffOutputDir, "main.go"))
	tu.AssertNoError(t, err, "save main.go")
}

func importGo(
	t *testing.T,
	outputDir, appName, pkgName string,
	files []string,
) {
	t.Helper()
	tu.AssertNoError(t, os.RemoveAll(outputDir), "remove output dir")
	// defer os.RemoveAll(out)

	err := kube.Import(
		kube.WithImportAppName(appName),
		kube.WithImportPackageName(pkgName),
		kube.WithImportOutputDirectory(outputDir),
		kube.WithImportManifestFiles(files),
		kube.WithImportSerializer(defaultSerializer()),
		kube.WithImportRemoveAppName(true),
		kube.WithImportGroupByKind(true),
	)
	tu.AssertNoError(t, err, "kube import")
}
