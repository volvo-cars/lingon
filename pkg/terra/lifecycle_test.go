// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"bytes"
	"fmt"
	"log/slog"

	"github.com/volvo-cars/lingon/pkg/internal/hcl"
)

func ExampleLifecycle() {
	type lifecycleBlock struct {
		Lifecycle *Lifecycle `hcl:"lifecycle,block"`
	}
	tagsRef := ReferenceAsString(
		ReferenceResource(&dummyResource{}).Append(
			"tags",
		),
	)
	d := lifecycleBlock{
		Lifecycle: &Lifecycle{
			CreateBeforeDestroy: Bool(true),
			PreventDestroy:      Bool(true),
			IgnoreChanges:       IgnoreChanges(tagsRef),
			ReplaceTriggeredBy:  ReplaceTriggeredBy(tagsRef),
		},
	}

	var b bytes.Buffer
	if err := hcl.EncodeRaw(&b, d); err != nil {
		slog.Error("encoding lifecycle", "err", err)
		return
	}

	fmt.Println(b.String())
	// Output:
	// lifecycle {
	//   create_before_destroy = true
	//   prevent_destroy       = true
	//   ignore_changes        = [tags]
	//   replace_triggered_by  = [dummy.dummy.tags]
	// }
}
