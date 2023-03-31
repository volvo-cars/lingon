// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/volvo-cars/lingon/example/platypus/pkg/infra/eks"
	"github.com/volvo-cars/lingon/example/platypus/pkg/infra/vpc"
	"github.com/volvo-cars/lingon/example/platypus/pkg/platform/awsauth"
	"github.com/volvo-cars/lingon/example/platypus/pkg/platform/grafana"
	"github.com/volvo-cars/lingon/example/platypus/pkg/platform/karpenter"
	"github.com/volvo-cars/lingon/example/platypus/pkg/terraclient"

	"github.com/volvo-cars/lingon/pkg/terra"

	"golang.org/x/exp/slog"
)

var (
	S = terra.String
	N = terra.Number
)

func main() {
	var apply bool
	var destroy bool
	var plan bool
	var migrate bool
	flag.BoolVar(
		&apply,
		"apply",
		false,
		"Apply the terraform changes (default: false)",
	)
	flag.BoolVar(
		&destroy,
		"destroy",
		false,
		"Destroy the terraform resources (default: false)",
	)
	flag.BoolVar(
		&plan,
		"plan",
		false,
		"Plan the terraform changes (default: false)",
	)
	flag.BoolVar(
		&migrate,
		"migrate",
		false,
		"Migrate Grafana (default: false)",
	)
	flag.Parse()

	ap := AWSParams{
		BackendS3Key: "terriyaki-tf-experiment",
		Region:       "eu-north-1",
		Profile:      "vcc-cdds-prod-legacy",
	}
	p := runParams{
		Apply:          apply,
		Destroy:        destroy,
		Plan:           plan,
		Migrate:        migrate,
		AWSParams:      ap,
		KubeconfigPath: "kubeconfig",
		ManifestPath:   ".kart/k8s",
		ClusterParams: ClusterParams{
			Name:    "platypus-1",
			Version: "1.24",
			ID:      1,
		},
		TFLabels: map[string]string{
			"environment": "dev",
			"terraform":   "true",
		},
		KLabels: map[string]string{
			"environment": "development",
		},
	}

	if err := run(p); err != nil {
		slog.Error("run", "err", err)
		os.Exit(1)
	}
	slog.Info("done")
}

type runParams struct {
	AWSParams      AWSParams
	KubeconfigPath string
	ManifestPath   string
	ClusterParams  ClusterParams
	TFLabels       map[string]string
	KLabels        map[string]string
	Apply          bool
	Destroy        bool
	Plan           bool
	Migrate        bool
}
type AWSParams struct {
	BackendS3Key string
	Region       string
	Profile      string
}
type ClusterParams struct {
	Name    string
	Version string
	ID      int
}

func run(p runParams) error {
	ctx := context.Background()
	uniqueName := p.ClusterParams.Name
	vpcOpts := vpc.Opts{
		Name: uniqueName,
		AZs: [3]string{
			"eu-north-1a", "eu-north-1b", "eu-north-1c",
		},
		CIDR: "10.0.0.0/16",
		PublicSubnetCIDRs: [3]string{
			"10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24",
		},
		PrivateSubnetCIDRs: [3]string{
			"10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24",
		},
	}

	tf := terraclient.NewClient(
		terraclient.WithDefaultPlan(p.Plan),
		terraclient.WithDefaultApply(p.Apply),
	)

	vpc := vpcStack{
		AWSStackConfig: newAWSStackConfig(uniqueName+"-vpc", p),
		AWSVPC:         *vpc.NewAWSVPC(vpcOpts),
	}

	if err := tf.Run(ctx, &vpc); err != nil {
		return fmt.Errorf("tfrun: handling vpc: %w", err)
	}
	if !vpc.IsStateComplete() {
		slog.Info("VPC state not in sync, finishing here")
		return finishAndDestroy(ctx, p, tf)
	}

	vpcState := vpc.AWSVPC.VPC.StateMust()
	privateSubnetIDs := [3]string{}
	for i, subnet := range vpc.AWSVPC.PrivateSubnets {
		privateSubnetIDs[i] = subnet.StateMust().Id
	}
	vpcID := vpcState.Id
	eks := eksStack{
		AWSStackConfig: newAWSStackConfig(uniqueName+"-eks", p),
		Cluster: *eks.NewEKSCluster(
			eks.ClusterOpts{
				Name:             p.ClusterParams.Name,
				Version:          p.ClusterParams.Version,
				VPCID:            vpcID,
				PrivateSubnetIDs: privateSubnetIDs,
			},
		),
	}
	if err := tf.Run(ctx, &eks); err != nil {
		return fmt.Errorf("tfrun: handling cluster: %w", err)
	}
	if !eks.IsStateComplete() {
		slog.Info("EKS cluster state not in sync, finishing here")
		return finishAndDestroy(ctx, p, tf)
	}

	eksSGID := eks.SecurityGroup.StateMust().Id
	eksState := eks.EKSCluster.StateMust()
	oidcState := eks.IAMOIDCProvider.StateMust()
	ks := karpenterStack{
		AWSStackConfig: newAWSStackConfig(uniqueName+"-karpenter", p),
		Infra: karpenter.NewInfra(
			karpenter.InfraOpts{
				Name:             eksState.Name + "-karpenter",
				ClusterName:      eksState.Name,
				ClusterARN:       eksState.Arn,
				PrivateSubnetIDs: privateSubnetIDs,
				OIDCProviderArn:  oidcState.Arn,
				OIDCProviderURL:  oidcState.Url,
			},
		),
	}
	if err := tf.Run(ctx, &ks); err != nil {
		return fmt.Errorf("terraforming karpenter: %w", err)
	}
	if !ks.IsStateComplete() {
		slog.Info(
			"stack state not in sync",
			slog.String("stack", ks.StackName()),
		)
		return finishAndDestroy(ctx, p, tf)
	}

	gmRDS, err := grafana.NewRDSPostgres(
		grafana.RDSOpts{
			Name:             uniqueName + "-grafana",
			VPCID:            vpcID,
			EKSSGID:          eksSGID,
			PrivateSubnetIDs: privateSubnetIDs,
		},
	)
	if err != nil {
		return fmt.Errorf("creating grafana rds infra: %w", err)
	}
	gs := grafanaStack{
		AWSStackConfig: newAWSStackConfig(uniqueName+"-grafana", p),
		RDSPostgres:    gmRDS,
	}
	if err := tf.Run(ctx, &gs); err != nil {
		return fmt.Errorf("terraforming grafana: %w", err)
	}
	if !gs.IsStateComplete() {
		slog.Info(
			"stack state not in sync",
			slog.String("stack", gs.StackName()),
		)
		return finishAndDestroy(ctx, p, tf)
	}

	slog.Info("getting kubeconfig from aws")
	if err := kubeconfigFromAWSCmd(
		ctx,
		p.AWSParams.Profile,
		p.ClusterParams.Name,
		p.AWSParams.Region,
		p.KubeconfigPath,
	); err != nil {
		return fmt.Errorf("kubeconfig from aws: %w", err)
	}

	k, err := NewClient(
		WithClientKubeconfig(p.KubeconfigPath),
		WithClientContext(p.ClusterParams.Name),
	)
	if err != nil {
		return fmt.Errorf("creating kubectl: %w", err)
	}

	clusterName := eks.EKSCluster.StateMust().Name
	clusterEndpoint := eks.EKSCluster.StateMust().Endpoint
	controllerIAMRoleArn := ks.Controller.Role.StateMust().Arn
	defaultInstanceProfile := ks.InstanceProfile.InstanceProfile.StateMust().Name
	interruptQueueName := ks.SimpleQueue.StateMust().Name
	kap := karpenter.New(
		karpenter.Opts{
			ClusterName:            clusterName,
			ClusterEndpoint:        clusterEndpoint,
			IAMRoleArn:             controllerIAMRoleArn,
			DefaultInstanceProfile: defaultInstanceProfile,
			InterruptQueue:         interruptQueueName,
		},
	)
	if err := k.Apply(ctx, kap); err != nil {
		return fmt.Errorf("applying karpenter app: %w", err)
	}
	// Wait for Karpenter to start before applying CRDs otherwise the webhooks fail
	objID := fmt.Sprintf(
		"%s/%s",
		kap.Deploy.TypeMeta.GetObjectKind().GroupVersionKind().GroupKind().String(),
		kap.Deploy.ObjectMeta.Name,
	)
	timeout := "5m"
	slog.Info(
		"waiting for karpenter deployment",
		slog.String("timeout", timeout),
	)
	if err := k.Cmd(
		ctx, "wait", "--namespace", kap.Deploy.Namespace, objID,
		"--for=condition=available",
		// Could take a while for the Fargate nodes to become available.
		// Usually it happens within 2 minutes, but just to be sure...
		"--timeout="+timeout,
	); err != nil {
		return fmt.Errorf("waiting for karpenter deployment: %w", err)
	}
	kapProvisioners := karpenter.NewProvisioners(
		karpenter.ProvisionersOpts{
			ClusterName:       clusterName,
			AvailabilityZones: vpcOpts.AZs,
		},
	)
	if err := k.Apply(ctx, &kapProvisioners); err != nil {
		return fmt.Errorf("applying karpenter provisioners app: %w", err)
	}

	db := gs.RDSPostgres.Postgres.StateMust()

	graf := grafana.New(
		grafana.AppOpts{
			Name:    grafana.AppName,
			Version: grafana.Version,
			Env:     "prod",
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

	kmNodeRoleARN := ks.InstanceProfile.IAMRole.StateMust().Arn
	kmFargateRoleARN := ks.FargateProfile.IAMRole.StateMust().Arn
	// Apply the aws-auth configmap
	awsAuth, err := awsauth.NewConfigMap(
		&awsauth.Data{
			MapRoles: karpenter.AWSAuthMapRoles(
				kmNodeRoleARN,
				kmFargateRoleARN,
			),
		},
	)
	if err != nil {
		return fmt.Errorf("creating aws-auth configmap: %w", err)
	}
	if err := k.Apply(
		ctx,
		awsAuth,
		// Required to become owner
		WithApplyForceConflicts(true),
	); err != nil {
		return fmt.Errorf("applying aws-auth: %w", err)
	}

	if p.Migrate {
		if err := migrateGrafana(ctx, p, tf, vpc, eks, gs); err != nil {
			slog.Error("migrate grafana", err)
			os.Exit(1)
		}
	}

	// This needs to come last, in case state is in sync but destroy flag was
	// passed
	if p.Destroy {
		return finishAndDestroy(ctx, p, tf)
	}

	fmt.Printf("\nTerriyaki Summary:\n")
	for _, mod := range tf.Stacks() {
		diff := "no plan"
		plan := mod.Plan()
		if plan != nil {
			diff = fmt.Sprintf(
				"add: %d, destroy: %d",
				len(plan.AddResources), len(plan.DestroyResources),
			)
		}
		fmt.Printf(
			"%s: resources: %s\n",
			mod.StackName(),
			diff,
		)
	}
	fmt.Println("")
	fmt.Println("")

	return nil
}

func kubectl(
	ctx context.Context,
	stdout io.Writer,
	stderr io.Writer,
	args ...string,
) error {
	cmd := exec.CommandContext(ctx, "kubectl", args...)
	cmd.Env = os.Environ() // inherit environment in case we need to use kubectl from a container

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	// waits for the command to exit and waits for any copying
	// to stdin or copying from stdout or stderr to complete
	return cmd.Wait()
}

func installCilium(
	ctx context.Context,
	co ClusterParams,
	kubeconfigPath string,
) error {
	cmd := exec.CommandContext(
		ctx, "cilium", "install",
		"--context", co.Name,
		"--cluster-name", co.Name,
		"--cluster-id", fmt.Sprintf("%d", co.ID),
		"--helm-set", "kubeProxyReplacement=strict",
		"--datapath-mode=aws-eni",
		"--version=v1.12.5",
		"--wait-duration=5m0s",
		"--wait",
	)
	path := os.Getenv("PATH")
	cmd.Env = append(cmd.Env, fmt.Sprintf("KUBECONFIG=%s", kubeconfigPath))
	cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=%s", path))

	fmt.Printf("%+v\n", cmd.Env)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}
