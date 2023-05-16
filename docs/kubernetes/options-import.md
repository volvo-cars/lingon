# Import Options for Kubernetes YAML to Go

Various settings to convert kubernetes YAML to Go.

- [Example with YAML files containing CustomResourceDefinition (CRDs)](#example-with-yaml-files-containing-customresourcedefinition-crds)
- [Example with io.Writer](#example-with-iowriter)

## Example with YAML files containing CustomResourceDefinition (CRDs)

```go
out := filepath.Join("gocode", "tekton")
_ = os.RemoveAll(out)
defer os.RemoveAll(out)
err := kube.Import(
	// name of the application
	kube.WithImportAppName("tekton"),
	// name of the Go package where the code is generated to
	kube.WithImportPackageName("tekton"),
	// the directory to write the generated code to
	kube.WithImportOutputDirectory(out),
	// the list of manifest files to read and convert
	kube.WithImportManifestFiles(
		[]string{filepath.Join(testdata, "tekton.yaml")},
	),
	// define the types for the CRDs
	kube.WithImportSerializer(
		func() runtime.Decoder {
			_ = apiextensions.AddToScheme(scheme.Scheme)
			_ = apiextensionsbeta.AddToScheme(scheme.Scheme)
			return scheme.Codecs.UniversalDeserializer()
		}(),
	),
	// will try to remove "tekton" from the name of the variable in the Go code, make them shorter
	kube.WithImportRemoveAppName(true),
	// group all the resources from the same kind into one file each
	// example: 10 ConfigMaps => 1 file "config-map.go" containing 10 variables storing ConfigMap, etc...
	kube.WithImportGroupByKind(true),
	// add convenience methods to the App struct
	kube.WithImportAddMethods(true),
	// do not print verbose information
	kube.WithImportVerbose(false),
	// do not ignore errors
	kube.WithImportIgnoreErrors(false),
	// just for example purposes
	// how to create a logger (see [golang.org/x/tools/slog](https://golang.org/x/tools/slog))
	// this has no effect with WithImportVerbose(false)
	kube.WithImportLogger(kube.Logger(os.Stderr)),
	// remove the status field and
	// other output-only fields from the manifest code before importing it.
	// Note that ConfigMap are not cleaned up as the comments will be lost.
	kube.WithImportCleanUp(true),
)
if err != nil {
	panic(fmt.Errorf("import: %w", err))
}
got, err := kubeutil.ListGoFiles(out)
if err != nil {
	panic(fmt.Errorf("list go files: %w", err))
}
// sort the files to make the output deterministic
sort.Strings(got)

for _, f := range got {
	fmt.Println(f)
}

// Output:
//
// gocode/tekton/app.go
// gocode/tekton/cluster-role-binding.go
// gocode/tekton/cluster-role.go
// gocode/tekton/config-map.go
// gocode/tekton/custom-resource-definition.go
// gocode/tekton/deployment.go
// gocode/tekton/horizontal-pod-autoscaler.go
// gocode/tekton/mutating-webhook-configuration.go
// gocode/tekton/namespace.go
// gocode/tekton/role-binding.go
// gocode/tekton/role.go
// gocode/tekton/secret.go
// gocode/tekton/service-account.go
// gocode/tekton/service.go
// gocode/tekton/validating-webhook-configuration.go
```

Another example:

```go
package crd_test

import (
	"fmt"
	"os"

	"github.com/volvo-cars/lingon/pkg/kube"
	istionetworkingv1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	secretsstorev1 "sigs.k8s.io/secrets-store-csi-driver/apis/v1"
)

func defaultSerializer() runtime.Decoder {
	// Add CRDs to scheme
	// This is needed to be able to import CRDs from kubernetes manifests files.
	_ = apiextensions.AddToScheme(scheme.Scheme)
	_ = apiextensionsv1.AddToScheme(scheme.Scheme)
	_ = secretsstorev1.AddToScheme(scheme.Scheme)
	_ = istionetworkingv1beta1.AddToScheme(scheme.Scheme)
	return scheme.Codecs.UniversalDeserializer()
}

// Example_import to shows how to import CRDs from kubernetes manifests files.
func Example_import() {
	// Remove previously generated output directory
	_ = os.RemoveAll("./out")

	if err := kube.Import(
		kube.WithImportAppName("team"),
		kube.WithImportPackageName("team"),
		kube.WithImportManifestFiles([]string{"./manifest.yaml"}),
		kube.WithImportOutputDirectory("./out"),
		kube.WithImportSerializer(defaultSerializer()),
		// do not print verbose information
		kube.WithImportVerbose(false),
		// do not ignore errors
		kube.WithImportIgnoreErrors(false),
		// use the default logger even if unused with WithImportVerbose(false)
		kube.WithImportLogger(kube.Logger(os.Stderr)),
	); err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Println("successfully imported CRDs from manifest.yaml")

	// Output:
	// successfully imported CRDs from manifest.yaml
}
```

## Example with io.Writer

If you need to manipulate the generated Go code, the option `WithImportWriter` allows for an `io.Writer`
to be passed and all the code will be written to it.

> NOTE: the output format is called [txtar](https://pkg.go.dev/golang.org/x/tools/txtar) format
>
> A txtar archive is zero or more comment lines and then a sequence of file entries.
> Each file entry begins with a file marker line of the form "-- FILENAME --" and
> is followed by zero or more file content lines making up the file data.
> The comment or file content ends at the next file marker line.
> The file marker line must begin with the three-byte sequence "-- " and
> end with the three-byte sequence " --", but the enclosed file name can be
> surrounding by additional white space, all of which is stripped.
>
> If the txtar file is missing a trailing newline on the final line,
> parsers should consider a final newline to be present anyway.
>
> There are no possible syntax errors in a txtar archive.

Code:

```go
filename := filepath.Join(testdata, "grafana.yaml")
file, _ := os.Open(filename)

var buf bytes.Buffer

err := kube.Import(
	kube.WithImportAppName("grafana"),
	kube.WithImportPackageName("grafana"),
	// will prefix all files with path "manifests/"
	kube.WithImportOutputDirectory("manifests/"),
	// we could just use kube.WithImportManifestFiles([]string{filename})
	// but this is just an example to show how to use WithImportReader
	// and WithImportWriter
	kube.WithImportReader(file),
	kube.WithImportWriter(&buf),
	// We don't want to group the resources by kind,
	// each file will contain a single resource
	kube.WithImportGroupByKind(false),
	// We rename the files to avoid name collisions.
	// Tip: use the Kind and Name of the resource to
	// create a unique name and avoid collision.
	//
	// Here, we didn't use WithImportGroupByKind,
	// each file will contain a single resource.
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
if err != nil {
	panic("failed to import")
}

// the output contained in bytes.Buffer is in the txtar format
// for more details, see https://pkg.go.dev/golang.org/x/tools/txtar
ar := txtar.Parse(buf.Bytes())

// sort the files to make the output deterministic
sort.SliceStable(
	ar.Files, func(i, j int) bool {
		return ar.Files[i].Name < ar.Files[j].Name
	},
)
for _, f := range ar.Files {
	fmt.Println(f.Name)
}
// Output:
//
// manifests/app.go
// manifests/clusterrole-grafana-clusterrole.go
// manifests/clusterrolebinding-grafana-clusterrolebinding.go
// manifests/configmap-grafana-dashboards-default.go
// manifests/configmap-grafana-test.go
// manifests/configmap-grafana.go
// manifests/deployment-grafana.go
// manifests/pod-grafana-test.go
// manifests/role-grafana.go
// manifests/rolebinding-grafana.go
// manifests/secret-grafana.go
// manifests/service-grafana.go
// manifests/serviceaccount-grafana-test.go
// manifests/serviceaccount-grafana.go
```
