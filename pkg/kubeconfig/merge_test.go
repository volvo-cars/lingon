// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeconfig_test

import (
	"path/filepath"
	"testing"

	k "github.com/volvo-cars/lingon/pkg/kubeconfig"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

const (
	testdataDir = "testdata"
)

func loadTestdata(t *testing.T, name string) *k.Config {
	t.Helper()
	config, err := k.LoadKubeConfig(filepath.Join(testdataDir, name))
	if err != nil {
		t.Fatal(err)
		return nil
	}
	return config
}

func makeInvalidConfig(user, cluster, context bool) *k.Config {
	c := k.New()
	c.APIVersion = "invalid"
	cc := &k.ClusterConfig{
		Name: "cluster",
		Cluster: k.Cluster{
			Server: "https://",
		},
	}
	u := &k.UserConfig{
		Name: "user",
		User: k.AuthInfo{
			Token: "",
		},
	}
	ctx := &k.ContextConfig{
		Name:    "context",
		Context: k.Context{Cluster: "cluster", User: "user"},
	}

	_ = c.AddCluster(nil)
	_ = c.AddCluster(&k.ClusterConfig{})
	if cluster {
		_ = c.AddCluster(cc)
	}

	_ = c.AddUser(nil)
	_ = c.AddUser(&k.UserConfig{})
	if user {
		_ = c.AddUser(u)
	}

	_ = c.AddContext(&k.ContextConfig{})
	_ = c.AddContext(nil)
	if context {
		_ = c.AddContext(ctx)
	}

	return c
}

func TestMerge(t *testing.T) {
	type TT struct {
		name string
		in   []*k.Config
		want *k.Config
		err  string
	}

	cluster1 := loadTestdata(t, "cluster1.yaml")
	cluster2 := loadTestdata(t, "cluster2.yaml")
	clusterMerged := loadTestdata(t, "cluster-merged.yaml")
	clusterMergedCopy := loadTestdata(t, "cluster-merged.yaml")
	cluster3 := loadTestdata(t, "cluster3.yaml")
	clusterMerged3 := loadTestdata(t, "cluster-merged3.yaml")

	tt := []TT{
		{
			name: "zero config",
			in:   []*k.Config{},
			err:  "no config to merge",
		},
		{
			name: "one empty config",
			in:   []*k.Config{k.New()},
			err:  "merge: config is empty",
		},
		{
			name: "one invalid config",
			in:   []*k.Config{makeInvalidConfig(false, true, false)},
			err:  "merge: no contexts defined + no users defined",
		},
		{
			name: "one valid config",
			in: []*k.Config{
				clusterMerged,
			},
			want: clusterMergedCopy,
		},
		{
			name: "two valid config",
			in: []*k.Config{
				cluster1,
				cluster2,
			},
			want: clusterMergedCopy,
		},
		{
			name: "two valid config & one invalid config",
			in: []*k.Config{
				cluster1,
				makeInvalidConfig(false, false, true),
				cluster2,
			},
			err: "merge: no users defined + no clusters defined",
		},
		{
			name: "two valid config & one nil config",
			in: []*k.Config{
				cluster1,
				nil,
				cluster2,
			},
			err: "merge: config is nil",
		},
		{
			name: "three valid config",
			in: []*k.Config{
				cluster1,
				cluster2,
				cluster3,
			},
			want: clusterMerged3,
		},
	}

	assert := func(t *testing.T, tt TT) {
		got, err := k.Merge(tt.in...)
		if err != nil {
			if err.Error() == tt.err {
				return
			}
			t.Errorf("got error %q, want %q", err, tt.err)
		}
		if diff := tu.Diff(got, tt.want); diff != "" {
			t.Error(tu.Callers(), diff)
		}
	}

	for _, tc := range tt {
		t.Run(
			tc.name, func(t *testing.T) {
				assert(t, tc)
			},
		)
	}
}
