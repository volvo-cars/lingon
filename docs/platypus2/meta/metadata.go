// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package meta

import (
	"errors"
	"fmt"

	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const MAX_CHAR = 64

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

	Img ContainerImg
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

var d = func(i int) string { return fmt.Sprintf("%d", i) }

func (b Metadata) ObjectMeta() metav1.ObjectMeta {
	if len(b.Name) > MAX_CHAR {
		panic("name is longer than " + d(MAX_CHAR) + " char: " + b.Name)
	}
	return metav1.ObjectMeta{
		Name:      b.Name,
		Namespace: b.Namespace,
		Labels:    b.Labels(),
	}
}

func (b Metadata) ObjectMetaNoNS() metav1.ObjectMeta {
	if len(b.Name) > MAX_CHAR {
		panic("name is longer than " + d(MAX_CHAR) + " char: " + b.Name)
	}
	return metav1.ObjectMeta{
		Name:   b.Name,
		Labels: b.Labels(),
	}
}

func (b Metadata) ObjectMetaAnnotations(annotations map[string]string) metav1.ObjectMeta {
	if len(b.Name) > MAX_CHAR {
		panic("name is longer than " + d(MAX_CHAR) + " char: " + b.Name)
	}
	return metav1.ObjectMeta{
		Name:        b.Name,
		Namespace:   b.Namespace,
		Labels:      b.Labels(),
		Annotations: annotations,
	}
}

func SetAnnotations(
	o metav1.ObjectMeta,
	a map[string]string,
) metav1.ObjectMeta {
	o.Annotations = a
	return o
}

func PatchLabels(o metav1.ObjectMeta, key, value string) metav1.ObjectMeta {
	o.Labels = ku.MergeLabels(o.Labels, map[string]string{key: value})
	return o
}

func (b Metadata) ObjectMetaNameSuffix(s string) metav1.ObjectMeta {
	n := b.Name + "-" + s
	if len(n) > MAX_CHAR {
		panic("name is longer than " + d(MAX_CHAR) + " char: " + n)
	}
	return metav1.ObjectMeta{
		Name:      n,
		Namespace: b.Namespace,
		Labels:    b.Labels(),
	}
}

func (b Metadata) ObjectMetaNameSuffixNoNS(s string) metav1.ObjectMeta {
	n := b.Name + "-" + s
	if len(n) > MAX_CHAR {
		panic("name is longer than " + d(MAX_CHAR) + " char: " + n)
	}
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

type ContainerImg struct {
	Registry string
	Image    string
	Sha      string
	Tag      string
}

func (c ContainerImg) URL() string {
	s, err := containerURL(c.Registry, c.Image, c.Sha, c.Tag)
	if err != nil {
		panic(fmt.Sprintf("%#v: %v", c, err))
	}
	return s
}

func containerURL(reg, img, sha, tag string) (string, error) {
	if img == "" {
		return "", errors.New("missing container image")
	}
	s := img
	if reg != "" {
		s = reg + "/" + s
	}

	// docker.io/nats:2.9.19@sha256:3ab6dc....
	// both tag and sha can be in the image
	// if a sha is defined, the tag is ignored

	if tag != "" {
		s = s + ":" + tag
	}
	if sha != "" {
		s = s + "@sha256:" + sha
	}

	return s, nil
}
