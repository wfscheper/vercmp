package semver

import (
	"fmt"
	"strconv"
	"strings"
)

const maxInt = int(^uint(0) >> 1)

var typeMap map[string]int = map[string]int{
	"a":  1,
	"b":  2,
	"rc": 3,
	"":   4,
}

type SemanticVersion struct {
	Major, Minor, Patch, PreRelease, DevCount int
	PreReleaseType                            string
}

func (s SemanticVersion) Keys() [7]int {
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

func New(v string) (*SemanticVersion, error) {
	s := new(SemanticVersion)
	parsed := strings.Split(strings.ToLower(strings.TrimSpace(v)), ".")
	if len(parsed) < 3 {
		return nil, fmt.Errorf("Invalid semantic version: %v\n", v)
	}

	part, parsed := pop(parsed)
	if major, err := strconv.Atoi(part); err == nil {
		s.Major = major
	} else {
		return nil, fmt.Errorf("Invalid major version: %v\n", v)
	}

	part, parsed = pop(parsed)
	if minor, err := strconv.Atoi(part); err == nil {
		s.Minor = minor
	} else {
		return nil, fmt.Errorf("Invalid minor version: %v\n", v)
	}

	part, parsed = pop(parsed)
	if patch, err := strconv.Atoi(part); err == nil {
		s.Patch = patch
	} else {
		return nil, fmt.Errorf("Invalid patch version: %v\n", v)
	}

	if len(parsed) > 0 {
		if parsed[0][0] == 'a' || parsed[0][0] == 'b' || parsed[0][0] == 'r' {
			// pre-release
			part, parsed = pop(parsed)
			if preReleaseType, preRelease, err := parsePreRelease(part); err == nil {
				s.PreReleaseType, s.PreRelease = preReleaseType, preRelease
			} else {
				return nil, fmt.Errorf("Invalid pre-release version: %v\n", v)
			}
		}
		if part, parsed = pop(parsed); part != "" {
			// dev part
			if len(part) < 4 || part[:3] != "dev" {
				return nil, fmt.Errorf("Invalid dev version: %v\n", v)
			}
			if devCount, err := strconv.Atoi(part[3:]); err == nil {
				s.DevCount = devCount
			} else {
				return nil, fmt.Errorf("Invalid dev version: %v\n", v)
			}
		}
	}
	return s, nil
}

func Vercmp(a, b *SemanticVersion) int {
	b_keys := b.Keys()
	for idx, a_key := range a.Keys() {
		b_key := b_keys[idx]
		if a_key != b_key {
			return a_key - b_key
		}
	}
	return 0
}

func parseBuffer(b string) interface{} {
	if r, err := strconv.Atoi(b); err == nil {
		return r
	} else {
		return b
	}
}

func parsePreRelease(s string) (string, int, error) {
	var preReleaseType string

	if s[0] == 'r' {
		preReleaseType, s = s[:2], s[2:]
	} else {
		preReleaseType, s = string(s[0]), s[1:]
	}
	if preRelease, err := strconv.Atoi(s); err == nil {
		return preReleaseType, preRelease, err
	} else {
		return "", 0, err
	}
}

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
