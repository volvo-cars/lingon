package kube

import (
	"fmt"

	"github.com/volvo-cars/lingon/pkg/kubeutil"
)

func ExampleNameVarFunc_addKind() {
	m := kubeutil.Metadata{
		Kind: "Deployment",
		Meta: kubeutil.Meta{Name: "super-duper-app"},
	}
	fmt.Println(NameVarFunc(m))
	// Output: SuperDuperAppDeploy
}

func ExampleNameVarFunc_kindSuffix() {
	m := kubeutil.Metadata{
		Kind: "Deployment",
		Meta: kubeutil.Meta{Name: "super-duper-deployment"},
	}
	fmt.Println(NameVarFunc(m))
	// Output: SuperDuperDeploy
}

func ExampleNameVarFunc_kindWithDash() {
	m := kubeutil.Metadata{
		Kind: "ClusterRole",
		Meta: kubeutil.Meta{Name: "argo-cluster-role"},
	}
	fmt.Println(NameVarFunc(m))
	// Output: ArgoCR
}

func ExampleNameFieldFunc_addKind() {
	m := kubeutil.Metadata{
		Kind: "Deployment",
		Meta: kubeutil.Meta{Name: "super-duper-app"},
	}
	fmt.Println(NameFieldFunc(m))
	// Output: SuperDuperAppDeploy
}

func ExampleNameFieldFunc_kindSuffix() {
	m := kubeutil.Metadata{
		Kind: "Deployment",
		Meta: kubeutil.Meta{Name: "super-duper-deployment"},
	}
	fmt.Println(NameFieldFunc(m))
	// Output: SuperDuperDeploy
}

func ExampleNameFieldFunc_kindWithDash() {
	m := kubeutil.Metadata{
		Kind: "ClusterRole",
		Meta: kubeutil.Meta{Name: "argo-cluster-role"},
	}
	fmt.Println(NameFieldFunc(m))
	// Output: ArgoCR
}

func ExampleNameFileFunc() {
	m := kubeutil.Metadata{
		Kind: "Deployment",
		Meta: kubeutil.Meta{Name: "super-duper-app"},
	}
	fmt.Println(NameFileFunc(m))
	// Output: super-duper-app_deploy.go
}

func ExampleDirectoryName_defaultNamespace() {
	fmt.Println(DirectoryName("output", "myapps", "Deployment"))
	// Output: output/myapps
}

func ExampleDirectoryName_clusterrole() {
	fmt.Println(DirectoryName("output", "", "ClusterRole"))
	// Output: output/_cluster/rbac
}

func ExampleDirectoryName_crd() {
	fmt.Println(DirectoryName("output", "", "CustomResourceDefinition"))
	// Output: output/_cluster/crd
}
