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

		// Unicode awareness: Treat as 1/2 code points, not 1/4 characters
		{"$€", "$E", 0.5},
		// (two errors vs. four characters in the proper title,
		// no substring matching any more)
		{"abcd (ft. pant)", "abcd (", 0.5},

		{"abcd (ft. banana pants)", "abcd", 1},
		{"abcd (ft banana pants)", "abcd", 1},
		{"abcd (feat. banana pants)", "abcd", 1},
		{"abcd (feat banana pants)", "abcd", 1},
		{"abcd (featuring banana pants)", "abcd", 1},
		{"abcd (featuring. banana pants)", "abcd", 1},

		// Ignore all parentheticals
		{"abcd (Taylor's Version)", "abcd", 1},
		{"abcd (Taylor's Version) (From The Vault)", "abcd", 1},
		{"abcd (Taylor's Version) [From The Vault]", "abcd", 1},
		{"abcd [Club Remix]", "abcd", 1},
		{"abcd", "abcd (Yes I know it's weird that I know this song)", 1},

		// But, if you guess the parenthetical perfectly, get some bonus points
		// (the parenthetical is 16 which is 4x the length of the title, so you get 2x as many points)
		{"abcd (ft. banana pants)", "abcd (ft. banana pants)", 3},
		{"abcd (ft. banana pants)", "ABCD (FT. BANANA PANTS)", 3},
		// 9 / 4 / 2 = extra 9/8 points
		{"abcd (ft. pants)", "abcd (ft. pants)", 2.125},
		// Getting ANY part of the parenthetical incorrect means no bonus points
		{"abcd (ft. pants)", "abcd (pants)", 1},
		{"abcd (ft. pants)", "abcd (ft. paints)", 1},
		// You can still get bonus points if you ONLY know the parenthetical
		{"abcd (ft. pants)", "I dunno (ft. pants)", 1.125},
		// Unicode awareness: Treat as 4 code points, not 8 characters
		{"abcd (€€€€)", "abcd (€€€€)", 1.5},

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
