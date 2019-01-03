package vercmp

import (
	"testing"
)

func TestMavenVerCmp(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"1", "2.0", -1},
		{"1", "1.0", 0},
		{"2.0", "1", 1},
	}

	var msg string
	var cmp bool
	for _, tt := range tests {
		got := MavenVerCmp(tt.a, tt.b)
		switch {
		case tt.want < 0:
			cmp = got < 0
			msg = "Expected MavenVerCmp(%s, %s) < 0, got %d"
		case tt.want == 0:
			cmp = got == 0
			msg = "Expected MavenVerCmp(%s, %s) == 0, got %d"
		default:
			cmp = got > 0
			msg = "Expected MavenVerCmp(%s, %s) > 0, got %d"
		}
		if !cmp {
			t.Errorf(msg, tt.a, tt.b, got)
		}
	}
}

func TestSemVerCmp(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"1.0.0", "2.0.0", -1},
		{"2.0.0", "1.0.0", 1},
		{"2.0.0", "2.0.0", 0},
	}

	var msg string
	var cmp bool
	for _, tt := range tests {
		got := SemVerCmp(tt.a, tt.b)
		switch {
		case tt.want < 0:
			cmp = got < 0
			msg = "Expected SemVerCmp(%s, %s) < 0, got %d"
		case tt.want == 0:
			cmp = got == 0
			msg = "Expected SemVerCmp(%s, %s) == 0, got %d"
		default:
			cmp = got > 0
			msg = "Expected SemVerCmp(%s, %s) > 0, got %d"
		}
		if !cmp {
			t.Errorf(msg, tt.a, tt.b, got)
		}
	}
}
