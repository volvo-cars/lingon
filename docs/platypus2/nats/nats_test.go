// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package nats

import (
	"os"
	"strings"
	"testing"

	"github.com/volvo-cars/lingon/pkg/kubeutil"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
	"github.com/volvo-cars/lingoneks/meta"
	corev1 "k8s.io/api/core/v1"
)

func TestNatsManifestExport(t *testing.T) {
	_ = os.RemoveAll("out")
	n := New()
	if err := n.Export("out"); err != nil {
		tu.AssertNoError(t, err, "nats")
	}
}

// TODO: THIS IS INTEGRATION and needs KWOK

// func TestNatsDeploy(t *testing.T) {
// 	n := New()
// 	if err := n.Apply(context.Background()); err != nil {
// 		tu.AssertNoError(t, err, "nats")
// 	}
// }

func TestConfig(t *testing.T) {
	tests := []struct {
		name string
		args Meta
		want string
	}{
		{
			name: "empty",
			args: Meta{
				Metadata:       meta.Metadata{},
				Config:         kubeutil.ConfigAndMount{},
				ConfigFile:     "",
				ConfigPath:     "",
				Client:         meta.NetPort{},
				Cluster:        meta.NetPort{},
				Monitor:        meta.NetPort{},
				Metrics:        meta.NetPort{},
				Leaf:           meta.NetPort{},
				Gw:             meta.NetPort{},
				pidVM:          corev1.VolumeMount{},
				pidPath:        "",
				storageDir:     "",
				StorageClass:   "",
				PvcName:        "",
				replicas:       0,
				ConfigReloader: meta.ContainerImg{},
				PromExporter:   meta.ContainerImg{},
			},
			want: `
# NATS Clients Port
port: 0
# PID file shared with configuration reloader.
pid_file: ""
###############
#             #
# Monitoring  #
#             #
###############

http: 0
server_name:$POD_NAME
server_tags: [
  4Gi
]

###################################
#                                 #
# NATS JetStream                  #
#                                 #
###################################

jetstream {
  max_mem:2G
  store_dir: ""
  max_file:10Gi
  unique_tag: "natsuniquetag"
}


###################################
#                                 #
# NATS Full Mesh Clustering Setup #
#                                 #
###################################

cluster {
  name: natscluster
  port: 0
  routes = [


  ]
  cluster_advertise: $CLUSTER_ADVERTISE
  connect_retries: 120
}


lame_duck_grace_period: 10s
lame_duck_duration: 30s
`,
		},
		{
			name: "simple",
			args: Meta{
				Metadata: meta.Metadata{
					Name:      "nats",
					Namespace: "nats",
				},
				Config:     kubeutil.ConfigAndMount{},
				ConfigFile: "nats.conf",
				ConfigPath: "/etc/nats-config",
				Client: meta.NetPort{
					Container: corev1.ContainerPort{
						Name:          "containerport",
						ContainerPort: 12345,
					},
				},
				Cluster: meta.NetPort{
					Container: corev1.ContainerPort{
						Name:          "clusterport",
						ContainerPort: 12346,
					},
					Service: corev1.ServicePort{Port: 9876},
				},
				Monitor: meta.NetPort{
					Container: corev1.ContainerPort{
						Name:          "monport",
						ContainerPort: 54321,
					},
				},
				pidVM:      corev1.VolumeMount{},
				pidPath:    "mypath",
				storageDir: "/data/storage",
				replicas:   3,
			},
			want: `
# NATS Clients Port
port: 12345
# PID file shared with configuration reloader.
pid_file: "mypath"
###############
#             #
# Monitoring  #
#             #
###############

http: 54321
server_name:$POD_NAME
server_tags: [
  4Gi
]

###################################
#                                 #
# NATS JetStream                  #
#                                 #
###################################

jetstream {
  max_mem:2G
  store_dir: "/data/storage"
  max_file:10Gi
  unique_tag: "natsuniquetag"
}


###################################
#                                 #
# NATS Full Mesh Clustering Setup #
#                                 #
###################################

cluster {
  name: natscluster
  port: 12346
  routes = [

    nats://nats-0.nats.nats.svc.cluster.local:9876
    nats://nats-1.nats.nats.svc.cluster.local:9876
    nats://nats-2.nats.nats.svc.cluster.local:9876

  ]
  cluster_advertise: $CLUSTER_ADVERTISE
  connect_retries: 120
}


lame_duck_grace_period: 10s
lame_duck_duration: 30s
`,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := strings.Split(Config(tt.args), "\n")
				want := strings.Split(tt.want, "\n")
				if diff := tu.Diff(got, want); diff != "" {
					t.Error(tu.Callers(), "diff", diff)
				}
			},
		)
	}
}
