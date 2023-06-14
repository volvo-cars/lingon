package nats

import (
	"fmt"

	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/timeseries"
	"github.com/K-Phoen/grabana/variable/interval"
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	"github.com/volvo-cars/lingoneks/monitoring"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	DashboardNatsCM = &corev1.ConfigMap{
		TypeMeta: ku.TypeConfigMapV1,
		ObjectMeta: metav1.ObjectMeta{
			Name:      "grafana-dashboard-nats",
			Namespace: monitoring.Namespace,
			Labels:    BaseLabels(), // TODO: check
		},
		Data: map[string]string{
			"nats.json": fmt.Sprintf("%s", DashboardNats),
		},
	}

	DashboardNats = DashMust(
		"NATS prometheus",
		dashboard.AutoRefresh("5s"),
		dashboard.Tags([]string{"generated"}),
		dashboard.VariableAsInterval(
			"interval",
			interval.Values(
				[]string{"30s", "1m", "5m", "10m", "30m", "1h", "6h", "12h"},
			),
		),
		dashboard.Row(
			"Prometheus",
			row.WithTimeSeries(
				"HTTP Rate",
				timeseries.DataSource("prometheus-default"),
				timeseries.WithPrometheusTarget(
					"rate(prometheus_http_requests_total[30s])",
					prometheus.Legend("{{handler}} - {{ code }}"),
				),
			),
		),
	)
)

func DashMust(name string, opts ...dashboard.Option) []byte {
	d, err := dashboard.New(name, opts...)
	if err != nil {
		panic(fmt.Sprintf("dashboard new: %s", err))
	}
	j, err := d.MarshalIndentJSON()
	if err != nil {
		panic(fmt.Sprintf("dashboard marshal: %s", err))
	}
	return j
}
