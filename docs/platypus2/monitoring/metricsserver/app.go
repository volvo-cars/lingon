// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package metricsserver

import (
	"context"
	"errors"
	"os"
	"os/exec"

	promoperator "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/volvo-cars/lingon/pkg/kube"
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	"github.com/volvo-cars/lingoneks/meta"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
)

// validate the struct implements the interface
var (
	_  kube.Exporter = (*MetricsServer)(nil)
	M                = Core()
	SA               = M.ServiceAccount()
)

func Core() Meta {
	n := "metrics-server"
	ns := "monitoring"
	ver := "0.6.3"

	return Meta{
		Metadata: meta.Metadata{
			Name:      n,
			Namespace: ns,
			Instance:  n + "-" + ns,
			Component: "metrics",
			PartOf:    n,
			Version:   ver,
			ManagedBy: "lingon",
			Img: meta.ContainerImg{
				Registry: "registry.k8s.io",
				Image:    "metrics-server/metrics-server",
				Tag:      "v" + ver,
			},
		},
		P: meta.NetPort{
			Container: corev1.ContainerPort{
				Name:          "https",
				ContainerPort: 10250,
				Protocol:      corev1.ProtocolTCP,
			},

			Service: corev1.ServicePort{
				Name:       "https",
				Protocol:   corev1.ProtocolTCP,
				Port:       443,
				TargetPort: intstr.FromString("https"),
			},
		},
		MetricsURL: ku.PathMetrics,
	}
}

type Meta struct {
	meta.Metadata
	P          meta.NetPort
	MetricsURL string
}

// MetricsServer contains kubernetes manifests
type MetricsServer struct {
	kube.App

	NS                        *corev1.Namespace
	AuthReaderRB              *rbacv1.RoleBinding
	Deploy                    *appsv1.Deployment
	SA                        *corev1.ServiceAccount
	SVC                       *corev1.Service
	ServiceMonitor            *promoperator.ServiceMonitor
	SystemAggregatedReaderCR  *rbacv1.ClusterRole
	SystemAuthDelegatorCRB    *rbacv1.ClusterRoleBinding
	SystemCR                  *rbacv1.ClusterRole
	SystemCRB                 *rbacv1.ClusterRoleBinding
	V1Beta1MetricsAPIServices *apiregistrationv1.APIService
}

type MetricsServerOption func(server *MetricsServer) *MetricsServer

// New creates a new MetricsServer
func New(opts ...MetricsServerOption) *MetricsServer {
	m := &MetricsServer{
		NS:                        ku.Namespace(M.Namespace, nil, nil),
		AuthReaderRB:              AuthReaderRB,
		Deploy:                    Deploy,
		SA:                        SA,
		SVC:                       SVC,
		ServiceMonitor:            ServiceMonitor,
		SystemAggregatedReaderCR:  SystemAggregatedReaderCR,
		SystemAuthDelegatorCRB:    SystemAuthDelegatorCRB,
		SystemCR:                  SystemCR,
		SystemCRB:                 SystemCRB,
		V1Beta1MetricsAPIServices: MetricsAPIServices,
	}

	for _, o := range opts {
		m = o(m)
	}

	return m
}

// Apply applies the kubernetes objects to the cluster
func (a *MetricsServer) Apply(ctx context.Context) error {
	return Apply(ctx, a)
}

// Export exports the kubernetes objects to YAML files in the given directory
func (a *MetricsServer) Export(dir string) error {
	return kube.Export(a, kube.WithExportOutputDirectory(dir))
}

// Apply applies the kubernetes objects contained in Exporter to the cluster
func Apply(ctx context.Context, km kube.Exporter) error {
	cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", "-")
	cmd.Env = os.Environ()        // inherit environment in case we need to use kubectl from a container
	stdin, err := cmd.StdinPipe() // pipe to pass data to kubectl
	if err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	go func() {
		defer func() {
			err = errors.Join(err, stdin.Close())
		}()
		if errEW := kube.Export(
			km,
			kube.WithExportWriter(stdin),
			kube.WithExportAsSingleFile("stdin"),
		); errEW != nil {
			err = errors.Join(err, errEW)
		}
	}()

	if errS := cmd.Start(); errS != nil {
		return errors.Join(err, errS)
	}

	// waits for the command to exit and waits for any copying
	// to stdin or copying from stdout or stderr to complete
	return errors.Join(err, cmd.Wait())
}

// P converts T to *T, useful for basic types
func P[T any](t T) *T {
	return &t
}

var MetricsAPIServices = &apiregistrationv1.APIService{
	TypeMeta: ku.TypeAPIServiceV1,
	ObjectMeta: metav1.ObjectMeta{
		Labels: M.Labels(),
		Name:   "v1beta1.metrics.k8s.io",
	},
	Spec: apiregistrationv1.APIServiceSpec{
		Group:                 "metrics.k8s.io",
		GroupPriorityMinimum:  int32(100),
		InsecureSkipTLSVerify: true,
		Service: &apiregistrationv1.ServiceReference{
			Name:      M.Name,
			Namespace: M.Namespace,
			Port:      &M.P.Service.Port,
		},
		Version:         "v1beta1",
		VersionPriority: int32(100),
	},
}

var ServiceMonitor = &promoperator.ServiceMonitor{
	ObjectMeta: M.ObjectMeta(),
	Spec: promoperator.ServiceMonitorSpec{
		Endpoints: []promoperator.Endpoint{
			{
				Interval:      promoperator.Duration("1m"),
				Path:          M.MetricsURL,
				Port:          M.P.Service.Name,
				Scheme:        "https",
				ScrapeTimeout: promoperator.Duration("10s"),
				TLSConfig: &promoperator.TLSConfig{
					SafeTLSConfig: promoperator.SafeTLSConfig{
						InsecureSkipVerify: true,
					},
				},
			},
		},
		JobLabel:          ku.AppLabelName,
		NamespaceSelector: promoperator.NamespaceSelector{MatchNames: []string{M.Namespace}},
		Selector: metav1.LabelSelector{
			MatchLabels: M.MatchLabels(),
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "monitoring.coreos.com/v1", Kind: "ServiceMonitor",
	},
}

var SVC = &corev1.Service{
	TypeMeta: ku.TypeServiceV1,
	ObjectMeta: metav1.ObjectMeta{
		Labels: ku.MergeLabels(
			M.Labels(),
			map[string]string{"kubernetes.io/cluster-service": "true"},
		),
		Name:      M.Name,
		Namespace: M.Namespace,
	},
	Spec: corev1.ServiceSpec{
		Ports:    []corev1.ServicePort{M.P.Service},
		Selector: M.MatchLabels(),
		Type:     corev1.ServiceTypeClusterIP,
	},
}
