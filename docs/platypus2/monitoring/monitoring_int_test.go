// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

//go:build inttest

package monitoring

import (
	"context"
	"flag"
	"os"
	"os/exec"
	"testing"
	"time"

	tu "github.com/volvo-cars/lingon/pkg/testutil"
	"github.com/volvo-cars/lingoneks/monitoring/metricsserver"
	"github.com/volvo-cars/lingoneks/monitoring/promcrd"
	"github.com/volvo-cars/lingoneks/monitoring/vmk8s"
	"github.com/volvo-cars/lingoneks/monitoring/vmop"
	"github.com/volvo-cars/lingoneks/monitoring/vmop/vmcrd"
	"sigs.k8s.io/kind/pkg/apis/config/v1alpha4"
	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cmd"
)

// How to run this test:
// 1. Docker and kind needs to be installed and working
// 2. in this folder, execute: `go test -v . -tags inttest`
//    note that without the -v, nothing appear on the screen
//
// to keep the cluster after the tests are done
// run : `go test -v . -keep-cluster=true -tags inttest`

var keepCluster = flag.Bool(
	"keep-cluster",
	false,
	"keep the cluster after the tests are done",
)

func TestKindDeploy(t *testing.T) {
	kubeConfigPath := "out/kind-test-monitoring.yaml"
	clusterName := "test-monitoring"

	provider := cluster.NewProvider(
		cluster.ProviderWithDocker(),
		cluster.ProviderWithLogger(cmd.NewLogger()),
	)
	t.Cleanup(
		func() {
			if *keepCluster {
				t.Logf(
					`

	# Set the kubeconfig (use full path)
		export KUBECONFIG=$(pwd)/%s

	# access the cluster
		kubectl get pod -A

	# to delete the cluster
		kind delete cluster --name %s

`, kubeConfigPath, clusterName,
				)
			} else {
				t.Log("deleting cluster", clusterName)
				if err := provider.Delete(
					clusterName,
					kubeConfigPath,
				); err != nil {
					t.Log("delete cluster", clusterName, "err:", err)
				}
			}
		},
	)

	// documented at https://kind.sigs.k8s.io/docs/user/configuration/
	kindConfig := &v1alpha4.Cluster{
		TypeMeta: v1alpha4.TypeMeta{
			Kind: "Cluster", APIVersion: "kind.x-k8s.io/v1alpha4",
		},
		// Name: clusterName, // will be overridden by provider.Create(clusterName)
		Nodes: []v1alpha4.Node{
			{Role: v1alpha4.ControlPlaneRole},
			{Role: v1alpha4.WorkerRole},
		},
	}
	err := provider.Create(
		clusterName,
		cluster.CreateWithV1Alpha4Config(kindConfig),
		cluster.CreateWithDisplayUsage(true),
		cluster.CreateWithDisplaySalutation(true),
		cluster.CreateWithWaitForReady(30*time.Second),
	)
	if err != nil {
		t.Errorf("unable to create kind test cluster: %v", err)
	}
	// If internal is true, this will contain the internal IP etc.
	// If internal is false, this will contain the host IP etc.
	// We want the host IP as we connect from outside the cluster.
	// internal := false
	err = provider.ExportKubeConfig(clusterName, kubeConfigPath, false)
	if err != nil {
		t.Errorf("unable to export test kube config: %v", err)
	}

	tu.AssertNoError(t, os.Setenv("KUBECONFIG", kubeConfigPath), "set env ")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// as it runs on kind, the certificates are self-signed,
	// it needs to be insecure but only in test/kind
	msInsec := func(m *metricsserver.MetricsServer) *metricsserver.MetricsServer {
		m.Deploy = metricsserver.PatchDeployInsecureTLS(m.Deploy)
		return m
	}

	type Applyer interface {
		Lingon()
		Apply(ctx2 context.Context) error
	}
	type deployApp struct {
		Name string
		App  Applyer
	}

	vm := vmk8s.New()
	tests := []deployApp{
		{Name: "promcrd", App: promcrd.New()},
		{Name: "vmcrd", App: vmcrd.New()},
		{Name: "vmop", App: vmop.New()},
		{Name: "metrics-server", App: metricsserver.New(msInsec)},
		{Name: "vmk8s", App: vm},
	}
	for _, da := range tests {
		t.Log("applying", da.Name)
		tu.AssertNoError(t, da.App.Apply(ctx), da.Name)
	}

	err = kubectl(
		t, ctx,
		"wait",
		"--namespace", vm.Grafana.Deploy.Namespace,
		// deploy/grafana
		vm.Grafana.Deploy.TypeMeta.Kind+"/"+vm.Grafana.Deploy.Name,
		"--for=condition=available",
		"--timeout=60s",
	)
	tu.AssertNoError(t, err, "kubectl wait for grafana deployment")
}

func kubectl(
	t *testing.T,
	ctx context.Context,
	args ...string,
) error {
	t.Helper()
	c := exec.CommandContext(ctx, "kubectl", args...)
	c.Env = os.Environ() // inherit environment in case we need to use kubectl from a container
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	err := c.Run()
	// o, err := c.CombinedOutput()
	// t.Log("kubectl output:", string(o))
	return err
}
