package kubeconfig_test

import (
	"testing"

	k "github.com/volvo-cars/go-terriyaki/pkg/kubeconfig"
	tu "github.com/volvo-cars/go-terriyaki/pkg/testutil"
)

func TestLoadKubeConfig(t *testing.T) {
	type TT struct {
		name string
		path string
		vf   []k.ValidationFunc
		want *k.Config
		err  string
	}

	tt := []TT{
		{
			name: "file does not exist",
			path: testdataDir + "/zero.yaml",
			err:  "load kubeconfig: file \"testdata/zero.yaml\" does not exist",
		},
		{
			name: "file is empty",
			path: testdataDir + "/empty.yaml",
			err:  "load kubeconfig: config \"testdata/empty.yaml\" invalid: config is empty",
		},
		{
			name: "file is empty and fails validation",
			path: testdataDir + "/empty.yaml",
			err:  "load kubeconfig: config \"testdata/empty.yaml\" invalid: config is empty",
		},
		{
			name: "file is not yaml",
			path: testdataDir + "/cluster1.json",
			err:  "load kubeconfig: file \"testdata/cluster1.json\" is not a yaml file",
		},
		{
			name: "file is not valid yaml",
			path: testdataDir + "/garbage.yaml",
			err:  "load kubeconfig: config \"testdata/garbage.yaml\" invalid: unmarshall json: json: cannot unmarshal string into Go value of type kubeconfig.Config", // nolint: lll
		},
		{
			name: "file is valid yaml but not kubeconfig",
			path: testdataDir + "/valid.yaml",
			err:  "load kubeconfig: config \"testdata/valid.yaml\" invalid: config is empty",
		},
		{
			name: "file is kubeconfig",
			path: testdataDir + "/cluster1.yaml",
			want: getCluster1(),
		},
		{
			name: "file is kubeconfig with valid contexts",
			path: testdataDir + "/cluster1.yaml",
			vf:   []k.ValidationFunc{k.WithValidContexts},
			want: getCluster1(),
		},
		{
			name: "file is kubeconfig with invalid contexts and fails validation",
			path: testdataDir + "/invalid-context.yaml",
			vf:   []k.ValidationFunc{k.WithValidContexts},
			err:  "load kubeconfig: config \"testdata/invalid-context.yaml\" invalid: context \"default\" references unknown user \"non-existent-user\"", // nolint: lll
		},
		{
			name: "file is kubeconfig with zero contexts",
			path: testdataDir + "/zero-context.yaml",
			vf:   []k.ValidationFunc{k.WithValidContexts},
			err:  "load kubeconfig: config \"testdata/zero-context.yaml\" invalid: no contexts defined",
		},
	}
	assert := func(t *testing.T, tt TT) {
		got, err := k.LoadKubeConfig(tt.path, tt.vf...)
		if tt.err != "" {
			if err != nil && err.Error() == tt.err {
				// t.Logf("got error: %v", err)
				return
			}
			t.Errorf("got error %q, want %q", err, tt.err)
			return
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

func getCluster1() *k.Config {
	c := k.New()
	_ = c.AddCluster(
		&k.ClusterConfig{
			Name: "default",
			Cluster: k.Cluster{
				CertificateAuthorityData: "LS0tLS1CRUdJTi",
				Server:                   "https://valid-cluster:6443",
			},
		},
	)
	_ = c.AddUser(
		&k.UserConfig{
			Name: "default",
			User: k.AuthInfo{
				ClientCertificateData: "LS0tLS1CRUdJTi",
				ClientKeyData:         "LS0tLS1CRUdJTi",
			},
		},
	)
	_ = c.AddContext(
		&k.ContextConfig{
			Name: "default",
			Context: k.Context{
				Cluster: "default",
				User:    "default",
			},
		},
	)

	return c
}
