# Export Options


{{ "Example_export" | example }}


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


