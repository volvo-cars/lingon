// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package terrajen

import (
	"fmt"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

type schemaAttribute struct {
	Name string
	tfjson.SchemaAttribute
}

type schemaBlockType struct {
	Name string
	tfjson.SchemaBlockType
}

func TestGraph_AttributeString(t *testing.T) {
	ex := schemaAttribute{
		Name: "string",
		SchemaAttribute: tfjson.SchemaAttribute{
			AttributeType: cty.String,
			Description:   "string",
			Required:      true,
		},
	}
	sb := tfjson.SchemaBlock{
		Attributes: map[string]*tfjson.SchemaAttribute{
			ex.Name: &ex.SchemaAttribute,
		},
	}

	g := newGraph(&sb)
	require.Len(t, g.attributes, 1)
	ac := g.attributes[0]
	assert.Equal(t, ac.name, ex.Name)
	assert.True(t, ex.AttributeType.Equals(ac.ctyType))
	assert.True(t, ac.isArg)
	assert.True(t, ac.isRequired)
}

func TestGraph_Attributes(t *testing.T) {
	count := 10
	attrs := make(map[string]*tfjson.SchemaAttribute, count)
	for i := 0; i < count; i++ {
		attrs[fmt.Sprintf("%d", i)] = &tfjson.SchemaAttribute{
			AttributeType: cty.String,
			Description:   "string",
			Required:      true,
		}
	}
	sb := tfjson.SchemaBlock{
		Attributes: attrs,
	}

	g := newGraph(&sb)
	assert.Len(t, g.attributes, 10)
}

func TestGraph_Blocks(t *testing.T) {
	exAttr := schemaAttribute{
		Name: "string",
		SchemaAttribute: tfjson.SchemaAttribute{
			AttributeType: cty.String,
			Description:   "string",
			Required:      true,
		},
	}
	ex := schemaBlockType{
		Name: "block",
		SchemaBlockType: tfjson.SchemaBlockType{
			NestingMode: tfjson.SchemaNestingModeSingle,
			Block: &tfjson.SchemaBlock{
				Attributes: map[string]*tfjson.SchemaAttribute{
					exAttr.Name: &exAttr.SchemaAttribute,
				},
			},
			// Make it required
			MinItems: 1,
			MaxItems: 1,
		},
	}
	sb := tfjson.SchemaBlock{
		NestedBlocks: map[string]*tfjson.SchemaBlockType{
			ex.Name: &ex.SchemaBlockType,
		},
	}

	g := newGraph(&sb)
	require.Len(t, g.children, 1)
	assert.Len(t, g.nodes, 1)

	ac := g.children[0]
	require.Len(t, ac.attributes, 1)
	assert.Equal(t, ex.Name, ac.name)
	assert.True(t, ac.isRequired)
	assert.True(t, ac.isArg)
	acc := ac.attributes[0]
	assert.Equal(t, acc.name, exAttr.Name)
	assert.True(t, exAttr.AttributeType.Equals(acc.ctyType))
	assert.True(t, acc.isArg)
	assert.True(t, acc.isRequired)
}
