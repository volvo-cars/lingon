// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"

	"github.com/volvo-cars/lingoneks/pkg/terraclient"
	"golang.org/x/exp/slog"
)

func finishAndDestroy(
	ctx context.Context, p runParams,
	runner *terraclient.Client,
) error {
	if !p.Destroy {
		return nil
	}
	stacks := runner.Stacks()
	// Iterate in reverse
	for i := len(stacks) - 1; i >= 0; i-- {
		stack := stacks[i]
		if err := runner.Run(
			ctx, stack,
			terraclient.WithRunDestroy(p.Destroy),
			terraclient.WithRunPlan(true),
			terraclient.WithRunApply(true),
		); err != nil {
			return fmt.Errorf(
				"destroying %s: %w",
				stack.StackName(), err,
			)
		}
	}
	slog.Info("EVERYTHING DESTROYED!!")
	return nil
}
