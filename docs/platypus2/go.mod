module github.com/volvo-cars/lingoneks

go 1.20

replace github.com/volvo-cars/lingon => ../../

require (
	github.com/VictoriaMetrics/metricsql v0.56.2
	github.com/VictoriaMetrics/operator/api v0.0.0-20230703080649-c199848d02f4
	github.com/ardanlabs/conf/v3 v3.1.6
	github.com/aws/karpenter v0.28.1
	github.com/aws/karpenter-core v0.28.1
	github.com/golingon/terraproviders/aws/5.6.2 v0.0.0-20230703111924-4b2c49f97f7c
	github.com/golingon/terraproviders/tls/4.0.4 v0.0.0-20230703111924-4b2c49f97f7c
	github.com/grafana/dashboard-linter v0.0.0-20230622143601-02e2cd156626
	github.com/hashicorp/terraform-exec v0.18.1
	github.com/hashicorp/terraform-json v0.17.1
	github.com/nats-io/nats.go v1.27.1
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.66.0
	github.com/prometheus/client_golang v1.16.0
	github.com/rogpeppe/go-internal v1.11.0
	github.com/tidwall/gjson v1.14.4
	github.com/volvo-cars/lingon v0.0.0-20230703105113-1bcac3444c58
	github.com/zeitlinger/conflate v0.0.0-20230622100834-279724abda8c
	go.uber.org/automaxprocs v1.5.2
	golang.org/x/exp v0.0.0-20230626212559-97b1e661b5df
	google.golang.org/grpc v1.56.1
	google.golang.org/protobuf v1.31.0
	k8s.io/api v0.27.3
	k8s.io/apiextensions-apiserver v0.27.3
	k8s.io/apimachinery v0.27.3
	k8s.io/kube-aggregator v0.27.3
	sigs.k8s.io/kind v0.20.0
	sigs.k8s.io/yaml v1.3.0
)

replace (
	k8s.io/api => k8s.io/api v0.26.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.23.1
	k8s.io/apimachinery => k8s.io/apimachinery v0.26.2
	k8s.io/client-go => k8s.io/client-go v0.26.2
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.26.2
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20230501164219-8b0f38b5fd1f
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.14.0
)

require (
	cloud.google.com/go v0.110.2 // indirect
	cloud.google.com/go/compute v1.19.3 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v1.0.1 // indirect
	cloud.google.com/go/storage v1.30.1 // indirect
	github.com/BurntSushi/toml v1.0.0 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/VictoriaMetrics/VictoriaMetrics v1.91.3 // indirect
	github.com/VictoriaMetrics/fasthttp v1.2.0 // indirect
	github.com/VictoriaMetrics/metrics v1.24.0 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/alessio/shellescape v1.4.1 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/aws/aws-sdk-go v1.44.294 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bmatcuk/doublestar/v4 v4.6.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dave/jennifer v1.6.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dennwc/varint v1.0.0 // indirect
	github.com/eidolon/wordwrap v0.0.0-20161011182207-e0f54129b8bb // indirect
	github.com/emicklei/go-restful/v3 v3.10.2 // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/gnostic v0.6.9 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/s2a-go v0.1.3 // indirect
	github.com/google/safetext v0.0.0-20220905092116-b49f7bc46da2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.3 // indirect
	github.com/googleapis/gax-go/v2 v2.8.0 // indirect
	github.com/grafana/regexp v0.0.0-20221122212121-6b5c0a4cb7fd // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/hcl/v2 v2.17.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.16.6 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nats-io/nats-server/v2 v2.9.17 // indirect
	github.com/nats-io/nkeys v0.4.4 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.4.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.11.0 // indirect
	github.com/prometheus/prometheus v0.44.0 // indirect
	github.com/samber/lo v1.38.1 // indirect
	github.com/spf13/cobra v1.7.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fastjson v1.6.4 // indirect
	github.com/valyala/fastrand v1.1.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/valyala/histogram v1.2.0 // indirect
	github.com/valyala/quicktemplate v1.7.0 // indirect
	github.com/veggiemonk/strcase v0.0.0-20230627213939-a882c834bcab // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/zclconf/go-cty v1.13.2 // indirect
	github.com/zegl/kube-score v1.16.1 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/goleak v1.2.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.10.0 // indirect
	golang.org/x/mod v0.11.0 // indirect
	golang.org/x/net v0.11.0 // indirect
	golang.org/x/oauth2 v0.9.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
	golang.org/x/term v0.9.0 // indirect
	golang.org/x/text v0.10.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/tools v0.10.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	gomodules.xyz/jsonpatch/v2 v2.3.0 // indirect
	google.golang.org/api v0.123.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230525234025-438c736192d0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230525234020-1aefcd67740a // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230629202037-9506855d4529 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/client-go v1.5.2 // indirect
	k8s.io/component-base v0.27.3 // indirect
	k8s.io/klog/v2 v2.100.1 // indirect
	k8s.io/kube-openapi v0.0.0-20230614213217-ba0abe644833 // indirect
	k8s.io/utils v0.0.0-20230505201702-9f6742963106 // indirect
	knative.dev/pkg v0.0.0-20230628105954-6eb4b40a9a30 // indirect
	mvdan.cc/gofumpt v0.5.0 // indirect
	sigs.k8s.io/controller-runtime v0.15.0 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)
