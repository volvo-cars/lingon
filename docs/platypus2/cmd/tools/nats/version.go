// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
)

var (
	// ldflags
	version = "dev"
	commit  = "dev"
	date    = "dev"
)

type V struct {
	Version   string          `json:"version"`
	Commit    string          `json:"commit"`
	Date      string          `json:"date"`
	BuildInfo debug.BuildInfo `json:"buildInfo"`
}

func VersionInfo(w http.ResponseWriter, r *http.Request) {
	// Set the CORS headers to the response.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set(
		"Access-Control-Allow-Methods",
		"GET, OPTIONS",
	)
	w.Header().Set(
		"Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
	)
	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		_, _ = fmt.Fprintln(os.Stderr, "error reading build-info")
		_, _ = fmt.Fprintln(w, "error reading build-info")
	}
	v := V{
		Version:   version,
		Commit:    commit,
		Date:      date,
		BuildInfo: *bi,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error encoding version: %v", err)
		_, _ = fmt.Fprintf(w, "error encoding version: %v", err)
	}
}
