// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package karpenter

import (
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var CertSecret = &corev1.Secret{
	TypeMeta: metav1.TypeMeta{
		Kind:       "Secret",
		APIVersion: "v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "karpenter-cert",
		Namespace: "karpenter",
		Labels:    commonLabels,
	},
	Data: map[string][]uint8{}, // Injected by karpenter-webhook
}

func GlobalSettings(
	opts Opts,
) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: kubeutil.TypeConfigMapV1,
		ObjectMeta: metav1.ObjectMeta{
			Name:      "karpenter-global-settings",
			Namespace: "karpenter",
			Labels:    commonLabels,
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
		Namespace: "karpenter",
		Labels:    commonLabels,
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
