// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package kubeconfig

import (
	"encoding/json"
	"fmt"

	"sigs.k8s.io/yaml"
)

func (c *Config) Unmarshal(b []byte) error {
	b, err := yaml.YAMLToJSON(b)
	if err != nil {
		return fmt.Errorf("unmarshall: convert yaml to json: %w", err)
	}
	err = json.Unmarshal(b, c)
	if err != nil {
		return fmt.Errorf("unmarshall json: %w", err)
	}
	return nil
}

func (c *Config) Marshal() ([]byte, error) {
	j, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("marshall: %w", err)
	}
	y, err := yaml.JSONToYAML(j)
	if err != nil {
		return nil, err
	}
	return y, nil
}
