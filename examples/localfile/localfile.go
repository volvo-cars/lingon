package localfile

//go:generate go run -mod=readonly github.com/volvo-cars/go-terriyaki/cmd/terragen -out . -pkg github.com/volvo-cars/go-terriyaki/examples/localfile -force -provider local=hashicorp/local:2.4.0

import (
	"github.com/volvo-cars/go-terriyaki/examples/localfile/providers/local"
	"github.com/volvo-cars/go-terriyaki/pkg/terra"
)

// NewLocalFileStack returns a new LocalFileStack which implements the terra.Exporter interface
// and can be exported into Terraform configuration
func NewLocalFileStack(filename string) *LocalFileStack {
	return &LocalFileStack{
		Backend: LocalBackend{
			Path: "terraform.tfstate",
		},
		Provider: local.NewProvider(local.ProviderArgs{}),
		File: local.NewFile(
			"file", local.FileArgs{
				Filename: terra.String(filename),
				Content:  terra.String("contents"),
			},
		),
	}
}

type LocalFileStack struct {
	terra.Stack

	Backend  LocalBackend    `validate:"required"`
	Provider *local.Provider `validate:"required"`
	File     *local.File     `validate:"required"`
}

var _ terra.Backend = (*LocalBackend)(nil)

// LocalBackend implements the Terraform local backend type.
// https://developer.hashicorp.com/terraform/language/settings/backends/local
type LocalBackend struct {
	Path string `hcl:"path,attr" validate:"required"`
}

// BackendType defines the type of the backend and implements the terra.Backend interface
func (b LocalBackend) BackendType() string {
	return "local"
}
