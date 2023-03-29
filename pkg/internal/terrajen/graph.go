// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terrajen

import (
	"fmt"
	"sort"

	"github.com/volvo-cars/lingon/pkg/internal/str"

	"github.com/dave/jennifer/jen"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"
)

func newGraph(schema *tfjson.SchemaBlock) *graph {
	g := graph{}

	for _, atName := range sortMapKeys(schema.Attributes) {
		attr := schema.Attributes[atName]
		// Attributes which are objects will be treated as children
		if _, ok := ctyTypeElementObject(attr.AttributeType); ok {
			g.children = append(
				g.children,
				g.traverseCtyType(
					nil,
					atName,
					attr.AttributeType,
					isAttributeArg(attr),
				),
			)
			continue
		}
		g.attributes = append(
			g.attributes, &attribute{
				name:       atName,
				ctyType:    attr.AttributeType,
				isArg:      isAttributeArg(attr),
				isRequired: attr.Required,
			},
		)
	}

	for _, blockName := range sortMapKeys(schema.NestedBlocks) {
		blockType := schema.NestedBlocks[blockName]
		g.children = append(
			g.children,
			g.traverseBlockType(nil, blockName, blockType),
		)
	}

	g.calculateUniqueNames()

	return &g
}

// graph is used to decouple the tfjson.SchemaBlock type from the code generator.
// A graph is created by supplying a tfjson.SchemaBlock.
type graph struct {
	// attributes are the top-level attributes for a terraform configuration
	attributes []*attribute
	// children contains all top-level nodes for a terraform configuration
	children []*node
	// nodes contains all nodes in the schema block
	nodes []*node
}

type node struct {
	name string
	// path is the list of nodes from the root (resource/data/provider)
	// node to this node.
	// E.g. if a resource "x" has a block "y" which contains this node's
	// block "z", then
	// name = "z"
	// path = ["x", "y"]
	path []string
	// uniqueName is a calculated unique name within the subpkg for a
	// resource/data/provider object.
	// If the name is already unique, then the uniqueName is the name.
	// If the name is not unique, take the last element from the path (if any)
	// and prefix is to the name until the name is unique
	uniqueName  string
	uniqueDepth int
	attributes  []*attribute

	// block values
	isArg bool
	// nestingPath is the path from one Go struct to it's child.
	// For Terraform Schema Blocks this is only ever going to be
	// length 0 (for single block) or 1 (for list/set/map).
	// For Cty types, this is a bit more complicated as an object
	// can be nested within a list/set/map with no limit to the depth,
	// hence why we use a list here, and that list represents the
	// path from object --> object. E.g. if we have the cty pseudo type
	// list(set(map(object))) then the nestingPath would be:
	// [list, set, map]
	nestingPath []nodeNestingMode
	isRequired  bool
	// minItems is taken from the HCL block type, and is not set
	// when using Cty types.
	minItems uint64
	// maxItems is taken from the HCL block type, and is not set
	// when using Cty types.
	maxItems uint64

	children []*node

	receiver string
}

func (n *node) argsStructName() string {
	return str.PascalCase(n.uniqueName)
}

func (n *node) attributesStructName() string {
	return str.PascalCase(n.uniqueName) + suffixAttributes
}

func (n *node) stateStructName() string {
	return str.PascalCase(n.uniqueName) + suffixState
}

func (n *node) isSingularArg() bool {
	return len(n.nestingPath) == 0 || n.maxItems == 1
}

func (n *node) isSingularState() bool {
	return len(n.nestingPath) == 0
}

func (n *node) comment() string {
	if n.isSingularArg() {
		required := "optional"
		if n.isRequired {
			required = "required"
		}
		return fmt.Sprintf("%s: %s", str.PascalCase(n.uniqueName), required)
	}
	return fmt.Sprintf(
		"%s: %s", str.PascalCase(n.uniqueName), nodeBlockListValidateTags(n),
	)
}

type nodeNestingMode int

const (
	nodeNestingModeList nodeNestingMode = 1
	nodeNestingModeSet  nodeNestingMode = 2
	nodeNestingModeMap  nodeNestingMode = 3
)

type attribute struct {
	name    string
	ctyType cty.Type
	// isArg is true if the attribute can be passed as an argument
	// to the node schema block, else it is false
	isArg bool
	// isRequired is true if the attribute can be passed as an argument
	// and is required, else it is false
	isRequired bool
}

func (a *attribute) comment() string {
	required := "optional"
	if a.isRequired {
		required = "required"
	}

	return fmt.Sprintf(
		"%s: %s, %s",
		str.PascalCase(a.name),
		a.ctyType.FriendlyName(),
		required,
	)
}

func (g *graph) isEmpty() bool {
	return len(g.nodes) == 0
}

func (g *graph) traverseBlockType(
	path []string,
	name string,
	blockType *tfjson.SchemaBlockType,
) *node {
	n := node{
		name:        name,
		path:        path,
		uniqueName:  name,
		nestingPath: blockNodeNestingMode(blockType),
		isRequired:  isArgBlockRequired(blockType),
		minItems:    blockType.MinItems,
		maxItems:    blockType.MaxItems,
		// Blocks can always be arguments
		isArg:    true,
		receiver: structReceiverFromName(name),
	}
	g.nodes = append(g.nodes, &n)

	// First handle the attributes
	for _, atName := range sortMapKeys(blockType.Block.Attributes) {
		attr := blockType.Block.Attributes[atName]
		// Attributes which are objects will be treated as children
		if _, ok := ctyTypeElementObject(attr.AttributeType); ok {
			n.children = append(
				n.children,
				g.traverseCtyType(
					appendPath(path, name),
					atName,
					attr.AttributeType,
					isAttributeArg(attr),
				),
			)
			continue
		}
		n.attributes = append(
			n.attributes, &attribute{
				name:       atName,
				ctyType:    attr.AttributeType,
				isArg:      isAttributeArg(attr),
				isRequired: attr.Required,
			},
		)
	}

	for _, blockName := range sortMapKeys(blockType.Block.NestedBlocks) {
		blockType := blockType.Block.NestedBlocks[blockName]
		n.children = append(
			n.children,
			g.traverseBlockType(appendPath(path, name), blockName, blockType),
		)
	}
	return &n
}

func (g *graph) traverseCtyType(
	path []string,
	name string,
	ct cty.Type,
	isArg bool,
) *node {
	n := node{
		name:        name,
		path:        path,
		uniqueName:  name,
		nestingPath: ctyNodeNestingMode(ct),
		isRequired:  false,
		isArg:       isArg,
		receiver:    structReceiverFromName(name),
	}
	g.nodes = append(g.nodes, &n)

	// Get the underlying object
	obj, _ := ctyTypeElementObject(ct)
	for _, atName := range sortMapKeys(obj.AttributeTypes()) {
		at := obj.AttributeType(atName)
		// If there are objects within the attributes of this object, traverse those objects
		// and make them children of this object
		if _, ok := ctyTypeElementObject(at); ok {
			n.children = append(
				n.children,
				g.traverseCtyType(appendPath(path, name), atName, at, isArg),
			)
			continue
		}
		n.attributes = append(
			n.attributes, &attribute{
				name:       atName,
				ctyType:    at,
				isArg:      isArg,
				isRequired: false,
			},
		)
	}
	return &n
}

func (g *graph) calculateUniqueNames() {
	// We hate this function
	uniq := false
	for !uniq {
		uniq = true
		dict := make(map[string][]*node)
		for _, n := range g.nodes {
			dict[n.uniqueName] = append(dict[n.uniqueName], n)
		}

		for _, nodes := range dict {
			if len(nodes) == 1 {
				continue
			}
			uniq = false

			for _, n := range nodes {
				if len(n.path) > n.uniqueDepth {
					pathIndex := len(n.path) - n.uniqueDepth - 1
					n.uniqueName = n.path[pathIndex] + "." + n.uniqueName
					n.uniqueDepth += 1
				}
			}
		}
	}
}

func appendPath(path []string, name string) []string {
	newPath := make([]string, len(path)+1)
	copy(newPath, path)
	newPath[len(path)] = name
	return newPath
}

func ctyTypeElementObject(ct cty.Type) (cty.Type, bool) {
	switch {
	case ct.IsObjectType():
		return ct, true
	case ct.IsCollectionType():
		return ctyTypeElementObject(ct.ElementType())
	default:
		return cty.NilType, false
	}
}

func blockNodeNestingMode(block *tfjson.SchemaBlockType) []nodeNestingMode {
	switch block.NestingMode {
	case tfjson.SchemaNestingModeSingle, tfjson.SchemaNestingModeGroup:
		return nil
	case tfjson.SchemaNestingModeList, tfjson.SchemaNestingModeMap:
		// Unintuitively, tfjson.SchemaNestingModeMap is not actually a map, just a list,
		// but they get keyed by the block labels into a Map.
		// For our use case, we therefore treat it like a list.
		return []nodeNestingMode{nodeNestingModeList}
	case tfjson.SchemaNestingModeSet:
		return []nodeNestingMode{nodeNestingModeSet}
	default:
		panic(
			fmt.Sprintf(
				"unsupported SchemaNestingMode: %s",
				block.NestingMode,
			),
		)
	}
}

// ctyNodeNestingMode calculates the nesting mode for a given cty.Type that is
// known to contain a cty.Object.
// The nesting path is from parent to the child cty.Object,
// which may contain lists, sets or maps along the way.
// I.e. what is the path from parent --> cty.Object (child)
func ctyNodeNestingMode(ct cty.Type) []nodeNestingMode {
	if ct.IsObjectType() {
		return nil
	}
	switch {
	case ct.IsListType():
		return append(
			[]nodeNestingMode{nodeNestingModeList},
			ctyNodeNestingMode(ct.ElementType())...,
		)
	case ct.IsSetType():
		return append(
			[]nodeNestingMode{nodeNestingModeSet},
			ctyNodeNestingMode(ct.ElementType())...,
		)
	case ct.IsMapType():
		ct.IsCollectionType()
		return append(
			[]nodeNestingMode{nodeNestingModeMap},
			ctyNodeNestingMode(ct.ElementType())...,
		)
	default:
		panic(
			fmt.Sprintf(
				"unsupported cty.Type nesting mode: %s",
				ct.FriendlyName(),
			),
		)
	}
}

func isArgBlockRequired(block *tfjson.SchemaBlockType) bool {
	if block.MaxItems == 1 && block.MinItems == 1 {
		return true
	}
	return false
}

func isAttributeArg(attr *tfjson.SchemaAttribute) bool {
	// If it is not computed, then it must be an argument.
	// If it is computed, and it is optional, then it can be an argument.
	// If it is required, then it must be an argument.
	switch {
	case !attr.Computed:
		return true
	case attr.Computed && attr.Optional:
		return true
	case attr.Required:
		return true
	default:
		return false
	}
}

// returnTypeFromNestingPath returns the jen statement for the type that
// node represents, e.g.
//
//	terra.ListValue[emrcluster.StepRef]
func returnTypeFromNestingPath(
	nestingPath []nodeNestingMode,
	qual *jen.Statement,
) *jen.Statement {
	fullQual := qual.Clone()
	// Iterate over nestingPath backwards to get the correct order.
	// E.g. if List-->Set-->Map is the nestingPath, then we can reverse this to
	// easily construct List[Set[Map[qual]]]
	for i := len(nestingPath) - 1; i >= 0; i-- {
		path := nestingPath[i]
		switch path {
		case nodeNestingModeList:
			fullQual = qualListValue().Types(fullQual.Clone())
		case nodeNestingModeSet:
			fullQual = qualSetValue().Types(fullQual.Clone())
		case nodeNestingModeMap:
			fullQual = qualMapValue().Types(fullQual.Clone())
		default:
			panic(fmt.Sprintf("unsupport node nesting type: %d", nestingPath))
		}
	}
	return fullQual
}

// jenNodeReturnValue returns the jen statement for the return value that
// node represents, e.g.
//
//	terra.ReferenceList[emrcluster.StepRef](terra.Reference("a","b","c")
func jenNodeReturnValue(
	n *node, childQual *jen.Statement,
) *jen.Statement {
	// If the child is not nested in any lists/sets/maps then simply return
	// a reference to the single child qualifier
	if len(n.nestingPath) == 0 {
		return qualReferenceSingle().Types(childQual)
	}

	first := n.nestingPath[0]
	remainder := n.nestingPath[1:]
	subType := returnTypeFromNestingPath(remainder, childQual)

	var fullQual *jen.Statement
	switch first {
	case nodeNestingModeList:
		fullQual = qualReferenceList().Types(subType)
	case nodeNestingModeSet:
		fullQual = qualReferenceSet().Types(subType)
	case nodeNestingModeMap:
		fullQual = qualReferenceMap().Types(subType)
	default:
		panic(fmt.Sprintf("unsupport node nesting type: %d", n.nestingPath))
	}
	return fullQual
}

func sortMapKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
