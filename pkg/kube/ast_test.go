// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"os"
	"strings"
	"testing"

	tu "github.com/volvo-cars/lingon/pkg/testutil"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestKube2GoJen(t *testing.T) {
	type TT struct {
		name     string
		manifest string
		golden   string
		redact   bool
	}
	tests := []TT{
		{
			name:     "deployment",
			manifest: "testdata/golden/deployment.yaml",
			golden:   "testdata/golden/deployment.golden",
		},
		{
			name:     "service",
			manifest: "testdata/golden/service.yaml",
			golden:   "testdata/golden/service.golden",
		},
		{
			name:     "secret",
			manifest: "testdata/golden/secret.yaml",
			golden:   "testdata/golden/secret.golden",
			redact:   true,
		},
		{
			name:     "empty configmap",
			manifest: "testdata/golden/configmap.yaml",
			golden:   "testdata/golden/configmap.golden",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				obj := objectFromManifest(t, tt.manifest)
				got := convert(t, obj, tt.redact)
				want := readGolden(t, tt.golden)
				if diff := tu.Diff(got, want); diff != "" {
					t.Error(tu.Callers(), diff)
				}
			},
		)
	}
}

func readGolden(t *testing.T, path string) string {
	t.Helper()
	golden, err := os.ReadFile(path)
	tu.AssertNoError(t, err, "read golden file")
	return string(golden)
}

func convert(t *testing.T, obj runtime.Object, redact bool) string {
	t.Helper()
	j := jamel{o: option{RedactSecrets: redact}}
	code := j.kube2GoJen(obj)
	var b strings.Builder
	err := code.Render(&b)
	tu.AssertNoError(t, err, "render code")
	return b.String()
}

func objectFromManifest(t *testing.T, path string) runtime.Object {
	t.Helper()
	data, err := os.ReadFile(path)
	tu.AssertNoError(t, err, "read manifest")

	serializer := scheme.Codecs.UniversalDeserializer()
	obj, _, err := serializer.Decode(data, nil, nil)
	tu.AssertNoError(t, err, "decode manifest")
	return obj
}
