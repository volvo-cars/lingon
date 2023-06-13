// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rogpeppe/go-internal/txtar"
	"github.com/volvo-cars/lingon/pkg/kube/testdata/go/tekton"

	// "github.com/volvo-cars/lingon/docs/kubernetes/kube/out/tekton"
	"github.com/volvo-cars/lingon/pkg/kube"
)

func Example() {
	tk := tekton.New()

	out := filepath.Join("out", "export")
	fmt.Printf("exporting to %s\n", out)

	_ = os.RemoveAll(out)
	defer os.RemoveAll(out)

	err := kube.Export(tk, kube.WithExportOutputDirectory(out))
	if err != nil {
		panic(err)
	}

	// or use io.Writer
	var buf bytes.Buffer
	_ = kube.Export(tk, kube.WithExportWriter(&buf))

	ar := txtar.Parse(buf.Bytes())

	fmt.Printf("\nexported %d manifests\n\n", len(ar.Files))
	fmt.Println("\t>>> first manifest <<<")
	if len(ar.Files) > 2 {
		// avoiding diff due to character invisible to the naked eye 😅
		l := strings.ReplaceAll(string(ar.Files[2].Data), "\n", "\n\t")
		fmt.Printf("\t%s\n", l)
	}

	// Output:
	// exporting to out/export
	//
	// exported 65 manifests
	//
	//	>>> first manifest <<<
	//	apiVersion: rbac.authorization.k8s.io/v1
	//	kind: ClusterRole
	//	metadata:
	//	  labels:
	//	    app.kubernetes.io/instance: default
	//	    app.kubernetes.io/part-of: tekton-pipelines
	//	    rbac.authorization.k8s.io/aggregate-to-admin: "true"
	//	    rbac.authorization.k8s.io/aggregate-to-edit: "true"
	//	  name: tekton-aggregate-edit
	//	rules:
	//	- apiGroups:
	//	  - tekton.dev
	//	  resources:
	//	  - tasks
	//	  - taskruns
	//	  - pipelines
	//	  - pipelineruns
	//	  - pipelineresources
	//	  - runs
	//	  - customruns
	//	  verbs:
	//	  - create
	//	  - delete
	//	  - deletecollection
	//	  - get
	//	  - list
	//	  - patch
	//	  - update
	//	  - watch
	//
	//
}
