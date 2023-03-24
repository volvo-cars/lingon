// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package terra

// Exporter represents a Terraform Stack.
// Embed the Stack struct into your struct
// to implement the interface, e.g.
//
// type EKSCluster struct {
//   	terra.Stack

//		IAMRole    aws.IamRole
//		EKSCluster aws.EksCluster
//		...
//	}
type Exporter interface {
	// Terriyaki is the original name of this project when it was being built, and is used to
	// explicitly mark a struct as implementing the Exporter interface
	Terriyaki()
}

var _ Exporter = (*Stack)(nil)

// Stack minimally implements the Exporter interface and can be embedded into user-defined stacks
// to make them implement the Exporter interface
type Stack struct{}

func (*Stack) Terriyaki() {}
