# Export Options


```go
package crd_test

import (
	"fmt"
	"os"
	"path/filepath"

	team "github.com/volvo-cars/lingon/docs/kubernetes/crd/out"
	"github.com/volvo-cars/lingon/pkg/kube"
)

var defaultOut = "out"

// Example_export to shows how to export to kubernetes manifests files in YAML.
func Example_export() {
	out := filepath.Join(defaultOut, "manifests")

	// Remove previously generated output directory
	_ = os.RemoveAll(out)

	tm := team.New()
	if err := kube.Export(
		tm,
		// directory where to export the manifests
		kube.WithExportOutputDirectory(out),

		// // Write the manifests to a bytes.Buffer.
		// // Note that it will be written as txtar format.
		// kube.WithExportWriter(buf),
		// // Add a kustomization.yaml file next to the manifests.
		// kube.WithExportKustomize(true),

		// // Write the manifest as a single file.
		// kube.WithExportAsSingleFile("manifest.yaml"),

		// // Write the manifests as JSON instead of YAML.
		// // Written as an JSON array if written as a single file.
		// kube.WithExportOutputJSON(true),

		// // Write to standard output instead of a file.
		// kube.WithExportStdOut(),

		// // Write the manifests as it would appear on the cluster.
		// // 1 resource per file and one folder per namespace.
		// kube.WithExportExplodeManifests(true),

		// // Define a hook to be called before exporting a secret.
		// // Note that this will remove the secret from the output.
		// kube.WithExportSecretHook(
		// 	func(s *corev1.Secret) error {
		// 		// do something with the secret
		// 		return nil
		// 	}),
	); err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Println("successfully exported CRDs to manifests")

	// Output:
	// successfully exported CRDs to manifests
}
```


## Note with io.Writer

If you need to manipulate the manifest once marshaled from Go, the option `WithExportWriter` allows for an `io.Writer`
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


