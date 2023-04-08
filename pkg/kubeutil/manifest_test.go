// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeutil

import (
	"bytes"
	"os"
	"testing"

	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestManifestSplit(t *testing.T) {
	golden, err := os.ReadFile("testdata/empty.yaml")
	tu.AssertNoError(t, err, "read golden file")
	got, err := ManifestSplit(bytes.NewReader(golden))
	tu.AssertNoError(t, err, "split manifest")
	want := []string{
		`apiVersion: v1
kind: Service
metadata:
  name: grafana
  namespace: monitoring
  labels:
    helm.sh/chart: grafana-6.50.7
    app.kubernetes.io/name: grafana
    app.kubernetes.io/instance: grafana
    app.kubernetes.io/version: "9.3.6"
    app.kubernetes.io/managed-by: Helm
spec:
  type: ClusterIP
  ports:
    - name: service
      port: 80
      protocol: TCP
      targetPort: 3000
  selector:
    app.kubernetes.io/name: grafana
    app.kubernetes.io/instance: grafana
`, `apiVersion: v1
data:
  _example: |
    ################################
    #                              #
    #    EXAMPLE CONFIGURATION     #
    #                              #
    ################################
    # This is an example config file highlighting the most common options.
    # this is particularly annoying as --- is kind of important in YAML.
    # ---------------------------------------
    # Settings Category
    # ---------------------------------------
    # some settings here
    bla: "48h"
    blabla: true
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/component: controller
    app.kubernetes.io/name: thename
    app.kubernetes.io/version: 1.8.0
  name: config---gc  # Oh you don't want that
  namespace: thenamespace
`,
	}
	tu.AssertEqualSlice(t, want, got)
}
