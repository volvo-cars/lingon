package promstack

import (
	"os"
	"path/filepath"
	"testing"

	tu "github.com/volvo-cars/lingon/pkg/testutil"
	corev1 "k8s.io/api/core/v1"
)

func TestExportDash(t *testing.T) {
	tu.AssertNoError(t, os.RemoveAll("out"))
	tu.AssertNoError(t, os.MkdirAll("out", os.ModePerm))
	cms := []*corev1.ConfigMap{
		KubePromtheusStackGrafanaCM,
		KubePromtheusStackGrafanaCM,
		KubePromtheusStackGrafanaConfigDashboardsCM,
		KubePromtheusStackGrafanaTestCM,
		KubePromtheusStackKubeAlertmanagerOverviewCM,
		KubePromtheusStackKubeApiserverCM,
		KubePromtheusStackKubeClusterTotalCM,
		KubePromtheusStackKubeControllerManagerCM,
		KubePromtheusStackKubeEtcdCM,
		KubePromtheusStackKubeGrafanaDatasourceCM,
		KubePromtheusStackKubeGrafanaOverviewCM,
		KubePromtheusStackKubeK8SCorednsCM,
		KubePromtheusStackKubeK8SResourcesClusterCM,
		KubePromtheusStackKubeK8SResourcesMulticlusterCM,
		KubePromtheusStackKubeK8SResourcesNamespaceCM,
		KubePromtheusStackKubeK8SResourcesNodeCM,
		KubePromtheusStackKubeK8SResourcesPodCM,
		KubePromtheusStackKubeK8SResourcesWorkloadCM,
		KubePromtheusStackKubeK8SResourcesWorkloadsNamespaceCM,
		KubePromtheusStackKubeKubeletCM,
		KubePromtheusStackKubeNamespaceByPodCM,
		KubePromtheusStackKubeNamespaceByWorkloadCM,
		KubePromtheusStackKubeNodeClusterRsrcUseCM,
		KubePromtheusStackKubeNodeRsrcUseCM,
		KubePromtheusStackKubeNodesCM,
		KubePromtheusStackKubeNodesDarwinCM,
		KubePromtheusStackKubePersistentvolumesusageCM,
		KubePromtheusStackKubePodTotalCM,
		KubePromtheusStackKubePrometheusCM,
		KubePromtheusStackKubeProxyCM,
		KubePromtheusStackKubeSchedulerCM,
		KubePromtheusStackKubeWorkloadTotalCM,
	}

	for _, cm := range cms {
		for k, v := range cm.Data {
			tu.AssertNoError(
				t,
				os.WriteFile(
					filepath.Join("out", k),
					[]byte(v),
					os.ModePerm,
				),
			)
		}
	}
}
