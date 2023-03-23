package kubeutil

import (
	"github.com/volvo-cars/lingon/pkg/meta"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SimpleCRB creates a ClusterRoleBinding from a service account to a cluster role
func SimpleCRB(
	sa *corev1.ServiceAccount,
	cr *rbacv1.ClusterRole,
) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta:   meta.TypeClusterRoleBindingV1,
		ObjectMeta: metav1.ObjectMeta{Name: sa.Name},
		Subjects:   RoleSubject(sa.Name, sa.Namespace),
		RoleRef:    ClusterRoleRef(cr.Name),
	}
}

// BindClusterRole binds a cluster role to a service account
func BindClusterRole(
	name string,
	sa *corev1.ServiceAccount,
	cr *rbacv1.ClusterRole,
	labels map[string]string,
) *rbacv1.ClusterRoleBinding {
	if name == "" {
		name = sa.Name + "-" + cr.Name
	}
	return &rbacv1.ClusterRoleBinding{
		TypeMeta:   meta.TypeClusterRoleBindingV1,
		ObjectMeta: metav1.ObjectMeta{Name: name, Labels: labels},
		Subjects:   RoleSubject(sa.Name, sa.Namespace),
		RoleRef:    ClusterRoleRef(cr.Name),
	}
}

// BindRole binds a role to a service account inside the Role's namespace
func BindRole(
	name string,
	sa *corev1.ServiceAccount,
	r *rbacv1.Role,
	labels map[string]string,
) *rbacv1.RoleBinding {
	if name == "" {
		name = sa.Name + "-" + r.Name
	}
	return &rbacv1.RoleBinding{
		TypeMeta: meta.TypeRoleBindingV1,
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: r.Namespace, // A RoleBinding may reference any Role in the same namespace
			Labels:    labels,
		},
		Subjects: RoleSubject(sa.Name, sa.Namespace),
		RoleRef:  RoleRef(r.Name),
	}
}

func ClusterRoleRef(name string) rbacv1.RoleRef {
	return rbacv1.RoleRef{
		APIGroup: meta.TypeClusterRoleV1.GroupVersionKind().Group,
		Kind:     meta.TypeClusterRoleV1.GroupVersionKind().Kind,
		Name:     name,
	}
}

func RoleRef(name string) rbacv1.RoleRef {
	return rbacv1.RoleRef{
		APIGroup: meta.TypeRoleV1.GroupVersionKind().Group,
		Kind:     meta.TypeRoleV1.GroupVersionKind().Kind,
		Name:     name,
	}
}

func RoleSubject(name, namespace string) []rbacv1.Subject {
	return []rbacv1.Subject{
		{
			Kind:      meta.TypeServiceAccountV1.GroupVersionKind().Kind,
			Name:      name,
			Namespace: namespace,
		},
	}
}
