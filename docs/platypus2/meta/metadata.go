// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package meta

import (
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type NetPort struct {
	Container corev1.ContainerPort
	Service   corev1.ServicePort
}

type Metadata struct {
	Name      string
	Namespace string
	Instance  string
	Component string
	PartOf    string
	Version   string
	ManagedBy string
	Registry  string
	Image     string
	Tag       string
}

func (b Metadata) Labels() map[string]string {
	return map[string]string{
		"app":                b.Name,
		ku.AppLabelName:      b.Name,
		ku.AppLabelInstance:  b.Instance,
		ku.AppLabelComponent: b.Component,
		ku.AppLabelPartOf:    b.PartOf,
		ku.AppLabelVersion:   b.Version,
		ku.AppLabelManagedBy: b.ManagedBy,
	}
}

func (b Metadata) MatchLabels() map[string]string {
	return map[string]string{
		ku.AppLabelName:     b.Name,
		ku.AppLabelInstance: b.Instance,
	}
}

func (b Metadata) ObjectMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      b.Name,
		Namespace: b.Namespace,
		Labels:    b.Labels(),
	}
}

func (b Metadata) ObjectMetaNoNS() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:   b.Name,
		Labels: b.Labels(),
	}
}

func (b Metadata) ObjectMetaAnnotations(annotations map[string]string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:        b.Name,
		Namespace:   b.Namespace,
		Labels:      b.Labels(),
		Annotations: annotations,
	}
}

func (b Metadata) ObjectMetaNameSuffix(s string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      b.Name + "-" + s,
		Namespace: b.Namespace,
		Labels:    b.Labels(),
	}
}

func (b Metadata) ObjectMetaNameSuffixNoNS(s string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:   b.Name + "-" + s,
		Labels: b.Labels(),
	}
}

func (b Metadata) NS() *corev1.Namespace {
	return &corev1.Namespace{
		TypeMeta:   ku.TypeNamespaceV1,
		ObjectMeta: b.ObjectMetaNoNS(),
		Spec:       corev1.NamespaceSpec{},
	}
}

func (b Metadata) ServiceAccount() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta:   ku.TypeServiceAccountV1,
		ObjectMeta: b.ObjectMeta(),
	}
}

func (b Metadata) ContainerURL() string {
	if b.Image == "" {
		panic("missing container image for: " + b.Name)
	}
	s := b.Image
	if b.Registry != "" {
		s = b.Registry + "/" + s
	}
	if b.Tag != "" {
		s = s + ":" + b.Tag
	}

	return s
}

func (b *Metadata) Service(
	port, targetPort int,
	portName string,
) *corev1.Service {
	return &corev1.Service{
		TypeMeta:   ku.TypeServiceV1,
		ObjectMeta: b.ObjectMeta(),
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       portName,
					Port:       int32(port),
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(targetPort),
				},
			},
			Selector: b.MatchLabels(),
			Type:     corev1.ServiceTypeClusterIP,
		},
	}
}
