// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package grafana

import (
	"fmt"

	"github.com/volvo-cars/lingon/pkg/kube"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

type AppOpts struct {
	Name    string `validate:"required"`
	Version string `validate:"required"`
	Env     string `validate:"required"`
}

func (a AppOpts) NameEnv() string {
	return fmt.Sprintf(
		"%s-%s",
		a.Name,
		a.Env,
	)
}

type KubeOpts struct {
	Name         string
	Namespace    string
	CommonLabels map[string]string

	PostgresHost   string
	PostgresDBName string
	PostgresUser   string
	// TODO: cannot actually store password here...
	PostgresPassword string
}

var _ kube.Exporter = (*KubeApp)(nil)

type KubeApp struct {
	kube.App
	Ns                    *corev1.Namespace
	Svc                   *corev1.Service
	Sa                    *corev1.ServiceAccount
	Cm                    *corev1.ConfigMap
	ClusterroleCr         *rbacv1.ClusterRole
	ClusterrolebindingCrb *rbacv1.ClusterRoleBinding
	Rb                    *rbacv1.RoleBinding
	Secret                *corev1.Secret
	DashboardsDefaultCm   *corev1.ConfigMap
	Role                  *rbacv1.Role
	Deploy                *appsv1.Deployment
}

const (
	AppName = "grafana"
	Version = "9.3.6"
)

func commonLabels(opts AppOpts) map[string]string {
	return map[string]string{
		kubeutil.AppLabelName:      opts.Name,
		kubeutil.AppLabelInstance:  opts.Name,
		kubeutil.AppLabelVersion:   opts.Version,
		kubeutil.AppLabelManagedBy: "lingon",
	}
}

func New(opts AppOpts, kOpts KubeOpts) *KubeApp {
	kOpts.Name = opts.NameEnv()
	kOpts.Namespace = opts.NameEnv()
	kOpts.CommonLabels = commonLabels(opts)

	SA := kubeutil.ServiceAccount(
		AppName,
		kOpts.Namespace,
		kOpts.CommonLabels,
		nil,
	)
	CR := kubeutil.ClusterRole(kOpts.Name, kOpts.CommonLabels, nil)

	Role := kubeutil.Role(
		AppName,
		kOpts.Namespace,
		kOpts.CommonLabels,
		RoleRules,
	)

	return &KubeApp{
		Ns: kubeutil.Namespace(
			kOpts.Name,
			kOpts.CommonLabels,
			nil,
		),
		Secret:              Secret(kOpts),
		Cm:                  Config(kOpts),
		DashboardsDefaultCm: DashboardsDefaultCm(kOpts),
		ClusterroleCr:       CR,
		ClusterrolebindingCrb: kubeutil.BindClusterRole(
			kOpts.Name+"-crb",
			SA,
			CR,
			kOpts.CommonLabels,
		),
		Role:   Role,
		Rb:     kubeutil.BindRole(AppName+"-rb", SA, Role, kOpts.CommonLabels),
		Deploy: Deployment(kOpts),
		Sa:     SA,
		Svc:    Service(kOpts),
	}
}
