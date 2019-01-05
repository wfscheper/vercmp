// Package semver implements parsing and comparing semantic versions
//
// Semantic versions conform to the Semantic Versioning 3.0.0 standard
// described at http://docs.openstack.org/developer/pbr/semver.html.
package semver

import (
	"fmt"
	"strconv"
	"strings"
)

const maxInt = int(^uint(0) >> 1)

var typeMap = map[string]int{
	"a":  1,
	"b":  2,
	"rc": 3,
	"":   4,
}

// Version represents a Semantic Version string.
type Version struct {
	Major, Minor, Patch, PreRelease, DevCount int
	PreReleaseType                            string
}

func (s Version) keys() [7]int {
	k := [7]int{}
	k[0], k[1], k[2] = s.Major, s.Minor, s.Patch
	if s.DevCount != 0 && s.PreReleaseType == "" {
		k[3] = 0
	} else {
		k[3] = 1
	}
	k[4] = typeMap[s.PreReleaseType]
	k[5] = s.PreRelease
	if s.DevCount != 0 {
		k[6] = s.DevCount
	} else {
		k[6] = maxInt
	}
	return k
}

func (s Version) String() string {
	str := fmt.Sprintf("%d.%d.%d", s.Major, s.Minor, s.Patch)
	if s.PreReleaseType != "" {
		str = str + fmt.Sprintf(".%s%d", s.PreReleaseType, s.PreRelease)
	}
	if s.DevCount != 0 {
		str = str + fmt.Sprintf(".dev%d", s.DevCount)
	}
	return str
}

// New parses a semantic version, per Semantic Versioning 3.0.0.
func New(v string) (*Version, error) {
	s := new(Version)
	parsed := strings.Split(strings.ToLower(strings.TrimSpace(v)), ".")
	if len(parsed) < 3 {
		return nil, fmt.Errorf("Invalid semantic version: %v", v)
	}

	part, parsed := pop(parsed)
	if major, err := strconv.Atoi(part); err == nil {
		s.Major = major
	} else {
		return nil, fmt.Errorf("Invalid major version: %v", v)
	}

	part, parsed = pop(parsed)
	if minor, err := strconv.Atoi(part); err == nil {
		s.Minor = minor
	} else {
		return nil, fmt.Errorf("Invalid minor version: %v", v)
	}

	part, parsed = pop(parsed)
	if patch, err := strconv.Atoi(part); err == nil {
		s.Patch = patch
	} else {
		return nil, fmt.Errorf("Invalid patch version: %v", v)
	}

	if len(parsed) > 0 {
		if parsed[0][0] == 'a' || parsed[0][0] == 'b' || parsed[0][0] == 'r' {
			// pre-release
			part, parsed = pop(parsed)
			if preReleaseType, preRelease, err := parsePreRelease(part); err == nil {
				s.PreReleaseType, s.PreRelease = preReleaseType, preRelease
			} else {
				return nil, fmt.Errorf("Invalid pre-release version: %v", v)
			}
		}
		if part, parsed = pop(parsed); part != "" {
			// dev part
			if len(part) < 4 || part[:3] != "dev" {
				return nil, fmt.Errorf("Invalid dev version: %v", v)
			}
			if devCount, err := strconv.Atoi(part[3:]); err == nil {
				s.DevCount = devCount
			} else {
				return nil, fmt.Errorf("Invalid dev version: %v", v)
			}
		}
	}
	return s, nil
}

// Vercmp compares two semantic versions and returns an integer less than 0
// if a is older than b, 0 if a and b are the same, and an integer greater than
// 0 if a is newer than b.
func Vercmp(a, b interface{}) int {
	var aVer, bVer *Version
	var err error

	switch a := a.(type) {
	default:
		panic(fmt.Sprintf("Unparsable type %T", a))
	case string:
		aVer, err = New(a)
		if err != nil {
			panic(fmt.Sprint(err))
		}
	case Version:
		aVer = &a
	case *Version:
		aVer = a
	}
	switch b := b.(type) {
	default:
		panic(fmt.Sprintf("Unparsable type %T", b))
	case string:
		bVer, err = New(b)
		if err != nil {
			panic(fmt.Sprint(err))
		}
	case Version:
		bVer = &b
	case *Version:
		bVer = b
	}

	bKeys := bVer.keys()
	for idx, aKey := range aVer.keys() {
		bKey := bKeys[idx]
		if aKey != bKey {
			return aKey - bKey
		}
	}
	return 0
}

// parseBuffer converts a numeric string to an interger, otherwise it returns
// the string unmodified.
func parseBuffer(b string) interface{} {
	if r, err := strconv.Atoi(b); err == nil {
		return r
	}
	return b
}

// parsePreRelease parses s and returns the pre-release type and the
// pre-release version.
func parsePreRelease(s string) (string, int, error) {
	var preReleaseType string

	if s[0] == 'r' {
		preReleaseType, s = s[:2], s[2:]
	} else {
		preReleaseType, s = string(s[0]), s[1:]
	}
	if preRelease, err := strconv.Atoi(s); err == nil {
		return preReleaseType, preRelease, err
	}
	return "", 0, nil
}

// pop removes the first element from the slice, and returns it and the
// remainder of the slice. If the slice contains only one element, then pop
// returns the element and an empty slice. If the slice is empty, pop returns
// an empty string and an empty slice.
func pop(s []string) (string, []string) {
	n := len(s)
	if n > 1 {
		return s[0], s[1:]
	} else if n == 1 {
		return s[0], []string{}
	} else {
		return "", s
	}
}
