package awsauth

import (
	"fmt"

	"github.com/invopop/yaml"
	"github.com/volvo-cars/lingon/pkg/kube"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigMap is an application to manage the aws-auth ConfigMap.
// The AWS EKS kube-system/aws-auth ConfigMap manages access to the Kubernetes
// cluster and AWS Roles and Users need to be added to grant access
type ConfigMap struct {
	kube.App

	ConfigMap *corev1.ConfigMap
}

func NewConfigMap(data *Data) (*ConfigMap, error) {
	mapRoles, err := yaml.Marshal(data.MapRoles)
	if err != nil {
		return nil, fmt.Errorf("marshalling mapRoles: %w", err)
	}
	mapUsers, err := yaml.Marshal(data.MapUsers)
	if err != nil {
		return nil, fmt.Errorf("marshalling mapUsers: %w", err)
	}
	return &ConfigMap{
		ConfigMap: &corev1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "ConfigMap",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "aws-auth",
				Namespace: "kube-system",
			},
			Data: map[string]string{
				"mapRoles": string(mapRoles),
				"mapUsers": string(mapUsers),
			},
		},
	}, nil
}

// Data represents the data of the aws-auth configmap
type Data struct {
	MapRoles []*RolesAuth `json:"mapRoles"`
	MapUsers []*UsersAuth `json:"mapUsers"`
}

// RolesAuth is the basic structure of a mapRoles authentication object
type RolesAuth struct {
	RoleARN  string   `json:"rolearn"`
	Username string   `json:"username"`
	Groups   []string `json:"groups,omitempty"`
}

// UsersAuth is the basic structure of a mapUsers authentication object
type UsersAuth struct {
	UserARN  string   `json:"userarn"`
	Username string   `json:"username"`
	Groups   []string `json:"groups,omitempty"`
}
