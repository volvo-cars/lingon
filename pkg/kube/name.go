// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/veggiemonk/strcase"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	"github.com/volvo-cars/lingon/pkg/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

//
// func outputFileNameYAML(outDir string, m meta.Metadata) string {
// 	return filepath.Join(
// 		DirectoryName(outDir, m.Meta.Namespace, m.Kind),
// 		fmt.Sprintf(
// 			"%d_%s.yaml",
// 			rankOfKind(m.Kind),
// 			BasicName(m.Meta.Name, m.Kind),
// 		),
// 	)
// }

// NameVarFunc returns the name of the variable containing the imported kubernetes object
// TIP: ALWAYS put the kind somewhere in the name to avoid collisions
func NameVarFunc(m kubeutil.Metadata) string {
	bn := BasicName(m.Meta.Name, m.Kind)
	b, a, found := strings.Cut(bn, "_")
	if found {
		if len(a) <= 4 && strings.ToLower(a) != "role" {
			return strcase.Pascal(b) + strings.ToUpper(a)
		}
		return strcase.Pascal(b + "_" + a)
	}
	return strcase.Pascal(bn)
}

// NameFieldFunc returns the name of the field in the App struct
func NameFieldFunc(m kubeutil.Metadata) string {
	bn := BasicName(m.Meta.Name, m.Kind)
	b, a, found := strings.Cut(bn, "_")
	if found {
		if len(a) <= 4 && strings.ToLower(a) != "role" {
			return strcase.Pascal(b) + strings.ToUpper(a)
		}
		return strcase.Pascal(b + "_" + a)
	}
	return strcase.Pascal(bn)
}

func NameFileObjFunc(m kubeutil.Metadata) string {
	return BasicName(m.Meta.Name, m.Kind) + ".go"
}

func isRuneAlphaNumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func RemoveAppName(name, appName string) string {
	res := strings.ReplaceAll(name, appName, "")

	first := string(name[0])
	if strings.ToUpper(first) == first {
		res = strings.ReplaceAll(res, strcase.Pascal(appName), "")
	}
	if res == "" {
		return name
	}
	// remove the non-alphanumeric characters in the beginning
	for i, c := range []rune(res) {
		if isRuneAlphaNumeric(c) {
			return res[i:]
		}
	}
	return name
}

func BasicName(name, kind string) string {
	sk := ShortKind(kind)

	// when the short kind is already the suffix: i.e. for podsecuritypolicy: webapp_psp
	if strings.HasSuffix(strings.ToLower(name), strings.ToLower(sk)) {
		n := strings.TrimSuffix(name, sk)
		// remove the last dash
		n = strings.TrimSuffix(n, "-")
		return n + "_" + sk
	}

	// when dash in kind name: i.e. service-account, cluster-role
	li := strings.LastIndex(name, "-")
	if li > 0 && len(name)-len(kind) < li {
		nn := name[:li] + name[li+1:] // remove the dash
		if strings.HasSuffix(strings.ToLower(nn), strings.ToLower(kind)) {
			n := nn[:len(nn)-len(kind)-1] // remove the kind
			// fmt.Printf(
			// 	"removing kind %q from name %q to get %q\n",
			// 	kind,
			// 	name,
			// 	n,
			// )
			return n + "_" + sk
		}
		return name + "_" + sk
	}

	// replace the kind by the short suffixed in the name: i.e. podsecuritypolicy: webapp_psp
	if strings.HasSuffix(strings.ToLower(name), strings.ToLower(kind)) {
		n := name[:len(name)-len(kind)-1]
		// fmt.Printf("removing kind %q from name %q to get %q\n", kind, name, n)
		return n + "_" + sk
	}

	// just suffix the short kind: i.e. deployment: webapp_deploy
	return name + "_" + sk
}

func ShortKind(s string) string {
	o, ok := meta.KAPI.ByKind(s)
	if !ok || o.ShortName == "" {
		return s
	}
	return o.ShortName
}

func rank(o runtime.Object) string {
	if o == nil || o.GetObjectKind() == nil {
		return ""
	}
	kind := o.GetObjectKind().GroupVersionKind().Kind
	return strconv.Itoa(rankOfKind(kind))
}

// rankOfKind returns an int denoting the position of the given kind
// in the partial ordering of Kubernetes resources, according to which
// kinds depend on which (derived by hand).
// Code taken from FluxCD.
func rankOfKind(kind string) int {
	switch strings.ToLower(kind) {
	// namespaces need to be created first
	case "namespace":
		return 0
	// not namespaced or don't depend on anything else
	case "customresourcedefinition", "serviceaccount", "clusterrole", "role", "persistentvolume", "service":
		return 1
	// These depend on something above, but not each other
	case "resourcequota", "limitrange", "secret", "configmap", "rolebinding", "clusterrolebinding", "persistentvolumeclaim", "ingress": // nolint: lll
		return 2
	// These depend on something above, but not each other
	case "daemonset", "deployment", "replicationcontroller", "replicaset", "job", "cronjob", "statefulset":
		return 3
	// best effort: no dependency
	default:
		return 4
	}
}

const (
	notNamespaced      = "_not-namespaced"
	clusterResourceDir = "_cluster"
	dirNS              = "namespace"
	dirCRD             = "crd"
	dirRBAC            = "rbac"
	dirWH              = "webhook"
	dirST              = "storage"
)

func DirectoryName(out, ns, kind string) string {
	ko, ok := meta.KAPI.ByKind(kind)
	if !ok || ko.Namespaced {
		if ns == "" {
			ns = notNamespaced
		}
		return filepath.Join(out, ns)
	}

	// cluster scoped
	switch kind {
	case "Namespace":
		return filepath.Join(out, clusterResourceDir, dirNS)
	case "ClusterRole", "ClusterRoleBinding":
		return filepath.Join(out, clusterResourceDir, dirRBAC)
	case "MutatingWebhookConfiguration", "ValidatingWebhookConfiguration":
		return filepath.Join(out, clusterResourceDir, dirWH)
	case "PersistentVolume":
		return filepath.Join(out, clusterResourceDir, dirST)
	case "CustomResourceDefinition":
		return filepath.Join(out, clusterResourceDir, dirCRD)
	default:
		return filepath.Join(out, clusterResourceDir)
	}
}
