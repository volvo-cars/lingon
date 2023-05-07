// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"path/filepath"
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
		name     string
		km       kube.Exporter
		opts     []kube.ExportOption
		outFiles []string
		err      error
		isErr    bool
	}
	outEBDS := filepath.Join(defaultExportOutputDir, "embeddedstruct")
	outTekton := filepath.Join(defaultExportOutputDir, "tekton")
	TT := []args{
		{
			name: "export embedded struct",
			km:   newEmbeddedStruct(),
			opts: []kube.ExportOption{
				kube.WithExportKustomize(true),
				kube.WithExportOutputDirectory(outEBDS),
			},
			outFiles: []string{
				"out/export/embeddedstruct/1_iamcr.yaml",
				"out/export/embeddedstruct/1_iamsa.yaml",
				"out/export/embeddedstruct/2_iamcrb.yaml",
				"out/export/embeddedstruct/3_depl.yaml",
				"out/export/embeddedstruct/3_iamdepl.yaml",
				"out/export/embeddedstruct/kustomization.yaml",
			},
		},
		{
			name: "export embedded struct explode",
			km:   newEmbeddedStruct(),
			opts: []kube.ExportOption{
				kube.WithExportExplodeManifests(true),
				kube.WithExportOutputDirectory(outEBDS),
			},
			outFiles: []string{
				"out/export/embeddedstruct/_cluster/rbac/1_iamcr.yaml",
				"out/export/embeddedstruct/_cluster/rbac/2_iamcrb.yaml",
				"out/export/embeddedstruct/defaultns/1_iamsa.yaml",
				"out/export/embeddedstruct/defaultns/3_depl.yaml",
				"out/export/embeddedstruct/defaultns/3_iamdepl.yaml",
			},
		},
		{
			name: "export embedded struct with name file func",
			km:   newEmbeddedStruct(),
			opts: []kube.ExportOption{
				kube.WithExportOutputDirectory(outEBDS),
				kube.WithExportNameFileFunc(
					func(name *kubeutil.Metadata) string {
						return strings.ToLower(name.Kind) + "_" + name.Meta.Name + ".yaml"
					},
				),
			},
			outFiles: []string{
				"out/export/embeddedstruct/clusterrole_imthename.yaml",
				"out/export/embeddedstruct/clusterrolebinding_imthename.yaml",
				"out/export/embeddedstruct/deployment_anotherimthename.yaml",
				"out/export/embeddedstruct/deployment_imthename.yaml",
				"out/export/embeddedstruct/serviceaccount_imthename.yaml",
			},
		},
		{
			name: "export embedded struct with explode and name file func",
			km:   newEmbeddedStruct(),
			opts: []kube.ExportOption{
				kube.WithExportOutputDirectory(outEBDS),
				kube.WithExportExplodeManifests(true),
				kube.WithExportNameFileFunc(
					func(name *kubeutil.Metadata) string {
						return strings.ToLower(name.Kind) + "_" + name.Meta.Name + ".yaml"
					},
				),
			},
			outFiles: []string{
				"out/export/embeddedstruct/_cluster/rbac/clusterrole_imthename.yaml",
				"out/export/embeddedstruct/_cluster/rbac/clusterrolebinding_imthename.yaml",
				"out/export/embeddedstruct/defaultns/deployment_anotherimthename.yaml",
				"out/export/embeddedstruct/defaultns/deployment_imthename.yaml",
				"out/export/embeddedstruct/defaultns/serviceaccount_imthename.yaml",
			},
		},
		{
			name: "export embedded struct with explode and name file func as JSON",
			km:   newEmbeddedStruct(),
			opts: []kube.ExportOption{
				kube.WithExportOutputDirectory(outEBDS),
				kube.WithExportExplodeManifests(true),
				kube.WithExportOutputJSON(true),
				kube.WithExportNameFileFunc(
					func(name *kubeutil.Metadata) string {
						return strings.ToLower(name.Kind) + "_" + name.Meta.Name + ".json"
					},
				),
			},
			outFiles: []string{
				"out/export/embeddedstruct/_cluster/rbac/clusterrole_imthename.json",
				"out/export/embeddedstruct/_cluster/rbac/clusterrolebinding_imthename.json",
				"out/export/embeddedstruct/defaultns/deployment_anotherimthename.json",
				"out/export/embeddedstruct/defaultns/deployment_imthename.json",
				"out/export/embeddedstruct/defaultns/serviceaccount_imthename.json",
			},
		},
		{
			name: "export incompatible options explode with single file",
			km:   newEmbeddedStruct(),
			opts: []kube.ExportOption{
				kube.WithExportOutputDirectory(outEBDS),
				kube.WithExportExplodeManifests(true),
				kube.WithExportAsSingleFile("single.yaml"),
				kube.WithExportKustomize(true),
			},
			err:   kube.ErrIncompatibleOptions,
			isErr: true,
		},
		{
			name: "export incompatible options json with kustomize",
			km:   newEmbeddedStruct(),
			opts: []kube.ExportOption{
				kube.WithExportOutputDirectory(outEBDS),
				kube.WithExportExplodeManifests(true),
				kube.WithExportOutputJSON(true),
				kube.WithExportKustomize(true),
			},
			err:   kube.ErrIncompatibleOptions,
			isErr: true,
		},
		{
			name: "export tekton",
			km:   tekton.New(),
			opts: []kube.ExportOption{
				kube.WithExportOutputDirectory(outTekton),
			},
			outFiles: []string{
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
			name: "export remove secrets",
			km:   tekton.New(),
			opts: []kube.ExportOption{
				kube.WithExportOutputDirectory(outTekton),
				kube.WithExportKustomize(true),
				kube.WithExportExplodeManifests(true),
				kube.WithExportSecretHook(
					func(s *corev1.Secret) error {
						return nil
					},
				),
			},
			outFiles: []string{
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
				"out/export/tekton/kustomization.yaml",
			},
		},
	}

	for _, tt := range TT {
		tc := tt
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()
				var buf bytes.Buffer
				tc.opts = append(
					tc.opts,
					kube.WithExportWriter(&buf),
				)
				err := kube.Export(tc.km, tc.opts...)
				if tc.isErr {
					if !errors.Is(err, tc.err) {
						t.Fatalf("expected error %v but got %v", tc.err, err)
					}
					return
				}
				tu.AssertNoError(t, tc.err, "failed to check error")
				tu.AssertNoError(t, err, "failed to export")
				got := txtar.Parse(buf.Bytes())
				tu.AssertEqualSlice(t, tu.Filenames(got), tc.outFiles)
				f := filepath.Join(
					"testdata",
					"golden",
					exportGoldenFileName(tc.name),
				)
				want, err := txtar.ParseFile(f)
				tu.AssertNoError(t, err, "failed to parse expected txtar")
				if diff := tu.DiffTxtar(got, want); diff != "" {
					t.Fatal(tu.Callers(), diff)
				}
			},
		)
	}
}

func exportGoldenFileName(s string) string {
	return strings.ReplaceAll(s, " ", "_") + ".txt"
}

func TestExport_SingleFileJSON(t *testing.T) {
	var buf bytes.Buffer

	err := kube.Export(
		tekton.New(),
		kube.WithExportOutputDirectory("out/export"),
		kube.WithExportOutputJSON(true),
		kube.WithExportAsSingleFile("tekton.json"),
		kube.WithExportWriter(&buf),
	)
	tu.AssertNoError(t, err, "failed to import")

	var got []map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &got)
	tu.AssertNoError(t, err, "failed to unmarshal json")
	tu.IsEqual(t, 65, len(got))
}

func TestExport_SingleFileYAML(t *testing.T) {
	var buf bytes.Buffer

	err := kube.Export(
		tekton.New(),
		kube.WithExportWriter(&buf),
		kube.WithExportOutputDirectory("out/export"),
		kube.WithExportAsSingleFile("tekton.yaml"),
	)
	tu.AssertNoError(t, err, "failed to import")
	got, err := kubeutil.ManifestSplit(&buf)
	tu.AssertNoError(t, err, "failed to split manifest")
	tu.IsEqual(t, 65, len(got))
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
				APIGroups: []string{""},
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
