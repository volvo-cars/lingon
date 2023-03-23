package kubeutil

import (
	"github.com/volvo-cars/lingon/pkg/meta"
	corev1 "k8s.io/api/core/v1"
)

// Labels for namespaces
// taken from https://kubernetes.io/docs/tasks/configure-pod-container/enforce-standards-namespace-labels/
const (
	NSLabelPodSecurityEnforce        = "pod-security.kubernetes.io/enforce"
	NSLabelPodSecurityEnforceVersion = "pod-security.kubernetes.io/enforce-version"
	NSLabelPodSecurityAudit          = "pod-security.kubernetes.io/audit"
	NSLabelPodSecurityAuditVersion   = "pod-security.kubernetes.io/audit-version"
	NSLabelPodSecurityWarn           = "pod-security.kubernetes.io/warn"
	NSLabelPodSecurityWarnVersion    = "pod-security.kubernetes.io/warn-version"
)

const (
	NSValuePodSecurityPrivileged = "privileged"
	NSValuePodSecurityRestricted = "restricted"
	NSValuePodSecurityBaseline   = "baseline"
)

func Namespace(
	name string,
	labels, annotations map[string]string,
) *corev1.Namespace {
	return &corev1.Namespace{
		TypeMeta:   meta.TypeNamespaceV1,
		ObjectMeta: meta.ObjectMeta(name, "", labels, annotations),
	}
}
