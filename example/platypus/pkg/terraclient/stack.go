// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terraclient

import (
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/volvo-cars/lingon/pkg/terra"
)

type StateMode int

const (
	// StateModeUnknown the state mode has not been determined yet (no plan/apply)
	StateModeUnknown StateMode = 0
	// StateModeEmpty there is no state (no apply yet)
	StateModeEmpty StateMode = 1
	// StateModePartial there is a state, but there are resources in the Stack
	// that are not in the state yet (need to be applied)
	StateModePartial StateMode = 2
	// StateModeComplete the Stack and Terraform state are in complete sync
	StateModeComplete StateMode = 3
)

var _ terra.Exporter = (Stacker)(nil)

// Stacker represents a Terraform Stack.
// Embed the Stack struct into your struct
// to implement the interface, e.g.
//
// type EKSCluster struct {
//   	terra.Stack

//		IAMRole    aws.IamRole
//		EKSCluster aws.EksCluster
//		...
//	}
type Stacker interface {
	terra.Exporter
	StackName() string
	SetPlan(*tfjson.Plan)
	SetStateMode(StateMode)
	IsStateComplete() bool
	Plan() *Plan
}

var _ Stacker = (*Stack)(nil)

type Stack struct {
	// Name is the unique name of the Stack.
	// It is used for the working directory where the Terraform code is
	// generated and the Terraform CLI is executed.
	Name      string       `lingon:"-" validate:"required"`
	stateMode StateMode    `lingon:"-"`
	plan      *Plan        `lingon:"-" `
	tfplan    *tfjson.Plan `lingon:"-" `
}

func (*Stack) Terriyaki() {}

func (r *Stack) StackName() string {
	return r.Name
}

func (r *Stack) SetPlan(tfplan *tfjson.Plan) {
	r.plan = parseTfPlan(tfplan)
	r.tfplan = tfplan
}

func (r *Stack) SetStateMode(sm StateMode) {
	r.stateMode = sm
}

func (r *Stack) StateMode() StateMode {
	return r.stateMode
}

func (r *Stack) IsStateComplete() bool {
	return r.stateMode == StateModeComplete
}

func (r *Stack) Plan() *Plan {
	return r.plan
}
