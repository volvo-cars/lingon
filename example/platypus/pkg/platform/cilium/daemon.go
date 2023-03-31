package cilium

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var Daemon = &appsv1.DaemonSet{
	TypeMeta: kubeutil.TypeMeta("DaemonSet"),
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cilium",
		Namespace: "kube-system",
		Labels:    map[string]string{"k8s-app": "cilium"},
	},
	Spec: appsv1.DaemonSetSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"k8s-app": "cilium",
			},
		},
		Template: v1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{"k8s-app": "cilium"},
				Annotations: map[string]string{
					"container.apparmor.security.beta.kubernetes.io/apply-sysctl-overwrites": "unconfined",
					"container.apparmor.security.beta.kubernetes.io/cilium-agent":            "unconfined",
					"container.apparmor.security.beta.kubernetes.io/clean-cilium-state":      "unconfined",
					"container.apparmor.security.beta.kubernetes.io/mount-cgroup":            "unconfined",
				},
			},
			Spec: v1.PodSpec{
				Volumes: []v1.Volume{
					{
						Name: "cilium-run",
						VolumeSource: v1.VolumeSource{
							HostPath: &v1.HostPathVolumeSource{
								Path: "/var/run/cilium",
								Type: P(v1.HostPathDirectoryOrCreate),
							},
						},
					},
					{
						Name: "bpf-maps",
						VolumeSource: v1.VolumeSource{
							HostPath: &v1.HostPathVolumeSource{
								Path: "/sys/fs/bpf",
								Type: P(v1.HostPathDirectoryOrCreate),
							},
						},
					},
					{
						Name: "hostproc",
						VolumeSource: v1.VolumeSource{
							HostPath: &v1.HostPathVolumeSource{
								Path: "/proc",
								Type: P(v1.HostPathDirectory),
							},
						},
					},
					{
						Name: "cilium-cgroup",
						VolumeSource: v1.VolumeSource{
							HostPath: &v1.HostPathVolumeSource{
								Path: "/run/cilium/cgroupv2",
								Type: P(v1.HostPathDirectoryOrCreate), // v1.HostPathType("DirectoryOrCreate")),
							},
						},
					},
					{
						Name: "cni-path",
						VolumeSource: v1.VolumeSource{
							HostPath: &v1.HostPathVolumeSource{
								Path: "/opt/cni/bin",
								Type: P(v1.HostPathDirectoryOrCreate),
							},
						},
					},
					{
						Name: "etc-cni-netd",
						VolumeSource: v1.VolumeSource{
							HostPath: &v1.HostPathVolumeSource{
								Path: "/etc/cni/net.d",
								Type: P(v1.HostPathDirectoryOrCreate),
							},
						},
					},
					{
						Name: "lib-modules",
						VolumeSource: v1.VolumeSource{
							HostPath: &v1.HostPathVolumeSource{Path: "/lib/modules"},
						},
					},
					{
						Name: "xtables-lock",
						VolumeSource: v1.VolumeSource{
							HostPath: &v1.HostPathVolumeSource{
								Path: "/run/xtables.lock",
								Type: P(v1.HostPathFileOrCreate),
							},
						},
					},
					{
						Name: "clustermesh-secrets",
						VolumeSource: v1.VolumeSource{
							Secret: &v1.SecretVolumeSource{
								SecretName:  "cilium-clustermesh",
								DefaultMode: P(int32(256)),
								Optional:    P(true),
							},
						},
					},
					{
						Name: "cilium-config-path",
						VolumeSource: v1.VolumeSource{
							ConfigMap: &v1.ConfigMapVolumeSource{
								LocalObjectReference: v1.LocalObjectReference{
									Name: "cilium-config",
								},
							},
						},
					},
					{
						Name: "host-proc-sys-net",
						VolumeSource: v1.VolumeSource{
							HostPath: &v1.HostPathVolumeSource{
								Path: "/proc/sys/net",
								Type: P(v1.HostPathDirectory),
							},
						},
					},
					{
						Name: "host-proc-sys-kernel",
						VolumeSource: v1.VolumeSource{
							HostPath: &v1.HostPathVolumeSource{
								Path: "/proc/sys/kernel",
								Type: P(v1.HostPathDirectory),
							},
						},
					},
					{
						Name: "hubble-tls",
						VolumeSource: v1.VolumeSource{
							Projected: &v1.ProjectedVolumeSource{
								Sources: []v1.VolumeProjection{
									{
										Secret: &v1.SecretProjection{
											LocalObjectReference: v1.LocalObjectReference{
												Name: "hubble-server-certs",
											},
											Items: []v1.KeyToPath{
												{
													Key:  "ca.crt",
													Path: "client-ca.crt",
												},
												{
													Key:  "tls.crt",
													Path: "server.crt",
												},
												{
													Key:  "tls.key",
													Path: "server.key",
												},
											},
											Optional: P(true),
										},
									},
								},
								DefaultMode: P(int32(256)),
							},
						},
					},
				},
				InitContainers: []v1.Container{
					{
						Name:  "mount-cgroup",
						Image: "quay.io/cilium/cilium:v1.12.4@sha256:4b074fcfba9325c18e97569ed1988464309a5ebf64bbc79bec6f3d58cafcb8cf",
						Command: []string{
							"sh",
							"-ec",
							ciliumMountScript,
						},
						Env: []v1.EnvVar{
							{
								Name:  "CGROUP_ROOT",
								Value: "/run/cilium/cgroupv2",
							},
							{
								Name:  "BIN_PATH",
								Value: "/opt/cni/bin",
							},
						},
						VolumeMounts: []v1.VolumeMount{
							{
								Name:      "hostproc",
								MountPath: "/hostproc",
							},
							{
								Name:      "cni-path",
								MountPath: "/hostbin",
							},
						},
						TerminationMessagePolicy: v1.TerminationMessagePolicy("FallbackToLogsOnError"),
						ImagePullPolicy:          v1.PullPolicy("IfNotPresent"),
						SecurityContext: &v1.SecurityContext{
							Capabilities: &v1.Capabilities{
								Add: []v1.Capability{
									v1.Capability("SYS_ADMIN"),
									v1.Capability("SYS_CHROOT"),
									v1.Capability("SYS_PTRACE"),
								},
								Drop: []v1.Capability{v1.Capability("ALL")},
							},
							SELinuxOptions: &v1.SELinuxOptions{
								Type:  "spc_t",
								Level: "s0",
							},
						},
					},
					{
						Name:  "apply-sysctl-overwrites",
						Image: "quay.io/cilium/cilium:v1.12.4@sha256:4b074fcfba9325c18e97569ed1988464309a5ebf64bbc79bec6f3d58cafcb8cf",
						Command: []string{
							"sh",
							"-ec",
							ciliumSysctlScript,
						},
						Env: []v1.EnvVar{
							{
								Name:  "BIN_PATH",
								Value: "/opt/cni/bin",
							},
						},
						VolumeMounts: []v1.VolumeMount{
							{
								Name:      "hostproc",
								MountPath: "/hostproc",
							},
							{
								Name:      "cni-path",
								MountPath: "/hostbin",
							},
						},
						TerminationMessagePolicy: v1.TerminationMessagePolicy("FallbackToLogsOnError"),
						ImagePullPolicy:          v1.PullPolicy("IfNotPresent"),
						SecurityContext: &v1.SecurityContext{
							Capabilities: &v1.Capabilities{
								Add: []v1.Capability{
									v1.Capability("SYS_ADMIN"),
									v1.Capability("SYS_CHROOT"),
									v1.Capability("SYS_PTRACE"),
								},
								Drop: []v1.Capability{v1.Capability("ALL")},
							},
							SELinuxOptions: &v1.SELinuxOptions{
								Type:  "spc_t",
								Level: "s0",
							},
						},
					},
					{
						Name:  "mount-bpf-fs",
						Image: "quay.io/cilium/cilium:v1.12.4@sha256:4b074fcfba9325c18e97569ed1988464309a5ebf64bbc79bec6f3d58cafcb8cf",
						Command: []string{
							"/bin/bash",
							"-c",
							"--",
						},
						Args: []string{`mount | grep "/sys/fs/bpf type bpf" || mount -t bpf bpf /sys/fs/bpf`},
						VolumeMounts: []v1.VolumeMount{
							{
								Name:             "bpf-maps",
								MountPath:        "/sys/fs/bpf",
								MountPropagation: P(v1.MountPropagationMode("Bidirectional")),
							},
						},
						TerminationMessagePolicy: v1.TerminationMessagePolicy("FallbackToLogsOnError"),
						ImagePullPolicy:          v1.PullPolicy("IfNotPresent"),
						SecurityContext:          &v1.SecurityContext{Privileged: P(true)},
					},
					{
						Name:    "clean-cilium-state",
						Image:   "quay.io/cilium/cilium:v1.12.4@sha256:4b074fcfba9325c18e97569ed1988464309a5ebf64bbc79bec6f3d58cafcb8cf",
						Command: []string{"/init-container.sh"},
						Env: []v1.EnvVar{
							{
								Name: "CILIUM_ALL_STATE",
								ValueFrom: &v1.EnvVarSource{
									ConfigMapKeyRef: &v1.ConfigMapKeySelector{
										LocalObjectReference: v1.LocalObjectReference{
											Name: "cilium-config",
										},
										Key:      "clean-cilium-state",
										Optional: P(true),
									},
								},
							},
							{
								Name: "CILIUM_BPF_STATE",
								ValueFrom: &v1.EnvVarSource{
									ConfigMapKeyRef: &v1.ConfigMapKeySelector{
										LocalObjectReference: v1.LocalObjectReference{Name: "cilium-config"},
										Key:                  "clean-cilium-bpf-state",
										Optional:             P(true),
									},
								},
							},
						},
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								v1.ResourceName("cpu"):    resource.MustParse("100m"),
								v1.ResourceName("memory"): resource.MustParse("100Mi"),
							},
						},
						VolumeMounts: []v1.VolumeMount{
							{
								Name:      "bpf-maps",
								MountPath: "/sys/fs/bpf",
							},
							{
								Name:             "cilium-cgroup",
								MountPath:        "/run/cilium/cgroupv2",
								MountPropagation: P(v1.MountPropagationMode("HostToContainer")),
							},
							{
								Name:      "cilium-run",
								MountPath: "/var/run/cilium",
							},
						},
						TerminationMessagePolicy: v1.TerminationMessagePolicy("FallbackToLogsOnError"),
						ImagePullPolicy:          v1.PullPolicy("IfNotPresent"),
						SecurityContext: &v1.SecurityContext{
							Capabilities: &v1.Capabilities{
								Add: []v1.Capability{
									v1.Capability("NET_ADMIN"),
									v1.Capability("SYS_MODULE"),
									v1.Capability("SYS_ADMIN"),
									v1.Capability("SYS_RESOURCE"),
								},
								Drop: []v1.Capability{v1.Capability("ALL")},
							},
							SELinuxOptions: &v1.SELinuxOptions{
								Type:  "spc_t",
								Level: "s0",
							},
						},
					},
				},
				Containers: []v1.Container{
					{
						Name:    "cilium-agent",
						Image:   "quay.io/cilium/cilium:v1.12.4@sha256:4b074fcfba9325c18e97569ed1988464309a5ebf64bbc79bec6f3d58cafcb8cf",
						Command: []string{"cilium-agent"},
						Args:    []string{"--config-dir=/tmp/cilium/config-map"},
						Env: []v1.EnvVar{
							{
								Name: "K8S_NODE_NAME",
								ValueFrom: &v1.EnvVarSource{
									FieldRef: &v1.ObjectFieldSelector{
										APIVersion: "v1",
										FieldPath:  "spec.nodeName",
									},
								},
							},
							{
								Name: "CILIUM_K8S_NAMESPACE",
								ValueFrom: &v1.EnvVarSource{
									FieldRef: &v1.ObjectFieldSelector{
										APIVersion: "v1",
										FieldPath:  "metadata.namespace",
									},
								},
							},
							{
								Name:  "CILIUM_CLUSTERMESH_CONFIG",
								Value: "/var/lib/cilium/clustermesh/",
							},
							{
								Name: "CILIUM_CNI_CHAINING_MODE",
								ValueFrom: &v1.EnvVarSource{
									ConfigMapKeyRef: &v1.ConfigMapKeySelector{
										LocalObjectReference: v1.LocalObjectReference{Name: "cilium-config"},
										Key:                  "cni-chaining-mode",
										Optional:             P(true),
									},
								},
							},
							{
								Name: "CILIUM_CUSTOM_CNI_CONF",
								ValueFrom: &v1.EnvVarSource{
									ConfigMapKeyRef: &v1.ConfigMapKeySelector{
										LocalObjectReference: v1.LocalObjectReference{Name: "cilium-config"},
										Key:                  "custom-cni-conf",
										Optional:             P(true),
									},
								},
							},
						},
						VolumeMounts: []v1.VolumeMount{
							{
								Name:      "host-proc-sys-net",
								MountPath: "/host/proc/sys/net",
							},
							{
								Name:      "host-proc-sys-kernel",
								MountPath: "/host/proc/sys/kernel",
							},
							{
								Name:             "bpf-maps",
								MountPath:        "/sys/fs/bpf",
								MountPropagation: P(v1.MountPropagationMode("HostToContainer")),
							},
							{
								Name:      "cilium-run",
								MountPath: "/var/run/cilium",
							},
							{
								Name:      "cni-path",
								MountPath: "/host/opt/cni/bin",
							},
							{
								Name:      "etc-cni-netd",
								MountPath: "/host/etc/cni/net.d",
							},
							{
								Name:      "clustermesh-secrets",
								ReadOnly:  true,
								MountPath: "/var/lib/cilium/clustermesh",
							},
							{
								Name:      "cilium-config-path",
								ReadOnly:  true,
								MountPath: "/tmp/cilium/config-map",
							},
							{
								Name:      "lib-modules",
								ReadOnly:  true,
								MountPath: "/lib/modules",
							},
							{
								Name:      "xtables-lock",
								MountPath: "/run/xtables.lock",
							},
							{
								Name:      "hubble-tls",
								ReadOnly:  true,
								MountPath: "/var/lib/cilium/tls/hubble",
							},
						},
						LivenessProbe: &v1.Probe{
							ProbeHandler: v1.ProbeHandler{
								HTTPGet: &v1.HTTPGetAction{
									Path:   "/healthz",
									Port:   intstr.IntOrString{IntVal: 9879},
									Host:   "127.0.0.1",
									Scheme: v1.URIScheme("HTTP"),
									HTTPHeaders: []v1.HTTPHeader{
										{
											Name:  "brief",
											Value: "true",
										},
									},
								},
							},
							TimeoutSeconds:   5,
							PeriodSeconds:    30,
							SuccessThreshold: 1,
							FailureThreshold: 10,
						},
						ReadinessProbe: &v1.Probe{
							ProbeHandler: v1.ProbeHandler{
								HTTPGet: &v1.HTTPGetAction{
									Path:   "/healthz",
									Port:   intstr.IntOrString{IntVal: 9879},
									Host:   "127.0.0.1",
									Scheme: v1.URIScheme("HTTP"),
									HTTPHeaders: []v1.HTTPHeader{
										{
											Name:  "brief",
											Value: "true",
										},
									},
								},
							},
							TimeoutSeconds:   5,
							PeriodSeconds:    30,
							SuccessThreshold: 1,
							FailureThreshold: 3,
						},
						StartupProbe: &v1.Probe{
							ProbeHandler: v1.ProbeHandler{
								HTTPGet: &v1.HTTPGetAction{
									Path:   "/healthz",
									Port:   intstr.IntOrString{IntVal: 9879},
									Host:   "127.0.0.1",
									Scheme: v1.URIScheme("HTTP"),
									HTTPHeaders: []v1.HTTPHeader{
										{
											Name:  "brief",
											Value: "true",
										},
									},
								},
							},
							PeriodSeconds:    2,
							SuccessThreshold: 1,
							FailureThreshold: 105,
						},
						Lifecycle: &v1.Lifecycle{
							PostStart: &v1.LifecycleHandler{
								Exec: &v1.ExecAction{
									Command: []string{
										"/cni-install.sh",
										"--enable-debug=false",
										"--cni-exclusive=true",
										"--log-file=/var/run/cilium/cilium-cni.log",
									},
								},
							},
							PreStop: &v1.LifecycleHandler{Exec: &v1.ExecAction{Command: []string{"/cni-uninstall.sh"}}},
						},
						TerminationMessagePolicy: v1.TerminationMessagePolicy("FallbackToLogsOnError"),
						ImagePullPolicy:          v1.PullPolicy("IfNotPresent"),
						SecurityContext: &v1.SecurityContext{
							Capabilities: &v1.Capabilities{
								Add: []v1.Capability{
									v1.Capability("CHOWN"),
									v1.Capability("KILL"),
									v1.Capability("NET_ADMIN"),
									v1.Capability("NET_RAW"),
									v1.Capability("IPC_LOCK"),
									v1.Capability("SYS_MODULE"),
									v1.Capability("SYS_ADMIN"),
									v1.Capability("SYS_RESOURCE"),
									v1.Capability("DAC_OVERRIDE"),
									v1.Capability("FOWNER"),
									v1.Capability("SETGID"),
									v1.Capability("SETUID"),
								},
								Drop: []v1.Capability{v1.Capability("ALL")},
							},
							SELinuxOptions: &v1.SELinuxOptions{
								Type:  "spc_t",
								Level: "s0",
							},
						},
					},
				},
				RestartPolicy:                 v1.RestartPolicy("Always"),
				TerminationGracePeriodSeconds: P(int64(1)),
				NodeSelector:                  map[string]string{"kubernetes.io/os": "linux"},
				ServiceAccountName:            "cilium",
				HostNetwork:                   true,
				Affinity: &v1.Affinity{
					PodAntiAffinity: &v1.PodAntiAffinity{
						RequiredDuringSchedulingIgnoredDuringExecution: []v1.PodAffinityTerm{
							{
								LabelSelector: &metav1.LabelSelector{MatchLabels: map[string]string{"k8s-app": "cilium"}},
								TopologyKey:   "kubernetes.io/hostname",
							},
						},
					},
				},
				Tolerations:       []v1.Toleration{{Operator: v1.TolerationOperator("Exists")}},
				PriorityClassName: "system-node-critical",
			},
		},
		UpdateStrategy: appsv1.DaemonSetUpdateStrategy{
			Type:          appsv1.DaemonSetUpdateStrategyType("RollingUpdate"),
			RollingUpdate: &appsv1.RollingUpdateDaemonSet{MaxUnavailable: &intstr.IntOrString{IntVal: 2}},
		},
	},
}

var ciliumMountScript = `cp /usr/bin/cilium-mount /hostbin/cilium-mount;
nsenter --cgroup=/hostproc/1/ns/cgroup --mount=/hostproc/1/ns/mnt "${BIN_PATH}/cilium-mount" $CGROUP_ROOT;
rm /hostbin/cilium-mount
`

var ciliumSysctlScript = `cp /usr/bin/cilium-sysctlfix /hostbin/cilium-sysctlfix;
nsenter --mount=/hostproc/1/ns/mnt "${BIN_PATH}/cilium-sysctlfix";
rm /hostbin/cilium-sysctlfix
`
