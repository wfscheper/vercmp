package maven

import (
	"fmt"
	"reflect"
	"testing"
)

func TestZip(t *testing.T) {
	tests := []struct {
		title string
		a, b  []interface{}
		want  []interfaceTuple
	}{
		{
			"One element each",
			[]interface{}{0},
			[]interface{}{1},
			[]interfaceTuple{interfaceTuple{0, 1}},
		},
		{
			"Two elements each",
			[]interface{}{0, "0"},
			[]interface{}{"1", 1},
			[]interfaceTuple{interfaceTuple{0, "1"}, interfaceTuple{"0", 1}},
		},
		{
			"One and two elements",
			[]interface{}{0},
			[]interface{}{"1", 1},
			[]interfaceTuple{interfaceTuple{0, "1"}, interfaceTuple{nil, 1}},
		},
		{
			"Two and one elements",
			[]interface{}{0, "0"},
			[]interface{}{"1"},
			[]interfaceTuple{interfaceTuple{0, "1"}, interfaceTuple{"0", nil}},
		},
		{
			"Empty slices",
			[]interface{}{},
			[]interface{}{},
			[]interfaceTuple{},
		},
		{
			"Empty slice a",
			[]interface{}{},
			[]interface{}{0, 1},
			[]interfaceTuple{interfaceTuple{nil, 0}, interfaceTuple{nil, 1}},
		},
		{
			"Empty slice b",
			[]interface{}{0, 1},
			[]interface{}{},
			[]interfaceTuple{interfaceTuple{0, nil}, interfaceTuple{1, nil}},
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			got := zip(tt.a, tt.b)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseBuffer(t *testing.T) {
	tests := []struct {
		b            string
		digitFollows bool
		want         interface{}
	}{
		{"1", false, 1},
		{"10", false, 10},
		{"1", true, 1},
		{"10", true, 10},
		{"a", false, "a"},
		{"b", false, "b"},
		{"m", false, "m"},
		{"a", true, "alpha"},
		{"b", true, "beta"},
		{"m", true, "milestone"},
		{"ga", false, ""},
		{"final", false, ""},
		{"rc", false, "rc"},
		{"cr", false, "rc"},
		{"ga", true, ""},
		{"final", true, ""},
		{"rc", true, "rc"},
		{"cr", true, "rc"},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(fmt.Sprintf("parseBuffer(%v, %v)", tt.b, tt.digitFollows), func(t *testing.T) {
			got := parseBuffer(tt.b, tt.digitFollows)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		s    []interface{}
		want []interface{}
	}{
		{
			[]interface{}{0},
			[]interface{}{},
		},
		{
			[]interface{}{""},
			[]interface{}{},
		},
		{
			[]interface{}{&[]interface{}{}},
			[]interface{}{},
		},
		{
			[]interface{}{1, 0},
			[]interface{}{1},
		},
		{
			[]interface{}{0, 1},
			[]interface{}{0, 1},
		},
		{
			[]interface{}{0, &[]interface{}{1}},
			[]interface{}{[]interface{}{1}},
		},
		{
			[]interface{}{1, ""},
			[]interface{}{1},
		},
		{
			[]interface{}{1, &[]interface{}{}},
			[]interface{}{1},
		},
		{
			[]interface{}{1, &[]interface{}{1}},
			[]interface{}{1, []interface{}{1}},
		},
		{
			[]interface{}{1, &[]interface{}{2, &[]interface{}{}}},
			[]interface{}{1, []interface{}{2}},
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(fmt.Sprintf("normalize(%v)", tt.s), func(t *testing.T) {
			got := make([]interface{}, len(tt.s))
			copy(got, tt.s)
			normalize(&got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSlice(t *testing.T) {
	initial := []interface{}{1}
	n := newSlice(&initial)
	want := []interface{}{1, n}
	if !reflect.DeepEqual(initial, want) {
		t.Fatalf("Got %v, want %v", initial, want)
	}
	if !reflect.DeepEqual(initial[1], n) {
		t.Fatalf("Got %v, want %v", initial[1], n)
	}
}

func TestNewVersion(t *testing.T) {
	tests := []Version{
		// weird versions
		Version{".1", []interface{}{0, 1}},
		Version{"-1", []interface{}{[]interface{}{1}}},
		// test some major.minor.tiny parsing
		Version{"1", []interface{}{1}},
		Version{"1.0", []interface{}{1}},
		Version{"1.0.0", []interface{}{1}},
		Version{"1.0.0.0", []interface{}{1}},
		Version{"11", []interface{}{11}},
		Version{"11.0", []interface{}{11}},
		Version{"1-1", []interface{}{1, []interface{}{1}}},
		Version{"1-1-1", []interface{}{1, []interface{}{1, []interface{}{1}}}},
		Version{" 1 ", []interface{}{1}},
		// test qualifeirs
		Version{"1.0-ALPHA", []interface{}{1, []interface{}{"alpha"}}},
		Version{"1-alpha", []interface{}{1, []interface{}{"alpha"}}},
		Version{"1.0ALPHA", []interface{}{1, []interface{}{"alpha"}}},
		Version{"1-alpha", []interface{}{1, []interface{}{"alpha"}}},
		Version{"1.0-A", []interface{}{1, []interface{}{"a"}}},
		Version{"1-a", []interface{}{1, []interface{}{"a"}}},
		Version{"1.0A", []interface{}{1, []interface{}{"a"}}},
		Version{"1a", []interface{}{1, []interface{}{"a"}}},
		Version{"1.0-BETA", []interface{}{1, []interface{}{"beta"}}},
		Version{"1-beta", []interface{}{1, []interface{}{"beta"}}},
		Version{"1.0-B", []interface{}{1, []interface{}{"b"}}},
		Version{"1-b", []interface{}{1, []interface{}{"b"}}},
		Version{"1.0B", []interface{}{1, []interface{}{"b"}}},
		Version{"1b", []interface{}{1, []interface{}{"b"}}},
		Version{"1.0-MILESTONE", []interface{}{1, []interface{}{"milestone"}}},
		Version{"1.0-milestone", []interface{}{1, []interface{}{"milestone"}}},
		Version{"1-M", []interface{}{1, []interface{}{"m"}}},
		Version{"1.0-m", []interface{}{1, []interface{}{"m"}}},
		Version{"1M", []interface{}{1, []interface{}{"m"}}},
		Version{"1m", []interface{}{1, []interface{}{"m"}}},
		Version{"1.0-RC", []interface{}{1, []interface{}{"rc"}}},
		Version{"1-rc", []interface{}{1, []interface{}{"rc"}}},
		Version{"1.0-SNAPSHOT", []interface{}{1, []interface{}{"snapshot"}}},
		Version{"1.0-snapshot", []interface{}{1, []interface{}{"snapshot"}}},
		Version{"1-SP", []interface{}{1, []interface{}{"sp"}}},
		Version{"1.0-sp", []interface{}{1, []interface{}{"sp"}}},
		Version{"1-GA", []interface{}{1}},
		Version{"1-ga", []interface{}{1}},
		Version{"1.0-FINAL", []interface{}{1}},
		Version{"1-final", []interface{}{1}},
		Version{"1.0-CR", []interface{}{1, []interface{}{"rc"}}},
		Version{"1-cr", []interface{}{1, []interface{}{"rc"}}},
		// test some transistion
		Version{"1.0-alpha1", []interface{}{1, []interface{}{"alpha", []interface{}{1}}}},
		Version{"1.0-alpha2", []interface{}{1, []interface{}{"alpha", []interface{}{2}}}},
		Version{"1.0.0alpha1", []interface{}{1, []interface{}{"alpha", []interface{}{1}}}},
		Version{"1.0-beta1", []interface{}{1, []interface{}{"beta", []interface{}{1}}}},
		Version{"1-beta2", []interface{}{1, []interface{}{"beta", []interface{}{2}}}},
		Version{"1.0.0beta1", []interface{}{1, []interface{}{"beta", []interface{}{1}}}},
		Version{"1.0-BETA1", []interface{}{1, []interface{}{"beta", []interface{}{1}}}},
		Version{"1-BETA2", []interface{}{1, []interface{}{"beta", []interface{}{2}}}},
		Version{"1.0.0BETA1", []interface{}{1, []interface{}{"beta", []interface{}{1}}}},
		Version{"1.0-milestone1", []interface{}{1, []interface{}{"milestone", []interface{}{1}}}},
		Version{"1.0-milestone2", []interface{}{1, []interface{}{"milestone", []interface{}{2}}}},
		Version{"1.0.0milestone1", []interface{}{1, []interface{}{"milestone", []interface{}{1}}}},
		Version{"1.0-MILESTONE1", []interface{}{1, []interface{}{"milestone", []interface{}{1}}}},
		Version{"1.0-milestone2", []interface{}{1, []interface{}{"milestone", []interface{}{2}}}},
		Version{"1.0.0MILESTONE1", []interface{}{1, []interface{}{"milestone", []interface{}{1}}}},
		Version{"1.0-alpha2snapshot", []interface{}{1, []interface{}{"alpha", []interface{}{2, []interface{}{"snapshot"}}}}},
	}

	t.Parallel()
	for _, want := range tests {
		t.Run(want.unparsed, func(t *testing.T) {
			got := New(want.unparsed)
			if !reflect.DeepEqual(got, &want) {
				t.Errorf("New(%v): got %v, want %v", want.unparsed, got, want)
			}
		})
	}
}

func TestVersionQualifiers(t *testing.T) {
	qualifiers := []string{"1-alpha2snapshot", "1-alpha2", "1-alpha-123",
		"1-beta-2", "1-beta123", "1-m2", "1-m11", "1-rc", "1-cr2", "1-rc123",
		"1-SNAPSHOT", "1", "1-sp", "1-sp2", "1-sp123", "1-abc", "1-def",
		"1-pom-1", "1-1-snapshot", "1-1", "1-2", "1-123"}

	t.Parallel()
	for i, low := range qualifiers[:len(qualifiers)-1] {
		for _, high := range qualifiers[i+1:] {
			t.Run(low+" < "+high, func(t *testing.T) {
				if !assertVersionOrder(low, high) {
					t.Error("got false")
				}
			})
		}
	}
}

func TestVersionNumbers(t *testing.T) {
	numbers := []string{"2.0", "2-1", "2.0.a", "2.0.0.a", "2.0.2", "2.0.123",
		"2.1.0", "2.1-a", "2.1b", "2.1-x", "2.1-1", "2.1.0.1", "2.2", "2.123",
		"11.a2", "11.a11", "11.b2", "11.b11", "11.m2", "11.m11", "11", "11.a",
		"11b", "11c", "11m"}

	t.Parallel()
	for i, low := range numbers[:len(numbers)-1] {
		for _, high := range numbers[i+1:] {
			t.Run(low+" < "+high, func(t *testing.T) {
				if !assertVersionOrder(low, high) {
					t.Error("got false")
				}
			})
		}
	}
}

func TestVersionEquality(t *testing.T) {
	tests := []struct {
		a, b string
	}{
		{"1", "1"},
		{"1", "1.0"},
		{"1", "1.0.0"},
		{"1.0", "1.0.0"},
		{"1", "1-0"},
		{"1", "1.0-0"},
		{"1.0", "1.0-0"},

		// no separator between number and character
		{"1a", "1-a"},
		{"1a", "1.0-a"},
		{"1a", "1.0.0-a"},
		{"1.0a", "1-a"},
		{"1.0.0a", "1-a"},
		{"1x", "1-x"},
		{"1x", "1.0-x"},
		{"1x", "1.0.0-x"},
		{"1.0x", "1-x"},
		{"1.0.0x", "1-x"},

		// aliases
		{"1ga", "1"},
		{"1final", "1"},
		{"1cr", "1rc"},

		// special "aliases" a, b and m for alpha, beta and milestone
		{"1a1", "1-alpha-1"},
		{"1b2", "1-beta-2"},
		{"1m3", "1-milestone-3"},

		// case insensitive
		{"1X", "1x"},
		{"1A", "1a"},
		{"1B", "1b"},
		{"1M", "1m"},
		{"1Ga", "1"},
		{"1GA", "1"},
		{"1Final", "1"},
		{"1FinaL", "1"},
		{"1FINAL", "1"},
		{"1Cr", "1Rc"},
		{"1cR", "1rC"},
		{"1m3", "1Milestone3"},
		{"1m3", "1MileStone3"},
		{"1m3", "1MILESTONE3"},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.a+" == "+tt.b, func(t *testing.T) {
			if !assertVersionEquality(tt.a, tt.b) {
				t.Error("got false")
			}
		})
	}
}

func TestVersionCompare(t *testing.T) {
	tests := []struct{ low, high string }{
		{"1", "2"},
		{"1.5", "2"},
		{"1", "2.5"},
		{"1.0", "1.1"},
		{"1.1", "1.2"},
		{"1.0.0", "1.1"},
		{"1.0.1", "1.1"},
		{"1.1", "1.2.0"},
		{"1.0-alpha-1", "1.0"},
		{"1.0-alpha-1", "1.0-alpha-2"},
		{"1.0-alpha-1", "1.0-beta-1"},
		{"1.0-beta-1", "1.0-SNAPSHOT"},
		{"1.0-SNAPSHOT", "1.0"},
		{"1.0-alpha-1-SNAPSHOT", "1.0-alpha-1"},
		{"1.0", "1.0-1"},
		{"1.0-1", "1.0-2"},
		{"1.0.0", "1.0-1"},
		{"2.0-1", "2.0.1"},
		{"2.0.1-klm", "2.0.1-lmn"},
		{"2.0.1", "2.0.1-xyz"},
		{"2.0.1", "2.0.1-123"},
		{"2.0.1-xyz", "2.0.1-123"},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.low+" < "+tt.high, func(t *testing.T) {
			if !assertVersionOrder(tt.low, tt.high) {
				t.Error("got false")
			}
		})
	}
}

func BenchmarkVercmp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Vercmp("1.2.3-milestone.1", "1.2.3-milestone.2")
	}
}

func BenchmarkVercmpVersion(b *testing.B) {
	v1 := New("1.2.3-milestone.1")
	v2 := New("1.2.3-milestone.2")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Vercmp(v1, v2)
	}
}

func assertVersionEquality(v1, v2 string) bool {
	if Vercmp(v1, v2) != 0 {
		return false
	}
	if Vercmp(v2, v1) != 0 {
		return false
	}
	return true
}

func assertVersionOrder(low, high string) bool {
	if Vercmp(low, high) >= 0 {
		return false
	}
	if Vercmp(high, low) <= 0 {
		return false
	}
	return true
}
