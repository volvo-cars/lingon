// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeutil

const (
	PortNameMetrics   = "metrics"
	PortNameProfiling = "profiling"
	PortNameProbes    = "probes"
	PortNameHTTP      = "http"
)

const (
	PathMetrics   = "/metrics"
	PathProfiling = "/debug/pprof"
	PathProbes    = "/healthz"
	PathReadiness = "/readiness"
)
