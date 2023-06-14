module github.com/volvo-cars/lingoneks

go 1.20

replace github.com/volvo-cars/lingon => ../../

require (
	github.com/K-Phoen/grabana v0.21.18
	github.com/ardanlabs/conf/v3 v3.1.5
	github.com/aws/karpenter v0.27.5
	github.com/aws/karpenter-core v0.27.5
	github.com/golingon/terraproviders/aws/5.0.1 v0.0.0-20230527233228-68663550bae0
	github.com/golingon/terraproviders/tls/4.0.4 v0.0.0-20230527233228-68663550bae0
	github.com/hashicorp/terraform-exec v0.18.1
	github.com/hashicorp/terraform-json v0.17.0
	github.com/nats-io/nats.go v1.26.0
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.65.1
	github.com/prometheus/client_golang v1.14.0
	github.com/rogpeppe/go-internal v1.10.0
	github.com/tidwall/gjson v1.14.4
	github.com/volvo-cars/lingon v0.0.0-20230529113525-2f8eb8598205
	go.uber.org/automaxprocs v1.5.1
	golang.org/x/exp v0.0.0-20230522175609-2e198f4a06a1
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.30.0
	k8s.io/api v0.27.2
	k8s.io/apiextensions-apiserver v0.27.2
	k8s.io/apimachinery v0.27.2
	k8s.io/kube-aggregator v0.27.2
	sigs.k8s.io/yaml v1.3.0
)

require (
	github.com/K-Phoen/sdk v0.12.2 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/aws/aws-sdk-go v1.44.271 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dave/jennifer v1.6.1 // indirect
	github.com/eidolon/wordwrap v0.0.0-20161011182207-e0f54129b8bb // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/gosimple/slug v1.13.1 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/hcl/v2 v2.17.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nats-io/nats-server/v2 v2.9.17 // indirect
	github.com/nats-io/nkeys v0.4.4 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.39.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/samber/lo v1.38.1 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
	github.com/veggiemonk/strcase v0.0.0-20230526161048-ad38aa882cb5 // indirect
	github.com/zclconf/go-cty v1.13.2 // indirect
	github.com/zegl/kube-score v1.16.1 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/mod v0.10.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	golang.org/x/tools v0.9.1 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/client-go v0.27.2 // indirect
	k8s.io/klog/v2 v2.100.1 // indirect
	k8s.io/utils v0.0.0-20230505201702-9f6742963106 // indirect
	knative.dev/pkg v0.0.0-20230525143525-9bda38b21643 // indirect
	mvdan.cc/gofumpt v0.5.0 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)
