package externalsecrets

import (
	"github.com/volvo-cars/lingon/pkg/kubeutil"
)

var CertControllerLabels = map[string]string{
	kubeutil.AppLabelName:      certControllerName,
	kubeutil.AppLabelInstance:  AppName,
	kubeutil.AppLabelVersion:   "v" + Version,
	kubeutil.AppLabelManagedBy: "lingon",
	// "helm.sh/chart":                "external-secrets-0.7.2",
}

var WebhookLabels = map[string]string{
	kubeutil.AppLabelInstance:  AppName,
	kubeutil.AppLabelManagedBy: "lingon",
	kubeutil.AppLabelVersion:   "v" + Version,
	"app.kubernetes.io/name":   "external-secrets-webhook",
	// "helm.sh/chart":                "external-secrets-0.7.2",
}

var ESLabels = map[string]string{
	kubeutil.AppLabelInstance:  AppName,
	kubeutil.AppLabelName:      AppName,
	kubeutil.AppLabelManagedBy: "lingon",
	kubeutil.AppLabelVersion:   "v" + Version,
	// "helm.sh/chart":                "external-secrets-0.7.2",
}

var WebhookMatchLabels = map[string]string{
	kubeutil.AppLabelInstance: AppName,
	kubeutil.AppLabelName:     webhookName,
}

var ESMatchLabels = map[string]string{
	kubeutil.AppLabelInstance: AppName,
	kubeutil.AppLabelName:     AppName,
}

var CertControllerMatchLabels = map[string]string{
	kubeutil.AppLabelInstance: AppName,
	kubeutil.AppLabelName:     certControllerName,
}
