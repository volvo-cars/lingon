// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeutil

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClusterRole(t *testing.T) {
	type args struct {
		name   string
		labels map[string]string
		rules  []rbacv1.PolicyRule
	}
	tests := []struct {
		name string
		args args
		want *rbacv1.ClusterRole
	}{
		{
			name: "cr",
			args: args{
				name:   "cr",
				labels: map[string]string{"l": "v"},
				rules: []rbacv1.PolicyRule{
					{
						Verbs:           nil,
						APIGroups:       nil,
						Resources:       nil,
						ResourceNames:   nil,
						NonResourceURLs: nil,
					},
				},
			},
			want: &rbacv1.ClusterRole{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ClusterRole",
					APIVersion: "rbac.authorization.k8s.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:   "cr",
					Labels: map[string]string{"l": "v"},
				},
				Rules: []rbacv1.PolicyRule{{}},
				// AggregationRule: &rbacv1.AggregationRule{
				// 	ClusterRoleSelectors: []metav1.LabelSelector{
				// 		{
				// 			MatchLabels:      nil,
				// 			MatchExpressions: nil,
				// 		},
				// 	},
				// },
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := ClusterRole(
					tt.args.name,
					tt.args.labels,
					tt.args.rules,
				)
				tu.AssertEqual(t, tt.want, got)
			},
		)
	}
}

func TestClusterRoleRef(t *testing.T) {
	tests := []struct {
		name  string
		rname string
		want  rbacv1.RoleRef
	}{
		{
			name:  "ref",
			rname: "rr",
			want: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     "rr",
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tu.AssertEqual(t, tt.want, ClusterRoleRef(tt.rname))
			},
		)
	}
}

func TestMergeLabels(t *testing.T) {
	tests := []struct {
		name   string
		labels []map[string]string
		want   map[string]string
	}{
		{
			name: "merge",
			labels: []map[string]string{
				{"key1": "val1"},
				{"key2": "val2", "key3": "val3"},
			},
			want: map[string]string{
				"key1": "val1",
				"key2": "val2",
				"key3": "val3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if diff := cmp.Diff(
					tt.want,
					MergeLabels(tt.labels...),
				); diff != "" {
					t.Error(diff)
				}
			},
		)
	}
}

func TestNamespace(t *testing.T) {
	type args struct {
		name        string
		labels      map[string]string
		annotations map[string]string
	}
	tests := []struct {
		name string
		args args
		want *corev1.Namespace
	}{
		{
			name: "ns",
			args: args{
				name:        "testns",
				labels:      map[string]string{"mylabel": "labelvalue"},
				annotations: map[string]string{"annot": "tation"},
			},
			want: &corev1.Namespace{
				TypeMeta: TypeNamespaceV1,
				ObjectMeta: metav1.ObjectMeta{
					Name:        "testns",
					Labels:      map[string]string{"mylabel": "labelvalue"},
					Annotations: map[string]string{"annot": "tation"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := Namespace(
					tt.args.name,
					tt.args.labels,
					tt.args.annotations,
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Namespace() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestResources(t *testing.T) {
	type args struct {
		cpuWant string
		memWant string
		cpuMax  string
		memMax  string
	}
	tests := []struct {
		name string
		args args
		want corev1.ResourceRequirements
	}{
		// TODO: Add test cases.
		{
			name: "ram cpu",
			args: args{
				cpuWant: "2",
				memWant: "2Gi",
				cpuMax:  "4",
				memMax:  "4Gi",
			},
			want: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("2"),
					corev1.ResourceMemory: resource.MustParse("2Gi"),
				},
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("4"),
					corev1.ResourceMemory: resource.MustParse("4Gi"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := Resources(
					tt.args.cpuWant,
					tt.args.memWant,
					tt.args.cpuMax,
					tt.args.memMax,
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Resources() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestRole(t *testing.T) {
	type args struct {
		name      string
		namespace string
		labels    map[string]string
		rules     []rbacv1.PolicyRule
	}
	tests := []struct {
		name string
		args args
		want *rbacv1.Role
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := Role(
					tt.args.name,
					tt.args.namespace,
					tt.args.labels,
					tt.args.rules,
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Role() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestRoleRef(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want rbacv1.RoleRef
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := RoleRef(tt.args.name); !reflect.DeepEqual(
					got,
					tt.want,
				) {
					t.Errorf("RoleRef() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestRoleSubject(t *testing.T) {
	type args struct {
		name      string
		namespace string
	}
	tests := []struct {
		name string
		args args
		want []rbacv1.Subject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := RoleSubject(
					tt.args.name,
					tt.args.namespace,
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("RoleSubject() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestSecret(t *testing.T) {
	type args struct {
		name      string
		namespace string
		data      map[string][]byte
	}
	tests := []struct {
		name string
		args args
		want *corev1.Secret
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := Secret(
					tt.args.name,
					tt.args.namespace,
					tt.args.data,
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Secret() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestServiceAccount(t *testing.T) {
	type args struct {
		name        string
		namespace   string
		labels      map[string]string
		annotations map[string]string
	}
	tests := []struct {
		name string
		args args
		want *corev1.ServiceAccount
	}{
		{
			name: "sa",
			args: args{
				name:        "sa",
				namespace:   "ns",
				labels:      nil,
				annotations: nil,
			},
			want: SimpleSA("sa", "ns"),
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := ServiceAccount(
					tt.args.name,
					tt.args.namespace,
					tt.args.labels,
					tt.args.annotations,
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("ServiceAccount() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestSetDeploySA(t *testing.T) {
	type args struct {
		deploy *appsv1.Deployment
		saName string
	}
	tests := []struct {
		name string
		args args
		want *appsv1.Deployment
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := SetDeploySA(
					tt.args.deploy,
					tt.args.saName,
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("SetDeploySA() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestSimpleCRB(t *testing.T) {
	type args struct {
		sa *corev1.ServiceAccount
		cr *rbacv1.ClusterRole
	}
	tests := []struct {
		name string
		args args
		want *rbacv1.ClusterRoleBinding
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := SimpleCRB(
					tt.args.sa,
					tt.args.cr,
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("SimpleCRB() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestSimpleDeployment(t *testing.T) {
	type args struct {
		name      string
		namespace string
		labels    map[string]string
		replicas  int32
		image     string
	}
	tests := []struct {
		name string
		args args
		want *appsv1.Deployment
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := SimpleDeployment(
					tt.args.name,
					tt.args.namespace,
					tt.args.labels,
					tt.args.replicas,
					tt.args.image,
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("SimpleDeployment() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestSimpleSA(t *testing.T) {
	type args struct {
		name      string
		namespace string
	}
	tests := []struct {
		name string
		args args
		want *corev1.ServiceAccount
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := SimpleSA(
					tt.args.name,
					tt.args.namespace,
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("SimpleSA() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
