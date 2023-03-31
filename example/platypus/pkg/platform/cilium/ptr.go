// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package cilium

// P returns a pointer to the given value.
func P[T any](t T) *T {
	return &t
}
