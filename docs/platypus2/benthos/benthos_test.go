// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package benthos

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/rogpeppe/go-internal/txtar"
	"github.com/tidwall/gjson"
	"github.com/volvo-cars/lingon/pkg/kube"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
	corev1 "k8s.io/api/core/v1"

	"sigs.k8s.io/yaml"
)

func TestBenthos_LintConfig(t *testing.T) {
	outDir := filepath.Join("out", "config")
	tu.AssertNoError(t, os.RemoveAll("out"))
	testCases := []struct {
		name string
		in   BenthosArgs
	}{
		{
			name: "sane_defaults",
			in:   BenthosArgs{},
		},
		{
			name: "nats_cluster",
			in: BenthosArgs{
				Name:      "benthos-a",
				Namespace: "nats",
				Version:   "",
				Image:     "",
				Config:    "",
				Port:      Port{},
				Replicas:  0,
				EnvVar: &[]corev1.EnvVar{
					{Name: "NATS_URL", Value: "nats://nats.nats.svc.cluster.local:4222"},
				},
				Resource: corev1.ResourceRequirements{},
			},
		},
		{
			name: "nats_topic",
			in: BenthosArgs{
				Name:      "benthos-a",
				Namespace: "nats",
				Version:   "",
				Image:     "",
				Config:    "",
				Port:      Port{},
				Replicas:  0,
				EnvVar: &[]corev1.EnvVar{
					{Name: "NATS_URL", Value: "nats://nats.nats.svc.cluster.local:4222"},
					{Name: "NATS_TOPIC", Value: "demo.output"},
				},
				Resource: corev1.ResourceRequirements{},
			},
		},
	}

	tu.AssertNoError(t, os.MkdirAll(outDir, os.ModePerm))
	configFileOut := filepath.Join(outDir, "2_config_cm.yaml")

	for _, tc := range testCases {
		var buf bytes.Buffer
		var cm []byte

		b := New(tc.in)
		tu.AssertNoError(t, kube.Export(b,
			kube.WithExportWriter(&buf),
			kube.WithExportOutputDirectory(outDir)))

		ar := txtar.Parse(buf.Bytes())
		for _, f := range ar.Files {
			if f.Name == configFileOut {
				cm = f.Data
			}
		}
		cm, _ = yaml.YAMLToJSON(cm)
		s := gjson.GetBytes(cm, `data.benthos\.yaml`).String()
		cmfile := filepath.Join(outDir, "config_"+tc.name+".yaml")
		tu.AssertNoError(t, os.WriteFile(cmfile, []byte(s), os.ModePerm))

		// lint config /!\ require benthos to be in the path!
		o, err := exec.Command("benthos", "lint", cmfile).CombinedOutput()
		tu.AssertNoError(t, err, "output=", string(o), "name=", tc.name)
	}
}
