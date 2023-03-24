package terraform_test

import (
	"bytes"
	"fmt"

	"github.com/volvo-cars/lingon/pkg/terra"
	"golang.org/x/exp/slog"
)

type MinimalStack struct {
	terra.Stack
}

func ExampleMinimalStack() {
	// Initialise the minimal stack
	stack := MinimalStack{}
	// Export the stack to Terraform HCL
	var b bytes.Buffer
	if err := terra.ExportWriter(&stack, &b); err != nil {
		slog.Error("exporting stack", "err", err)
		return
	}
	fmt.Println(b.String())

	// Output:
	// terraform {
	// }
}
