package vmk8s

import (
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Single = &Metadata{
	Name:      "victoria-metrics",
	Namespace: namespace,
	Instance:  "victoria-metrics-" + namespace,
	Component: "tsmdb",
	PartOf:    appName,
	Version:   version,
	ManagedBy: "lingon",
}

var VMOp = &Metadata{
	Name:      "victoria-metrics-operator",
	Namespace: namespace,
	Instance:  "victoria-metrics-operator-" + namespace,
	Component: "operator",
	PartOf:    appName,
	Version:   OperatorVersion,
	ManagedBy: "lingon",
}

type Metadata struct {
	Name      string
	Namespace string
	Instance  string
	Component string
	PartOf    string
	Version   string
	ManagedBy string
}

func (b *Metadata) Labels() map[string]string {
	return map[string]string{
		"app":                b.Name,
		ku.AppLabelName:      b.Name,
		ku.AppLabelInstance:  b.Instance,
		ku.AppLabelComponent: b.Component,
		ku.AppLabelPartOf:    b.PartOf,
		ku.AppLabelVersion:   b.Version,
		ku.AppLabelManagedBy: b.ManagedBy,
	}
}

func (b *Metadata) MatchLabels() map[string]string {
	return map[string]string{
		ku.AppLabelName:     b.Name,
		ku.AppLabelInstance: b.Instance,
	}
}

func (b *Metadata) ObjectMeta() v1.ObjectMeta {
	return v1.ObjectMeta{
		Name:      b.Name,
		Namespace: b.Namespace,
		Labels:    b.Labels(),
	}
}

func (b *Metadata) ObjectMetaNoNS() v1.ObjectMeta {
	return v1.ObjectMeta{
		Name:   b.Name,
		Labels: b.Labels(),
	}
}

func (b *Metadata) ObjectMetaAnnotations(annotations map[string]string) v1.ObjectMeta {
	return v1.ObjectMeta{
		Name:        b.Name,
		Namespace:   b.Namespace,
		Labels:      b.Labels(),
		Annotations: annotations,
	}
}

func (b *Metadata) ObjectMetaNameSuffix(s string) v1.ObjectMeta {
	return v1.ObjectMeta{
		Name:      b.Name + "-" + s,
		Namespace: b.Namespace,
		Labels:    b.Labels(),
	}
}

func (b *Metadata) ObjectMetaNameSuffixNoNS(s string) v1.ObjectMeta {
	return v1.ObjectMeta{
		Name:   b.Name + "-" + s,
		Labels: b.Labels(),
	}
}
