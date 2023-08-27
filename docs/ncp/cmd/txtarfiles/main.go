package main

import (
	"flag"
	"ncp/bla"
	"os"
	"strings"

	"log/slog"

	"golang.org/x/tools/txtar"
)

func main() {
	var (
		dir string
		in  fileList
		out string
	)

	flag.StringVar(&dir, "dir", "", "relative directory for the input files")
	flag.Var(&in, "in", "file to include in the txtar")
	flag.StringVar(&out, "out", "", "name of the output file")
	flag.Parse()

	if len(in) == 0 {
		slog.Error("input file name is required")
		os.Exit(1)
	}
	if out == "" {
		slog.Error("output file name is required")
		os.Exit(1)
	}

	ar, err := bla.PackTxtar(dir, in)
	if err != nil {
		slog.Error("creating txtar", "error", err)
		os.Exit(1)
	}
	if err := os.WriteFile(out, txtar.Format(ar), os.ModePerm); err != nil {
		slog.Error("writing txtar", "error", err)
		os.Exit(1)
	}
}

type fileList []string

func (i *fileList) String() string {
	return strings.Join(*i, ", ")
}

func (i *fileList) Set(value string) error {
	*i = append(*i, value)
	return nil
}
