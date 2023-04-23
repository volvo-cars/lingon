// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package testutil_test

import (
	"bytes"
	"fmt"
	"os"

	"github.com/volvo-cars/lingon/pkg/testutil"
)

func ExampleScore() {
	f, err := os.Open("testdata/deployment.yaml")
	if err != nil {
		panic(err)
	}
	card, err := testutil.Score(f)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = testutil.RenderScoreCard2(card, &buf, false)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", buf.String())

	// Output:
	//
	// apps/v1/Deployment webapp
	//     [CRITICAL] Container Ephemeral Storage Request and Limit
	//         · nginx -> Ephemeral Storage limit is not set
	//             Resource limits are recommended to avoid resource DDOS. Set resources.limits.ephemeral-storage
	//     [CRITICAL] Container Image Tag
	//         · nginx -> Image with latest tag
	//             Using a fixed tag is recommended to avoid accidental upgrades
	//     [CRITICAL] Container Security Context ReadOnlyRootFilesystem
	//         · nginx -> Container has no configured security context
	//             Set securityContext to run the container in a more secure context.
	//     [CRITICAL] Container Security Context User Group ID
	//         · nginx -> Container has no configured security context
	//             Set securityContext to run the container in a more secure context.
	//     [CRITICAL] Pod NetworkPolicy
	//         · The pod does not have a matching NetworkPolicy
	//             Create a NetworkPolicy that targets this pod to control who/what can communicate with this pod.
	//             Note, this feature needs to be supported by the CNI implementation used in the Kubernetes cluster to
	//             have an effect.
}
