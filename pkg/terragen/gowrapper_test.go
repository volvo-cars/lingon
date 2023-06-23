// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terragen

import (
	"testing"

	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestParseProvider(t *testing.T) {
	type test struct {
		providerStr string
		provider    Provider
		expectErr   bool
		errmsg      string
	}

	tests := []test{
		{
			providerStr: "aws=hashicorp/aws:4.60.0",
			provider: Provider{
				Name:    "aws",
				Source:  "hashicorp/aws",
				Version: "4.60.0",
			},
			expectErr: false,
		},
		{
			providerStr: "awshashicorp/aws:4.60.0",
			expectErr:   true,
			errmsg:      "provider format incorrect: missing `=`",
		},
		{
			providerStr: "aws=hashicorp/aws",
			expectErr:   true,
			errmsg:      "provider format incorrect: missing `:` in `source:version`",
		},
		{
			providerStr: "aws=hashicorp/aws",
			expectErr:   true,
			errmsg:      "provider format incorrect: missing `:` in `source:version`",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.providerStr, func(t *testing.T) {
				ap, err := ParseProvider(tt.providerStr)
				if tt.expectErr {
					tu.AssertErrorMsg(t, err, tt.errmsg)
				} else {
					tu.AssertNoError(t, err, "parsing provider")
					tu.AssertEqual(t, tt.provider, ap)
				}
			},
		)
	}
}
