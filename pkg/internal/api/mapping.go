// Copyright (c) Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

func ImportKubernetesPkgAlias(file *jen.File) {
	for pkgPath, pkgAlias := range ImportsToPkgAlias {
		file.ImportAlias(pkgPath, pkgAlias)
	}
}

func PkgPathFromAPIVersion(apiVersion string) (string, error) {
	alias, ok := VersionToPkgAlias[apiVersion]
	if !ok {
		return "", fmt.Errorf("unknown api version: %s", apiVersion)
	}

	pkgPath, ok := PkgAliasToImport[alias]
	if !ok {
		return "", fmt.Errorf(
			"unknown pkg alias from APIVersion (%s): %s",
			apiVersion,
			alias,
		)
	}
	return pkgPath, nil
}

var ImportsToPkgAlias = map[string]string{
	"k8s.io/api/admission/v1":                                  "admissionv1",
	"k8s.io/api/admission/v1beta1":                             "admissionv1beta1",
	"k8s.io/api/admissionregistration/v1":                      "admissionregistrationv1",
	"k8s.io/api/admissionregistration/v1beta1":                 "admissionregistrationv1beta1",
	"k8s.io/api/apps/v1":                                       "appsv1",
	"k8s.io/api/authentication/v1":                             "authenticationv1",
	"k8s.io/api/authentication/v1beta1":                        "authenticationv1beta1",
	"k8s.io/api/authorization/v1":                              "authorizationv1",
	"k8s.io/api/authorization/v1beta":                          "authorizationv1beta",
	"k8s.io/api/autoscaling/v1":                                "autoscalingv1",
	"k8s.io/api/autoscaling/v2":                                "autoscalingv2",
	"k8s.io/api/autoscaling/v2beta1":                           "autoscalingv2beta1",
	"k8s.io/api/autoscaling/v2beta2":                           "autoscalingv2beta2",
	"k8s.io/api/batch/v1":                                      "batchv1",
	"k8s.io/api/batch/v1beta1":                                 "batchv1beta1",
	"k8s.io/api/certificates/v1beta1":                          "certificatesv1beta1",
	"k8s.io/api/coordination/v1":                               "coordinationv1",
	"k8s.io/api/core/v1":                                       "corev1",
	"k8s.io/api/discovery/v1":                                  "discoveryv1",
	"k8s.io/api/discovery/v1beta1":                             "discoveryv1beta1",
	"k8s.io/api/events/v1":                                     "eventsv1",
	"k8s.io/api/events/v1beta1":                                "eventsv1beta1",
	"k8s.io/api/extensions/v1beta1":                            "extensionsv1beta1",
	"k8s.io/api/networking/v1":                                 "networkingv1",
	"k8s.io/api/networking/v1beta1":                            "networkingv1beta1",
	"k8s.io/api/node/v1":                                       "nodev1",
	"k8s.io/api/node/v1beta1":                                  "nodev1beta1",
	"k8s.io/api/policy/v1":                                     "policyv1",
	"k8s.io/api/policy/v1beta1":                                "policyv1beta1",
	"k8s.io/api/rbac/v1":                                       "rbacv1",
	"k8s.io/api/scheduling/v1":                                 "schedulingv1",
	"k8s.io/api/scheduling/v1beta1":                            "schedulingv1beta1",
	"k8s.io/api/storage/v1":                                    "storagev1",
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1": "apiextensionsv1",
	"k8s.io/apimachinery/pkg/apis/meta/v1":                     "metav1",
}

var PkgAliasToImport = map[string]string{
	"admissionregistrationv1":      "k8s.io/api/admissionregistration/v1",
	"admissionregistrationv1beta1": "k8s.io/api/admissionregistration/v1beta1",
	"admissionv1":                  "k8s.io/api/admission/v1",
	"admissionv1beta1":             "k8s.io/api/admission/v1beta1",
	"apiextensionsv1":              "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1",
	"appsv1":                       "k8s.io/api/apps/v1",
	"authenticationv1":             "k8s.io/api/authentication/v1",
	"authenticationv1beta1":        "k8s.io/api/authentication/v1beta1",
	"authorizationv1":              "k8s.io/api/authorization/v1",
	"authorizationv1beta":          "k8s.io/api/authorization/v1beta",
	"autoscalingv1":                "k8s.io/api/autoscaling/v1",
	"autoscalingv2":                "k8s.io/api/autoscaling/v2",
	"autoscalingv2beta1":           "k8s.io/api/autoscaling/v2beta1",
	"autoscalingv2beta2":           "k8s.io/api/autoscaling/v2beta2",
	"batchv1":                      "k8s.io/api/batch/v1",
	"batchv1beta1":                 "k8s.io/api/batch/v1beta1",
	"certificatesv1beta1":          "k8s.io/api/certificates/v1beta1",
	"coordinationv1":               "k8s.io/api/coordination/v1",
	"corev1":                       "k8s.io/api/core/v1",
	"discoveryv1":                  "k8s.io/api/discovery/v1",
	"discoveryv1beta1":             "k8s.io/api/discovery/v1beta1",
	"eventsv1":                     "k8s.io/api/events/v1",
	"eventsv1beta1":                "k8s.io/api/events/v1beta1",
	"extensionsv1beta1":            "k8s.io/api/extensions/v1beta1",
	"metav1":                       "k8s.io/apimachinery/pkg/apis/metaObject/v1",
	"networkingv1":                 "k8s.io/api/networking/v1",
	"networkingv1beta1":            "k8s.io/api/networking/v1beta1",
	"nodev1":                       "k8s.io/api/node/v1",
	"nodev1beta1":                  "k8s.io/api/node/v1beta1",
	"policyv1":                     "k8s.io/api/policy/v1",
	"policyv1beta1":                "k8s.io/api/policy/v1beta1",
	"rbacv1":                       "k8s.io/api/rbac/v1",
	"schedulingv1":                 "k8s.io/api/scheduling/v1",
	"schedulingv1beta1":            "k8s.io/api/scheduling/v1beta1",
	"storagev1":                    "k8s.io/api/storage/v1",
}

var VersionToPkgAlias = map[string]string{
	"admissionregistration.k8s.io/v1":       "admissionregistrationv1",
	"admissionregistration.k8s.io/v1alpha1": "admissionregistrationv1alpha1",
	"admissionregistration.k8s.io/v1beta1":  "admissionregistrationv1beta1",
	"apiextensions.k8s.io/v1":               "apiextensionsv1",
	"apiextensions.k8s.io/v1beta1":          "apiextensionsv1beta1",
	"apiregistration.k8s.io/v1":             "apiregistrationv1",
	"apiregistration.k8s.io/v1beta1":        "apiregistrationv1beta1",
	"apps/v1":                               "appsv1",
	"apps/v1beta1":                          "appsv1beta1",
	"apps/v1beta2":                          "appsv1beta2",
	"authentication.k8s.io/v1":              "authenticationv1",
	"authentication.k8s.io/v1alpha1":        "authenticationv1alpha1",
	"authentication.k8s.io/v1beta1":         "authenticationv1beta1",
	"authorization.k8s.io/v1":               "authorizationv1",
	"authorization.k8s.io/v1beta1":          "authorizationv1beta1",
	"autoscaling/v1":                        "autoscalingv1",
	"autoscaling/v2":                        "autoscalingv2",
	"autoscaling/v2beta1":                   "autoscalingv2beta1",
	"autoscaling/v2beta2":                   "autoscalingv2beta2",
	"batch/v1":                              "batchv1",
	"batch/v1beta1":                         "batchv1beta1",
	"certificates.k8s.io/v1":                "certificatesv1",
	"certificates.k8s.io/v1beta1":           "certificatesv1beta1",
	"coordination.k8s.io/v1":                "coordinationv1",
	"coordination.k8s.io/v1beta1":           "coordinationv1beta1",
	"discovery.k8s.io/v1":                   "discoveryv1",
	"discovery.k8s.io/v1beta1":              "discoveryv1beta1",
	"events.k8s.io/v1":                      "eventsv1",
	"events.k8s.io/v1beta1":                 "eventsv1beta1",
	"extensions/v1beta1":                    "extensionsv1beta1",
	"networking.k8s.io/v1":                  "networkingv1",
	"networking.k8s.io/v1beta1":             "networkingv1beta1",
	"node.k8s.io/v1":                        "nodev1",
	"node.k8s.io/v1alpha1":                  "nodev1alpha1",
	"node.k8s.io/v1beta1":                   "nodev1beta1",
	"policy/v1":                             "policyv1",
	"policy/v1beta1":                        "policyv1beta1",
	"rbac.authorization.k8s.io/v1":          "rbacv1",
	"rbac.authorization.k8s.io/v1alpha1":    "rbacv1alpha",
	"rbac.authorization.k8s.io/v1beta1":     "rbacv1beta1",
	"scheduling.k8s.io/v1":                  "schedulingv1",
	"scheduling.k8s.io/v1beta1":             "schedulingv1beta1",
	"storage.k8s.io/v1":                     "storagev1",
	"storage.k8s.io/v1alpha1":               "storagev1alpha1",
	"storage.k8s.io/v1beta1":                "storagev1beta1",
	"v1":                                    "corev1",
}
