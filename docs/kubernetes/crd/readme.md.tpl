# Example CRD

All can be done in one step:

```bash
go generate -v ./...
```

## Import the manifest in Go

{{ "Example_import" | example }}


## Export the Go struct to YAML


{{ "Example_export" | example }}


## How to find CRD types

This is more of a heuristic than a rule.

- go to the repo of the project
- search for `AddToScheme` function

Often it is in the `pkg/apis` or `apis` folder of the project.

## List of well-known CRDs types

> This is best-effort

```go
import(
    "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/runtime/serializer"
    apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
    capiahelm "sigs.k8s.io/cluster-api-addon-provider-helm/api/v1alpha1"
    certmanager "github.com/cert-manager/cert-manager/pkg/api"
    clientgoscheme "k8s.io/client-go/kubernetes/scheme"
    externalsecretsv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
    gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
    gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
    helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
    imageautov1 "github.com/fluxcd/image-automation-controller/api/v1beta1"
    imagereflectv1 "github.com/fluxcd/image-reflector-controller/api/v1beta2"
    istionetworkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
    istionetworkingv1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
    istiosecurityv1beta1 "istio.io/client-go/pkg/apis/security/v1beta1"
    karpenterapi "github.com/aws/karpenter/pkg/apis"
    karpenterv1alpha5 "github.com/aws/karpenter-core/pkg/apis"
    knativecachingalpha1 "knative.dev/caching/pkg/apis/caching/v1alpha1"
    knativeservingv1 "knative.dev/serving/pkg/apis/serving/v1"
    kservev1alpha1 "github.com/kserve/kserve/pkg/apis/serving/v1alpha1"
    kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta2"
    metacontrolleralpha "github.com/metacontroller/metacontroller/pkg/apis/metacontroller/v1alpha1"
    notificationv1b2 "github.com/fluxcd/notification-controller/api/v1beta2"
    otelv1alpha1 "github.com/open-telemetry/opentelemetry-operator/apis/v1alpha1"
    profilev1 "github.com/kubeflow/kubeflow/components/profile-controller/api/v1"
    profilev1beta1 "github.com/kubeflow/kubeflow/components/profile-controller/api/v1beta1"
    secretsstorev1 "sigs.k8s.io/secrets-store-csi-driver/apis/v1"
    servingv1alpha1 "github.com/kserve/modelmesh-serving/apis/serving/v1alpha1"
    slothv1alpha1 "github.com/slok/sloth/pkg/kubernetes/api/sloth/v1"
    sourcev1 "github.com/fluxcd/source-controller/api/v1"
    tektonpipelinesv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
    tektontriggersv1alpha1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
    utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)
```