// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package testutil

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReadYAML(filePath string) ([]string, error) {
	e := filepath.Ext(filePath)
	if !contains([]string{".yaml", ".yml"}, e) {
		return nil, fmt.Errorf("not yaml file: %s", filePath)
	}
	yf, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read manifest %s: %w", filePath, err)
	}
	scanner := bufio.NewScanner(bytes.NewReader(yf))
	var content []string
	var buf bytes.Buffer

	for scanner.Scan() {
		txt := scanner.Text()
		if strings.Contains(txt, "---") {
			if buf.Len() > 0 {
				content = append(content, buf.String())
				buf.Reset()
			}
		} else {
			buf.WriteString(txt + "\n")
		}
	}
	content = append(content, buf.String())

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return content, nil
}
