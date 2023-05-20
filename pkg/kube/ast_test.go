// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/volvo-cars/lingon/pkg/kubeutil"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
	rbacv1 "k8s.io/api/rbac/v1"
)

func TestKube2GoJen(t *testing.T) {
	type TT struct {
		name     string
		manifest string
		golden   string
		redact   bool
	}
	tests := []TT{
		{
			name:     "deployment",
			manifest: "testdata/golden/deployment.yaml",
			golden:   "testdata/golden/deployment.golden",
		},
		{
			name:     "service",
			manifest: "testdata/golden/service.yaml",
			golden:   "testdata/golden/service.golden",
		},
		{
			name:     "secret",
			manifest: "testdata/golden/secret.yaml",
			golden:   "testdata/golden/secret.golden",
			redact:   true,
		},
		{
			name:     "empty configmap",
			manifest: "testdata/golden/configmap.yaml",
			golden:   "testdata/golden/configmap.golden",
		},
		{
			name:     "configmap with comments",
			manifest: "testdata/golden/cm-comment.yaml",
			golden:   "testdata/golden/cm-comment.golden",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := convert(t, tt.manifest, tt.redact)
				want := tu.ReadGolden(t, tt.golden)
				tu.AssertEqual(t, want, got)
			},
		)
	}
}

func convert(t *testing.T, path string, redact bool) string {
	t.Helper()
	data, err := os.ReadFile(path)
	tu.AssertNoError(t, err, "read manifest")
	j := jamel{
		o: importOption{
			Serializer:    defaultSerializer(),
			RedactSecrets: redact,
		},
	}
	m, err := kubeutil.ExtractMetadata(data)
	tu.AssertNoError(t, err, "extract metadata")
	code, err := j.yaml2GoJen(data, m)
	tu.AssertNoError(t, err, "convert yaml to go")
	var b strings.Builder
	err = code.Render(&b)
	tu.AssertNoError(t, err, "render code")
	return b.String()
}

func TestConvertValue(t *testing.T) {
	type TT struct {
		name string
		in   interface{}
		want string
	}
	type sarray struct {
		Names          []string
		Bin            []byte
		UINT           []uint8
		unexported     string
		alsoUnexported struct {
			Nope bool
		}
	}
	tests := []TT{
		{
			name: "string array",
			in:   []string{"bla", "oops", "yeah"},
			want: "[]string{\"bla\", \"oops\", \"yeah\"}",
		},
		{
			name: "bytes array",
			in:   []byte("ok"),
			want: "[]byte(\"ok\")",
		},
		{
			name: "string array empty",
			in:   []string{},
			want: `[]string{}`,
		},
		{
			name: "string array empty string",
			in:   []string{""},
			want: `[]string{""}`,
		},
		{
			name: "string array array",
			in:   [][]string{{"test"}, {"bla"}},
			want: "[][]string{[]string{\"test\"}, []string{\"bla\"}}",
		},
		{
			name: "bytes array array",
			in:   [][]byte{[]byte("ok"), []byte("yay")},
			want: `[][]byte{[]byte("ok"), []byte("yay")}`,
		},
		{
			name: "array of struct",
			in:   []sarray{{Names: []string{"s1"}}, {Names: []string{"s2"}}},
			want: `[]kube.sarray{kube.sarray{Names: []string{"s1"}}, kube.sarray{Names: []string{"s2"}}}`,
		},
		{
			name: "array of array of struct",
			in: [][]sarray{
				{
					{Names: []string{"s1"}},
					{Names: []string{"s2", "s3"}},
				},
			},
			want: `[][]kube.sarray{[]kube.sarray{kube.sarray{Names: []string{"s1"}}, kube.sarray{Names: []string{"s2", "s3"}}}}`,
		},
		{
			name: "array of struct of array of byte",
			in: []sarray{
				{
					Bin: []byte("bibin"),
				},
			},
			want: `[]kube.sarray{kube.sarray{Bin: []byte("bibin")}}`,
		},
		{
			name: "array of struct of array of uint8",
			in: []sarray{
				{
					UINT: []uint8{
						uint8(0x62),
						uint8(0x69),
						uint8(0x62),
						uint8(0x69),
						uint8(0x6e),
					}, // "bibin"
				},
			},
			want: `[]kube.sarray{kube.sarray{UINT: []byte("bibin")}}`,
		},
		{
			name: "map of array",
			in:   map[string][]string{"key1": {"ok", "yay"}},
			want: `map[string][]string{"key1": []string{"ok", "yay"}}`,
		},
		{
			name: "map of map string",
			in: map[string]map[string]string{
				"key1": {"key2": "value1"},
			},
			want: `map[string]map[string]string{"key1": map[string]string{"key2": "value1"}}`,
		},
		{
			name: "map of struct of array",
			in: map[string]sarray{
				"key": {
					Names: []string{
						"one",
						"two",
					},
					unexported:     "not exported",
					alsoUnexported: struct{ Nope bool }{Nope: true},
				},
			},
			want: `map[string]kube.sarray{"key": kube.sarray{Names: []string{"one", "two"}}}`,
		},
		{
			name: "map of pointers to struct",
			in:   map[string]*sarray{"key1": {Names: []string{"hi"}}},
			want: `map[string]*kube.sarray{"key1": &kube.sarray{Names: []string{"hi"}}}`,
		},
		{
			name: "map of interface",
			in:   map[string]interface{}{},
			want: `map[string]interface{}{}`,
		},
		{
			name: "array of pointers to struct",
			in:   []*sarray{{Names: []string{"hi"}}},
			want: `[]*kube.sarray{&kube.sarray{Names: []string{"hi"}}}`,
		},
		{
			name: "array of bool - ignore zero value",
			in:   []bool{true, false, true, false},
			want: `[]bool{true, true}`,
		},
		{
			name: "array of interface",
			in:   []interface{}{},
			want: `[]interface{}{}`,
		},
		// {
		// 	name: "nil",
		// 	in:   nil,
		// 	want: "nil",
		// },
		{
			name: "string",
			in:   "test",
			want: `"test"`,
		},
		{
			name: "function - unsupported",
			in:   func() {},
			want: "nil",
		},
		{
			name: "complex - unsupported",
			in:   complex(10, 11),
			want: "nil",
		},
		{
			name: "int",
			in:   42,
			want: "42",
		},
		{
			name: "int16",
			in:   int16(10),
			want: "int16(10)",
		},
		{
			name: "int8",
			in:   int8(64),
			want: "int8(64)",
		},
		{
			name: "uint - converted to uint64",
			in:   uint(10),
			want: "uint64(0xa)",
		},
		{
			name: "uint64",
			in:   uint64(10),
			want: "uint64(0xa)",
		},
		{
			name: "uint32",
			in:   uint32(10),
			want: "uint32(0xa)",
		},
		{
			name: "uint16",
			in:   uint16(10),
			want: "uint16(0xa)",
		},
		{
			name: "uint8",
			in:   uint8(42),
			want: "uint8(0x2a)",
		},
		{
			name: "float32",
			in:   float32(10),
			want: "float32(10)",
		},
		{
			name: "float64",
			in:   float64(10),
			want: "10.0",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				j := jamel{
					crdPkgAlias: make(map[string]string, 0),
					o:           importOption{},
				}
				got := j.convertValue(reflect.ValueOf(tt.in))
				tu.AssertEqual(t, tt.want, fmt.Sprintf("%#v", got))
			},
		)
	}
}

func Test_convertPolicyRule(t *testing.T) {
	tests := []struct {
		name string
		in   rbacv1.PolicyRule
		want string
	}{
		{
			name: "list configmaps",
			in: rbacv1.PolicyRule{
				APIGroups: []string{""},
				Resources: []string{"configmaps"},
				Verbs:     []string{"list", "watch"},
			},
			want: `v1.PolicyRule{
	APIGroups: []string{""},
	Resources: []string{"configmaps"},
	Verbs:     []string{"list", "watch"},
}`,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				j := jamel{
					crdPkgAlias: make(map[string]string, 0),
					o:           importOption{},
				}
				got := j.convertValue(reflect.ValueOf(tt.in))
				tu.AssertEqual(t, tt.want, fmt.Sprintf("%#v", got))
			},
		)
	}
}

func Test_rawString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "simple",
			s:    "simple",
			want: `"simple"`,
		},
		{
			name: "empty",
			s:    "",
			want: `""`,
		},
		{
			name: "with double quote",
			s:    `this "double quote" in a string`,
			want: `"this \"double quote\" in a string"`,
		},
		{
			name: "with new line",
			s:    "this line \n and this line",
			want: "`" + `
this line 
 and this line
` + "`",
		},
		{
			name: "with backtick",
			s:    "this ` is a backtick",
			want: "\"this ` is a backtick\"",
		},
		{
			name: "backticks with new line",
			s: `fun 
stuff` + "`\"with backticks`\" and new lines",
			want: "`\nfun \nstuff\"\"with backticks\"\" and new lines\n`",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := rawString(tt.s)
				tu.AssertEqual(t, tt.want, fmt.Sprintf("%#v", got))
			},
		)
	}
}
