package vercmp

import (
	"github.com/wfscheper/vercmp/maven"
	"github.com/wfscheper/vercmp/semver"
)

// Vercmp compares two versions a and b, and returns an integer less than zero
// if a is older than b, 0 if a is the same as b, or an integer greater than
// zero if a is newer than b
func Vercmp(a, b string) int {
	if aVer, err := semver.New(a); err == nil {
		if bVer, err := semver.New(b); err == nil {
			return semver.Vercmp(aVer, bVer)
		}
	}
	aVer := maven.New(a)
	bVer := maven.New(b)
	return maven.Vercmp(aVer, bVer)
}
