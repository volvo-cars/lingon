// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"bytes"

	"github.com/rogpeppe/go-internal/txtar"
)

func Txtar2YAML(ar *txtar.Archive) []byte {
	var buf bytes.Buffer
	for _, f := range ar.Files {
		buf.WriteString("\n\n---\n")
		buf.WriteString("# " + f.Name + "\n")
		buf.Write(f.Data)
	}
	return buf.Bytes()
}
