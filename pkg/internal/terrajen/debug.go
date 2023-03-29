// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terrajen

import (
	"errors"
	"fmt"
	"go/scanner"
	"strings"
)

func JenDebug(err error) {
	es := err.Error()
	var se scanner.ErrorList
	if errors.As(err, &se) {
		for _, e := range se {
			// fmt.Println("line:", e.Pos.Line)
			// fmt.Println("column:", e.Pos.Column)
			// fmt.Println("message:", e.Msg)
			_, a, found := strings.Cut(es, e.Msg)
			if !found {
				continue
			}
			ss := strings.Split(a, "\n")
			for i, l := range ss {
				if i == e.Pos.Line {
					fmt.Printf("%d: %s  <== /!\\ ERROR: %s\n", i, l, e.Msg)
					fmt.Println(strings.Repeat(" ", e.Pos.Column-5) + "___^")
				} else {
					fmt.Printf("%d: %s\n", i, l)
				}
			}
		}
	}
}
