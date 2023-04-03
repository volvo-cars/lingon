// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/rogpeppe/go-internal/txtar"
	"github.com/volvo-cars/lingon/pkg/kube"
	"github.com/volvo-cars/lingon/pkg/kube/testdata/go/tekton"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const defaultExportOutputDir = "out/export"

func TestExport(t *testing.T) {
	type args struct {
		Name     string
		km       kube.Exporter
		Opts     []kube.ExportOption
		OutFiles []string
	}
	outEBDS := filepath.Join(defaultExportOutputDir, "embeddedstruct")
	outTekton := filepath.Join(defaultExportOutputDir, "tekton")
	TT := []args{
		{
			Name: "embedded struct",
			km:   newEmbeddedStruct(),
			Opts: []kube.ExportOption{
				kube.WithExportKustomize(true),
				kube.WithExportOutputDirectory(outEBDS),
			},
			OutFiles: []string{
				"out/export/embeddedstruct/1_iamcr.yaml",
				"out/export/embeddedstruct/1_iamsa.yaml",
				"out/export/embeddedstruct/2_iamcrb.yaml",
				"out/export/embeddedstruct/3_depl.yaml",
				"out/export/embeddedstruct/3_iamdepl.yaml",
				"out/export/embeddedstruct/kustomization.yaml",
			},
		},
		{
			Name: "embedded struct explode",
			km:   newEmbeddedStruct(),
			Opts: []kube.ExportOption{
				kube.WithExportExplodeManifests(true),
				kube.WithExportOutputDirectory(outEBDS),
			},
			OutFiles: []string{
				"out/export/embeddedstruct/_cluster/rbac/1_iamcr.yaml",
				"out/export/embeddedstruct/_cluster/rbac/2_iamcrb.yaml",
				"out/export/embeddedstruct/defaultns/1_iamsa.yaml",
				"out/export/embeddedstruct/defaultns/3_depl.yaml",
				"out/export/embeddedstruct/defaultns/3_iamdepl.yaml",
			},
		},
		{
			Name: "embedded struct with name file func",
			km:   newEmbeddedStruct(),
			Opts: []kube.ExportOption{
				kube.WithExportOutputDirectory(outEBDS),
				kube.WithExportNameFileFunc(
					func(name *kubeutil.Metadata) string {
						return strings.ToLower(name.Kind) + "_" + name.Meta.Name + ".yaml"
					},
				),
			},
			OutFiles: []string{
				"out/export/embeddedstruct/clusterrole_imthename.yaml",
				"out/export/embeddedstruct/clusterrolebinding_imthename.yaml",
				"out/export/embeddedstruct/deployment_anotherimthename.yaml",
				"out/export/embeddedstruct/deployment_imthename.yaml",
				"out/export/embeddedstruct/serviceaccount_imthename.yaml",
			},
		},
		{
			Name: "embedded struct with explode and name file func",
			km:   newEmbeddedStruct(),
			Opts: []kube.ExportOption{
				kube.WithExportOutputDirectory(outEBDS),
				kube.WithExportExplodeManifests(true),
				kube.WithExportNameFileFunc(
					func(name *kubeutil.Metadata) string {
						return strings.ToLower(name.Kind) + "_" + name.Meta.Name + ".yaml"
					},
				),
			},
			OutFiles: []string{
				"out/export/embeddedstruct/_cluster/rbac/clusterrole_imthename.yaml",
				"out/export/embeddedstruct/_cluster/rbac/clusterrolebinding_imthename.yaml",
				"out/export/embeddedstruct/defaultns/deployment_anotherimthename.yaml",
				"out/export/embeddedstruct/defaultns/deployment_imthename.yaml",
				"out/export/embeddedstruct/defaultns/serviceaccount_imthename.yaml",
			},
		},
		{
			Name: "tekton",
			km:   tekton.New(),
			Opts: []kube.ExportOption{
				kube.WithExportOutputDirectory(outTekton),
			},
			OutFiles: []string{
				"out/export/tekton/0_pipelines_ns.yaml",
				"out/export/tekton/0_pipelines_resolvers_ns.yaml",
				"out/export/tekton/1_aggregate_edit_cr.yaml",
				"out/export/tekton/1_aggregate_view_cr.yaml",
				"out/export/tekton/1_clustertasks_dev_crd.yaml",
				"out/export/tekton/1_customruns_dev_crd.yaml",
				"out/export/tekton/1_pipelineresources_dev_crd.yaml",
				"out/export/tekton/1_pipelineruns_dev_crd.yaml",
				"out/export/tekton/1_pipelines_controller_cluster_access_cr.yaml",
				"out/export/tekton/1_pipelines_controller_role.yaml",
				"out/export/tekton/1_pipelines_controller_sa.yaml",
				"out/export/tekton/1_pipelines_controller_svc.yaml",
				"out/export/tekton/1_pipelines_controller_tenant_access_cr.yaml",
				"out/export/tekton/1_pipelines_dev_crd.yaml",
				"out/export/tekton/1_pipelines_info_role.yaml",
				"out/export/tekton/1_pipelines_leader_election_role.yaml",
				"out/export/tekton/1_pipelines_resolvers_namespace_rbac_role.yaml",
				"out/export/tekton/1_pipelines_resolvers_resolution_request_updates_cr.yaml",
				"out/export/tekton/1_pipelines_resolvers_sa.yaml",
				"out/export/tekton/1_pipelines_webhook_cluster_access_cr.yaml",
				"out/export/tekton/1_pipelines_webhook_role.yaml",
				"out/export/tekton/1_pipelines_webhook_sa.yaml",
				"out/export/tekton/1_pipelines_webhook_svc.yaml",
				"out/export/tekton/1_resolutionrequests_resolution_dev_crd.yaml",
				"out/export/tekton/1_runs_dev_crd.yaml",
				"out/export/tekton/1_taskruns_dev_crd.yaml",
				"out/export/tekton/1_tasks_dev_crd.yaml",
				"out/export/tekton/1_verificationpolicies_dev_crd.yaml",
				"out/export/tekton/2_bundleresolver_config_cm.yaml",
				"out/export/tekton/2_cluster_resolver_config_cm.yaml",
				"out/export/tekton/2_config_artifact_bucket_cm.yaml",
				"out/export/tekton/2_config_artifact_pvc_cm.yaml",
				"out/export/tekton/2_config_defaults_cm.yaml",
				"out/export/tekton/2_config_leader_election_cm.yaml",
				"out/export/tekton/2_config_leader_election_cm1.yaml",
				"out/export/tekton/2_config_logging_cm.yaml",
				"out/export/tekton/2_config_logging_cm2.yaml",
				"out/export/tekton/2_config_observability_cm.yaml",
				"out/export/tekton/2_config_observability_cm3.yaml",
				"out/export/tekton/2_config_registry_cert_cm.yaml",
				"out/export/tekton/2_config_spire_cm.yaml",
				"out/export/tekton/2_config_trusted_resources_cm.yaml",
				"out/export/tekton/2_feature_flags_cm.yaml",
				"out/export/tekton/2_git_resolver_config_cm.yaml",
				"out/export/tekton/2_hubresolver_config_cm.yaml",
				"out/export/tekton/2_pipelines_controller_cluster_access_crb.yaml",
				"out/export/tekton/2_pipelines_controller_leaderelection_rb.yaml",
				"out/export/tekton/2_pipelines_controller_rb.yaml",
				"out/export/tekton/2_pipelines_controller_tenant_access_crb.yaml",
				"out/export/tekton/2_pipelines_info_cm.yaml",
				"out/export/tekton/2_pipelines_info_rb.yaml",
				"out/export/tekton/2_pipelines_resolvers_crb.yaml",
				"out/export/tekton/2_pipelines_resolvers_namespace_rbac_rb.yaml",
				"out/export/tekton/2_pipelines_webhook_cluster_access_crb.yaml",
				"out/export/tekton/2_pipelines_webhook_leaderelection_rb.yaml",
				"out/export/tekton/2_pipelines_webhook_rb.yaml",
				"out/export/tekton/2_resolvers_feature_flags_cm.yaml",
				"out/export/tekton/2_webhook_certs_secrets.yaml",
				"out/export/tekton/3_pipelines_controller_deploy.yaml",
				"out/export/tekton/3_pipelines_remote_resolvers_deploy.yaml",
				"out/export/tekton/3_pipelines_webhook_deploy.yaml",
				"out/export/tekton/4_config_webhook_pipeline_dev_validatingwebhookconfigurations.yaml",
				"out/export/tekton/4_pipelines_webhook_hpa.yaml",
				"out/export/tekton/4_validation_webhook_pipeline_dev_validatingwebhookconfigurations.yaml",
				"out/export/tekton/4_webhook_pipeline_dev_mutatingwebhookconfigurations.yaml",
			},
		},
		{
			Name: "tekton",
			km:   tekton.New(),
			Opts: []kube.ExportOption{
				kube.WithExportOutputDirectory(outTekton),
				kube.WithExportKustomize(true),
				kube.WithExportExplodeManifests(true),
				kube.WithExportSecretHook(
					func(s *corev1.Secret) error {
						return nil
					},
				),
			},
			OutFiles: []string{
				// should not have any secrets since we use the secret hook
				"out/export/tekton/_cluster/crd/1_clustertasks_dev_crd.yaml",
				"out/export/tekton/_cluster/crd/1_customruns_dev_crd.yaml",
				"out/export/tekton/_cluster/crd/1_pipelineresources_dev_crd.yaml",
				"out/export/tekton/_cluster/crd/1_pipelineruns_dev_crd.yaml",
				"out/export/tekton/_cluster/crd/1_pipelines_dev_crd.yaml",
				"out/export/tekton/_cluster/crd/1_resolutionrequests_resolution_dev_crd.yaml",
				"out/export/tekton/_cluster/crd/1_runs_dev_crd.yaml",
				"out/export/tekton/_cluster/crd/1_taskruns_dev_crd.yaml",
				"out/export/tekton/_cluster/crd/1_tasks_dev_crd.yaml",
				"out/export/tekton/_cluster/crd/1_verificationpolicies_dev_crd.yaml",
				"out/export/tekton/_cluster/namespace/0_pipelines_ns.yaml",
				"out/export/tekton/_cluster/namespace/0_pipelines_resolvers_ns.yaml",
				"out/export/tekton/_cluster/rbac/1_aggregate_edit_cr.yaml",
				"out/export/tekton/_cluster/rbac/1_aggregate_view_cr.yaml",
				"out/export/tekton/_cluster/rbac/1_pipelines_controller_cluster_access_cr.yaml",
				"out/export/tekton/_cluster/rbac/1_pipelines_controller_tenant_access_cr.yaml",
				"out/export/tekton/_cluster/rbac/1_pipelines_resolvers_resolution_request_updates_cr.yaml",
				"out/export/tekton/_cluster/rbac/1_pipelines_webhook_cluster_access_cr.yaml",
				"out/export/tekton/_cluster/rbac/2_pipelines_controller_cluster_access_crb.yaml",
				"out/export/tekton/_cluster/rbac/2_pipelines_controller_tenant_access_crb.yaml",
				"out/export/tekton/_cluster/rbac/2_pipelines_resolvers_crb.yaml",
				"out/export/tekton/_cluster/rbac/2_pipelines_webhook_cluster_access_crb.yaml",
				"out/export/tekton/_cluster/webhook/4_config_webhook_pipeline_dev_validatingwebhookconfigurations.yaml",
				"out/export/tekton/_cluster/webhook/4_validation_webhook_pipeline_dev_validatingwebhookconfigurations.yaml",
				"out/export/tekton/_cluster/webhook/4_webhook_pipeline_dev_mutatingwebhookconfigurations.yaml",
				"out/export/tekton/kustomization.yaml",
				"out/export/tekton/tekton-pipelines-resolvers/1_pipelines_resolvers_namespace_rbac_role.yaml",
				"out/export/tekton/tekton-pipelines-resolvers/1_pipelines_resolvers_sa.yaml",
				"out/export/tekton/tekton-pipelines-resolvers/2_bundleresolver_config_cm.yaml",
				"out/export/tekton/tekton-pipelines-resolvers/2_cluster_resolver_config_cm.yaml",
				"out/export/tekton/tekton-pipelines-resolvers/2_config_leader_election_cm1.yaml",
				"out/export/tekton/tekton-pipelines-resolvers/2_config_logging_cm2.yaml",
				"out/export/tekton/tekton-pipelines-resolvers/2_config_observability_cm3.yaml",
				"out/export/tekton/tekton-pipelines-resolvers/2_git_resolver_config_cm.yaml",
				"out/export/tekton/tekton-pipelines-resolvers/2_hubresolver_config_cm.yaml",
				"out/export/tekton/tekton-pipelines-resolvers/2_pipelines_resolvers_namespace_rbac_rb.yaml",
				"out/export/tekton/tekton-pipelines-resolvers/2_resolvers_feature_flags_cm.yaml",
				"out/export/tekton/tekton-pipelines-resolvers/3_pipelines_remote_resolvers_deploy.yaml",
				"out/export/tekton/tekton-pipelines/1_pipelines_controller_role.yaml",
				"out/export/tekton/tekton-pipelines/1_pipelines_controller_sa.yaml",
				"out/export/tekton/tekton-pipelines/1_pipelines_controller_svc.yaml",
				"out/export/tekton/tekton-pipelines/1_pipelines_info_role.yaml",
				"out/export/tekton/tekton-pipelines/1_pipelines_leader_election_role.yaml",
				"out/export/tekton/tekton-pipelines/1_pipelines_webhook_role.yaml",
				"out/export/tekton/tekton-pipelines/1_pipelines_webhook_sa.yaml",
				"out/export/tekton/tekton-pipelines/1_pipelines_webhook_svc.yaml",
				"out/export/tekton/tekton-pipelines/2_config_artifact_bucket_cm.yaml",
				"out/export/tekton/tekton-pipelines/2_config_artifact_pvc_cm.yaml",
				"out/export/tekton/tekton-pipelines/2_config_defaults_cm.yaml",
				"out/export/tekton/tekton-pipelines/2_config_leader_election_cm.yaml",
				"out/export/tekton/tekton-pipelines/2_config_logging_cm.yaml",
				"out/export/tekton/tekton-pipelines/2_config_observability_cm.yaml",
				"out/export/tekton/tekton-pipelines/2_config_registry_cert_cm.yaml",
				"out/export/tekton/tekton-pipelines/2_config_spire_cm.yaml",
				"out/export/tekton/tekton-pipelines/2_config_trusted_resources_cm.yaml",
				"out/export/tekton/tekton-pipelines/2_feature_flags_cm.yaml",
				"out/export/tekton/tekton-pipelines/2_pipelines_controller_leaderelection_rb.yaml",
				"out/export/tekton/tekton-pipelines/2_pipelines_controller_rb.yaml",
				"out/export/tekton/tekton-pipelines/2_pipelines_info_cm.yaml",
				"out/export/tekton/tekton-pipelines/2_pipelines_info_rb.yaml",
				"out/export/tekton/tekton-pipelines/2_pipelines_webhook_leaderelection_rb.yaml",
				"out/export/tekton/tekton-pipelines/2_pipelines_webhook_rb.yaml",
				"out/export/tekton/tekton-pipelines/3_pipelines_controller_deploy.yaml",
				"out/export/tekton/tekton-pipelines/3_pipelines_webhook_deploy.yaml",
				"out/export/tekton/tekton-pipelines/4_pipelines_webhook_hpa.yaml",
			},
		},
	}

	for _, tt := range TT {
		tc := tt
		t.Run(
			tt.Name, func(t *testing.T) {
				t.Parallel()
				var buf bytes.Buffer
				//nolint:gocritic
				tc.Opts = append(
					tc.Opts,
					kube.WithExportWriter(&buf),
				)
				err := kube.Export(tc.km, tc.Opts...)
				tu.AssertNoError(t, err, "failed to import")
				ar := txtar.Parse(buf.Bytes())
				got := make([]string, 0, len(ar.Files))
				for _, f := range ar.Files {
					got = append(got, f.Name)
				}
				sort.Strings(got)
				want := tc.OutFiles
				if diff := tu.Diff(got, want); diff != "" {
					t.Error(diff)
				}
			},
		)
	}
}

func TestExport_SingleFile(t *testing.T) {
	var buf bytes.Buffer
	err := kube.Export(
		tekton.New(),
		kube.WithExportWriter(&buf),
		kube.WithExportOutputDirectory("out/export"),
		kube.WithExportAsSingleFile("tekton.yaml"),
	)
	tu.AssertNoError(t, err, "failed to import")
	got, err := splitManifest(&buf)
	tu.AssertNoError(t, err, "failed to split manifest")
	if len(got) != 66 {
		t.Errorf("expected 66 manifests, got %d", len(got))
	}
}

func splitManifest(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var content []string
	var buf bytes.Buffer

	for scanner.Scan() {
		txt := scanner.Text()
		switch {
		// Skip comments
		case strings.HasPrefix(txt, "#"):
			continue
		// Split by '---'
		case strings.Contains(txt, "---"):
			if buf.Len() > 0 {
				content = append(content, buf.String())
				buf.Reset()
			}
		default:
			buf.WriteString(txt + "\n")
		}
	}

	s := buf.String()
	if len(s) > 0 { // if a manifest ends with '---', don't add it
		content = append(content, s)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("spliting manifests: %w", err)
	}
	return content, nil
}

type IAM struct {
	Sa   *corev1.ServiceAccount
	Crb  *rbacv1.ClusterRoleBinding
	Cr   *rbacv1.ClusterRole
	Depl *appsv1.Deployment
}

type EmbedStruct struct {
	kube.App

	IAM
	Depl *appsv1.Deployment
}

var appName = "imthename"

var labels = map[string]string{
	"app": appName,
}

func newEmbeddedStruct() *EmbedStruct {
	sa := kubeutil.SimpleSA(appName, "defaultns")
	sa.Labels = labels

	cr := &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{Name: sa.Name, Labels: labels},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			},
		},
	}
	crb := kubeutil.SimpleCRB(sa, cr)
	crb.Labels = labels

	iam := IAM{
		Sa:  sa,
		Crb: crb,
		Cr:  cr,
		Depl: kubeutil.SimpleDeployment(
			"another"+appName,
			sa.Namespace,
			labels,
			int32(1),
			"nginx:latest",
		),
	}

	return &EmbedStruct{
		Depl: kubeutil.SimpleDeployment(
			appName,
			sa.Namespace,
			labels,
			int32(1),
			"nginx:latest",
		),
		IAM: iam,
	}
}
