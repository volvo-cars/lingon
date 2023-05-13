// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"io"
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
)

// Logger returns a logger that writes to w [io.Writer]. If w is nil, os.Stderr is used.
// Timestamp is removed and directory from the source's filename is shown.
func Logger(w io.Writer) *slog.Logger {
	if w == nil {
		w = os.Stderr
	}
	return slog.New(
		slog.NewTextHandler(
			w,
			&slog.HandlerOptions{
				AddSource:   true,
				ReplaceAttr: logReplace,
			},
		).WithAttrs(
			[]slog.Attr{slog.String("app", defaultAppName)},
		),
	)
}

func logReplace(groups []string, a slog.Attr) slog.Attr {
	// Remove time.
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return slog.Attr{}
	}
	// Remove the directory from the source's filename.
	if a.Key == slog.SourceKey {
		source := a.Value.Any().(*slog.Source)
		source.File = filepath.Base(source.File)
	}
	return a
}
