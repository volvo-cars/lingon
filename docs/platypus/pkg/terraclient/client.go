// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terraclient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/volvo-cars/lingon/pkg/terra"
)

const (
	tfSuffix   = ".tf"
	tfExec     = "terraform"
	tfPlanFile = "tfplan"
	tfWorkDir  = ".terra"
)

var (
	ErrNoStackName        = errors.New("no stack name")
	ErrDuplicateStackName = errors.New("duplicate stack name")
)

func NewClient(opts ...func(*clientOpts)) *Client {
	var rOpt clientOpts
	for _, opt := range opts {
		opt(&rOpt)
	}
	return &Client{
		stacks: make([]Stacker, 0),
		opts:   rOpt,
	}
}

func WithDefaultPlan(b bool) func(o *clientOpts) {
	return func(o *clientOpts) {
		o.plan = b
	}
}

func WithDefaultApply(b bool) func(o *clientOpts) {
	return func(o *clientOpts) {
		o.apply = b
	}
}

func WithDefaultDestroy(b bool) func(o *clientOpts) {
	return func(o *clientOpts) {
		o.destroy = b
	}
}

type clientOpts struct {
	plan    bool
	apply   bool
	destroy bool
}

// Client runs Terraform stacks and keeps a record of the runs in order to provide a summary of
// the changes
type Client struct {
	mu     sync.Mutex
	stacks []Stacker
	opts   clientOpts
}

func WithRunPlan(b bool) func(o *runOpts) {
	return func(o *runOpts) {
		o.plan = b
	}
}

func WithRunApply(b bool) func(o *runOpts) {
	return func(o *runOpts) {
		o.apply = b
	}
}

func WithRunDestroy(b bool) func(o *runOpts) {
	return func(o *runOpts) {
		o.destroy = b
	}
}

type runOpts struct {
	plan    bool
	apply   bool
	destroy bool
}

func (r *Client) Run(
	ctx context.Context,
	stack Stacker,
	opts ...func(*runOpts),
) error {
	rOpts := runOpts{
		plan:    r.opts.plan,
		apply:   r.opts.apply,
		destroy: r.opts.destroy,
	}

	for _, opt := range opts {
		opt(&rOpts)
	}

	if err := r.addStackToRunner(stack); err != nil {
		return err
	}
	if err := r.initStack(ctx, stack, rOpts); err != nil {
		return fmt.Errorf(
			"initializing stack %s: %w",
			stack.StackName(), err,
		)
	}
	var diff bool
	if rOpts.plan || rOpts.apply {
		var err error
		diff, err = r.planStack(ctx, stack, rOpts)
		if err != nil {
			return fmt.Errorf(
				"planning stack %s: %w", stack.StackName(), err,
			)
		}
	}
	if diff && rOpts.apply {
		if err := r.applyStack(ctx, stack); err != nil {
			return fmt.Errorf(
				"applying stack %s: %w", stack.StackName(), err,
			)
		}
	}
	if err := r.showStack(ctx, stack); err != nil {
		return fmt.Errorf(
			"getting state for stack %s: %w",
			stack.StackName(), err,
		)
	}

	return nil
}

func (r *Client) Stacks() []Stacker {
	return r.stacks
}

func (r *Client) planStack(
	ctx context.Context, stack Stacker,
	opts runOpts,
) (bool, error) {
	tf, err := r.newTerraform(stack)
	if err != nil {
		return false, fmt.Errorf("creating terraform runtime")
	}
	slog.Info(
		"Running Terraform Plan",
		slog.String("working_dir", tf.WorkingDir()),
		slog.String("out", tfPlanFile),
	)
	tf.SetStdout(os.Stdout)
	diff, err := tf.Plan(
		ctx,
		tfexec.Out(tfPlanFile),
		tfexec.Destroy(opts.destroy),
	)
	if err != nil {
		return false, err
	}
	tf.SetStdout(io.Discard)

	plan, err := tf.ShowPlanFile(ctx, tfPlanFile)
	if err != nil {
		return false, err
	}
	r.setPlan(stack, plan)

	return diff, nil
}

func (r *Client) applyStack(ctx context.Context, stack Stacker) error {
	tf, err := r.newTerraform(stack)
	if err != nil {
		return err
	}

	slog.Info(
		"Running Terraform Apply",
		slog.String("working_dir", tf.WorkingDir()),
		slog.String("plan", tfPlanFile),
	)
	tf.SetStdout(os.Stdout)
	if err := tf.Apply(ctx, tfexec.DirOrPlan(tfPlanFile)); err != nil {
		return fmt.Errorf("terraform apply command failed: %w", err)
	}
	tf.SetStdout(io.Discard)

	return nil
}

func (r *Client) showStack(ctx context.Context, stack Stacker) error {
	tf, err := r.newTerraform(stack)
	if err != nil {
		return err
	}

	slog.Info(
		"Importing Terraform state into stack",
		slog.String("stack", stack.StackName()),
		slog.String("working_dir", tf.WorkingDir()),
	)
	tfState, err := tf.Show(ctx)
	if err != nil {
		return fmt.Errorf("terraform show command failed: %w", err)
	}

	if len(tfState.Values.RootModule.Resources) == 0 {
		stack.SetStateMode(StateModeEmpty)
		return nil
	}

	fullState, err := terra.StackImportState(stack, tfState)
	if err != nil {
		return fmt.Errorf("importing state: %w", err)
	}
	if fullState {
		stack.SetStateMode(StateModeComplete)
	} else {
		stack.SetStateMode(StateModePartial)
	}

	return nil
}

func (r *Client) initStack(
	ctx context.Context,
	stack Stacker,
	opts runOpts,
) error {
	if stack.StackName() == "" {
		return ErrNoStackName
	}

	if err := terra.Export(
		stack, terra.WithExportOutputDirectory(r.workingDir(stack)),
	); err != nil {
		return fmt.Errorf("terra export: %w", err)
	}
	tf, err := r.newTerraform(stack)
	if err != nil {
		return fmt.Errorf("creating terraform runtime: %w", err)
	}

	if err = tf.Init(ctx, tfexec.Upgrade(true)); err != nil {
		return fmt.Errorf("terraform init command failed: %w", err)
	}

	tfv, err := tf.Validate(ctx)
	if err != nil {
		return fmt.Errorf("terraform validate command failed: %w", err)
	}
	if !tfv.Valid {
		return fmt.Errorf(
			"terraform stack is not valid: %+v",
			tfv.Diagnostics,
		)
	}

	return nil
}

// addStackToRunner appends the given stack to the runner's internal list
// This ensures there are not stacks with duplicate names being
// run and the list can be used by the caller, if needed, for things like
// destroying/infrastructure in reverse order
func (r *Client) addStackToRunner(stack Stacker) error {
	for _, stk := range r.stacks {
		if stk.StackName() == stack.StackName() {
			// If it's the same stack being run twice then it's ok
			if stk == stack {
				return nil
			}
			return fmt.Errorf(
				"%w: %s",
				ErrDuplicateStackName,
				stack.StackName(),
			)
		}
	}
	r.mu.Lock()
	r.stacks = append(
		r.stacks, stack,
	)
	r.mu.Unlock()
	return nil
}

func (r *Client) setPlan(stack Stacker, plan *tfjson.Plan) {
	stack.SetPlan(plan)
}

//	func NewTerraform(s Stacker) (*tfexec.Terraform, error) {
//		wd := filepath.Join(tfWorkDir, s.StackName())
//		tf, err := tfexec.NewTerraform(wd, tfExec)
//	}
func (r *Client) newTerraform(stack Stacker) (*tfexec.Terraform, error) {
	workingDir := r.workingDir(stack)
	tf, err := tfexec.NewTerraform(workingDir, tfExec)
	if err != nil {
		return nil, fmt.Errorf("creating terraform runtime: %w", err)
	}
	return tf, nil
}

func (r *Client) workingDir(stack Stacker) string {
	return filepath.Join(tfWorkDir, stack.StackName())
}
