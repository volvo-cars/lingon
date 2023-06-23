// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package vmop

import (
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

var VictoriaMetricsOperatorPDB = &policyv1beta1.PodDisruptionBudget{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmop",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-operator",
			"helm.sh/chart":                "victoria-metrics-operator-0.23.1",
		},
		Name:      "vmop-victoria-metrics-operator",
		Namespace: "monitoring",
	},
	Spec: policyv1beta1.PodDisruptionBudgetSpec{
		MinAvailable: &intstr.IntOrString{IntVal: int32(1)},
		Selector: &metav1.LabelSelector{MatchLabels: map[string]string{
			"app.kubernetes.io/instance": "vmop",
			"app.kubernetes.io/name":     "victoria-metrics-operator",
		}},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "policy/v1beta1",
		Kind:       "PodDisruptionBudget",
	},
}
