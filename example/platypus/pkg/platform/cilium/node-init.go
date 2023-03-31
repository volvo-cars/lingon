package cilium

import (
	"github.com/hexops/valast"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var NodeInit = &appsv1.DaemonSet{
	TypeMeta: kubeutil.TypeDaemonSetV1,
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cilium-node-init",
		Namespace: "kube-system",
		Labels:    map[string]string{"app": "cilium-node-init"},
	},
	Spec: appsv1.DaemonSetSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app": "cilium-node-init",
			},
		},
		Template: v1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels:      map[string]string{"app": "cilium-node-init"},
				Annotations: map[string]string{"container.apparmor.security.beta.kubernetes.io/node-init": "unconfined"},
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "node-init",
						Image: "quay.io/cilium/startup-script:d69851597ea019af980891a4628fb36b7880ec26",
						Env: []v1.EnvVar{
							{
								Name:  "STARTUP_SCRIPT",
								Value: startUpScript,
							},
						},
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								v1.ResourceName("cpu"):    resource.MustParse("100m"),
								v1.ResourceName("memory"): resource.MustParse("100Mi"),
							},
						},
						Lifecycle: &v1.Lifecycle{
							PostStart: &v1.LifecycleHandler{
								Exec: &v1.ExecAction{
									Command: []string{
										"nsenter",
										"--target=1",
										"--mount",
										"--",
										"/bin/bash",
										"-c",
										editIPTableScript,
									},
								},
							},
						},
						TerminationMessagePolicy: v1.TerminationMessagePolicy("FallbackToLogsOnError"),
						ImagePullPolicy:          v1.PullPolicy("IfNotPresent"),
						SecurityContext: &v1.SecurityContext{
							Capabilities: &v1.Capabilities{
								Add: []v1.Capability{
									v1.Capability("SYS_MODULE"),
									v1.Capability("NET_ADMIN"),
									v1.Capability("SYS_ADMIN"),
									v1.Capability("SYS_CHROOT"),
									v1.Capability("SYS_PTRACE"),
								},
							},
							Privileged: valast.Addr(false).(*bool),
							SELinuxOptions: &v1.SELinuxOptions{
								Type:  "spc_t",
								Level: "s0",
							},
						},
					},
				},
				NodeSelector:      map[string]string{"kubernetes.io/os": "linux"},
				HostNetwork:       true,
				HostPID:           true,
				Tolerations:       []v1.Toleration{{Operator: v1.TolerationOperator("Exists")}},
				PriorityClassName: "system-node-critical",
			},
		},
		UpdateStrategy: appsv1.DaemonSetUpdateStrategy{Type: appsv1.DaemonSetUpdateStrategyType("RollingUpdate")},
	},
}

var startUpScript = `#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

echo "Link information:"
ip link

echo "Routing table:"
ip route

echo "Addressing:"
ip -4 a
ip -6 a
mkdir -p "/tmp/cilium-bootstrap.d"
date > "/tmp/cilium-bootstrap.d/cilium-bootstrap-time"
echo "Node initialization complete"
`

var editIPTableScript = `#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

# When running in AWS ENI mode, it's likely that 'aws-node' has
# had a chance to install SNAT iptables rules. These can result
# in dropped traffic, so we should attempt to remove them.
# We do it using a 'postStart' hook since this may need to run
# for nodes which might have already been init'ed but may still
# have dangling rules. This is safe because there are no
# dependencies on anything that is part of the startup script
# itself, and can be safely run multiple times per node (e.g. in
# case of a restart).
if [[ "$(iptables-save | grep -c AWS-SNAT-CHAIN)" != "0" ]];
then
    echo 'Deleting iptables rules created by the AWS CNI VPC plugin'
    iptables-save | grep -v AWS-SNAT-CHAIN | iptables-restore
fi
echo 'Done!'
`
