// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package externalsecrets

import (
	"fmt"

	"github.com/volvo-cars/lingon/pkg/kubeutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var metricsPort = corev1.ContainerPort{
	ContainerPort: int32(portMetric),
	Name:          "metrics",
	Protocol:      corev1.ProtocolTCP,
}

var Deploy = &appsv1.Deployment{
	TypeMeta:   kubeutil.TypeDeploymentV1,
	ObjectMeta: kubeutil.ObjectMeta(AppName, Namespace, ESLabels, nil),
	Spec: appsv1.DeploymentSpec{
		Replicas: P(int32(1)),
		Selector: &metav1.LabelSelector{
			MatchLabels: ESMatchLabels,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: ESMatchLabels,
			},
			Spec: corev1.PodSpec{
				AutomountServiceAccountToken: P(true),
				Containers: []corev1.Container{
					{
						Args:            []string{"--concurrent=1"},
						Image:           containerImage,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Name:            AppName,
						Ports:           []corev1.ContainerPort{metricsPort},
					},
				},
				ServiceAccountName: "external-secrets", // Patched by kubeutil.SetDeploySA
			},
		},
	},
}

var CertControllerDeploy = &appsv1.Deployment{
	TypeMeta: kubeutil.TypeDeploymentV1,
	ObjectMeta: kubeutil.ObjectMeta(
		certControllerName,
		Namespace,
		CertControllerLabels,
		nil,
	),
	Spec: appsv1.DeploymentSpec{
		Replicas: P(int32(1)),
		Selector: &metav1.LabelSelector{
			MatchLabels: CertControllerMatchLabels,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: CertControllerMatchLabels,
			},
			Spec: corev1.PodSpec{
				AutomountServiceAccountToken: P(true),
				Containers: []corev1.Container{
					{
						Args: []string{
							"certcontroller",
							"--crd-requeue-interval=5m",
							// "--service-name=external-secrets-webhook",
							F("--service-name=%s", webhookName),
							// "--service-namespace=external-secrets",
							F("--service-namespace=%s", Namespace),
							// "--secret-name=external-secrets-webhook",
							F("--secret-name=%s", webhookName),
							// "--secret-namespace=external-secrets",
							F("--secret-namespace=%s", Namespace),
						},
						Image:           containerImage,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Name:            "cert-controller",
						Ports:           []corev1.ContainerPort{metricsPort},
						ReadinessProbe: &corev1.Probe{
							InitialDelaySeconds: int32(20),
							PeriodSeconds:       int32(5),
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/readyz",
									Port: intstr.IntOrString{IntVal: int32(8081)},
								},
							},
						},
					},
				},
				ServiceAccountName: "external-secrets-cert-controller", // Patched by kubeutil.SetDeploySA
			},
		},
	},
}
var F = fmt.Sprintf

var WebhookDeploy = &appsv1.Deployment{
	TypeMeta:   kubeutil.TypeDeploymentV1,
	ObjectMeta: kubeutil.ObjectMeta(webhookName, Namespace, WebhookLabels, nil),
	Spec: appsv1.DeploymentSpec{
		Replicas: P(int32(1)),
		Selector: &metav1.LabelSelector{
			MatchLabels: WebhookMatchLabels,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: WebhookMatchLabels,
			},
			Spec: corev1.PodSpec{
				AutomountServiceAccountToken: P(true),
				Containers: []corev1.Container{
					{
						Args: []string{
							"webhook",
							// "--port=10250",
							F("--port=%d", webhookPort),
							// "--dns-name=external-secrets-webhook.external-secrets.svc",
							F("--dns-name=%s.%s.svc", webhookName, Namespace),
							// "--cert-dir=/tmp/certs",
							F("--cert-dir=%s", webhookSecretMountPath),
							"--check-interval=5m",
							// "--metrics-addr=:8080",
							F("--metrics-addr=:%d", portMetric),
							// "--healthz-addr=:8081",
							F("--healthz-addr=:%d", healthzPort),
						},
						Image:           containerImage,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Name:            "webhook",
						Ports: []corev1.ContainerPort{
							metricsPort,
							{
								ContainerPort: int32(webhookPort),
								Name:          "webhook",
								Protocol:      corev1.ProtocolTCP,
							},
						},
						ReadinessProbe: &corev1.Probe{
							InitialDelaySeconds: int32(20),
							PeriodSeconds:       int32(5),
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/readyz",
									Port: intstr.IntOrString{IntVal: int32(healthzPort)},
								},
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: webhookSecretMountPath,
								Name:      "certs",
								ReadOnly:  true,
							},
						},
					},
				},
				ServiceAccountName: "external-secrets-webhook", // Patched by kubeutil.SetDeploySA
				Volumes: []corev1.Volume{
					{
						Name: "certs",
						// VolumeSource: corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{SecretName: "external-secrets-webhook"}},
						VolumeSource: corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{SecretName: webhookName}},
					},
				},
			},
		},
	},
}
