// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package grafana

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	aws "github.com/golingon/terraproviders/aws/4.60.0"
	"github.com/golingon/terraproviders/aws/4.60.0/securitygroup"

	"github.com/volvo-cars/lingon/pkg/terra"
)

var (
	B = terra.Bool
	N = terra.Number
	S = terra.String
)

type RDSOpts struct {
	Name             string    `validate:"required"`
	VPCID            string    `validate:"required"`
	EKSSGID          string    `validate:"required"`
	PrivateSubnetIDs [3]string `validate:"required,dive,required"`
	// SnapshotID (optional) to restore the database from
	SnapshotID string
}

func NewRDSPostgres(
	opts RDSOpts,
) (RDSPostgres, error) {
	if err := validator.New().Struct(opts); err != nil {
		return RDSPostgres{}, fmt.Errorf("validating opts: %w", err)
	}
	sg := aws.NewSecurityGroup(
		"grafana", aws.SecurityGroupArgs{
			VpcId:       S(opts.VPCID),
			Name:        S(opts.Name),
			Description: S("Security group for access to Grafana RDS"),
			Ingress: []securitygroup.Ingress{
				{
					FromPort:       N(5432),
					ToPort:         N(5432),
					Protocol:       S("tcp"),
					Description:    S("Allow access from Grafana pods to RDS"),
					SecurityGroups: terra.Set(S(opts.EKSSGID)),
				},
			},
		},
	)
	dbSubnetGroup := aws.NewDbSubnetGroup(
		"grafana", aws.DbSubnetGroupArgs{
			Name:        S(opts.Name),
			Description: S("Grafana RDS subnet group"),
			SubnetIds:   terra.SetString(opts.PrivateSubnetIDs[:]...),
		},
	)

	rds := aws.NewDbInstance(
		"grafana", aws.DbInstanceArgs{
			Identifier:              S(opts.Name),
			Engine:                  S("postgres"),
			EngineVersion:           S("14.4"),
			AutoMinorVersionUpgrade: B(false),

			InstanceClass:       S("db.t4g.micro"),
			AllocatedStorage:    N(20),
			MaxAllocatedStorage: N(50),

			// TODO: need to do something better here
			DbName:   S("grafana"),
			Username: S("grafana"),
			Password: S("platypusgrafana"),

			DbSubnetGroupName:   dbSubnetGroup.Attributes().Id(),
			VpcSecurityGroupIds: terra.Set(sg.Attributes().Id()),

			PubliclyAccessible: B(false),
			MultiAz:            B(false),

			SkipFinalSnapshot: B(true),
		},
	)
	if opts.SnapshotID != "" {
		rds.Args.SnapshotIdentifier = S(opts.SnapshotID)
	}
	return RDSPostgres{
		SecurityGroup: sg,
		SubnetGroup:   dbSubnetGroup,
		Postgres:      rds,
	}, nil
}

type RDSPostgres struct {
	SecurityGroup *aws.SecurityGroup `validate:"required"`
	SubnetGroup   *aws.DbSubnetGroup `validate:"required"`
	Postgres      *aws.DbInstance    `validate:"required"`
}
