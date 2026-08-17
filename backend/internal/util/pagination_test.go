package util

import "testing"

func TestParsePage(t *testing.T) {
	cases := []struct {
		raw  string
		want int
	}{
		{"", 1},
		{"0", 1},
		{"-3", 1},
		{"abc", 1},
		{"3", 3},
	}
	for _, tc := range cases {
		if got := ParsePage(tc.raw); got != tc.want {
			t.Errorf("ParsePage(%q) = %d, want %d", tc.raw, got, tc.want)
		}
	}
}

func TestParsePageSize(t *testing.T) {
	cases := []struct {
		raw  string
		want int
	}{
		{"", 10},
		{"5", 5},
		{"999", 200},
		{"abc", 10},
	}
	for _, tc := range cases {
		if got := ParsePageSize(tc.raw); got != tc.want {
			t.Errorf("ParsePageSize(%q) = %d, want %d", tc.raw, got, tc.want)
		}
	}
}
