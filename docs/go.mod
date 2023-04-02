module github.com/volvo-cars/lingon/docs

go 1.20

replace (
	github.com/volvo-cars/lingon => ../
	github.com/volvo-cars/lingon/docs => ./
)

require (
	github.com/Masterminds/semver/v3 v3.2.0
	github.com/aws/karpenter v0.27.1
	github.com/aws/karpenter-core v0.27.1
	github.com/go-playground/validator/v10 v10.12.0
	github.com/golingon/terraproviders/aws/4.60.0 v0.0.0-20230331133707-9ae37d0b1bd3
	github.com/golingon/terraproviders/tls/4.0.4 v0.0.0-20230331133707-9ae37d0b1bd3
	github.com/google/go-containerregistry v0.14.0
	github.com/hashicorp/hcl/v2 v2.16.2
	github.com/hashicorp/terraform-exec v0.18.1
	github.com/hashicorp/terraform-json v0.16.0
	github.com/hexops/valast v1.4.3
	github.com/invopop/yaml v0.2.0
	github.com/rogpeppe/go-internal v1.10.0
	github.com/stretchr/testify v1.8.2
	github.com/volvo-cars/lingon v0.0.0-20230331132342-46c7f66ade5e
	golang.org/x/exp v0.0.0-20230321023759-10a507213a29
	istio.io/api v0.0.0-20230217221049-9d422bf48675
	istio.io/client-go v1.17.1
	k8s.io/api v0.26.3
	k8s.io/apiextensions-apiserver v0.26.3
	k8s.io/apimachinery v0.26.3
	k8s.io/client-go v0.26.3
	sigs.k8s.io/secrets-store-csi-driver v1.3.2
)

require (
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/aws/aws-sdk-go v1.44.195 // indirect
	github.com/dave/jennifer v1.6.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/cli v23.0.1+incompatible // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v23.0.1+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.7.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.16.0 // indirect
	github.com/leodido/go-urn v1.2.2 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/samber/lo v1.37.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/veggiemonk/strcase v0.0.0-20230325182039-9fa4e7cee676 // indirect
	github.com/zclconf/go-cty v1.13.1 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/mod v0.9.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	golang.org/x/tools v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20221018160656-63c7b68cfc55 // indirect
	google.golang.org/protobuf v1.29.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/klog/v2 v2.90.1 // indirect
	k8s.io/utils v0.0.0-20230313181309-38a27ef9d749 // indirect
	knative.dev/pkg v0.0.0-20221123154742-05b694ec4d3a // indirect
	mvdan.cc/gofumpt v0.4.0 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
