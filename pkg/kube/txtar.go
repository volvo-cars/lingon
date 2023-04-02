package kube

import (
	"bytes"

	"github.com/rogpeppe/go-internal/txtar"
)

func Txtar2YAML(ar *txtar.Archive) []byte {
	// ar := txtar.Parse(b)
	// txtar.Format(ar)
	var buf bytes.Buffer
	for _, f := range ar.Files {
		buf.WriteString("\n\n---\n")
		buf.WriteString("# " + f.Name + "\n")
		buf.Write(f.Data)
	}
	return buf.Bytes()
}
