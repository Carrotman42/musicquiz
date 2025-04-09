package mqgame

import (
	"math"
	"testing"
)

func TestCalcSimilarity(t *testing.T) {
	for _, row := range []struct {
		target, guess string
		want          float64
	}{
		{"abc", "abc", 1},
		{"Abc", "abC", 1},
		{"abcd", "BC", 0.5},

		{"abcd (ft. pants)", "abcd (ft. pants)", 1},
		// (two errors vs. four characters in the proper title,
		// no substring matching any more)
		{"abcd (ft. pant)", "abcd (", 0.5},

		{"abcd (ft. banana pants)", "abcd", 1},
		{"abcd (ft banana pants)", "abcd", 1},
		{"abcd (feat. banana pants)", "abcd", 1},
		{"abcd (feat banana pants)", "abcd", 1},
		{"abcd (featuring banana pants)", "abcd", 1},
		{"abcd (featuring. banana pants)", "abcd", 1},

		// Spaces on the inside matter, spaces on the outside don't
		{"abcd", "abcd ", 1},
		{"abcd", " abcd", 1},
		{" abcd", "abcd", 1},
		{"abcd ", "abcd", 1},
		{"abcd", "ab cd", 0.75},
		{"ab cd", "abcd", 0.8},

		{"abcd", "abed", 0.75},
		{"abcd", "abd", 0.75},
		// One wrong out of four characters (even though they *also*
		// got four right) is still 75%
		{"abcd", "abxcd", 0.75},
		{"abcd", "efgh", 0},
		// (controversially, this does have "...a...d" in it, but is so
		// wrong I think it deserves to not have those counted at all)
		{"abcd", "this takes more than 4 deletions", 0},
		{"abcd", "", 0},
		{"abcd", "c", 0.25},
		{"abcd", "  ac", 0.5},
	} {
		got := calcSimilarity(row.target, row.guess)
		if diff := got - row.want; math.Abs(diff) > 0.0001 {
			t.Errorf("calcSimilarity(%q, %q) = %v, want %v", row.target, row.guess, got, row.want)
		}
	}
}
