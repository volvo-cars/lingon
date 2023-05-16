# Import Options for Kubernetes YAML to Go

Various settings to convert kubernetes YAML to Go.

- [Example with YAML files containing CustomResourceDefinition (CRDs)](#example-with-yaml-files-containing-customresourcedefinition-crds)
- [Example with io.Writer](#example-with-iowriter)

## Example with YAML files containing CustomResourceDefinition (CRDs)

{{ "ExampleImport_withManifest" | example }}

Another example:

{{ "Example_import" | example }}

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

{{ "ExampleImport_withWriter" | example }}
