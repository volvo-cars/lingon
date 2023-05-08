// Copyright 2023 Volvo Car Corporation
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

	// AppLabelName defines the name of the application i.e. mysql
	AppLabelName = "app.kubernetes.io/name"
	// AppLabelInstance defines a unique name identifying the instance of an application i.e. mysql-abcxyz
	AppLabelInstance = "app.kubernetes.io/instance"
	// AppLabelVersion defines the current version of the application i.e. 5.7.21
	AppLabelVersion = "app.kubernetes.io/version"
	// AppLabelComponent defines the component withing the architecture i.e. database
	AppLabelComponent = "app.kubernetes.io/component"
	// AppLabelPartOf defines the name of a higher level application this one is part of i.e. WordPress
	AppLabelPartOf = "app.kubernetes.io/part-of"
	// AppLabelManagedBy defines the tool being used to manage the operation of an application i.e. lingon
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

// MergeLabels merges multiple label maps into one.
func MergeLabels(labels ...map[string]string) map[string]string {
	result := map[string]string{}
	for _, l := range labels {
		for k, v := range l {
			result[k] = v
		}
	}
	return result
}
