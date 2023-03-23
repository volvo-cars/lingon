package hack

//go:generate echo "\n> BUILDING\n"
//go:generate go build -mod=readonly -o ../bin/kygo github.com/volvo-cars/go-terriyaki/cmd/kygo
//go:generate go build -mod=readonly -o ../bin/explode github.com/volvo-cars/go-terriyaki/cmd/explode
//go:generate echo "\n> LINTING\n"
//go:generate golangci-lint -v run ../...
//go:generate echo "\n> TESTING\n"
//go:generate go test -mod=readonly ../...
//go:generate echo "\n> VULNERABILITIES\n"
//go:generate go run golang.org/x/vuln/cmd/govulncheck@latest ../...
