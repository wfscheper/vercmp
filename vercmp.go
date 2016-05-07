package vercmp

import (
	"github.com/wfscheper/vercmp/maven"
	"github.com/wfscheper/vercmp/semver"
)

//SemVerCmp compares two semantic version strings, and returns a negative
// integer if a is older than b, 0 if a is the same as b, and a positive
// integer if a is newer than b
func SemVerCmp(a, b string) int {
	return semver.Vercmp(a, b)
}

// MavenVerCmp compares to maven version strings, and returns a negative
// integer if a is older than b, 0 if a is the same as b, and a positive
// integer if a is newer than b
func MavenVerCmp(a, b string) int {
	return maven.Vercmp(a, b)
}
