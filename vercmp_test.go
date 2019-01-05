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

	for _, tt := range tests {
		got := MavenVerCmp(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("MavenVerCmp(%s, %s): got %d, want %d", tt.a, tt.b, got, tt.want)
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

	for _, tt := range tests {
		got := SemVerCmp(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("SemVerCmp(%s, %s): got %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}
