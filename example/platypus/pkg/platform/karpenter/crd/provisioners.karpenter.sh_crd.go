// CODE GENERATED.

package crd

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// P returns a pointer to the given value.
func P[T any](t T) *T {
	return &t
}

var ProvisionersKarpenterShCRD = &apiextensionsv1.CustomResourceDefinition{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{"controller-gen.kubebuilder.io/version": "v0.11.3"},
		Name:        "provisioners.karpenter.sh",
	},
	Spec: apiextensionsv1.CustomResourceDefinitionSpec{
		Group: "karpenter.sh",
		Names: apiextensionsv1.CustomResourceDefinitionNames{
			Categories: []string{"karpenter"},
			Kind:       "Provisioner",
			ListKind:   "ProvisionerList",
			Plural:     "provisioners",
			Singular:   "provisioner",
		},
		Scope: apiextensionsv1.ResourceScope("Cluster"),
		Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
			apiextensionsv1.CustomResourceDefinitionVersion{
				Name: "v1alpha5",
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Description: "Provisioner is the Schema for the Provisioners API",
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"apiVersion": apiextensionsv1.JSONSchemaProps{
								Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
								Type:        "string",
							},
							"kind": apiextensionsv1.JSONSchemaProps{
								Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
								Type:        "string",
							},
							"metadata": apiextensionsv1.JSONSchemaProps{Type: "object"},
							"spec": apiextensionsv1.JSONSchemaProps{
								Description: "ProvisionerSpec is the top level provisioner specification. Provisioners launch nodes in response to pods that are unschedulable. A single provisioner is capable of managing a diverse set of nodes. Node properties are determined from a combination of provisioner and pod scheduling constraints.",
								Properties: map[string]apiextensionsv1.JSONSchemaProps{
									"annotations": apiextensionsv1.JSONSchemaProps{
										AdditionalProperties: &apiextensionsv1.JSONSchemaPropsOrBool{
											Allows: true,
											Schema: &apiextensionsv1.JSONSchemaProps{Type: "string"},
										},
										Description: "Annotations are applied to every node.",
										Type:        "object",
									},
									"consolidation": apiextensionsv1.JSONSchemaProps{
										Description: "Consolidation are the consolidation parameters",
										Properties: map[string]apiextensionsv1.JSONSchemaProps{
											"enabled": apiextensionsv1.JSONSchemaProps{
												Description: "Enabled enables consolidation if it has been set",
												Type:        "boolean",
											},
										},
										Type: "object",
									},
									"kubeletConfiguration": apiextensionsv1.JSONSchemaProps{
										Description: "KubeletConfiguration are options passed to the kubelet when provisioning nodes",
										Properties: map[string]apiextensionsv1.JSONSchemaProps{
											"clusterDNS": apiextensionsv1.JSONSchemaProps{
												Description: "clusterDNS is a list of IP addresses for the cluster DNS server. Note that not all providers may use all addresses.",
												Items:       &apiextensionsv1.JSONSchemaPropsOrArray{Schema: &apiextensionsv1.JSONSchemaProps{Type: "string"}},
												Type:        "array",
											},
											"containerRuntime": apiextensionsv1.JSONSchemaProps{
												Description: "ContainerRuntime is the container runtime to be used with your worker nodes.",
												Type:        "string",
											},
											"cpuCFSQuota": apiextensionsv1.JSONSchemaProps{
												Description: "CPUCFSQuota enables CPU CFS quota enforcement for containers that specify CPU limits.",
												Type:        "boolean",
											},
											"evictionHard": apiextensionsv1.JSONSchemaProps{
												AdditionalProperties: &apiextensionsv1.JSONSchemaPropsOrBool{
													Allows: true,
													Schema: &apiextensionsv1.JSONSchemaProps{Type: "string"},
												},
												Description: "EvictionHard is the map of signal names to quantities that define hard eviction thresholds",
												Type:        "object",
											},
											"evictionMaxPodGracePeriod": apiextensionsv1.JSONSchemaProps{
												Description: "EvictionMaxPodGracePeriod is the maximum allowed grace period (in seconds) to use when terminating pods in response to soft eviction thresholds being met.",
												Format:      "int32",
												Type:        "integer",
											},
											"evictionSoft": apiextensionsv1.JSONSchemaProps{
												AdditionalProperties: &apiextensionsv1.JSONSchemaPropsOrBool{
													Allows: true,
													Schema: &apiextensionsv1.JSONSchemaProps{Type: "string"},
												},
												Description: "EvictionSoft is the map of signal names to quantities that define soft eviction thresholds",
												Type:        "object",
											},
											"evictionSoftGracePeriod": apiextensionsv1.JSONSchemaProps{
												AdditionalProperties: &apiextensionsv1.JSONSchemaPropsOrBool{
													Allows: true,
													Schema: &apiextensionsv1.JSONSchemaProps{Type: "string"},
												},
												Description: "EvictionSoftGracePeriod is the map of signal names to quantities that define grace periods for each eviction signal",
												Type:        "object",
											},
											"imageGCHighThresholdPercent": apiextensionsv1.JSONSchemaProps{
												Description: "ImageGCHighThresholdPercent is the percent of disk usage after which image garbage collection is always run. The percent is calculated by dividing this field value by 100, so this field must be between 0 and 100, inclusive. When specified, the value must be greater than ImageGCLowThresholdPercent.",
												Format:      "int32",
												Maximum:     P(100.0),
												Type:        "integer",
											},
											"imageGCLowThresholdPercent": apiextensionsv1.JSONSchemaProps{
												Description: "ImageGCLowThresholdPercent is the percent of disk usage before which image garbage collection is never run. Lowest disk usage to garbage collect to. The percent is calculated by dividing this field value by 100, so the field value must be between 0 and 100, inclusive. When specified, the value must be less than imageGCHighThresholdPercent",
												Format:      "int32",
												Maximum:     P(100.0),
												Type:        "integer",
											},
											"kubeReserved": apiextensionsv1.JSONSchemaProps{
												AdditionalProperties: &apiextensionsv1.JSONSchemaPropsOrBool{
													Allows: true,
													Schema: &apiextensionsv1.JSONSchemaProps{
														AnyOf: []apiextensionsv1.JSONSchemaProps{
															apiextensionsv1.JSONSchemaProps{Type: "integer"},
															apiextensionsv1.JSONSchemaProps{Type: "string"},
														},
														Pattern:      "^(\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))))?$",
														XIntOrString: true,
													},
												},
												Description: "KubeReserved contains resources reserved for Kubernetes system components.",
												Type:        "object",
											},
											"maxPods": apiextensionsv1.JSONSchemaProps{
												Description: "MaxPods is an override for the maximum number of pods that can run on a worker node instance.",
												Format:      "int32",
												Type:        "integer",
											},
											"podsPerCore": apiextensionsv1.JSONSchemaProps{
												Description: "PodsPerCore is an override for the number of pods that can run on a worker node instance based on the number of cpu cores. This value cannot exceed MaxPods, so, if MaxPods is a lower value, that value will be used.",
												Format:      "int32",
												Type:        "integer",
											},
											"systemReserved": apiextensionsv1.JSONSchemaProps{
												AdditionalProperties: &apiextensionsv1.JSONSchemaPropsOrBool{
													Allows: true,
													Schema: &apiextensionsv1.JSONSchemaProps{
														AnyOf: []apiextensionsv1.JSONSchemaProps{
															apiextensionsv1.JSONSchemaProps{Type: "integer"},
															apiextensionsv1.JSONSchemaProps{Type: "string"},
														},
														Pattern:      "^(\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))))?$",
														XIntOrString: true,
													},
												},
												Description: "SystemReserved contains resources reserved for OS system daemons and kernel memory.",
												Type:        "object",
											},
										},
										Type: "object",
									},
									"labels": apiextensionsv1.JSONSchemaProps{
										AdditionalProperties: &apiextensionsv1.JSONSchemaPropsOrBool{
											Allows: true,
											Schema: &apiextensionsv1.JSONSchemaProps{Type: "string"},
										},
										Description: "Labels are layered with Requirements and applied to every node.",
										Type:        "object",
									},
									"limits": apiextensionsv1.JSONSchemaProps{
										Description: "Limits define a set of bounds for provisioning capacity.",
										Properties: map[string]apiextensionsv1.JSONSchemaProps{
											"resources": apiextensionsv1.JSONSchemaProps{
												AdditionalProperties: &apiextensionsv1.JSONSchemaPropsOrBool{
													Allows: true,
													Schema: &apiextensionsv1.JSONSchemaProps{
														AnyOf: []apiextensionsv1.JSONSchemaProps{
															apiextensionsv1.JSONSchemaProps{Type: "integer"},
															apiextensionsv1.JSONSchemaProps{Type: "string"},
														},
														Pattern:      "^(\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))))?$",
														XIntOrString: true,
													},
												},
												Description: "Resources contains all the allocatable resources that Karpenter supports for limiting.",
												Type:        "object",
											},
										},
										Type: "object",
									},
									"provider": apiextensionsv1.JSONSchemaProps{
										Description:            "Provider contains fields specific to your cloudprovider.",
										Type:                   "object",
										XPreserveUnknownFields: P(true),
									},
									"providerRef": apiextensionsv1.JSONSchemaProps{
										Description: "ProviderRef is a reference to a dedicated CRD for the chosen provider, that holds additional configuration options",
										Properties: map[string]apiextensionsv1.JSONSchemaProps{
											"apiVersion": apiextensionsv1.JSONSchemaProps{
												Description: "API version of the referent",
												Type:        "string",
											},
											"kind": apiextensionsv1.JSONSchemaProps{
												Description: "Kind of the referent; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds\"",
												Type:        "string",
											},
											"name": apiextensionsv1.JSONSchemaProps{
												Description: "Name of the referent; More info: http://kubernetes.io/docs/user-guide/identifiers#names",
												Type:        "string",
											},
										},
										Required: []string{"name"},
										Type:     "object",
									},
									"requirements": apiextensionsv1.JSONSchemaProps{
										Description: "Requirements are layered with Labels and applied to every node.",
										Items: &apiextensionsv1.JSONSchemaPropsOrArray{
											Schema: &apiextensionsv1.JSONSchemaProps{
												Description: "A node selector requirement is a selector that contains values, a key, and an operator that relates the key and values.",
												Properties: map[string]apiextensionsv1.JSONSchemaProps{
													"key": apiextensionsv1.JSONSchemaProps{
														Description: "The label key that the selector applies to.",
														Type:        "string",
													},
													"operator": apiextensionsv1.JSONSchemaProps{
														Description: "Represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists, DoesNotExist. Gt, and Lt.",
														Type:        "string",
													},
													"values": apiextensionsv1.JSONSchemaProps{
														Description: "An array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. If the operator is Gt or Lt, the values array must have a single element, which will be interpreted as an integer. This array is replaced during a strategic merge patch.",
														Items:       &apiextensionsv1.JSONSchemaPropsOrArray{Schema: &apiextensionsv1.JSONSchemaProps{Type: "string"}},
														Type:        "array",
													},
												},
												Required: []string{"key", "operator"},
												Type:     "object",
											},
										},
										Type: "array",
									},
									"startupTaints": apiextensionsv1.JSONSchemaProps{
										Description: "StartupTaints are taints that are applied to nodes upon startup which are expected to be removed automatically within a short period of time, typically by a DaemonSet that tolerates the taint. These are commonly used by daemonsets to allow initialization and enforce startup ordering.  StartupTaints are ignored for provisioning purposes in that pods are not required to tolerate a StartupTaint in order to have nodes provisioned for them.",
										Items: &apiextensionsv1.JSONSchemaPropsOrArray{
											Schema: &apiextensionsv1.JSONSchemaProps{
												Description: "The node this Taint is attached to has the \"effect\" on any pod that does not tolerate the Taint.",
												Properties: map[string]apiextensionsv1.JSONSchemaProps{
													"effect": apiextensionsv1.JSONSchemaProps{
														Description: "Required. The effect of the taint on pods that do not tolerate the taint. Valid effects are NoSchedule, PreferNoSchedule and NoExecute.",
														Type:        "string",
													},
													"key": apiextensionsv1.JSONSchemaProps{
														Description: "Required. The taint key to be applied to a node.",
														Type:        "string",
													},
													"timeAdded": apiextensionsv1.JSONSchemaProps{
														Description: "TimeAdded represents the time at which the taint was added. It is only written for NoExecute taints.",
														Format:      "date-time",
														Type:        "string",
													},
													"value": apiextensionsv1.JSONSchemaProps{
														Description: "The taint value corresponding to the taint key.",
														Type:        "string",
													},
												},
												Required: []string{"effect", "key"},
												Type:     "object",
											},
										},
										Type: "array",
									},
									"taints": apiextensionsv1.JSONSchemaProps{
										Description: "Taints will be applied to every node launched by the Provisioner. If specified, the provisioner will not provision nodes for pods that do not have matching tolerations. Additional taints will be created that match pod tolerations on a per-node basis.",
										Items: &apiextensionsv1.JSONSchemaPropsOrArray{
											Schema: &apiextensionsv1.JSONSchemaProps{
												Description: "The node this Taint is attached to has the \"effect\" on any pod that does not tolerate the Taint.",
												Properties: map[string]apiextensionsv1.JSONSchemaProps{
													"effect": apiextensionsv1.JSONSchemaProps{
														Description: "Required. The effect of the taint on pods that do not tolerate the taint. Valid effects are NoSchedule, PreferNoSchedule and NoExecute.",
														Type:        "string",
													},
													"key": apiextensionsv1.JSONSchemaProps{
														Description: "Required. The taint key to be applied to a node.",
														Type:        "string",
													},
													"timeAdded": apiextensionsv1.JSONSchemaProps{
														Description: "TimeAdded represents the time at which the taint was added. It is only written for NoExecute taints.",
														Format:      "date-time",
														Type:        "string",
													},
													"value": apiextensionsv1.JSONSchemaProps{
														Description: "The taint value corresponding to the taint key.",
														Type:        "string",
													},
												},
												Required: []string{"effect", "key"},
												Type:     "object",
											},
										},
										Type: "array",
									},
									"ttlSecondsAfterEmpty": apiextensionsv1.JSONSchemaProps{
										Description: "TTLSecondsAfterEmpty is the number of seconds the controller will wait before attempting to delete a node, measured from when the node is detected to be empty. A Node is considered to be empty when it does not have pods scheduled to it, excluding daemonsets. \n Termination due to no utilization is disabled if this field is not set.",
										Format:      "int64",
										Type:        "integer",
									},
									"ttlSecondsUntilExpired": apiextensionsv1.JSONSchemaProps{
										Description: "TTLSecondsUntilExpired is the number of seconds the controller will wait before terminating a node, measured from when the node is created. This is useful to implement features like eventually consistent node upgrade, memory leak protection, and disruption testing. \n Termination due to expiration is disabled if this field is not set.",
										Format:      "int64",
										Type:        "integer",
									},
									"weight": apiextensionsv1.JSONSchemaProps{
										Description: "Weight is the priority given to the provisioner during scheduling. A higher numerical weight indicates that this provisioner will be ordered ahead of other provisioners with lower weights. A provisioner with no weight will be treated as if it is a provisioner with a weight of 0.",
										Format:      "int32",
										Maximum:     P(100.0),
										Minimum:     P(1.0),
										Type:        "integer",
									},
								},
								Type: "object",
							},
							"status": apiextensionsv1.JSONSchemaProps{
								Description: "ProvisionerStatus defines the observed state of Provisioner",
								Properties: map[string]apiextensionsv1.JSONSchemaProps{
									"conditions": apiextensionsv1.JSONSchemaProps{
										Description: "Conditions is the set of conditions required for this provisioner to scale its target, and indicates whether or not those conditions are met.",
										Items: &apiextensionsv1.JSONSchemaPropsOrArray{
											Schema: &apiextensionsv1.JSONSchemaProps{
												Description: "Condition defines a readiness condition for a Knative resource. See: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties",
												Properties: map[string]apiextensionsv1.JSONSchemaProps{
													"lastTransitionTime": apiextensionsv1.JSONSchemaProps{
														Description: "LastTransitionTime is the last time the condition transitioned from one status to another. We use VolatileTime in place of metav1.Time to exclude this from creating equality.Semantic differences (all other things held constant).",
														Type:        "string",
													},
													"message": apiextensionsv1.JSONSchemaProps{
														Description: "A human readable message indicating details about the transition.",
														Type:        "string",
													},
													"reason": apiextensionsv1.JSONSchemaProps{
														Description: "The reason for the condition's last transition.",
														Type:        "string",
													},
													"severity": apiextensionsv1.JSONSchemaProps{
														Description: "Severity with which to treat failures of this type of condition. When this is not specified, it defaults to Error.",
														Type:        "string",
													},
													"status": apiextensionsv1.JSONSchemaProps{
														Description: "Status of the condition, one of True, False, Unknown.",
														Type:        "string",
													},
													"type": apiextensionsv1.JSONSchemaProps{
														Description: "Type of condition.",
														Type:        "string",
													},
												},
												Required: []string{"status", "type"},
												Type:     "object",
											},
										},
										Type: "array",
									},
									"lastScaleTime": apiextensionsv1.JSONSchemaProps{
										Description: "LastScaleTime is the last time the Provisioner scaled the number of nodes",
										Format:      "date-time",
										Type:        "string",
									},
									"resources": apiextensionsv1.JSONSchemaProps{
										AdditionalProperties: &apiextensionsv1.JSONSchemaPropsOrBool{
											Allows: true,
											Schema: &apiextensionsv1.JSONSchemaProps{
												AnyOf: []apiextensionsv1.JSONSchemaProps{
													apiextensionsv1.JSONSchemaProps{Type: "integer"},
													apiextensionsv1.JSONSchemaProps{Type: "string"},
												},
												Pattern:      "^(\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))))?$",
												XIntOrString: true,
											},
										},
										Description: "Resources is the list of resources that have been provisioned.",
										Type:        "object",
									},
								},
								Type: "object",
							},
						},
						Type: "object",
					},
				},
				Served:       true,
				Storage:      true,
				Subresources: &apiextensionsv1.CustomResourceSubresources{},
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "apiextensions.k8s.io/v1",
		Kind:       "CustomResourceDefinition",
	},
}
