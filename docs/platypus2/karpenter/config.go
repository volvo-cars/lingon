// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package karpenter

import (
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var CertSecret = &corev1.Secret{
	TypeMeta:   kubeutil.TypeSecretV1,
	ObjectMeta: KA.ObjectMetaNameSuffix("cert"),
	Data:       map[string][]byte{}, // Injected by karpenter-webhook
}

func GlobalSettings(opts Opts) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: kubeutil.TypeConfigMapV1,
		ObjectMeta: metav1.ObjectMeta{
			Name:      KA.ConfigName,
			Namespace: KA.Namespace,
			Labels:    KA.Labels(),
		},
		Data: map[string]string{
			"aws.clusterEndpoint":            opts.ClusterEndpoint,
			"aws.clusterName":                opts.ClusterName,
			"aws.defaultInstanceProfile":     opts.DefaultInstanceProfile,
			"aws.enableENILimitedPodDensity": "true",
			"aws.enablePodENI":               "false",
			"aws.interruptionQueueName":      opts.InterruptQueue,
			"aws.isolatedVPC":                "false",
			"aws.nodeNameConvention":         "ip-name",
			"aws.vmMemoryOverheadPercent":    "0.075",
			"batchIdleDuration":              "1s",
			"batchMaxDuration":               "10s",
		},
	}
}

var LoggingConfig = &corev1.ConfigMap{
	TypeMeta: kubeutil.TypeConfigMapV1,
	ObjectMeta: metav1.ObjectMeta{
		Name:      "config-logging",
		Namespace: KA.Namespace,
		Labels:    KA.Labels(),
	},
	Data: map[string]string{
		"loglevel.webhook":  "error",
		"zap-logger-config": ZapLoggerConfig,
	},
}

var ZapLoggerConfig = `{
  "level": "debug",
  "development": false,
  "disableStacktrace": true,
  "disableCaller": true,
  "sampling": {
    "initial": 100,
    "thereafter": 100
  },
  "outputPaths": ["stdout"],
  "errorOutputPaths": ["stderr"],
  "encoding": "console",
  "encoderConfig": {
    "timeKey": "time",
    "levelKey": "level",
    "nameKey": "logger",
    "callerKey": "caller",
    "messageKey": "message",
    "stacktraceKey": "stacktrace",
    "levelEncoder": "capital",
    "timeEncoder": "iso8601"
  }
}
`
