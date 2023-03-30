// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import "github.com/hashicorp/hcl/v2"

func newRef(address ...string) Reference {
	tr := make(hcl.Traversal, len(address))
	for i, s := range address {
		if i == 0 {
			tr[i] = hcl.TraverseRoot{Name: s}
		} else {
			tr[i] = hcl.TraverseAttr{Name: s}
		}
	}
	return Reference{
		tr: tr,
	}
}
