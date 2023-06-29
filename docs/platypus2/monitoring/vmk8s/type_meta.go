// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package vmk8s

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

var TypeVMRuleV1Beta1 = metav1.TypeMeta{
	APIVersion: "operator.victoriametrics.com/v1beta1",
	Kind:       "VMRule",
}

var TypeVMServiceScrapeV1Beta1 = metav1.TypeMeta{
	APIVersion: "operator.victoriametrics.com/v1beta1",
	Kind:       "VMServiceScrape",
}

var TypeVMSingleV1Beta1 = metav1.TypeMeta{
	APIVersion: "operator.victoriametrics.com/v1beta1",
	Kind:       "VMSingle",
}

var TypeVMNodeScrape = metav1.TypeMeta{
	APIVersion: "operator.victoriametrics.com/v1beta1",
	Kind:       "VMNodeScrape",
}
