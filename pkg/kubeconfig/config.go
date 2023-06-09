// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeconfig

import (
	"fmt"
	"reflect"
)

const (
	version = "v1"
	kind    = "Config"
)

// New creates a new kubeconfig.
func New() *Config {
	return &Config{
		APIVersion: version,
		Kind:       kind,
		Clusters:   []*ClusterConfig{},
		Users:      []*UserConfig{},
		Contexts:   []*ContextConfig{},
	}
}

// Config is a kubeconfig.
type Config struct {
	Kind           string           `json:"kind"`
	APIVersion     string           `json:"apiVersion"`
	Preferences    Preferences      `json:"preferences"`
	CurrentContext string           `json:"current-context"`
	Clusters       []*ClusterConfig `json:"clusters"`
	Contexts       []*ContextConfig `json:"contexts"`
	Users          []*UserConfig    `json:"users"`
}

// UserConfig is a user in a kubeconfig.
type UserConfig struct {
	Name string   `json:"name"`
	User AuthInfo `json:"user"`
}

// ContextConfig is a context in a kubeconfig.
type ContextConfig struct {
	Name    string  `json:"name"`
	Context Context `json:"context"`
}

// ClusterConfig is a cluster in a kubeconfig.
type ClusterConfig struct {
	Name    string  `json:"name"`
	Cluster Cluster `json:"cluster"`
}

// AddCluster adds a cluster to the kubeconfig.
func (c *Config) AddCluster(cluster *ClusterConfig) error {
	if cluster == nil {
		return fmt.Errorf("add cluster: %w", errIsNil)
	}
	if reflect.ValueOf(cluster.Cluster).IsZero() {
		return fmt.Errorf("add cluster: %w", errIsEmpty)
	}

	c.Clusters = append(c.Clusters, cluster)
	return nil
}

// AddUser adds a user to the kubeconfig.
func (c *Config) AddUser(user *UserConfig) error {
	if user == nil {
		return fmt.Errorf("add user: %w", errIsNil)
	}
	if reflect.ValueOf(user.User).IsZero() {
		return fmt.Errorf("add user: %w", errIsEmpty)
	}
	c.Users = append(c.Users, user)
	return nil
}

// AddContext adds a context to the kubeconfig.
func (c *Config) AddContext(context *ContextConfig) error {
	if context == nil {
		return fmt.Errorf("add context: %w", errIsNil)
	}
	if reflect.ValueOf(context.Context).IsZero() {
		return fmt.Errorf("add context: %w", errIsEmpty)
	}

	c.Contexts = append(c.Contexts, context)
	return nil
}
