// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeutil

// adapted from: https://github.com/bwplotka/mimic/blob/prometheus-kubernetes-example/abstractions/kubernetes/volumes/volumes.go

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DataConfigMap creates a ConfigMap with the given data.
func DataConfigMap(
	name, namespace string,
	labels, annotations, data map[string]string,
) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: TypeConfigMapV1,
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

// ConfigMap creates a ConfigMap from the ConfigAndMount
func (m ConfigAndMount) ConfigMap() *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta:   TypeConfigMapV1,
		ObjectMeta: m.ObjectMeta,
		Data:       m.Data,
	}
}

// VolumeAndMount creates a VolumeAndMount from the ConfigAndMount
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

// HashEnv creates an environment variable with the hash of the ConfigMap data.
func (m ConfigAndMount) HashEnv(name string) corev1.EnvVar {
	h := sha256.New()
	if err := json.NewEncoder(h).Encode(m.Data); err != nil {
		panic(fmt.Sprintf("failed to JSON encode & hash configMap data for %s, err: %v",
			m.VolumeMount.Name,
			err))
	}

	return corev1.EnvVar{
		Name:  name,
		Value: base64.URLEncoding.EncodeToString(h.Sum(nil)),
	}
}

// Hash returns the hash of the ConfigMap data.
func (m ConfigAndMount) Hash() string {
	h := sha256.New()
	if err := json.NewEncoder(h).Encode(m.Data); err != nil {
		panic(fmt.Sprintf("failed to JSON encode & hash configMap data for %s, err: %v",
			m.VolumeMount.Name,
			err))
	}

	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

// VolumeAndMount is a helper struct to create a Volume and a VolumeSource
type VolumeAndMount struct {
	corev1.VolumeMount
	// corev1.Volume has just Name and VolumeSource.
	// A name field is already present in the VolumeMount,
	// so we just add the VolumeSource field here directly.
	VolumeSource corev1.VolumeSource
}

// Volume creates a Volume from the VolumeAndMount
func (vam VolumeAndMount) Volume() corev1.Volume {
	return corev1.Volume{
		Name:         vam.Name,
		VolumeSource: vam.VolumeSource,
	}
}

// VolumesAndMounts is a helper struct to create a list of Volumes and a list of VolumeSource
type VolumesAndMounts []VolumeAndMount

// Volumes creates a list of Volumes from the VolumesAndMounts
func (vams VolumesAndMounts) Volumes() []corev1.Volume {
	volumes := make([]corev1.Volume, 0, len(vams))
	for _, vam := range vams {
		volumes = append(volumes, vam.Volume())
	}
	return volumes
}

// VolumeMounts creates a list of VolumesMount from the VolumesAndMounts
func (vams VolumesAndMounts) VolumeMounts() []corev1.VolumeMount {
	mounts := make([]corev1.VolumeMount, 0, len(vams))
	for _, vam := range vams {
		mounts = append(mounts, vam.VolumeMount)
	}
	return mounts
}
