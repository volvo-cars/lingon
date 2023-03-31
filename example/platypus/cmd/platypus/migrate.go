package main

import (
	"context"
	"fmt"

	"github.com/volvo-cars/lingon/example/platypus/gen/providers/aws"
	"github.com/volvo-cars/lingon/example/platypus/pkg/platform/grafana"
	"github.com/volvo-cars/lingon/example/platypus/pkg/terraclient"

	"github.com/volvo-cars/lingon/pkg/terra"
	"golang.org/x/exp/slog"
)

func migrateGrafana(
	ctx context.Context,
	p runParams,
	tf *terraclient.Client,
	vpcProd vpcStack,
	eksProd eksStack,
	grafanaProd grafanaStack,
) error {
	// E.g. some PR
	env := "test-pr-12345"
	uniqueName := p.ClusterParams.Name + "-" + env

	vpcState := vpcProd.VPC.StateMust()
	dbState := grafanaProd.Postgres.StateMust()
	sgState := eksProd.SecurityGroup.StateMust()
	privateSubnetIDs := [3]string{
		vpcProd.PrivateSubnets[0].StateMust().Id,
		vpcProd.PrivateSubnets[1].StateMust().Id,
		vpcProd.PrivateSubnets[2].StateMust().Id,
	}

	snapshot := rdsSnapshotStack{
		AWSStackConfig: newAWSStackConfig(uniqueName+"-grafana-snapshot", p),
		Snapshot: aws.NewDbSnapshot(
			"grafana", aws.DbSnapshotArgs{
				DbInstanceIdentifier: terra.String(dbState.Identifier),
				DbSnapshotIdentifier: terra.String(uniqueName + "-grafana"),
			},
		),
	}

	if err := tf.Run(ctx, &snapshot); err != nil {
		return fmt.Errorf("handling snapshot: %w", err)
	}

	if !snapshot.IsStateComplete() {
		slog.Info("snapshot not sync'd")
		return nil
	}

	snapshotID := snapshot.Snapshot.StateMust().Id
	//
	// Create new Grafana instance
	//
	gmRDS, err := grafana.NewRDSPostgres(
		grafana.RDSOpts{
			Name:             uniqueName + "-grafana",
			VPCID:            vpcState.Id,
			EKSSGID:          sgState.Id,
			PrivateSubnetIDs: privateSubnetIDs,
			SnapshotID:       snapshotID,
		},
	)
	if err != nil {
		return fmt.Errorf("creating grafana rds infra: %w", err)
	}
	gm := grafanaStack{
		AWSStackConfig: newAWSStackConfig(uniqueName+"-grafana", p),
		RDSPostgres:    gmRDS,
	}
	if err := tf.Run(ctx, &gm); err != nil {
		return fmt.Errorf("handling grafana: %w", err)
	}
	if !gm.IsStateComplete() {
		slog.Info("grafana not sync'd")
		return nil
	}
	// return nil

	k, err := NewClient(
		WithClientKubeconfig(p.KubeconfigPath),
		WithClientContext(p.ClusterParams.Name),
	)
	if err != nil {
		return fmt.Errorf("creating kubectl client: %w", err)
	}

	db := gm.RDSPostgres.Postgres.StateMust()
	graf := grafana.New(
		grafana.AppOpts{
			Name:    uniqueName + "-grafana",
			Version: "9.3.8",
			Env:     env,
		},
		grafana.KubeOpts{
			PostgresHost:     db.Address,
			PostgresDBName:   db.DbName,
			PostgresUser:     db.Username,
			PostgresPassword: db.Password,
		},
	)
	if err := k.Apply(ctx, graf); err != nil {
		return fmt.Errorf("applying grafana app: %w", err)
	}

	return nil
}
