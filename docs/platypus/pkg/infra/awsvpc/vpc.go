// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package awsvpc

import (
	"fmt"

	aws "github.com/golingon/terraproviders/aws/4.60.0"

	"github.com/volvo-cars/lingon/pkg/terra"
)

var (
	S        = terra.String
	N        = terra.Number
	B        = terra.Bool
	Anywhere = S("0.0.0.0/0")
)

type Opts struct {
	Name               string
	AZs                [3]string
	CIDR               string
	PublicSubnetCIDRs  [3]string
	PrivateSubnetCIDRs [3]string
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

const (
	TagManagedBy      = "ManagedBy"
	TagManagedByValue = "Lingon"
	// TagName human-readable resource name. Note that the AWS Console UI displays the case-sensitive "Name" tag.
	TagName = "Name"
	// TagAppID is a tag specifying the application identifier, application using the resource.
	TagAppID = "app-id"
	// TagAppRole is a tag specifying the resource's technical function, e.g. webserver, database, etc.
	TagAppRole = "app-role"
	// TagPurpose  is a tag specifying the resource's business purpose, e.g. "frontend ui", "payment processor", etc.
	TagPurpose = "purpose"
	// TagEnv is a tag specifying the environment.
	TagEnv = "environment"
	// TagProject is a tag specifying the project.
	TagProject = "project"
	// TagOwner is a tag specifying the person of contact.
	TagOwner = "owner"
	// TagCostCenter is a tag specifying the cost center that will receive the bill.
	TagCostCenter = "cost-center"
	// TagAutomationExclude is a tag specifying if the resource should be excluded from automation.
	// Value: true/false
	TagAutomationExclude = "automation-exclude"
	// TagPII is a tag specifying if the resource contains Personally Identifiable Information.
	// Value: true/false
	TagPII = "pii"
)

func stags(ss ...string) terra.MapValue[terra.StringValue] {
	sv := make(map[string]terra.StringValue, 0)
	for i := 0; i < len(ss); i += 2 {
		if i+1 >= len(ss) {
			panic("odd number of strings")
		}
		sv[ss[i]] = S(ss[i+1])
	}
	sv[TagManagedBy] = S(TagManagedByValue)

	return terra.Map(sv)
}

func ttags(m map[string]string) terra.MapValue[terra.StringValue] {
	sv := make(map[string]terra.StringValue, 0)
	for k, v := range m {
		sv[k] = S(v)
	}
	sv[TagManagedBy] = S(TagManagedByValue)
	return terra.Map(sv)
}

func mergeTags(m ...map[string]string) terra.MapValue[terra.StringValue] {
	sv := make(map[string]terra.StringValue, 0)
	for _, mm := range m {
		for k, v := range mm {
			sv[k] = S(v)
		}
	}
	sv[TagManagedBy] = S(TagManagedByValue)
	return terra.Map(sv)
}

func mergeSTags(m map[string]string, ss ...string) terra.MapValue[terra.StringValue] {
	sv := make(map[string]terra.StringValue, 0)
	for k, v := range m {
		sv[k] = S(v)
	}
	for i := 0; i < len(ss); i += 2 {
		if i+1 >= len(ss) {
			sv[ss[i]] = S("")
			break
		}
		sv[ss[i]] = S(ss[i+1])
	}
	sv[TagManagedBy] = S(TagManagedByValue)
	return terra.Map(sv)
}

func NewAWSVPC(opts Opts) *AWSVPC {
	name := opts.Name

	tags := func(name string, tags ...string) terra.MapValue[terra.StringValue] {
		ss := []string{TagName, name}
		ss = append(ss, tags...)
		return mergeSTags(opts.CommonTags, ss...)
	}

	vpc := aws.NewVpc(
		name, aws.VpcArgs{
			CidrBlock: S(opts.CIDR),
			// Tags:             ttags(map[string]string{TagName: opts.Name}),
			InstanceTenancy:  S("default"),
			EnableDnsSupport: B(true),
			Tags:             tags(opts.Name),
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
				Vpc:  B(true),
				Tags: tags("nat-" + opts.AZs[i]),
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
				Tags: mergeSTags(opts.CommonTags,
					TagName, name+"-private",
					"karpenter.sh/discovery", "platypus-1",
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
