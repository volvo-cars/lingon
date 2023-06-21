// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeutil

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Secret creates a Secret with the given name, namespace, labels, annotations and data.
func Secret(name, namespace string, data map[string][]byte) *v1.Secret {
	return &v1.Secret{
		TypeMeta: TypeSecretV1,
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
}

// HashSecret returns the hash of the secret content
func HashSecret(s *v1.Secret) string {
	h := sha256.New()
	if err := json.NewEncoder(h).Encode(s.Data); err != nil {
		panic(
			fmt.Sprintf(
				"failed to JSON encode & hash secret data for %s, err: %v",
				s.Name, err,
			),
		)
	}
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func SecretEnvVar(varName, key, secretName string) v1.EnvVar {
	return v1.EnvVar{
		Name: varName,
		ValueFrom: &v1.EnvVarSource{
			SecretKeyRef: &v1.SecretKeySelector{
				Key:                  key,
				LocalObjectReference: v1.LocalObjectReference{Name: secretName},
			},
		},
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
