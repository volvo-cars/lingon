// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package kubeutil

import (
	"github.com/volvo-cars/lingon/pkg/meta"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Secret(name, namespace string, data map[string][]byte) *v1.Secret {
	return &v1.Secret{
		TypeMeta: meta.TypeMeta("Secret"),
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
}

// see https://kubernetes.io/docs/concepts/configuration/secret/
const (
	// TypeSecretSAToken ServiceAccount token
	TypeSecretSAToken = "kubernetes.io/service-account-token" // nolint:gosec
	// TypeSecretDockerCfg serialized ~/.dockercfg file
	TypeSecretDockerCfg = "kubernetes.io/dockercfg" // nolint:gosec
	// TypeSecretDockerJSON serialized ~/.docker/config.json file
	TypeSecretDockerJSON = "kubernetes.io/dockerconfigjson" // nolint:gosec
	// TypeSecretBasicAuth credentials for basic authentication
	TypeSecretBasicAuth = "kubernetes.io/basic-auth" // nolint:gosec
	// TypeSecretSSH credentials for SSH authentication
	TypeSecretSSH = "kubernetes.io/ssh-auth" // nolint:gosec
	// TypeSecretTLS data for a TLS client or server
	TypeSecretTLS = "kubernetes.io/tls" // nolint:gosec
	// TypeSecretBootstrapToken bootstrap token data
	TypeSecretBootstrapToken = "bootstrap.kubernetes.io/token" // nolint:gosec
)
