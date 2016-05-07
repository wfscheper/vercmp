package vercmp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMavenVerCmp(t *testing.T) {
	is := assert.New(t)
	tests := []struct {
		a, b     string
		expected int
	}{
		{"1", "2.0", -1},
		{"1", "1.0", 0},
		{"2.0", "1", 1},
	}

	var msg string
	var cmp bool
	for _, test := range tests {
		actual := MavenVerCmp(test.a, test.b)
		if test.expected < 0 {
			cmp = actual < 0
			msg = "Expected MavenVerCmp(%s, %s) < 0, got %d"
		} else if test.expected == 0 {
			cmp = actual == 0
			msg = "Expected MavenVerCmp(%s, %s) == 0, got %d"
		} else {
			cmp = actual > 0
			msg = "Expected MavenVerCmp(%s, %s) > 0, got %d"
		}
		is.True(cmp, msg, test.a, test.b, actual)
	}
}

func TestSemVerCmp(t *testing.T) {
	is := assert.New(t)
	tests := []struct {
		a, b     string
		expected int
	}{
		{"1.0.0", "2.0.0", -1},
		{"2.0.0", "1.0.0", 1},
		{"2.0.0", "2.0.0", 0},
	}

	var msg string
	var cmp bool
	for _, test := range tests {
		actual := SemVerCmp(test.a, test.b)
		if test.expected < 0 {
			cmp = actual < 0
			msg = "Expected SemVerCmp(%s, %s) < 0, got %d"
		} else if test.expected == 0 {
			cmp = actual == 0
			msg = "Expected SemVerCmp(%s, %s) == 0, got %d"
		} else {
			cmp = actual > 0
			msg = "Expected SemVerCmp(%s, %s) > 0, got %d"
		}
		is.True(cmp, msg, test.a, test.b, actual)
	}
}
