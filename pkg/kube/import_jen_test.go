// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"bytes"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/volvo-cars/lingon/pkg/internal/api"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestStmtStruct(t *testing.T) {
	f := jen.NewFile("test")
	sc := make(map[string]*jen.Statement)

	pkgPath, _ := api.PkgPathFromAPIVersion("apps/v1")
	sc["Depl"] = jen.Qual(pkgPath, "Deployment")
	pkgPath, _ = api.PkgPathFromAPIVersion("apps/v1")
	sc["CM"] = jen.Qual(pkgPath, "ConfigMap")

	f.Add(stmtStruct("MyApp", sc))

	var buf bytes.Buffer
	tu.AssertNoError(t, f.Render(&buf), "render")
	want := `package test

import (
	kube "github.com/volvo-cars/lingon/pkg/kube"
	v1 "k8s.io/api/apps/v1"
)

type MyApp struct {
	kube.App

	CM   *v1.ConfigMap
	Depl *v1.Deployment
}
`
	tu.AssertEqual(t, buf.String(), want)
}

func TestImportAddMethods(t *testing.T) {
	f := jen.NewFile("test")
	addMethods(f, "NAME")
	var buf bytes.Buffer
	tu.AssertNoError(t, f.Render(&buf), "render")
	want := `package test

import (
	"context"
	kube "github.com/volvo-cars/lingon/pkg/kube"
)

// Apply applies the kubernetes objects to the cluster
func (a *NAME) Apply(ctx context.Context) error {
	return Apply(ctx, a)
}

// Export exports the kubernetes objects to YAML files in the given directory
func (a *NAME) Export(dir string) error {
	return kube.Export(a, kube.WithExportOutputDirectory(dir))
}
`
	tu.AssertEqual(t, buf.String(), want)
}

func TestStmtApplyFunc(t *testing.T) {
	f := jen.NewFile("test")
	f.Add(stmtApplyFunc())
	var buf bytes.Buffer
	tu.AssertNoError(t, f.Render(&buf), "render")
	want := `package test

import (
	"context"
	"errors"
	kube "github.com/volvo-cars/lingon/pkg/kube"
	"os"
	"os/exec"
)

// Apply applies the kubernetes objects contained in Exporter to the cluster
func Apply(ctx context.Context, km kube.Exporter) error {
	cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", "-")
	cmd.Env = os.Environ()        // inherit environment in case we need to use kubectl from a container
	stdin, err := cmd.StdinPipe() // pipe to pass data to kubectl
	if err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	go func() {
		defer func() {
			err = errors.Join(err, stdin.Close())
		}()
		if errEW := kube.Export(km, kube.WithExportWriter(stdin), kube.WithExportAsSingleFile("stdin")); errEW != nil {
			err = errors.Join(err, errEW)
		}
	}()

	if errS := cmd.Start(); errS != nil {
		return errors.Join(err, errS)
	}

	// waits for the command to exit and waits for any copying
	// to stdin or copying from stdout or stderr to complete
	return errors.Join(err, cmd.Wait())
}
`

	tu.AssertEqual(t, buf.String(), want)
}
