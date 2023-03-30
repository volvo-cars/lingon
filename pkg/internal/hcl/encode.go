// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package hcl

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

const (
	tagHCL = "hcl"
)

type EncodeArgs struct {
	Backend       *Backend
	Providers     []Provider
	DataResources []DataResource
	Resources     []Resource
}

type Backend struct {
	Type          string
	Configuration interface{}
}

type Provider struct {
	LocalName     string
	Source        string
	Version       string
	Configuration interface{}
}

type DataResource struct {
	DataSource    string
	LocalName     string
	Configuration interface{}
}

type Resource struct {
	Type          string
	LocalName     string
	Configuration interface{}
}

type Tokenizer interface {
	// InternalTokens returns the HCL tokens that are rendered in the Terraform
	// configuration when a Terraform stack is exported.
	//
	// Internal: users should **not** use this!
	InternalTokens() hclwrite.Tokens
}

// Encode writes the HCL encoded from the given stack
// into the given io.Writer.
func Encode(wr io.Writer, args EncodeArgs) error {
	file := hclwrite.NewEmptyFile()
	fileBody := file.Body()

	// Encode terraform block
	tfBody := fileBody.AppendNewBlock("terraform", nil).Body()
	if err := encodeBackend(tfBody, args); err != nil {
		return fmt.Errorf("encoding backend: %w", err)
	}
	tfBody.AppendNewline()

	encodeRequiredProviders(tfBody, args)
	fileBody.AppendNewline()

	// Encode provider blocks
	if len(args.Providers) > 0 {
		fileBody.AppendUnstructuredTokens(hclwrite.TokensForIdentifier("// Provider blocks"))
		fileBody.AppendNewline()
	}
	for _, provider := range args.Providers {
		providerBlock := fileBody.AppendNewBlock(
			"provider",
			[]string{provider.LocalName},
		)
		rv := reflect.ValueOf(provider.Configuration)
		if err := encodeStruct(
			rv,
			providerBlock,
			providerBlock.Body(),
		); err != nil {
			return fmt.Errorf(
				"encoding provider %s: %w",
				provider.LocalName,
				err,
			)
		}
		fileBody.AppendNewline()
	}
	// Encode data blocks
	if len(args.DataResources) > 0 {
		fileBody.AppendUnstructuredTokens(hclwrite.TokensForIdentifier("// Data blocks"))
		fileBody.AppendNewline()
	}
	for _, data := range args.DataResources {
		dataBlock := fileBody.AppendNewBlock(
			"data",
			[]string{data.DataSource, data.LocalName},
		)
		rv := reflect.ValueOf(data.Configuration)
		if err := encodeStruct(
			rv,
			dataBlock,
			dataBlock.Body(),
		); err != nil {
			return fmt.Errorf(
				"encoding data resource %s.%s: %w",
				data.DataSource,
				data.LocalName,
				err,
			)
		}
		fileBody.AppendNewline()
	}
	// Encode resource blocks
	if len(args.Resources) > 0 {
		fileBody.AppendUnstructuredTokens(hclwrite.TokensForIdentifier("// Resource blocks"))
		fileBody.AppendNewline()
	}
	for _, resource := range args.Resources {
		resourceBlock := fileBody.AppendNewBlock(
			"resource",
			[]string{resource.Type, resource.LocalName},
		)
		rv := reflect.ValueOf(resource.Configuration)
		if err := encodeStruct(
			rv,
			resourceBlock,
			resourceBlock.Body(),
		); err != nil {
			return fmt.Errorf(
				"encoding resource %s.%s: %w",
				resource.Type,
				resource.LocalName,
				err,
			)
		}
		fileBody.AppendNewline()
	}

	if _, err := file.WriteTo(wr); err != nil {
		return fmt.Errorf("writing hcl: %w", err)
	}
	return nil
}

// EncodeRaw takes an empty Go interface and attempts to encode it
// using reflection and hcl tags in the provided Go struct.
// This should be used for edge cases only, and better to rely on
// Encode which takes a Exporter
func EncodeRaw(wr io.Writer, val interface{}) error {
	file := hclwrite.NewEmptyFile()
	rv := reflect.ValueOf(val)

	if err := encodeStruct(rv, nil, file.Body()); err != nil {
		return err
	}

	if _, err := file.WriteTo(wr); err != nil {
		return fmt.Errorf("writing hcl: %w", err)
	}
	return nil
}

func encodeStruct(
	rv reflect.Value,
	block *hclwrite.Block,
	body *hclwrite.Body,
) error {
	if rv.Kind() != reflect.Struct {
		if rv.Kind() == reflect.Pointer {
			return encodeStruct(rv.Elem(), block, body)
		}
		return fmt.Errorf("cannot encode non-struct type: %s", rv.Kind())
	}
	labels := make([]string, 0)

	for i := 0; i < rv.NumField(); i++ {
		sf := rv.Type().Field(i)
		fv := rv.Field(i)

		if sf.Anonymous {
			if !sf.IsExported() {
				continue
			}
			if err := encodeStruct(fv, block, body); err != nil {
				return err
			}
			continue
		}

		hclTag, ok := sf.Tag.Lookup(tagHCL)
		// Ignore fields without an HCL tag or unconventional format
		if hclTag == "" || !ok {
			continue
		}
		tagName, tagKind := splitStructTag(hclTag)
		switch tagKind {
		// If tag kind is missing, it defaults to attr
		case "", "attr":
			switch v := fv.Interface().(type) {
			case Tokenizer:
				tokens := v.InternalTokens()
				// Make sure that tokens is not nil because we don't want to
				// write empty attributes
				if tokens != nil {
					body.SetAttributeRaw(tagName, tokens)
				}
			default:
				if fv.CanInterface() && fv.Interface() != nil {
					ctyVal, err := impliedCtyValue(fv)
					if err != nil {
						return err
					}
					body.SetAttributeRaw(
						tagName,
						hclwrite.TokensForValue(ctyVal),
					)
				}
			}
		case "block":
			if !sf.IsExported() {
				return fmt.Errorf("cannot encode private field: %s", sf.Name)
			}
			if err := encodeBlock(fv, tagName, body); err != nil {
				return err
			}
		case "label":
			if fv.Kind() != reflect.String {
				return fmt.Errorf(
					"hcl `,label` tag found on non-string field: %s of type %s",
					sf.Name,
					sf.Type,
				)
			}
			label := fv.String()
			if label == "" {
				return fmt.Errorf("hcl label is empty for field: %s", sf.Name)
			}
			labels = append(labels, label)
		case "remain":
			if err := encodeRemainBody(fv, body); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown hcl label: %s", tagKind)
		}
	}

	if len(labels) > 0 {
		// When working against the top-level Go struct, no HCL block exists, so `hcl:",label"` tags
		// are not allowed here.
		// Only on the root Go struct is the block nil.
		if block == nil {
			return fmt.Errorf("cannot set hcl label tag on struct without a block")
		}
		block.SetLabels(labels)
	}

	return nil
}

func encodeBlock(
	rv reflect.Value,
	tagName string,
	body *hclwrite.Body,
) error {
	if rv.CanInterface() && rv.Interface() == nil {
		return nil
	}
	switch rv.Kind() {
	case reflect.Interface:
		// Get the underlying value of the interface
		iVal := reflect.ValueOf(rv.Interface())
		return encodeBlock(iVal, tagName, body)
	case reflect.Struct:
		newBlock := body.AppendNewBlock(tagName, nil)
		return encodeStruct(rv, newBlock, newBlock.Body())
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			if err := encodeBlock(rv.Index(i), tagName, body); err != nil {
				return err
			}
		}
	case reflect.Pointer:
		if rv.IsNil() {
			return nil
		}
		return encodeBlock(rv.Elem(), tagName, body)
	default:
		if rv.IsNil() {
			return nil
		}
		return fmt.Errorf(
			"supported type for \",block\" HCL tag: %s",
			rv.Kind(),
		)
	}
	return nil
}

func encodeRemainBody(rv reflect.Value, body *hclwrite.Body) error {
	switch rv.Kind() {
	case reflect.Map:
		iter := rv.MapRange()
		for iter.Next() {
			key := iter.Key().String()
			ctyVal, err := impliedCtyValue(iter.Value())
			if err != nil {
				return err
			}
			body.SetAttributeRaw(key, hclwrite.TokensForValue(ctyVal))
		}
	case reflect.Pointer:
		return encodeRemainBody(rv.Elem(), body)
	default:
		return fmt.Errorf(
			"supported type for \",remain\" HCL tag: %s",
			rv.Kind(),
		)
	}
	return nil
}

func splitStructTag(tag string) (string, string) {
	split := strings.Split(tag, ",")
	if len(split) == 2 {
		return split[0], split[1]
	}
	return split[0], ""
}

func impliedCtyValue(rv reflect.Value) (cty.Value, error) {
	ctyType, err := gocty.ImpliedType(rv.Interface())
	if err != nil {
		return cty.NilVal, fmt.Errorf(
			"getting implied cty type for %s: %w",
			rv,
			err,
		)
	}
	ctyVal, err := gocty.ToCtyValue(rv.Interface(), ctyType)
	if err != nil {
		return cty.NilVal, fmt.Errorf(
			"getting cty value from implied type for %s: %w",
			rv,
			err,
		)
	}
	return ctyVal, nil
}
