// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package main

//go:generate terragen -out providers/aws/4.66.1 -pkg github.com/volvo-cars/lingoneks/providers/aws/5.0.1 -provider aws=hashicorp/aws:4.66.1
//go:generate terragen -out providers/tls/4.0.4 -pkg github.com/volvo-cars/lingoneks/providers/tls/4.0.4 -provider tls=hashicorp/tls:4.0.4
