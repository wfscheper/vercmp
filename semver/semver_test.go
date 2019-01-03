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
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		v    string
		want *SemanticVersion
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

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.v, func(t *testing.T) {
			got, err := New(tt.v)
			if err != nil {
				t.Errorf("got %v, want nil", err)
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

var versionEqualityTests = []string{
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

func TestVersionEquality(t *testing.T) {
	t.Parallel()
	for _, v := range versionEqualityTests {
		t.Run(v+" == "+v, func(t *testing.T) {
			if !assertVersionEqual(v, v) {
				t.Error("got false, want true")
			}
		})
	}
}

func TestVersionOrdering(t *testing.T) {
	t.Parallel()
	for _, pairs := range combinations(versionEqualityTests, 2) {
		left, right := pairs[0], pairs[1]
		t.Run(left+" < "+right, func(t *testing.T) {
			l_pos, r_pos := index(versionEqualityTests, left), index(versionEqualityTests, right)
			if l_pos < r_pos {
				if !assertVersionOrder(left, right) {
					t.Errorf("Expected %v < %v", left, right)
				}
			} else {
				if !assertVersionOrder(right, left) {
					t.Errorf("Expected %v < %v", right, left)
				}
			}
		})
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
}
