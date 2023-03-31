// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package main

//go:generate go run -mod=readonly github.com/volvo-cars/lingon/cmd/terragen -out ../gen/providers/aws -pkg github.com/volvo-cars/lingon/example/platypus/gen/providers/aws -force -provider aws=hashicorp/aws:4.49.0
//go:generate go run -mod=readonly github.com/volvo-cars/lingon/cmd/terragen -out ../gen/providers/tls -pkg github.com/volvo-cars/lingon/example/platypus/gen/providers/tls -force -provider tls=hashicorp/tls:4.0.4
