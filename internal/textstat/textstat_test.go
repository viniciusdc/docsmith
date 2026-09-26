package textstat

import (
	"strings"
	"testing"
)

func TestCount(t *testing.T) {
	cases := []struct {
		in   string
		want Stats
	}{
		{"", Stats{}},
		{"one\n", Stats{Lines: 1, Words: 1, Bytes: 4}},
		{"two words\nand three more\n", Stats{Lines: 2, Words: 5, Bytes: 25}},
		{"no newline", Stats{Lines: 0, Words: 2, Bytes: 10}},
		{"  spaced\t\tout  \n", Stats{Lines: 1, Words: 2, Bytes: 16}},
		{"héllo wörld\n", Stats{Lines: 1, Words: 2, Bytes: 14}},
	}
	for _, c := range cases {
		got, err := Count(strings.NewReader(c.in))
		if err != nil {
			t.Fatalf("Count(%q): %v", c.in, err)
		}
		if got != c.want {
			t.Errorf("Count(%q) = %+v, want %+v", c.in, got, c.want)
		}
	}
}

func TestAdd(t *testing.T) {
	got := Stats{1, 2, 3}.Add(Stats{4, 5, 6})
	if want := (Stats{5, 7, 9}); got != want {
		t.Errorf("Add = %+v, want %+v", got, want)
	}
}
