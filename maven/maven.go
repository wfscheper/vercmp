package maven

import (
	"fmt"
	"strconv"
	"strings"
)

var aliases map[string]string = map[string]string{
	"ga":    "",
	"final": "",
	"cr":    "rc",
}

var qualifiers [7]string = [7]string{"alpha", "beta", "milestone", "rc", "snapshot", "", "sp"}

type MavenVersion struct {
	unparsed string
	parsed   []interface{}
}

func (m *MavenVersion) String() string {
	return m.unparsed
}

// New returns a new MavenVersion parsed from the version string v.
func New(v string) *MavenVersion {
	parsed := make([]interface{}, 0, 10)
	currentSlice := &parsed
	start := 0
	isDigit := false

	buf := strings.ToLower(strings.TrimSpace(v))
	for idx, ch := range buf {
		if ch == '.' {
			if idx == start {
				*currentSlice = append(*currentSlice, 0)
			} else {
				*currentSlice = append(*currentSlice,
					parseBuffer(buf[start:idx], false))
			}
			start = idx + 1
		} else if ch == '-' {
			if idx == start {
				*currentSlice = append(*currentSlice, 0)
			} else {
				*currentSlice = append(*currentSlice,
					parseBuffer(buf[start:idx], false))
			}
			start = idx + 1
			currentSlice = newSlice(currentSlice)
		} else if _, err := strconv.Atoi(string(ch)); err == nil {
			if !isDigit && idx > start {
				*currentSlice = append(*currentSlice,
					parseBuffer(buf[start:idx], true))
				currentSlice = newSlice(currentSlice)
				start = idx
			}
			isDigit = true
		} else {
			if isDigit && idx > start {
				*currentSlice = append(*currentSlice,
					parseBuffer(buf[start:idx], false))
				currentSlice = newSlice(currentSlice)
				start = idx
			}
			isDigit = false
		}
	}
	if len(buf) > start {
		*currentSlice = append(*currentSlice, parseBuffer(buf[start:], false))
	}
	normalize(&parsed)
	return &MavenVersion{v, parsed}
}

func Vercmp(a, b *MavenVersion) int {
	return compareSlice(a.parsed, b.parsed)
}

func compare(a, b interface{}) int {
	switch a := a.(type) {
	case int:
		return compareInt(a, b)
	case string:
		return compareString(a, b)
	case []interface{}:
		return compareSlice(a, b)
	default:
		return 1
	}
}

func compareInt(a int, b interface{}) int {
	switch b := b.(type) {
	default:
		panic(fmt.Sprintf("Unkown type %t", b))
	case int:
		return a - b
	case string, []interface{}, nil:
		return a
	}
}

func compareSlice(a []interface{}, b interface{}) int {
	switch b := b.(type) {
	default:
		panic(fmt.Sprintf("Unkown type %t", b))
	case nil:
		if len(a) == 0 {
			return 0
		}
		return compare(a[0], b)
	case int:
		return -1
	case string:
		return 1
	case []interface{}:
		var result int
		for _, pair := range zip(a, b) {
			if pair.Left == nil {
				if pair.Right == nil {
					result = 0
				} else {
					result = -1 * compare(pair.Right, pair.Left)
				}
			} else {
				result = compare(pair.Left, pair.Right)
			}
			if result != 0 {
				return result
			}
		}
		return 0
	}
}

func compareString(a string, b interface{}) int {
	switch b := b.(type) {
	default:
		panic(fmt.Sprintf("Unkown type %t", b))
	case int, []interface{}:
		return -1
	case nil:
		return compareString(a, "")
	case string:
		a_value := stringValue(a)
		b_value := stringValue(b)
		if a_value < b_value {
			return -1
		} else if a_value == b_value {
			return 0
		} else {
			return 1
		}
	}
}

// appendSlicePtr appends to a slice poitner
func appendSlicePtr(sPtr *[]interface{}, elem ...interface{}) {
	s := *sPtr
	s = append(s, elem...)
	*sPtr = s
}

// parseBuffer determines if the string b is an integer or a string
func parseBuffer(b string, digitFollows bool) interface{} {
	if r, err := strconv.Atoi(b); err == nil {
		return r
	}
	if digitFollows && len(b) == 1 {
		switch b {
		case "a":
			b = "alpha"
		case "b":
			b = "beta"
		case "m":
			b = "milestone"
		}
	}
	if r, ok := aliases[b]; ok {
		return r
	}
	return b
}

// newSlice appends a new slice to s and returns the new slice
func newSlice(sPtr *[]interface{}) *[]interface{} {
	n := make([]interface{}, 0)
	*sPtr = append(*sPtr, &n)
	return &n
}

// normalize removes zero-value elements from the end of s until it encounters
// a non-zero value that isn't a slice and returns the modified slice
func normalize(sPtr *[]interface{}) {
	s := *sPtr
	for i := len(s) - 1; i >= 0; i-- {
		switch e := s[i].(type) {
		case int:
			if e == 0 {
				s = append(s[:i], s[i+1:]...)
			} else {
				*sPtr = s
				return
			}
		case string:
			if e == "" {
				s = append(s[:i], s[i+1:]...)
			} else {
				*sPtr = s
				return
			}
		case *[]interface{}:
			normalize(e)
			if len(*e) == 0 {
				s = append(s[:i], s[i+1:]...)
			} else {
				s[i] = *e
			}
		}
	}
	*sPtr = s
}

func stringValue(s string) string {
	for idx, q := range qualifiers {
		if s == q {
			return strconv.FormatInt(int64(idx+1), 10)
		}
	}
	return fmt.Sprintf("%d-%s", len(qualifiers), s)
}

type interfaceTuple struct {
	Left, Right interface{}
}

func zip(a, b []interface{}) []interfaceTuple {
	var r []interfaceTuple

	len_a := len(a)
	len_b := len(b)

	if len_a < len_b {
		r = make([]interfaceTuple, len_b)
		for i, e := range a {
			r[i] = interfaceTuple{Left: e, Right: b[i]}
		}
		for i, e := range b[len_a:] {
			r[len_a+i] = interfaceTuple{Left: nil, Right: e}
		}
	} else {
		r = make([]interfaceTuple, len_a)
		for i, e := range b {
			r[i] = interfaceTuple{Left: a[i], Right: e}
		}
		for i, e := range a[len_b:] {
			r[len_b+i] = interfaceTuple{Left: e, Right: nil}
		}
	}
	return r
}
