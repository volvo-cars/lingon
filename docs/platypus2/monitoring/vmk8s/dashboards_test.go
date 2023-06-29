// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

//go:build inttest

package vmk8s

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

const (
	ghuc      = "https://raw.githubusercontent.com"
	vmRepo    = "/VictoriaMetrics/VictoriaMetrics/master"
	dotdcRepo = "/dotdc/grafana-dashboards-kubernetes/master"
)

var srcDash = []DashSource{
	// VICTORIA METRICS DASHBOARDS URLS
	{
		Name: "backupmanager.json",
		URL:  ghuc + vmRepo + "/dashboards/backupmanager.json",
	},
	{
		Name: "victoriametrics.json",
		URL:  ghuc + vmRepo + "/dashboards/victoriametrics.json",
	},
	{
		Name: "vmagent.json",
		URL:  ghuc + vmRepo + "/dashboards/vmagent.json",
	},
	{
		Name: "victoriametrics-cluster.json",
		URL:  ghuc + vmRepo + "/dashboards/victoriametrics-cluster.json",
	},
	{
		Name: "vmalert.json",
		URL:  ghuc + vmRepo + "/dashboards/vmalert.json",
	},
	{
		Name: "vm-operator.json",
		URL:  ghuc + vmRepo + "/dashboards/operator.json",
	},
	// KUBERNETES DASHBOARDS URLS
	{
		Name: "k8s-system-api-server.json",
		URL:  ghuc + dotdcRepo + "/dashboards/k8s-system-api-server.json",
	},
	{
		Name: "k8s-system-coredns.json",
		URL:  ghuc + dotdcRepo + "/dashboards/k8s-system-coredns.json",
	},
	{
		Name: "k8s-views-global.json",
		URL:  ghuc + dotdcRepo + "/dashboards/k8s-views-global.json",
	},
	{
		Name: "k8s-views-namespaces.json",
		URL:  ghuc + dotdcRepo + "/dashboards/k8s-views-namespaces.json",
	},
	{
		Name: "k8s-views-nodes.json",
		URL:  ghuc + dotdcRepo + "/dashboards/k8s-views-nodes.json",
	},
	{
		Name: "k8s-views-pods.json",
		URL:  ghuc + dotdcRepo + "/dashboards/k8s-views-pods.json",
	},
}

func TestDashboardsDownload(t *testing.T) {
	c := http.Client{Timeout: 30 * time.Second}

	for _, src := range srcDash {
		resp, err := c.Get(src.URL)
		tu.AssertNoError(t, err, "url", src.URL)
		defer resp.Body.Close()
		file, err := os.Create(filepath.Join("dashboards", src.Name))
		tu.AssertNoError(t, err, "create file", src.Name)
		defer file.Close()
		_, err = io.Copy(file, resp.Body)
		tu.AssertNoError(t, err, "copying", src.Name)
	}
}
