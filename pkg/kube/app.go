// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package kube

// Exporter interfaces for kubernetes objects defined in a Go structs
type Exporter interface {
	Dummy()
}

// _ is a dummy variable to make sure that App implements Exporter
var _ Exporter = (*App)(nil)

// App struct is meant to be embedded in other structs
// to specify that they are a set of kubernetes manifests
type App struct{}

// dummy is a dummy method to make sure that App implements Exporter
func (a *App) Dummy() {}
