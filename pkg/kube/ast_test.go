package kube

import (
	"os"
	"strings"
	"testing"

	tu "github.com/volvo-cars/go-terriyaki/pkg/testutil"
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
			manifest: "testdata/deployment.yaml",
			golden:   "testdata/deployment.golden",
		},
		{
			name:     "service",
			manifest: "testdata/service.yaml",
			golden:   "testdata/service.golden",
		},
		{
			name:     "secret",
			manifest: "testdata/secret.yaml",
			golden:   "testdata/secret.golden",
			redact:   true,
		},
		{
			name:     "empty configmap",
			manifest: "testdata/configmap.yaml",
			golden:   "testdata/configmap.golden",
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
	want := string(golden)
	return want
}

func convert(t *testing.T, obj runtime.Object, redact bool) string {
	t.Helper()
	j := jamel{o: Option{RedactSecrets: redact}}
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
