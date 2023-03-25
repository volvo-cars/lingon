// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package kube_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/dave/jennifer/jen"
	"github.com/volvo-cars/lingon/pkg/kube"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
	// "github.com/volvo-cars/platypus/pkg/platform/tekton"
)

const (
	kubeAppPkgPath = "github.com/volvo-cars/lingon/pkg/kube"
)

var diffOutputDir = filepath.Join(defOutDir, "diff")

// func TestUpdate(t *testing.T) {
// 	out := filepath.Join(defOutDir, "update")
// 	tu.AssertNoError(t, os.RemoveAll(out), "remove output dir")
// 	diff, err := kube.DiffLatest(
// 		"tekton",
// 		"tekton",
// 		out,
// 		nil,
// 		tekton.New(),
// 		[]string{"testdata/tekton-updated.yaml"},
// 	)
// 	tu.AssertNoError(t, err, "update tekton app")
// 	tu.NotNil(t, diff)
// }

func TestDiff2YAML(t *testing.T) {
	tu.AssertNoError(t, os.RemoveAll(diffOutputDir), "remove output dir")
	defer tu.AssertNoError(t, os.RemoveAll(diffOutputDir), "remove output dir")

	appName := "tekton"
	oldPkgName := "old"
	newPkgName := "new"
	oldPkgDir := filepath.Join(diffOutputDir, oldPkgName)
	newPkgDir := filepath.Join(diffOutputDir, newPkgName)

	// export the tekton app we modified from Go to YAML files
	// then import it back to Go
	importGo(
		t,
		oldPkgDir,
		appName,
		oldPkgName,
		[]string{"testdata/tekton.yaml"},
	)
	importGo(
		t,
		newPkgDir,
		appName,
		newPkgName,
		[]string{"testdata/tekton-updated.yaml"},
	)
	generateMain(t, oldPkgName, newPkgName)
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
		kube.WithAppName(appName),
		kube.WithPackageName(pkgName),
		kube.WithOutputDirectory(outputDir),
		kube.WithManifestFiles(files),
		kube.WithSerializer(defaultSerializer()),
		kube.WithRemoveAppName(true),
		kube.WithGroupByKind(true),
	)
	tu.AssertNoError(t, err, "kube import")
}
