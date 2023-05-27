// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package nats

import (
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO: template the config

var natsConfig = map[string]string{
	"nats.conf": `
# NATS Clients Port
port: 4222
# PID file shared with configuration reloader.
pid_file: "/var/run/nats/nats.pid"
###############
#             #
# Monitoring  #
#             #
###############
http: 8222
server_name:$POD_NAME
server_tags: [
    "4GiB",
]
###################################
#                                 #
# NATS Full Mesh Clustering Setup #
#                                 #
###################################
cluster {
  port: 6222
  routes = [
    nats://nats-0.nats.nats.svc.cluster.local:6222,nats://nats-1.nats.nats.svc.cluster.local:6222,nats://nats-2.nats.nats.svc.cluster.local:6222,
  ]
  cluster_advertise: $CLUSTER_ADVERTISE
  connect_retries: 120
}
lame_duck_grace_period: 10s
lame_duck_duration: 30s

`,
}

var cm = ku.ConfigAndMount{
	Data: natsConfig,
	ObjectMeta: metav1.ObjectMeta{
		Labels:    BaseLabels(),
		Name:      "nats-config",
		Namespace: namespace,
	},
	VolumeMount: corev1.VolumeMount{
		Name:      "config-volume",
		MountPath: "/etc/nats-config",
	},
}