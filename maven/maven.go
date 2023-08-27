package maven

import (
	"fmt"
	"strconv"
	"strings"
)

var aliases = map[string]string{
	"ga":    "",
	"final": "",
	"cr":    "rc",
}

var qualifiers = [7]string{"alpha", "beta", "milestone", "rc", "snapshot", "", "sp"}

// Version repersents a parsed Maven 3 version string.
type Version struct {
	unparsed string
	parsed   []interface{}
}

// New returns a new Version parsed from the version string v.
func New(v string) *Version {
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
	return &Version{v, parsed}
}

// String returns the oringal Maven version.
func (m *Version) String() string {
	return m.unparsed
}

// Vercmp compares two Maven 3 versions, a and b, and returns 1 if a is newer
// than b, 0 if a and b are equal, or -1 if a is older than b. a and b an be
// either a string or a Version.
func Vercmp(a, b interface{}) int {
	var aVer, bVer *Version
	switch a := a.(type) {
	case string:
		aVer = New(a)
	case *Version:
		aVer = a
	case Version:
		aVer = &a
	}
	switch b := b.(type) {
	case string:
		bVer = New(b)
	case *Version:
		bVer = b
	case Version:
		bVer = &b
	}
	return compareSlice(aVer.parsed, bVer.parsed)
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
		aValue := stringValue(a)
		bValue := stringValue(b)
		if aValue < bValue {
			return -1
		} else if aValue == bValue {
			return 0
		} else {
			return 1
		}
	}
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

	aLen := len(a)
	bLen := len(b)

	if aLen < bLen {
		r = make([]interfaceTuple, bLen)
		for i, e := range a {
			r[i] = interfaceTuple{Left: e, Right: b[i]}
		}
		for i, e := range b[aLen:] {
			r[aLen+i] = interfaceTuple{Left: nil, Right: e}
		}
	} else {
		r = make([]interfaceTuple, aLen)
		for i, e := range b {
			r[i] = interfaceTuple{Left: a[i], Right: e}
		}
		for i, e := range a[bLen:] {
			r[bLen+i] = interfaceTuple{Left: e, Right: nil}
		}
	}
	return r
}
