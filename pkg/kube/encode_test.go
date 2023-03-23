package kube

import (
	"errors"
	"testing"

	"github.com/volvo-cars/lingon/pkg/kubeutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestEncode_EmptyField(t *testing.T) {
	i := newEmbedWithEmptyField()

	_, err := encodeApp(i)
	if err != nil {
		if !errors.Is(err, ErrFieldMissing) {
			t.Fatal(err)
		}
	}
}

func TestEncode_EmbeddedStruct(t *testing.T) {
	i := newEmbeddedStruct()
	names := map[string]struct{}{
		"1_IAMSa":   {},
		"1_IAMCr":   {},
		"2_IAMCrb":  {},
		"3_Depl":    {},
		"3_IAMDepl": {},
	}

	m, err := encodeApp(i)
	if err != nil {
		if !errors.Is(err, ErrFieldMissing) {
			t.Fatal(err)
		}
	}

	// all the fields are present
	for k := range m {
		if _, ok := names[k]; !ok {
			t.Errorf("unexpected key %q", k)
		}
	}

	// only the fields defined are present
	for k := range names {
		if _, ok := m[k]; !ok {
			t.Errorf("missing key %q", k)
		}
	}
}

type IAM struct {
	App
	Sa   *corev1.ServiceAccount
	Crb  *rbacv1.ClusterRoleBinding
	Cr   *rbacv1.ClusterRole
	Depl *appsv1.Deployment
}

type EmbedStruct struct {
	IAM
	Depl *appsv1.Deployment
}

type EmbedWithEmptyField struct {
	IAM
	Depl      *appsv1.Deployment
	EmptyDepl *appsv1.Deployment
}

var name = "fyaml"

var labels = map[string]string{
	"app": name,
}

func newEmbeddedStruct() *EmbedStruct {
	sa := kubeutil.SimpleSA(name, "defaultns")
	sa.Labels = labels

	cr := &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{Name: sa.Name, Labels: labels},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			},
		},
	}
	crb := kubeutil.SimpleCRB(sa, cr)
	crb.Labels = labels

	iam := IAM{
		Sa:  sa,
		Crb: crb,
		Cr:  cr,
		Depl: kubeutil.SimpleDeployment(
			"another"+name,
			sa.Namespace,
			labels,
			int32(1),
			"nginx:latest",
		),
	}

	return &EmbedStruct{
		Depl: kubeutil.SimpleDeployment(
			name,
			sa.Namespace,
			labels,
			int32(1),
			"nginx:latest",
		),
		IAM: iam,
	}
}

func newEmbedWithEmptyField() *EmbedWithEmptyField {
	sa := kubeutil.SimpleSA(name, "defaultns")
	sa.Labels = labels

	cr := &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{Name: sa.Name, Labels: labels},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			},
		},
	}
	crb := kubeutil.SimpleCRB(sa, cr)
	crb.Labels = labels

	iam := IAM{
		Sa:  sa,
		Crb: crb,
		Cr:  cr,
		Depl: kubeutil.SimpleDeployment(
			"another"+name,
			sa.Namespace,
			labels,
			int32(1),
			"nginx:latest",
		),
	}

	return &EmbedWithEmptyField{
		Depl: kubeutil.SimpleDeployment(
			name,
			sa.Namespace,
			labels,
			int32(1),
			"nginx:latest",
		),
		IAM: iam,
	}
}
