// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

// Package testutils
// Example:
//
//	func TestExample(t *testing.T) {
//		type TT struct {
//			got  string
//			want string
//		}
//
//		assert := func(t *testing.T, tt TT) {
//			if diff := Diff(tt.got, tt.want); diff != "" {
//				t.Error(Callers(), diff)
//			}
//		}
//
//		t.Run(
//			"1", func(t *testing.T) {
//				t.Parallel()
//				assert(t, TT{"lorem ipsum dolor amet", "lorem ipsum dolor sit amet"})
//			},
//		)
//
//		t.Run(
//			"2", func(t *testing.T) {
//				t.Parallel()
//				assert(t, TT{"the quick fox jumped over lazy dog", "the quick brown fox jumped over the lazy dog"})
//			},
//		)
//
//		t.Run(
//			"2", func(t *testing.T) {
//				t.Parallel()
//				assert(t, TT{"Sphinx of black quartz judge my vow", "Sphinx of black quartz, judge my vow"})
//			},
//		)
//	}
package testutil
