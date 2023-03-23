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
