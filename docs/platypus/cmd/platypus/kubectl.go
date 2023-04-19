// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/volvo-cars/lingon/pkg/kube"
	"golang.org/x/exp/slog"
)

func Kubectl(
	ctx context.Context,
	out io.Writer,
	errw io.Writer,
	args ...string,
) error {
	slog.Info("kubectl", slog.Any("args", args))
	cmd := exec.CommandContext(ctx, "kubectl", args...)

	cmd.Env = os.Environ()

	cmd.Stdin = os.Stdin
	cmd.Stdout = out
	cmd.Stderr = errw

	err := cmd.Start()
	if err != nil {
		return err
	}
	// waits for the command to exit and waits for any copying
	// to stdin or copying from stdout or stderr to complete
	return cmd.Wait()
}

func WithClientKubeconfig(kubeconfig string) func(o *clientOpts) {
	return func(o *clientOpts) {
		o.kubeconfig = kubeconfig
	}
}

func WithClientContext(context string) func(o *clientOpts) {
	return func(o *clientOpts) {
		o.context = context
	}
}

type clientOpts struct {
	kubeconfig string
	context    string
}

func NewClient(optParams ...func(o *clientOpts)) (*Client, error) {
	opts := clientOpts{}
	for _, opt := range optParams {
		opt(&opts)
	}
	// TODO: do we want to get kubeconfig/context from env vars?
	if opts.kubeconfig == "" {
		return nil, fmt.Errorf("kubeconfig required")
	}
	if opts.context == "" {
		return nil, fmt.Errorf("context required")
	}
	// TODO: can we validate the context and the kubeconfig?

	return &Client{
		opts: opts,
	}, nil
}

// Client rerepsents a kubectl client setup to communicate with a single Kubernetes cluster
type Client struct {
	opts clientOpts
}

func WithApplyForceConflicts(b bool) func(o *applyOpts) {
	return func(o *applyOpts) {
		o.forceConflicts = b
	}
}

type applyOpts struct {
	forceConflicts bool
}

// Apply performs a kubectl apply for the given manifest
func (k *Client) Apply(
	ctx context.Context,
	km kube.Exporter,
	opts ...func(o *applyOpts),
) error {
	ao := applyOpts{}
	for _, opt := range opts {
		opt(&ao)
	}

	args := k.baseArgs()
	args = append(
		args,
		"apply",
		"--server-side=true",
		"-f",
		"-",
	)

	if ao.forceConflicts {
		args = append(args, "--force-conflicts")
	}

	cmd := exec.CommandContext(
		ctx,
		"kubectl",
		args...,
	)
	cmd.Env = os.Environ()        // inherit environment in case we need to use kubectl from a container
	stdin, err := cmd.StdinPipe() // pipe to pass data to kubectl
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		var buf bytes.Buffer
		if err := kube.Export(
			km,
			kube.WithExportWriter(&buf),
			kube.WithExportAsSingleFile("karpenter.yaml"),
		); err != nil {
			log.Fatal("export", err)
		}
		log.Printf("kubectl apply: %s", buf.String())
		if _, err := io.Copy(stdin, &buf); err != nil {
			log.Fatal("copy", err)
		}
	}()

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}
	// waits for the command to exit and waits for any copying
	// to stdin or copying from stdout or stderr to complete
	return cmd.Wait()
}

// Diff performs a kubectl diff to see if the manifests have changed
func (k *Client) Diff(ctx context.Context, km kube.Exporter) error {
	args := k.baseArgs()
	args = append(
		args,
		"diff",
		"--server-side=true",
		"-f",
		"-",
	)
	cmd := exec.CommandContext(
		ctx,
		"kubectl",
		args...,
	)
	cmd.Env = os.Environ()        // inherit environment in case we need to use kubectl from a container
	stdin, err := cmd.StdinPipe() // pipe to pass data to kubectl
	if err != nil {
		log.Fatal(err)
	}

	if err := kube.Export(
		km,
		kube.WithExportOutputDirectory("karpenter"),
	); err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		if err := kube.Export(km,
			kube.WithExportWriter(stdin),
			kube.WithExportAsSingleFile("stdin"),
		); err != nil {
			log.Fatal(err)
		}
	}()

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}
	// waits for the command to exit and waits for any copying
	// to stdin or copying from stdout or stderr to complete
	return cmd.Wait()
}

func (k *Client) Cmd(
	ctx context.Context,
	args ...string,
) error {
	baseArgs := []string{
		"--kubeconfig", k.opts.kubeconfig, "--context", k.opts.context,
	}
	cmd := exec.CommandContext(
		ctx,
		"kubectl",
		append(baseArgs, args...)...,
	)
	cmd.Env = os.Environ() // inherit environment in case we need to use kubectl from a container

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}
	// waits for the command to exit and waits for any copying
	// to stdin or copying from stdout or stderr to complete
	return cmd.Wait()
}

func (k *Client) baseArgs() []string {
	return []string{
		"--kubeconfig",
		k.opts.kubeconfig,
		"--context",
		k.opts.context,
	}
}
