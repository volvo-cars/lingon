// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeutil

import (
	corev1 "k8s.io/api/core/v1"
)

// Recommended Kubernetes Labels
//
// https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
//
// Shared labels and annotations share a common prefix: `app.kubernetes.io`.
// TFLabels without a prefix are private to users.
// The shared prefix ensures that shared labels do not interfere with custom user labels.
const (
	AppLabelName      = "app.kubernetes.io/name"
	AppLabelInstance  = "app.kubernetes.io/instance"
	AppLabelVersion   = "app.kubernetes.io/version"
	AppLabelComponent = "app.kubernetes.io/component"
	AppLabelPartOf    = "app.kubernetes.io/part-of"
	AppLabelManagedBy = "app.kubernetes.io/managed-by"
)

const (
	// LabelServiceName is used to indicate the name of a Kubernetes service.
	LabelServiceName = "kubernetes.io/service-name"

	LabelInstanceTypeStable = "node.kubernetes.io/instance-type"
	LabelOSStable           = "kubernetes.io/os"
	LabelArchStable         = "kubernetes.io/arch"
	LabelHostname           = "kubernetes.io/hostname"
	LabelTopologyZone       = "topology.kubernetes.io/zone"
	LabelTopologyRegion     = "topology.kubernetes.io/region"
)

// RBAC aggregated cluster roles
//
// see https://kubernetes.io/docs/reference/access-authn-authz/rbac/#aggregated-clusterroles
const (
	LabelRbacAggregateToAdmin = "rbac.authorization.k8s.io/aggregate-to-admin"
	LabelRbacAggregateToEdit  = "rbac.authorization.k8s.io/aggregate-to-edit"
	LabelRbacAggregateToView  = "rbac.authorization.k8s.io/aggregate-to-view"
)

var NotInWindows = corev1.NodeSelectorTerm{
	MatchExpressions: []corev1.NodeSelectorRequirement{
		{
			Key:      LabelOSStable,
			Operator: corev1.NodeSelectorOpNotIn,
			Values:   []string{"windows"},
		},
	},
}

func MergeLabels(labels ...map[string]string) map[string]string {
	result := map[string]string{}
	for _, l := range labels {
		for k, v := range l {
			result[k] = v
		}
	}
	return result
}
