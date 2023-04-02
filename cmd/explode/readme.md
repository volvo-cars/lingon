# Explode

organize kubernetes manifests according to their kind

## Usage

```shell
explode -in=<dir> -out=<dir>
expllode -in=manifest.yaml  # default -out=out/
explode -out=<dir> < manifest.yaml  # default -in=- # read from stdin
```

## Example

At the root of the repo, run:

```shell
go build -o bin/explode ./cmd/explode && ./bin/explode -in=./pkg/kube/testdata
```

output:

```shell
$ ls -R out/
_cluster  tekton-pipelines  tekton-pipelines-resolvers

out/_cluster:
crd  namespace  rbac  webhook

out/_cluster/crd:
1_clustertasks.tekton.dev_crd.yaml       1_resolutionrequests.resolution.tekton.dev_crd.yaml
1_customruns.tekton.dev_crd.yaml         1_runs.tekton.dev_crd.yaml
1_pipelineresources.tekton.dev_crd.yaml  1_taskruns.tekton.dev_crd.yaml
1_pipelineruns.tekton.dev_crd.yaml       1_tasks.tekton.dev_crd.yaml
1_pipelines.tekton.dev_crd.yaml          1_verificationpolicies.tekton.dev_crd.yaml

out/_cluster/namespace:
0_tekton-pipelines-resolvers_ns.yaml  0_tekton-pipelines_ns.yaml

out/_cluster/rbac:
1_tekton-aggregate-edit_cr.yaml                                  1_tekton-pipelines-webhook-cluster-access_cr.yaml
1_tekton-aggregate-view_cr.yaml                                  2_tekton-pipelines-controller-cluster-access_crb.yaml
1_tekton-pipelines-controller-cluster-access_cr.yaml             2_tekton-pipelines-controller-tenant-access_crb.yaml
1_tekton-pipelines-controller-tenant-access_cr.yaml              2_tekton-pipelines-resolvers_crb.yaml
1_tekton-pipelines-resolvers-resolution-request-updates_cr.yaml  2_tekton-pipelines-webhook-cluster-access_crb.yaml

out/_cluster/webhook:
4_config.webhook.pipeline.tekton.dev_validatingwebhookconfigurations.yaml
4_validation.webhook.pipeline.tekton.dev_validatingwebhookconfigurations.yaml
4_webhook.pipeline.tekton.dev_mutatingwebhookconfigurations.yaml

out/tekton-pipelines:
1_tekton-pipelines-controller_role.yaml       2_config-registry-cert_cm.yaml
1_tekton-pipelines-controller_sa.yaml         2_config-spire_cm.yaml
1_tekton-pipelines-controller_svc.yaml        2_config-trusted-resources_cm.yaml
1_tekton-pipelines-info_role.yaml             2_feature-flags_cm.yaml
1_tekton-pipelines-leader-election_role.yaml  2_pipelines-info_cm.yaml
1_tekton-pipelines-webhook_role.yaml          2_tekton-pipelines-controller-leaderelection_rb.yaml
1_tekton-pipelines-webhook_sa.yaml            2_tekton-pipelines-controller_rb.yaml
1_tekton-pipelines-webhook_svc.yaml           2_tekton-pipelines-info_rb.yaml
2_config-artifact-bucket_cm.yaml              2_tekton-pipelines-webhook-leaderelection_rb.yaml
2_config-artifact-pvc_cm.yaml                 2_tekton-pipelines-webhook_rb.yaml
2_config-defaults_cm.yaml                     2_webhook-certs_secrets.yaml
2_config-leader-election_cm.yaml              3_tekton-pipelines-controller_deploy.yaml
2_config-logging_cm.yaml                      3_tekton-pipelines-webhook_deploy.yaml
2_config-observability_cm.yaml                4_tekton-pipelines-webhook_hpa.yaml

out/tekton-pipelines-resolvers:
1_tekton-pipelines-resolvers-namespace-rbac_role.yaml  2_config-observability_cm.yaml
1_tekton-pipelines-resolvers_sa.yaml                   2_git-resolver-config_cm.yaml
2_bundleresolver-config_cm.yaml                        2_hubresolver-config_cm.yaml
2_cluster-resolver-config_cm.yaml                      2_resolvers-feature-flags_cm.yaml
2_config-leader-election_cm.yaml                       2_tekton-pipelines-resolvers-namespace-rbac_rb.yaml
2_config-logging_cm.yaml                               3_tekton-pipelines-remote-resolvers_deploy.yaml

```
