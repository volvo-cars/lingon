// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/go-playground/validator/v10"
	tfjson "github.com/hashicorp/terraform-json"
)

// Resource represents a Terraform Resource.
// The generated Go structs from a Terraform provider resource will implement this interface
// and be used when exported to Terraform Configuration
type Resource interface {
	// Type returns the resource type, e.g. aws_iam_role
	Type() string
	// LocalName returns the unique name of the resource as it will be stored in the state
	LocalName() string
	// Configuration returns the arguments for the resource
	Configuration() interface{}
	// ImportState takes the given attributes value map (from a Terraform state) and imports it
	// into this resource
	ImportState(attributes io.Reader) error
}

// DataResource represents a Terraform DataResource.
// The generated Go structs from a Terraform provider data resource will implement this interface
type DataResource interface {
	DataSource() string
	LocalName() string
	Configuration() interface{}
}

// Provider represents a Terraform Provider.
// The generated Go structs from a Terraform provider configuration will implement this interface
type Provider interface {
	LocalName() string
	Source() string
	Version() string
	Configuration() interface{}
}

// Backend represents a Terraform Backend.
// Users will define their backends to implement this interface
type Backend interface {
	BackendType() string
}

// stackObjects contains all the blocks that are extracted from a user-defined stack
type stackObjects struct {
	Backend       Backend
	Providers     []Provider
	Resources     []Resource
	DataResources []DataResource
}

const (
	tagTerriyaki = "tki"
)

var (
	ErrNoBackendBlock        = errors.New("stack must have a backend block")
	ErrMultipleBackendBlocks = errors.New("stack cannot have multiple backend blocks")
	ErrNoProviderBlock       = errors.New("stack must have a provider block")
	ErrNotExportedField      = errors.New("stack has non-exported (private) field")
	ErrUnknownPublicField    = errors.New("unknown public field")
)

// StackImportState imports the Terraform state into the Terraform Stack.
// A bool is returned indicating whether all the resources have state. If the bool is true,
// every resource in the stack will have some state. If the bool is false, the state is
// incomplete meaning some resources may have state.
func StackImportState(stack Exporter, tfState *tfjson.State) (bool, error) {
	sb, err := objectsFromStack(stack)
	if err != nil {
		return false, err
	}
	isFullState := true
	stateResources := tfState.Values.RootModule.Resources
	// Iterate over the resources in the Stack and try to find the corresponding resource in the state.
	// If it exists, import the state into the Stack.
	for _, res := range sb.Resources {
		resFound := false
		for _, sr := range stateResources {
			// Find the resource in the state. It is the same resource if the resource type
			// and resource local name match because that is how Terraform uniquely identifies
			// resources in its state.
			if res.Type() == sr.Type && res.LocalName() == sr.Name {
				resFound = true
				var b bytes.Buffer
				if err := json.NewEncoder(&b).Encode(sr.AttributeValues); err != nil {
					return false, fmt.Errorf(
						"encoding attribute values for resource %s.%s: %w",
						res.Type(), res.LocalName(), err,
					)
				}
				if err := res.ImportState(&b); err != nil {
					return false, fmt.Errorf(
						"importing state into resource %s.%s: %w",
						res.Type(), res.LocalName(), err,
					)
				}
				break
			}
		}
		if !resFound {
			isFullState = false
		}
	}
	return isFullState, nil
}

func objectsFromStack(stack Exporter) (*stackObjects, error) {
	if err := validator.New().Struct(stack); err != nil {
		return nil, fmt.Errorf("stack validation failed: %w", err)
	}

	rv := reflect.ValueOf(stack)
	// Exporter is a pointer so expect a pointer
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	sb := stackObjects{}
	if err := parseStackStructFields(rv, &sb); err != nil {
		return nil, err
	}

	return &sb, nil
}

// parseStackStructFields takes a struct reflect.Value and appends any Terraform objects that it
// finds to the provider stackObjects
func parseStackStructFields(rv reflect.Value, sb *stackObjects) error {
	for i := 0; i < rv.NumField(); i++ {
		sf := rv.Type().Field(i)
		fv := rv.Field(i)
		// Handle embedded (i.e. anonymous) structs
		if sf.Anonymous {
			if !sf.IsExported() {
				return fmt.Errorf(
					"%w, PkgPath: %s, Name: %s, Type: %s",
					ErrNotExportedField,
					sf.PkgPath,
					sf.Name,
					sf.Type,
				)
			}
			// Proceed and parse the embedded struct
			if err := parseStackStructFields(fv, sb); err != nil {
				return err
			}
			continue
		}

		// Check if field has a terriyaki struct tag
		if tkiTag, ok := sf.Tag.Lookup(tagTerriyaki); ok {
			// Ignore fields with the "-" value
			if tkiTag == "-" {
				continue
			}
		}
		if !sf.IsExported() {
			return fmt.Errorf(
				"%w, PkgPath: %s, Name: %s, Type: %s",
				ErrNotExportedField,
				sf.PkgPath,
				sf.Name,
				sf.Type,
			)
		}

		tkiObjects := make([]interface{}, 0)
		// Handle slices and arrays
		switch fv.Kind() {
		case reflect.Array, reflect.Slice:
			// If it's an array or slice, iterate over the elements and process
			// them individually
			for j := 0; j < fv.Len(); j++ {
				tkiObjects = append(tkiObjects, fv.Index(j).Interface())
			}
		default:
			tkiObjects = append(tkiObjects, fv.Interface())
		}
		for _, obj := range tkiObjects {
			switch v := obj.(type) {
			case Resource:
				sb.Resources = append(sb.Resources, v)
			case DataResource:
				sb.DataResources = append(sb.DataResources, v)
			case Provider:
				sb.Providers = append(sb.Providers, v)
			case Backend:
				if sb.Backend != nil {
					return ErrMultipleBackendBlocks
				}
				sb.Backend = v
			default:
				// Not sure if this should be an error, but rather be explicit at this point
				return fmt.Errorf(
					"%w: %s, type: %s",
					ErrUnknownPublicField,
					sf.Name,
					sf.Type,
				)
			}
		}
	}
	return nil
}
