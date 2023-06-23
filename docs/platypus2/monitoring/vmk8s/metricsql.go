// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package vmk8s

import (
	"github.com/VictoriaMetrics/metricsql"
)

func MetricMust(s string) string {
	_, err := metricsql.Parse(s)
	if err != nil {
		panic("invalid query: " + s + "\t" + err.Error())
	}
	// fmt.Printf("%s\n", valast.String(expr))

	return s
}

// func MetricsConvertMQLToPromQL(mql string) string {
// 	// Convert mql to PromQL
// 	pql, err := metricsql.ExpandWithExprs(mql)
// 	if err != nil {
// 		panic("cannot expand with expressions: " + err.Error())
// 	}
// 	return pql
// }
