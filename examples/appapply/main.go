package main

//go:generate go run github.com/volvo-cars/go-terriyaki/cmd/kygo -in ../../pkg/kube/testdata/grafana.yaml -out out/grafana -app grafana -group=false -clean-name=false
//go:generate go run github.com/volvo-cars/go-terriyaki/cmd/explode -in ../../pkg/kube/testdata/grafana.yaml -out out/explode

import (
	"appapply/out/grafana"
	"fmt"

	"github.com/volvo-cars/go-terriyaki/pkg/kube"
)

func main() {
	g := grafana.New()
	if err := kube.Export(g, "export"); err != nil {
		fmt.Println(err)
		return
	}
}
