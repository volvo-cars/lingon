// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

//go:build inttest

package vmk8s

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/grafana/dashboard-linter/lint"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
	"github.com/zeitlinger/conflate"
)

const (
	ghuc      = "https://raw.githubusercontent.com"
	vmRepo    = "/VictoriaMetrics/VictoriaMetrics/master"
	dotdcRepo = "/dotdc/grafana-dashboards-kubernetes/master"
)

var srcDash = []DashSource{
	// VICTORIA METRICS DASHBOARDS URLS
	// {
	// 	Name: "backupmanager.json",
	// 	URL:  ghuc + vmRepo + "/dashboards/backupmanager.json",
	// },
	// {
	// 	Name: "victoriametrics.json",
	// 	URL:  ghuc + vmRepo + "/dashboards/victoriametrics.json",
	// },
	// {
	// 	Name: "vmagent.json",
	// 	URL:  ghuc + vmRepo + "/dashboards/vmagent.json",
	// },
	// {
	// 	Name: "victoriametrics-cluster.json",
	// 	URL:  ghuc + vmRepo + "/dashboards/victoriametrics-cluster.json",
	// },
	// {
	// 	Name: "vmalert.json",
	// 	URL:  ghuc + vmRepo + "/dashboards/vmalert.json",
	// },
	// {
	// 	Name: "vm-operator.json",
	// 	URL:  ghuc + vmRepo + "/dashboards/operator.json",
	// },
	// // KUBERNETES DASHBOARDS URLS
	// {
	// 	Name: "k8s-system-api-server.json",
	// 	URL:  ghuc + dotdcRepo + "/dashboards/k8s-system-api-server.json",
	// },
	// {
	// 	Name: "k8s-system-coredns.json",
	// 	URL:  ghuc + dotdcRepo + "/dashboards/k8s-system-coredns.json",
	// },
	// {
	// 	Name: "k8s-views-global.json",
	// 	URL:  ghuc + dotdcRepo + "/dashboards/k8s-views-global.json",
	// },
	// {
	// 	Name: "k8s-views-namespaces.json",
	// 	URL:  ghuc + dotdcRepo + "/dashboards/k8s-views-namespaces.json",
	// },
	// {
	// 	Name: "k8s-views-nodes.json",
	// 	URL:  ghuc + dotdcRepo + "/dashboards/k8s-views-nodes.json",
	// },
	// {
	// 	Name: "k8s-views-pods.json",
	// 	URL:  ghuc + dotdcRepo + "/dashboards/k8s-views-pods.json",
	// },
	// {
	// 	Name: "node-exporter-full.json",
	// 	URL:  "https://grafana.com/api/dashboards/1860/revisions/22/download",
	// },
	// Karpenter dashboards
	// {
	// 	Name:   "karpenter-performance-dashboard.json",
	// 	URL:    "https://raw.githubusercontent.com/aws/karpenter/main/website/content/en/v0.28/getting-started/getting-started-with-karpenter/karpenter-performance-dashboard.json",
	// },
	// {
	// 	Name:   "karpenter-controllers.json",
	// 	URL:    "https://raw.githubusercontent.com/aws/karpenter/main/website/content/en/v0.28/getting-started/getting-started-with-karpenter/karpenter-controllers.json",
	// },
	// {
	// 	Name:   "karpenter-controllers-allocation.json",
	// 	URL:    "https://raw.githubusercontent.com/aws/karpenter/main/website/content/en/v0.28/getting-started/getting-started-with-karpenter/karpenter-controllers-allocation.json",
	// },
	// {
	// 	Name:   "karpenter-capacity-dashboard.json",
	// 	URL:    "https://raw.githubusercontent.com/aws/karpenter/main/website/content/en/v0.28/getting-started/getting-started-with-karpenter/karpenter-capacity-dashboard.json",
	// },
	{
		Name: "nats-dash.json",
		URL:  "https://github.com/nats-io/prometheus-nats-exporter/raw/main/walkthrough/grafana-nats-dash.json",
	},
	{
		Name: "nats-jetstream-dash.json",
		URL:  "https://github.com/nats-io/prometheus-nats-exporter/raw/main/walkthrough/grafana-jetstream-dash.json",
	},
}

func TestDashboardsDownload(t *testing.T) {
	c := http.Client{Timeout: 30 * time.Second}

	for _, src := range srcDash {
		resp, err := c.Get(src.URL)
		tu.AssertNoError(t, err, "url", src.URL)
		// buf, err := io.ReadAll(resp.Body)
		tu.AssertNoError(t, err, "read body", src.URL)
		file, err := os.Create(filepath.Join("dashboards", src.Name))
		tu.AssertNoError(t, err, "create file", src.Name)
		_, err = io.Copy(file, resp.Body)
		tu.AssertNoError(t, err, "copying", src.Name)
		_ = file.Close()
		_ = resp.Body.Close()
	}
}

func LintDashboards(path string, buf []byte, autofix bool) error {
	dashboard, err := lint.NewDashboard(buf)
	if err != nil {
		return fmt.Errorf("parse dashboard %s: %w", path, err)
	}

	config := &lint.ConfigurationFile{
		Exclusions: map[string]*lint.ConfigurationRuleEntries{
			"target-promql-rule":           {},
			"panel-title-description-rule": {},
		},
		Warnings: map[string]*lint.ConfigurationRuleEntries{},
		Verbose:  true,
		Autofix:  autofix,
	}

	rules := lint.NewRuleSet()
	results, err := rules.Lint([]lint.Dashboard{dashboard})
	if err != nil {
		return fmt.Errorf("lint dashboard %s: %w", path, err)
	}

	if config.Autofix {
		changes := results.AutoFix(&dashboard)
		if changes > 0 {
			err = writeDash(dashboard, path, buf)
			if err != nil {
				return err
			}
		}
	}

	results.Configure(config)
	results.ReportByRule()

	if results.MaximumSeverity() >= lint.Warning {
		return fmt.Errorf("linting errors")
	}
	return nil
}

func writeDash(dashboard lint.Dashboard, filename string, old []byte) error {
	newBytes, err := dashboard.Marshal()
	if err != nil {
		return err
	}
	c := conflate.New()
	err = c.AddData(old, newBytes)
	if err != nil {
		return err
	}
	b, err := c.MarshalJSON()
	if err != nil {
		return err
	}
	json := strings.ReplaceAll(
		string(b),
		"\"options\": null,",
		"\"options\": [],",
	)

	return os.WriteFile(filename, []byte(json), 0o600)
}
