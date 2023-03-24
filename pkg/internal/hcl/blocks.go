// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package hcl

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func encodeBackend(body *hclwrite.Body, args EncodeArgs) error {
	if args.Backend == nil {
		return nil
	}
	be := args.Backend
	// Encode backend block
	backendBlock := body.AppendNewBlock(
		"backend",
		[]string{be.Type},
	)
	rv := reflect.ValueOf(be.Configuration)
	if err := encodeStruct(
		rv,
		backendBlock,
		backendBlock.Body(),
	); err != nil {
		return fmt.Errorf(
			"encoding backend %s: %w",
			be.Type,
			err,
		)
	}
	return nil
}

func encodeRequiredProviders(body *hclwrite.Body, args EncodeArgs) {
	if len(args.Providers) == 0 {
		return
	}
	body.AppendNewline()
	reqProvBody := body.AppendNewBlock(
		"required_providers",
		nil,
	).Body()
	for _, prov := range args.Providers {
		reqProvBody.SetAttributeValue(
			prov.LocalName, cty.MapVal(
				map[string]cty.Value{
					"source":  cty.StringVal(prov.Source),
					"version": cty.StringVal(prov.Version),
				},
			),
		)
	}
}
