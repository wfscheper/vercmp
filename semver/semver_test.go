// Copyright 2012 Red Hat, Inc.
// Copyright 2012-2013 Hewlett-Packard Development Company, L.P.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		v        string
		expected *SemanticVersion
	}{
		{"1.2.3.dev6", &SemanticVersion{
			Major:    1,
			Minor:    2,
			Patch:    3,
			DevCount: 6,
		}},
		{"1.2.3.dev7", &SemanticVersion{
			Major:    1,
			Minor:    2,
			Patch:    3,
			DevCount: 7,
		}},
		{"1.2.3.a4.dev12", &SemanticVersion{
			Major:          1,
			Minor:          2,
			Patch:          3,
			PreRelease:     4,
			PreReleaseType: "a",
			DevCount:       12,
		}},
		{"1.2.3.a4.dev13", &SemanticVersion{
			Major:          1,
			Minor:          2,
			Patch:          3,
			PreRelease:     4,
			PreReleaseType: "a",
			DevCount:       13,
		}},
		{"1.2.3.a4", &SemanticVersion{
			Major:          1,
			Minor:          2,
			Patch:          3,
			PreRelease:     4,
			PreReleaseType: "a",
		}},
		{"1.2.3.a5.dev1", &SemanticVersion{
			Major:          1,
			Minor:          2,
			Patch:          3,
			PreRelease:     5,
			PreReleaseType: "a",
			DevCount:       1,
		}},
		{"1.2.3.a5", &SemanticVersion{
			Major:          1,
			Minor:          2,
			Patch:          3,
			PreRelease:     5,
			PreReleaseType: "a",
		}},
		{"1.2.3.b3.dev1", &SemanticVersion{
			Major:          1,
			Minor:          2,
			Patch:          3,
			PreRelease:     3,
			PreReleaseType: "b",
			DevCount:       1,
		}},
		{"1.2.3.b3", &SemanticVersion{
			Major:          1,
			Minor:          2,
			Patch:          3,
			PreRelease:     3,
			PreReleaseType: "b",
		}},
		{"1.2.3.rc2.dev1", &SemanticVersion{
			Major:          1,
			Minor:          2,
			Patch:          3,
			PreRelease:     2,
			PreReleaseType: "rc",
			DevCount:       1,
		}},
		{"1.2.3.rc2", &SemanticVersion{
			Major:          1,
			Minor:          2,
			Patch:          3,
			PreRelease:     2,
			PreReleaseType: "rc",
		}},
		{"1.2.3.rc3.dev1", &SemanticVersion{
			Major:          1,
			Minor:          2,
			Patch:          3,
			PreRelease:     3,
			PreReleaseType: "rc",
			DevCount:       1,
		}},
		{"1.2.3", &SemanticVersion{
			Major: 1,
			Minor: 2,
			Patch: 3,
		}},
		{"1.2.4", &SemanticVersion{
			Major: 1,
			Minor: 2,
			Patch: 4,
		}},
		{"1.3.3", &SemanticVersion{
			Major: 1,
			Minor: 3,
			Patch: 3,
		}},
		{"2.2.3", &SemanticVersion{
			Major: 2,
			Minor: 2,
			Patch: 3,
		}},
	}

	for _, test := range tests {
		actual, err := New(test.v)
		if assert.Nil(t, err) {
			assert.Equal(t, test.expected, actual, "Expected New(%v) == %v, got %v",
				test.v, test.expected, actual)
		}
	}
}

func TestVersionOrdering(t *testing.T) {
	versions := []string{
		"1.2.3.dev6",
		"1.2.3.dev7",
		"1.2.3.a4.dev12",
		"1.2.3.a4.dev13",
		"1.2.3.a4",
		"1.2.3.a5.dev1",
		"1.2.3.a5",
		"1.2.3.b3.dev1",
		"1.2.3.b3",
		"1.2.3.rc2.dev1",
		"1.2.3.rc2",
		"1.2.3.rc3.dev1",
		"1.2.3",
		"1.2.4",
		"1.3.3",
		"2.2.3",
	}

	for _, v := range versions {
		assert.True(t, assertVersionEqual(v, v), "Expected %s == %s", v, v)
	}

	for _, pairs := range combinations(versions, 2) {
		left, right := pairs[0], pairs[1]
		l_pos := index(versions, left)
		r_pos := index(versions, right)
		if l_pos < r_pos {
			assert.True(t, assertVersionOrder(left, right),
				"Expected %v < %v", left, right)
		} else {
			assert.True(t, assertVersionOrder(right, left),
				"Expected %v < %v", right, left)
		}
	}
}

func BenchmarkVercmp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Vercmp("1.2.3.a5.dev6", "1.2.3.a5.dev7")
	}
}

func BenchmarkVercmpSemanticVersion(b *testing.B) {
	v1, _ := New("1.2.3.a5.dev6")
	v2, _ := New("1.2.3.a5.dev7")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Vercmp(v1, v2)
	}
}

func assertVersionEqual(v1, v2 interface{}) bool {
	if Vercmp(v1, v2) != 0 {
		return false
	}
	if Vercmp(v2, v1) != 0 {
		return false
	}
	return true
}

func assertVersionOrder(low, high interface{}) bool {
	if Vercmp(low, high) >= 0 {
		return false
	}
	if Vercmp(high, low) <= 0 {
		return false
	}
	return true
}

func index(s []string, v string) int {
	for idx, e := range s {
		if e == v {
			return idx
		}
	}
	return -1
}

func combinations(s []string, r int) [][]string {
	rv := make([][]string, 0)
	n := len(s)
	if r > n {
		return rv
	}

	indices := make([]int, r)
	for i, _ := range indices {
		indices[i] = i
	}

	c := make([]string, r)
	for _, i := range indices {
		c[i] = s[i]
	}
	rv = append(rv, c)

	var idx int
	for {
		for i := r - 1; i >= -1; i-- {
			idx = i
			if i == -1 {
				return rv
			} else if indices[i] != i+n-r {
				break
			}
		}
		indices[idx] = indices[idx] + 1
		for j := idx + 1; j < r; j++ {
			indices[j] = indices[j-1] + 1
		}
		c = make([]string, r)
		for i, p := range indices {
			c[i] = s[p]
		}
		rv = append(rv, c)
	}
	return rv
}
