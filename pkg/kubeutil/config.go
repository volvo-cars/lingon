// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package kubeutil

// adapted from: https://github.com/bwplotka/mimic/blob/prometheus-kubernetes-example/abstractions/kubernetes/volumes/volumes.go

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/volvo-cars/lingon/pkg/meta"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DataConfigMap(
	name, namespace string,
	labels, annotations, data map[string]string,
) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: meta.TypeConfigMapV1,
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Data: data,
	}
}

// ConfigAndMount is a helper struct to create a ConfigMap and a VolumeMount
type ConfigAndMount struct {
	metav1.ObjectMeta
	corev1.VolumeMount //nolint:govet
	Data               map[string]string
}

func (m ConfigAndMount) ConfigMap() corev1.ConfigMap {
	return corev1.ConfigMap{
		TypeMeta:   meta.TypeConfigMapV1,
		ObjectMeta: m.ObjectMeta,
		Data:       m.Data,
	}
}

func (m ConfigAndMount) VolumeAndMount() VolumeAndMount {
	return VolumeAndMount{
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: m.ObjectMeta.Name},
			},
		},
		VolumeMount: m.VolumeMount,
	}
}

func (m ConfigAndMount) HashEnv(name string) corev1.EnvVar {
	h := sha256.New()
	if err := json.NewEncoder(h).Encode(m.Data); err != nil {
		panic(
			fmt.Sprintf(
				"failed to JSON encode & hash configMap data for %s, err: %v",
				m.VolumeMount.Name,
				err,
			),
		)
	}

	return corev1.EnvVar{
		Name:  name,
		Value: base64.URLEncoding.EncodeToString(h.Sum(nil)),
	}
}
