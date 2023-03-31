// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractBlocks_Simple(t *testing.T) {
	type simpleStack struct {
		DummyStack
		DummyRes  *dummyResource     `validate:"required"`
		DummyData *dummyDataResource `validate:"required"`
	}
	dr := &dummyResource{}
	ddr := &dummyDataResource{}
	st := simpleStack{
		DummyStack: newDummyBaseStack(),
		DummyRes:   dr,
		DummyData:  ddr,
	}
	sb, err := objectsFromStack(&st)
	require.NoError(t, err)
	require.Len(t, sb.Resources, 1)
	assert.Equal(t, dr, sb.Resources[0])
	require.Len(t, sb.DataResources, 1)
	assert.Equal(t, ddr, sb.DataResources[0])
}

func TestExtractBlocks_Complex(t *testing.T) {
	type DummyModule struct {
		Resource *dummyResource     `validate:"required"`
		Data     *dummyDataResource `validate:"required"`
	}
	type complexStack struct {
		DummyStack
		DummyModule
		SliceRes []*dummyResource      `validate:"required,dive,required"`
		OneRes   [1]*dummyResource     `validate:"required,dive,required"`
		OneData  [1]*dummyDataResource `validate:"required,dive,required"`
	}
	dr := &dummyResource{}
	ddr := &dummyDataResource{}
	st := complexStack{
		DummyStack: newDummyBaseStack(),
		DummyModule: DummyModule{
			Resource: dr,
			Data:     ddr,
		},
		SliceRes: []*dummyResource{dr, dr},
		OneRes:   [1]*dummyResource{dr},
		OneData:  [1]*dummyDataResource{ddr},
	}
	sb, err := objectsFromStack(&st)
	require.NoError(t, err)
	assert.Len(t, sb.Resources, 4)
	assert.Len(t, sb.DataResources, 2)
}

func TestExtractBlocks_UnknownField(t *testing.T) {
	type simpleStack struct {
		DummyStack
		UnknownString string
	}
	st := simpleStack{
		DummyStack: newDummyBaseStack(),
	}
	_, err := objectsFromStack(&st)
	assert.ErrorIs(t, err, ErrUnknownPublicField)
}

func TestExtractBlocks_PrivateField(t *testing.T) {
	type simpleStack struct {
		DummyStack
		privateField string //nolint:structcheck,unused
	}
	st := simpleStack{
		DummyStack: newDummyBaseStack(),
	}
	_, err := objectsFromStack(&st)
	assert.ErrorIs(t, err, ErrNotExportedField)
}

func newDummyBaseStack() DummyStack {
	return DummyStack{
		Backend:  &dummyBackend{},
		Provider: &dummyProvider{},
	}
}

func TestExtractBlocks_IgnoredField(t *testing.T) {
	type simpleStack struct {
		DummyStack
		UnknownString string `lingon:"-"`
	}
	st := simpleStack{
		DummyStack: newDummyBaseStack(),
	}
	_, err := objectsFromStack(&st)
	require.NoError(t, err)
}

type DummyStack struct {
	Stack
	Backend  *dummyBackend
	Provider *dummyProvider
}

//
// Dummy Backend
//

var _ Backend = (*dummyBackend)(nil)

type dummyBackend struct{}

func (b dummyBackend) BackendType() string {
	return "dummy"
}

//
// Dummy Provider
//

var _ Provider = (*dummyProvider)(nil)

type dummyProvider struct{}

func (p *dummyProvider) LocalName() string {
	return "dummy"
}

func (p *dummyProvider) Source() string {
	return "dummy"
}

func (p *dummyProvider) Version() string {
	return "dummy"
}

func (p *dummyProvider) Configuration() interface{} {
	return dummyConfig
}

//
// Dummy Resources
//

var _ Resource = (*dummyResource)(nil)

type dummyResource struct{}

func (r *dummyResource) Type() string {
	return "dummy"
}

func (r *dummyResource) LocalName() string {
	return "dummy"
}

func (r *dummyResource) Configuration() interface{} {
	return dummyConfig
}

func (r *dummyResource) ImportState(av io.Reader) error {
	return nil
}

func (r *dummyResource) Dependencies() Dependencies {
	return nil
}

func (r *dummyResource) LifecycleManagement() *Lifecycle {
	return nil
}

//
// Dummy Data Resources
//

var _ DataResource = (*dummyDataResource)(nil)

type dummyDataResource struct{}

func (d *dummyDataResource) DataSource() string {
	return "dummy"
}

func (d *dummyDataResource) LocalName() string {
	return "dummy"
}

func (d *dummyDataResource) Configuration() interface{} {
	return dummyConfig
}

//
// Dummy Args / Configuration
//

var dummyConfig = dummyArgs{Name: "dummy"}

type dummyArgs struct {
	Name string `hcl:"name,attr"`
}
