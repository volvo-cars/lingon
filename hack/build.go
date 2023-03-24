// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package hack

// TODO: CONVERT THIS TO MAGEFILES

//go:generate echo "\n> BUILDING\n"
//go:generate go build -mod=readonly -o ../bin/kygo github.com/volvo-cars/lingon/cmd/kygo
//go:generate go build -mod=readonly -o ../bin/explode github.com/volvo-cars/lingon/cmd/explode
//go:generate go build -mod=readonly -o ../bin/terragen github.com/volvo-cars/lingon/cmd/terragen

//go:generate echo "\n> LINTING\n"
// Check license headers. Remove the --plan argument to apply any necessary
// changes
// go:generate go run github.com/hashicorp/copywrite@v0.16.3 headers --dirPath ./.. --config ./../.copywrite.hcl --plan
//go:generate go run mvdan.cc/gofumpt@latest -extra -w ../
//go:generate go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.2 -v run ../...

//go:generate echo "\n> TESTING\n"
//go:generate go test -mod=readonly ../...

//go:generate echo "\n> SBOM\n"
//go:generate go run github.com/anchore/syft/cmd/syft@latest packages dir:.. -o spdx-json --file ../sbom.json

//go:generate echo "\n> VULNERABILITIES\n"
//go:generate go run golang.org/x/vuln/cmd/govulncheck@latest ../...
//go:generate go run github.com/google/osv-scanner/cmd/osv-scanner@v1 --sbom sbom.json --recursive ..

//go:generate echo "\n> LICENSES\n"
//go:generate go run github.com/google/go-licenses@v1.6.0 check ../...
//go:generate go run github.com/google/go-licenses@v1.6.0 save ../... --save_path=../bin/licenses --force
