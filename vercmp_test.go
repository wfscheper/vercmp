package vercmp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVercmp(t *testing.T) {
	is := assert.New(t)
	tests := []struct {
		title, a, b string
		expected    int
	}{
		{"Semantic less than semantic", "1.0.0", "2.0.0", -1},
		{"Semantic equal to semantic", "2.0.0", "2.0.0", 0},
		{"Semantic greater than semantic", "2.0.0", "1.0.0", 1},
		{"Maven less than maven", "1", "2.0", -1},
		{"Maven equal to maven", "1", "1.0", 0},
		{"Maven greater than manve", "2.0", "1", 1},
		{"Semantic vs Maven", "1.0.0", "1", 0},
		{"Maven vs Semantic", "1", "1.0.0", 0},
		{"Semantic vs gibberish", "1.0.0", "flargh", 1},
	}

	for _, test := range tests {
		actual := Vercmp(test.a, test.b)
		if test.expected < 0 {
			is.True(actual < 0, "(%s) Expected Vercmp(%s, %s) < 0, got %d",
				test.title, test.a, test.b, actual)
		} else if test.expected == 0 {
			is.True(actual == 0, "(%s) Expected Vercmp(%s, %s) == 0, got %d",
				test.title, test.a, test.b, actual)
		} else {
			is.True(actual > 0, "(%s) Expected Vercmp(%s, %s) > 0, got %d",
				test.title, test.a, test.b, actual)
		}
	}
}
