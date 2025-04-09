package mqgame

import (
	"math"
	"testing"
)

func TestLongestMutualSubstring(t *testing.T) {
	for _, row := range []struct {
		inputA, inputB string
		want           string
	}{
		{"1234", "12", "12"},
		{"12", "1234", "12"},
		{"1234", "34", "34"},
		{"1234", "234", "234"},
		{"123444567", "3456", "456"},
		{"123444567", "34566", "456"},
		// "56" would be a valid substring,
		// but this func is defined to return the earlier match.
		{"123444567", "34556", "34"},
		{"12345", "6789", ""},
	} {
		got := longestMutualSubstring(row.inputA, row.inputB)

		if got != row.want {
			t.Errorf("longestMutualSubstring(%q, %q) = %q, want %q", row.inputA, row.inputB, got, row.want)
		}
	}
}

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
