// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package cilium

import (
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Config(data map[string]string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: kubeutil.TypeConfigMapV1,
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cilium-config",
			Namespace: "kube-system",
		},
		Data: data,
	}
}

var DefaultConfigData = map[string]string{
	// Identity allocation mode selects how identities are shared between cilium
	// nodes by setting how they are stored. The options are "crd" or "kvstore".
	// - "crd" stores identities in kubernetes as CRDs (custom resource definition).
	//   These can be queried with:
	//     kubectl get ciliumid
	// - "kvstore" stores identities in an etcd kvstore, that is
	//   configured below. Cilium versions before 1.6 supported only the kvstore
	//   backend. Upgrades from these older cilium versions should continue using
	//   the kvstore by commenting out the identity-allocation-mode below, or
	//   setting it to "kvstore".
	"identity-allocation-mode":      "crd",
	"cilium-endpoint-gc-interval":   "5m0s",
	"nodes-gc-interval":             "5m0s",
	"skip-cnp-status-startup-clean": "false",
	// Disable the usage of CiliumEndpoint CRD
	"disable-endpoint-crd": "false",

	// If you want to run cilium in debug mode change this value to true
	"debug": "false",
	// The agent can be put into the following three policy enforcement modes
	// default, always and never.
	// https://docs.cilium.io/en/latest/policy/intro///policy-enforcement-modes
	"enable-policy": "default",

	// Enable IPv4 addressing. If enabled, all endpoints are allocated an IPv4
	// address.
	"enable-ipv4": "true",

	// Enable IPv6 addressing. If enabled, all endpoints are allocated an IPv6
	// address.
	"enable-ipv6": "false",
	// Users who wish to specify their own custom CNI configuration file must set
	// custom-cni-conf to "true", otherwise Cilium may overwrite the configuration.
	"custom-cni-conf":        "false",
	"enable-bpf-clock-probe": "true",
	// If you want cilium monitor to aggregate tracing for packets, set this level
	// to "low", "medium", or "maximum". The higher the level, the less packets
	// that will be seen in monitor output.
	"monitor-aggregation": "medium",

	// The monitor aggregation interval governs the typical time between monitor
	// notification events for each allowed connection.
	//
	// Only effective when monitor aggregation is set to "medium" or higher.
	"monitor-aggregation-interval": "5s",

	// The monitor aggregation flags determine which TCP flags which, upon the
	// first observation, cause monitor notifications to be generated.
	//
	// Only effective when monitor aggregation is set to "medium" or higher.
	"monitor-aggregation-flags": "all",
	// Specifies the ratio (0.0-1.0) of total system memory to use for dynamic
	// sizing of the TCP CT, non-TCP CT, NAT and policy BPF maps.
	"bpf-map-dynamic-size-ratio": "0.0025",
	// bpf-policy-map-max specifies the maximum number of entries in endpoint
	// policy map (per endpoint)
	"bpf-policy-map-max": "16384",
	// bpf-lb-map-max specifies the maximum number of entries in bpf lb service,
	// backend and affinity maps.
	"bpf-lb-map-max": "65536",
	// bpf-lb-bypass-fib-lookup instructs Cilium to enable the FIB lookup bypass
	// optimization for nodeport reverse NAT handling.
	"bpf-lb-external-clusterip": "false",

	// Pre-allocation of map entries allows per-packet latency to be reduced, at
	// the expense of up-front memory allocation for the entries in the maps. The
	// default value below will minimize memory usage in the default installation;
	// users who are sensitive to latency may consider setting this to "true".
	//
	// This option was introduced in Cilium 1.4. Cilium 1.3 and earlier ignore
	// this option and behave as though it is set to "true".
	//
	// If this value is modified, then during the next Cilium startup the restore
	// of existing endpoints and tracking of ongoing connections may be disrupted.
	// As a result, reply packets may be dropped and the load-balancing decisions
	// for established connections may change.
	//
	// If this option is set to "false" during an upgrade from 1.3 or earlier to
	// 1.4 or later, then it may cause one-time disruptions during the upgrade.
	"preallocate-bpf-maps": "false",

	// Regular expression matching compatible Istio sidecar istio-proxy
	// container image names
	"sidecar-istio-proxy-image": "cilium/istio_proxy",

	// Name of the cluster. Only relevant when building a mesh of clusters.
	"cluster-name": "go-terriyaki-test-1",
	// Unique ID of the cluster. Must be unique across all conneted clusters and
	// in the range of 1 and 255. Only relevant when building a mesh of clusters.
	"cluster-id": "1",

	// Encapsulation mode for communication between nodes
	// Possible values:
	//   - disabled
	//   - vxlan (default)
	//   - geneve
	"tunnel":                           "disabled",
	"enable-endpoint-routes":           "true",
	"auto-create-cilium-node-resource": "true",
	"ec2-api-endpoint":                 "",
	"eni-tags":                         "{}",
	"subnet-ids-filter":                "",
	"subnet-tags-filter":               "",
	"instance-tags-filter":             "",
	// Enables L7 proxy for L7 policy enforcement and visibility
	"enable-l7-proxy": "true",

	"enable-ipv4-masquerade":       "true",
	"enable-ipv6-masquerade":       "true",
	"egress-masquerade-interfaces": "eth0",

	"enable-xt-socket-fallback":           "true",
	"install-iptables-rules":              "true",
	"install-no-conntrack-iptables-rules": "false",

	"auto-direct-node-routes":      "false",
	"enable-local-redirect-policy": "false",

	"kube-proxy-replacement":              "disabled",
	"bpf-lb-sock":                         "false",
	"enable-health-check-nodeport":        "true",
	"node-port-bind-protection":           "true",
	"enable-auto-protect-node-port-range": "true",
	"enable-svc-source-range-check":       "true",
	"enable-l2-neigh-discovery":           "true",
	"arping-refresh-period":               "30s",
	"enable-endpoint-health-checking":     "true",
	"enable-health-checking":              "true",
	"enable-well-known-identities":        "false",
	"enable-remote-node-identity":         "true",
	"synchronize-k8s-nodes":               "true",
	"operator-api-serve-addr":             "127.0.0.1:9234",
	// Enable Hubble gRPC service.
	"enable-hubble": "true",
	// UNIX domain socket for Hubble server to listen to.
	"hubble-socket-path": "/var/run/cilium/hubble.sock",
	// An additional address for Hubble server to listen to (e.g. ":4244").
	"hubble-listen-address":                   "4244",
	"hubble-disable-tls":                      "false",
	"hubble-tls-cert-file":                    "/var/lib/cilium/tls/hubble/server.crt",
	"hubble-tls-key-file":                     "/var/lib/cilium/tls/hubble/server.key",
	"hubble-tls-client-ca-files":              "/var/lib/cilium/tls/hubble/client-ca.crt",
	"ipam":                                    "eni",
	"disable-cnp-status-updates":              "true",
	"enable-vtep":                             "false",
	"vtep-endpoint":                           "",
	"vtep-cidr":                               "",
	"vtep-mask":                               "",
	"vtep-mac":                                "",
	"enable-bgp-control-plane":                "false",
	"procfs":                                  "/host/proc",
	"bpf-root":                                "/sys/fs/bpf",
	"cgroup-root":                             "/run/cilium/cgroupv2",
	"enable-k8s-terminating-endpoint":         "true",
	"remove-cilium-node-taints":               "true",
	"set-cilium-is-up-condition":              "true",
	"unmanaged-pod-watcher-interval":          "15",
	"tofqdns-dns-reject-response-code":        "refused",
	"tofqdns-enable-dns-compression":          "true",
	"tofqdns-endpoint-max-ip-per-hostname":    "50",
	"tofqdns-idle-connection-grace-period":    "0s",
	"tofqdns-max-deferred-connection-deletes": "10000",
	"tofqdns-min-ttl":                         "3600",
	"tofqdns-proxy-response-max-delay":        "100ms",
	"agent-not-ready-taint-key":               "node.cilium.io/agent-not-ready",
}
