package meta

// it needs a kubernetes cluster to extract all the metadata.
//go:generate go run -mod=readonly github.com/volvo-cars/lingon/cmd/tools/apisources -out apiresources_gen.go

import (
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

var ErrFieldMissing = errors.New("missing")

func ObjectMeta(
	name, namespace string,
	labels map[string]string,
	annotations map[string]string,
) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:        name,
		Namespace:   namespace,
		Labels:      labels,
		Annotations: annotations,
	}
}

// Metadata is a struct that holds the metadata of a kubernetes object.
type Metadata struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Meta       Meta   `json:"metadata"`
}

type Meta struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// ExtractMetadata returns the Metadata of a kubernetes manifest object.
func ExtractMetadata(data []byte) (*Metadata, error) {
	var m Metadata
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	if m.Meta.Name == "" || m.Kind == "" {
		return nil, fmt.Errorf("name or kind: %w", ErrFieldMissing)
	}
	if m.Meta.Namespace == "" {
		m.Meta.Namespace = "default"
	}
	if len(m.Meta.Labels) == 0 {
		m.Meta.Labels = map[string]string{}
	}
	if len(m.Meta.Annotations) == 0 {
		m.Meta.Annotations = map[string]string{}
	}
	return &m, nil
}
