// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package infra

import (
	"fmt"

	aws "github.com/golingon/terraproviders/aws/5.0.1"
	"github.com/volvo-cars/lingon/pkg/terra"
)

type Opts struct {
	Name               string
	AZs                [3]string
	CIDR               string
	PublicSubnetCIDRs  [3]string
	PrivateSubnetCIDRs [3]string
	KarpenterDiscovery string
	CommonTags         map[string]string
}

type AWSVPC struct {
	VPC *aws.Vpc `validate:"required"`

	PublicSubnets  [3]*aws.Subnet                `validate:"required,dive,required"`
	PublicRT       *aws.RouteTable               `validate:"required"`
	PublicRoute    *aws.Route                    `validate:"required"`
	PublicRTAssocs [3]*aws.RouteTableAssociation `validate:"required,dive,required"`

	PrivateSubnets  [3]*aws.Subnet                `validate:"required,dive,required"`
	PrivateRTs      [3]*aws.RouteTable            `validate:"required,dive,required"`
	PrivateRoutes   [3]*aws.Route                 `validate:"required,dive,required"`
	PrivateRTAssocs [3]*aws.RouteTableAssociation `validate:"required,dive,required"`

	InternetGateway *aws.InternetGateway `validate:"required"`
	EIPNat          [3]*aws.Eip          `validate:"required,dive,required"`
	NatGateways     [3]*aws.NatGateway   `validate:"required,dive,required"`
}

func NewAWSVPC(opts Opts) *AWSVPC {
	name := opts.Name

	tags := func(
		name string,
		tags ...string,
	) terra.MapValue[terra.StringValue] {
		ss := []string{TagName, name}
		ss = append(ss, tags...)
		return MergeSTags(opts.CommonTags, ss...)
	}

	vpc := aws.NewVpc(
		name, aws.VpcArgs{
			CidrBlock: S(opts.CIDR),
			// Tags:             ttags(map[string]string{TagName: opts.Name}),
			InstanceTenancy:    S("default"),
			EnableDnsSupport:   B(true), // must be enabled for EFS
			EnableDnsHostnames: B(true), // must be enabled for EFS
			Tags:               tags(opts.Name),
		},
	)

	igw := aws.NewInternetGateway(
		name, aws.InternetGatewayArgs{
			VpcId: vpc.Attributes().Id(),
			Tags:  tags(name + "-igw"),
		},
	)

	eipNats := [3]*aws.Eip{}
	for i := 0; i < 3; i++ {
		eipNats[i] = aws.NewEip(
			fmt.Sprintf("nats_%d", i), aws.EipArgs{
				// Vpc:    B(true), // deprecated
				Domain: S("vpc"),
				Tags:   tags("nat-" + opts.AZs[i]),
			},
		)
	}

	publicSubnets := [3]*aws.Subnet{}
	for i := 0; i < 3; i++ {
		publicSubnets[i] = aws.NewSubnet(
			fmt.Sprintf("public_%d", i), aws.SubnetArgs{
				VpcId:               vpc.Attributes().Id(),
				AvailabilityZone:    S(opts.AZs[i]),
				CidrBlock:           S(opts.PublicSubnetCIDRs[i]),
				MapPublicIpOnLaunch: terra.Bool(true),
				Tags:                tags(name + "-public"),
			},
		)
	}

	publicRT := aws.NewRouteTable(
		"public", aws.RouteTableArgs{
			VpcId: vpc.Attributes().Id(),
			Tags:  tags(name + "-public"),
		},
	)
	publicRoute := aws.NewRoute(
		"public", aws.RouteArgs{
			DestinationCidrBlock: Anywhere,
			RouteTableId:         publicRT.Attributes().Id(),
			GatewayId:            igw.Attributes().Id(),
		},
	)

	pubRTAssocs := [3]*aws.RouteTableAssociation{}
	for i := 0; i < 3; i++ {
		pubRTAssocs[i] = aws.NewRouteTableAssociation(
			fmt.Sprintf("public_%d", i), aws.RouteTableAssociationArgs{
				SubnetId:     publicSubnets[i].Attributes().Id(),
				RouteTableId: publicRT.Attributes().Id(),
			},
		)
	}

	natGateways := [3]*aws.NatGateway{}
	for i := 0; i < 3; i++ {
		ng := aws.NewNatGateway(
			fmt.Sprintf("nat_gateway_%d", i), aws.NatGatewayArgs{
				SubnetId:     publicSubnets[i].Attributes().Id(),
				AllocationId: eipNats[i].Attributes().Id(),
				Tags:         tags(fmt.Sprintf("ng-%d", i)),
			},
		)
		ng.DependsOn = terra.DependsOn(igw)
		natGateways[i] = ng
	}

	privateSubnets := [3]*aws.Subnet{}
	for i := 0; i < 3; i++ {
		privateSubnets[i] = aws.NewSubnet(
			fmt.Sprintf("private_%d", i), aws.SubnetArgs{
				VpcId:            vpc.Attributes().Id(),
				AvailabilityZone: S(opts.AZs[i]),
				CidrBlock:        S(opts.PrivateSubnetCIDRs[i]),
				Tags: MergeSTags(
					opts.CommonTags,
					TagName, name+"-private",
					KarpenterDiscoveryKey, opts.KarpenterDiscovery,
				),
			},
		)
	}

	privateRTs := [3]*aws.RouteTable{}
	for i := 0; i < 3; i++ {
		privateRTs[i] = aws.NewRouteTable(
			fmt.Sprintf("private_%d", i), aws.RouteTableArgs{
				VpcId: vpc.Attributes().Id(),
				Tags:  tags(fmt.Sprintf("platypus-private-%d", i)),
			},
		)
	}
	privateRoutes := [3]*aws.Route{}
	for i := 0; i < 3; i++ {
		privateRoutes[i] = aws.NewRoute(
			fmt.Sprintf("private_%d", i), aws.RouteArgs{
				RouteTableId:         privateRTs[i].Attributes().Id(),
				DestinationCidrBlock: Anywhere,
				NatGatewayId:         natGateways[i].Attributes().Id(),
			},
		)
	}

	privateRTAssocs := [3]*aws.RouteTableAssociation{}
	for i := 0; i < 3; i++ {
		privateRTAssocs[i] = aws.NewRouteTableAssociation(
			fmt.Sprintf("private_%d", i), aws.RouteTableAssociationArgs{
				SubnetId:     privateSubnets[i].Attributes().Id(),
				RouteTableId: privateRTs[i].Attributes().Id(),
			},
		)
	}
	return &AWSVPC{
		VPC:             vpc,
		InternetGateway: igw,
		EIPNat:          eipNats,
		PublicSubnets:   publicSubnets,
		PublicRT:        publicRT,
		PublicRoute:     publicRoute,
		PublicRTAssocs:  pubRTAssocs,

		NatGateways:     natGateways,
		PrivateSubnets:  privateSubnets,
		PrivateRTs:      privateRTs,
		PrivateRoutes:   privateRoutes,
		PrivateRTAssocs: privateRTAssocs,
	}
}
