package mqgame

import (
	"chowski3/common/automation/ytmgui"
	"cmp"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
	"unsafe"
)

type GuessRound struct {
	Song    SongInfo
	Guesses []SongGuess
	Passes  []*Player
}

type normmap map[string]*Player

func (nm normmap) normalize(p *Player) *Player {
	if existing, ok := nm[p.Name]; ok {
		return existing
	}
	nm[p.Name] = p
	return p
}

func (gr *GuessRound) normalizePlayers(players normmap) {
	for i, p := range gr.Passes {
		gr.Passes[i] = players.normalize(p)
	}
	for i := range gr.Guesses {
		gr.Guesses[i].normalizePlayers(players)
	}
}

// TODO: not a great dependency, I think.
type SongInfo = ytmgui.SongInfo

type SongGuess struct {
	Player *Player
	Delay  time.Duration
	Guess  string

	Votes []PlayerVote
}

func (sg *SongGuess) normalizePlayers(players normmap) {
	sg.Player = players.normalize(sg.Player)
	for i := range sg.Votes {
		sg.Votes[i].Player = players.normalize(sg.Votes[i].Player)
	}
}

// Safe to copy by value.
type PlayerVote struct {
	Player *Player
	Voted  bool
	Vote   bool
}

func (sg SongGuess) Correct(curSong SongInfo) (similarity float64, closeEnough bool) {
	d := sg.Description(curSong)
	return d.Score, d.Correct()
}

func (sg SongGuess) Clone() SongGuess {
	return SongGuess{
		sg.Player,
		sg.Delay,
		sg.Guess,
		slices.Clone(sg.Votes),
	}
}

// Primarily for the purposes of dealing with the template package, and its
// lack of handling multiple returned values.
type SongGuessDescription struct {
	SongGuess
	// Similarity score.
	Score              float64 // [0, 1]
	CorrectedByContest bool
}

func (sgd SongGuessDescription) ScorePercent() int {
	return int(sgd.Score * 100)
}

const minCorrectSongScore = 0.9

func (sgd SongGuessDescription) Correct() bool {
	return sgd.CorrectedByContest || sgd.Score >= minCorrectSongScore
}

func (sg SongGuess) Description(correct SongInfo) SongGuessDescription {
	correctByContest := false
	if len(sg.Votes) > 0 {
		yes, unvoted := 0, 0
		for _, vote := range sg.Votes {
			if !vote.Voted {
				unvoted++
			} else if vote.Vote {
				yes++
			}
		}
		_ = unvoted
		correctByContest = yes > len(sg.Votes)/2
	}
	score := calcSimilarity(correct.Title, sg.Guess)
	return SongGuessDescription{
		sg,
		score,
		correctByContest,
	}
}

// Semantically modifies the receiver, even though technically we don't need
// it to be a pointer.
func (sg *SongGuess) ContestVote(other *Player, valid bool) error {
	if sg.Player == other {
		return fmt.Errorf("ContestVote: invalid API use: cannot call ContestVote on own SongGuess, please call via PlayerContest")
	}
	if !sg.IsContested() {
		return fmt.Errorf("ContestVote: cannot vote because %s hasn't contested %q yet", sg.Player.Name, sg.Guess)
	}
	for i := range sg.Votes {
		x := &sg.Votes[i]
		if x.Player == other {
			x.Vote = valid
			x.Voted = true
			return nil
		}
	}
	return fmt.Errorf("ContestVote: %s can't vote on %s's %q: %[1]s not found in Votes list", other.Name, sg.Player.Name, sg.Guess)
}

func (sg SongGuess) ContestedString() string {
	if !sg.IsContested() {
		return "(not contested)"
	}
	yay, nay, na := 0, 0, 0
	for _, v := range sg.Votes {
		if !v.Voted {
			na++
		} else if v.Vote {
			yay++
		} else {
			nay++
		}
	}
	if na == 0 {
		return fmt.Sprintf("%d yes %d no", yay, nay)
	}
	return fmt.Sprintf("%d yes %d no (%d undecided)", yay, nay, na)
}

func (sg SongGuess) IsContested() bool {
	return len(sg.Votes) > 0
}

func (sg SongGuess) NeedContestVotesFrom() []*Player {
	var ret []*Player
	for _, v := range sg.Votes {
		if !v.Voted {
			ret = append(ret, v.Player)
		}
	}
	return ret
}

func (sg SongGuess) FindContestedVote(other *Player) (PlayerVote, error) {
	if !sg.IsContested() {
		return PlayerVote{}, fmt.Errorf("%s's guess of %q is not contested, tell them to contest it", sg.Player.Name, sg.Guess)
	}
	for _, x := range sg.Votes {
		if x.Player == other {
			return x, nil
		}
	}
	return PlayerVote{}, nil
	// TODO: this is a bad failure mode, but leaving it blank isn't great either :/
	//return PlayerVote{}, fmt.Errorf("can't find contested vote: player %q not found; did you just join?", sg.Player.Name, sg.Guess, other.Name)
}

func (sg SongGuess) NeedsContestVote(other *Player) bool {
	pv, err := sg.FindContestedVote(other)
	if err != nil {
		// error doesn't matter: it's mostly signaling to the
		// template package.
		return false
	}
	return !pv.Voted
}

type PlayerRoundState struct {
	Guesses  []SongGuess
	Correct  bool
	Passed   bool
	BestTime time.Duration // only set if Correct
}

func (prs PlayerRoundState) StillGuessing() bool {
	return !prs.Correct && !prs.Passed
}

// Helps protect the IsZero implementation
func init() {
	if unsafe.Sizeof(GuessRound{}) != 12*unsafe.Sizeof((*int)(nil)) {
		panic("BAD size for GuessRound!")
	}
}

func (gr GuessRound) IsZero() bool {
	return gr.Song == SongInfo{} && gr.Guesses == nil && gr.Passes == nil
}

func (gr GuessRound) Clone() GuessRound {
	guesses := make([]SongGuess, len(gr.Guesses))
	for i, old := range gr.Guesses {
		guesses[i] = old.Clone()
	}
	return GuessRound{
		gr.Song.Clone(),
		guesses,
		slices.Clone(gr.Passes),
	}
}

func (gr GuessRound) OutstandingContests() []string {
	var counts map[*Player]int
	for _, g := range gr.Guesses {
		for _, p := range g.NeedContestVotesFrom() {
			if counts == nil {
				counts = make(map[*Player]int)
			}
			counts[p]++
		}
	}
	if len(counts) == 0 {
		return nil
	}
	ret := make([]string, 0, len(counts))
	for _, p := range slices.SortedFunc(maps.Keys(counts), func(a, b *Player) int {
		return cmp.Compare(a.Name, b.Name)
	}) {
		ret = append(ret, fmt.Sprintf("%v (%d missing votes)", p.Name, counts[p]))
	}
	return ret
}

func (gr GuessRound) PlayerGuesses(p *Player) []SongGuess {
	if gr.IsZero() {
		panic("GuessRound is zero value - no round found!")
	}
	var ret []SongGuess
	for _, g := range gr.Guesses {
		if g.Player == p {
			ret = append(ret, g)
		}
	}
	return ret
}

func (gr GuessRound) PlayerState(p *Player) PlayerRoundState {
	if gr.IsZero() {
		panic("GuessRound is zero value - no round found!")
	}
	guesses := gr.PlayerGuesses(p)
	var (
		bestTime time.Duration
		correct  bool
	)
	for _, guess := range guesses {
		if _, ok := guess.Correct(gr.Song); ok {
			bestTime, correct = guess.Delay, true
			break
		}
	}
	// Normally a user wouldn't pass after having a correct guess, but in
	// the face of a contest that can happen. Having a contested (and then
	// approved) guess essentially overrides any skip signal.
	passed := !correct && slices.Contains(gr.Passes, p)

	return PlayerRoundState{
		Guesses:  guesses,
		Correct:  correct,
		Passed:   passed,
		BestTime: bestTime,
	}
}

func (gr GuessRound) AllPlayers() []*Player {
	ret := slices.Clone(gr.Passes)
	for _, g := range gr.Guesses {
		// number of players will be small, and testing equality is
		// just a pointer check, so a set is wasteful.
		if slices.Contains(ret, g.Player) {
			continue
		}
		ret = append(ret, g.Player)
	}
	return ret
}

func (gr GuessRound) FinishedPlayers() []*Player {
	ret := slices.Clone(gr.Passes)
	for _, g := range gr.Guesses {
		// number of players will be small, so a set is wasteful.
		if slices.Contains(ret, g.Player) {
			continue
		}
		if _, correct := g.Correct(gr.Song); correct {
			ret = append(ret, g.Player)
		}
	}
	return ret
}

func (gr GuessRound) PlayerSummary(p *Player) []string {
	guesses := gr.PlayerGuesses(p)
	ret := make([]string, len(guesses), len(guesses)+1)
	for i, g := range guesses {
		desc := g.Description(gr.Song)
		gotit := ""
		if desc.Correct() {
			gotit = " - got it!"
		}
		ret[i] = fmt.Sprintf("%q (%v, %v%%%s)", g.Guess, g.Delay.Truncate(time.Millisecond), desc.ScorePercent(), gotit)
	}
	if slices.Contains(gr.Passes, p) {
		ret = append(ret, "PASSED")
	}
	return ret
}

// Acts on *GuessRound since this operation semantically mutates the round.
// Passing in allPlayers is a bit weird, but prevents needing to pass them
// around later.  plus, someone shouldn't get to just create a bunch of new
// players just to vote, after the fact at least.
func (gr *GuessRound) PlayerContest(p *Player, title string, allPlayers []*Player) error {
	guesses := gr.Guesses
	var guess *SongGuess
	for i := range guesses {
		g := &guesses[i]
		if g.Player == p && g.Guess == title {
			guess = g
			break
		}
	}
	if guess == nil {
		return fmt.Errorf("couldn't find %s's guess of %q", p.Name, title)
	}
	if guess.Votes != nil {
		return fmt.Errorf("can't contest %s's guess of %q: it is already contested", p.Name, title)
	}
	votes := make([]PlayerVote, len(allPlayers))
	for i, nextPlayer := range allPlayers {
		self := p == nextPlayer
		votes[i] = PlayerVote{
			Player: nextPlayer,
			Voted:  self,
			Vote:   self,
		}
	}
	guess.Votes = votes
	return nil
}

// Acts on *GuessRound since this operation semantically mutates the round.
func (gr *GuessRound) PlayerContestVote(other *Player, guessPlayerName, title string, shouldCount bool) error {
	guesses := gr.Guesses
	var guess *SongGuess
	for i := range guesses {
		g := &guesses[i]
		if g.Player.Name == guessPlayerName && g.Guess == title {
			guess = g
			break
		}
	}
	if guess == nil {
		return fmt.Errorf("couldn't find %s's guess of %q", guessPlayerName, title)
	}
	if !guess.IsContested() {
		return fmt.Errorf("can't vote on %s's guess of %q: it is not contested, tell them to contest it", guessPlayerName, title)
	}
	return guess.ContestVote(other, shouldCount)
}

type Scoreboard struct {
	Ranking []Scorerank
	Rounds  []ScoredRound
}

func (sc Scoreboard) PlayerTotal(p *Player) float64 {
	for _, r := range sc.Ranking {
		if r.Player == p {
			return r.Total
		}
	}
	return 0 //math.NaN()
}

type Scorerank struct {
	Player *Player
	Total  float64
}

// ScoredGuess is the result of a scored round for a single player's guess,
// recording the score relative to other player's guesses.
type PlayerScore struct {
	Player *Player
	// [0, 1], see ScoreRound for details.
	Score float64

	Delay time.Duration
}

// Used in the rendered html template; easier than calling Truncate in the
// template, which would require pushing things like parseDuration, or
// arbitrary constants, into the template.
func (ps PlayerScore) DelayString() string {
	return ps.Delay.Truncate(100 * time.Millisecond).String()
}

type ScoredRound struct {
	Round  GuessRound
	Scores []PlayerScore
}

func (sr ScoredRound) ScoresByPlayer() map[*Player]PlayerScore {
	ret := make(map[*Player]PlayerScore, len(sr.Scores))
	for _, score := range sr.Scores {
		ret[score.Player] = score
	}
	return ret
}

// Each round is scored according to the following:
//   - 0.0 is for people who tried to guess but gave up, or who haven't finished
//     yet (i.e. round isn't over).
//   - 0.25 is for people who didn't guess at all and skipped.
//   - 0.5 is the slowest guesser who got it right, and 1.0 is the fastest;
//     between that, there is a sliding scale of points based on relative delay
func ScoreRound(round GuessRound) ScoredRound {
	var scores []PlayerScore

	type pprs struct {
		p   *Player
		prs PlayerRoundState
	}
	var pprss []pprs
	// Not super efficient (iterates through all guesses multiple times),
	// but that's ok: n is small!
	for _, player := range round.FinishedPlayers() {
		state := round.PlayerState(player)
		if state.Passed {
			score := 0.25
			if len(state.Guesses) > 0 {
				score = 0
			}
			scores = append(scores, PlayerScore{player, score, 0})
			continue
		}
		if !state.Correct {
			// Not finished with round: 0 points!
			scores = append(scores, PlayerScore{player, 0, 0})
			continue
		}
		pprss = append(pprss, pprs{player, state})
	}
	switch len(pprss) {
	case 0:
		return ScoredRound{round, scores}
	case 1:
		return ScoredRound{round, append(scores, PlayerScore{pprss[0].p, 1, pprss[0].prs.BestTime})}
	}

	slices.SortFunc(pprss, func(a, b pprs) int {
		return cmp.Compare(a.prs.BestTime, b.prs.BestTime)
	})
	bestTime := pprss[0].prs.BestTime
	worstTime := pprss[len(pprss)-1].prs.BestTime
	scoreRange := (worstTime - bestTime).Seconds()
	for _, cur := range pprss {
		thisDelay := cur.prs.BestTime
		// ratio is [0, 1], where bestTime is 0 and worstTime is 1.
		ratio := (thisDelay - bestTime).Seconds() / scoreRange
		score := 1 - ratio/2
		scores = append(scores, PlayerScore{cur.p, score, thisDelay})
	}
	return ScoredRound{round, scores}
}

func ScoreRounds(rounds []GuessRound) Scoreboard {
	m := make(map[*Player][]float64)
	scoredRounds := make([]ScoredRound, len(rounds))
	for i, round := range rounds {
		scored := ScoreRound(round)
		scoredRounds[i] = scored
		for _, ps := range scored.Scores {
			m[ps.Player] = append(m[ps.Player], ps.Score)
		}
	}

	rows := make([]Scorerank, len(m))
	i := 0
	for player, scores := range m {
		sum := 0.
		for _, score := range scores {
			sum += score
		}
		rows[i] = Scorerank{player, sum}
		i++
	}
	slices.SortFunc(rows, func(a, b Scorerank) int {
		return -cmp.Compare(a.Total, b.Total)
	})
	return Scoreboard{rows, scoredRounds}
}

type multiguessGame struct {
	mu            sync.Mutex
	rounds        []GuessRound
	curSong       SongInfo
	lastTimestamp time.Time
}

func (mg *multiguessGame) Scoreboard() Scoreboard {
	mg.mu.Lock()
	defer mg.mu.Unlock()

	return ScoreRounds(mg.rounds)
}

// Rounds returns a deep copy of all current guess rounds.
func (mg *multiguessGame) Rounds() []GuessRound {
	mg.mu.Lock()
	defer mg.mu.Unlock()

	ret := make([]GuessRound, len(mg.rounds))
	for i, old := range mg.rounds {
		ret[i] = old.Clone()
	}
	return ret
}

// NumRounds returns the number of rounds currently being tracked.
func (mg *multiguessGame) NumRounds() int {
	mg.mu.Lock()
	defer mg.mu.Unlock()
	return len(mg.rounds)
}

func (mg *multiguessGame) CurrentRound() GuessRound {
	mg.mu.Lock()
	defer mg.mu.Unlock()

	if len(mg.rounds) == 0 || mg.curSong == (SongInfo{}) {
		return GuessRound{}
	}
	return mg.rounds[len(mg.rounds)-1]
}

func (mg *multiguessGame) PriorRound() GuessRound {
	mg.mu.Lock()
	defer mg.mu.Unlock()

	idx := len(mg.rounds) - 1
	if mg.curSong != (SongInfo{}) {
		idx--
	}
	if idx < 0 {
		return GuessRound{}
	}
	return mg.rounds[idx]
}

func (mg *multiguessGame) PlayerGuess(p *Player, title string) (gotIt bool, _ error) {
	mg.mu.Lock()
	defer mg.mu.Unlock()

	if len(mg.rounds) == 0 || mg.curSong == (SongInfo{}) {
		return false, fmt.Errorf("can't guess: no round found")
	}
	round := &mg.rounds[len(mg.rounds)-1]
	if slices.Contains(round.Passes, p) {
		return false, fmt.Errorf("can't guess: player already passed")
	}
	sg := SongGuess{p, time.Since(mg.lastTimestamp), title, nil}
	_, gotIt = sg.Correct(mg.curSong)
	round.Guesses = append(round.Guesses, sg)
	return gotIt, nil
}

func (mg *multiguessGame) PlayerPass(p *Player) error {
	mg.mu.Lock()
	defer mg.mu.Unlock()

	if len(mg.rounds) == 0 || mg.curSong == (SongInfo{}) {
		return fmt.Errorf("can't pass: no round found")
	}
	round := &mg.rounds[len(mg.rounds)-1]
	if slices.Contains(round.Passes, p) {
		return fmt.Errorf("can't pass: player already passed")
	}
	round.Passes = append(round.Passes, p)
	return nil
}

func (mg *multiguessGame) PlayerContest(p *Player, title string, allPlayers []*Player) error {
	mg.mu.Lock()
	defer mg.mu.Unlock()

	if len(mg.rounds) == 0 {
		return fmt.Errorf("can't contest %q: no round exists at all", title)
	}
	return mg.rounds[len(mg.rounds)-1].PlayerContest(p, title, allPlayers)
}

func (mg *multiguessGame) PlayerContestVote(other *Player, guessPlayerName, title string, shouldCount bool) error {
	mg.mu.Lock()
	defer mg.mu.Unlock()

	if len(mg.rounds) == 0 {
		return fmt.Errorf("can't contest %q: no round exists at all", title)
	}
	return mg.rounds[len(mg.rounds)-1].PlayerContestVote(other, guessPlayerName, title, shouldCount)
}

// Pass an empty SongInfo if you want to stop the round but not start a new one.
func (mg *multiguessGame) SetCurrentSong(info SongInfo) {
	mg.mu.Lock()
	defer mg.mu.Unlock()

	mg.curSong = info
	mg.lastTimestamp = time.Now()
	if info != (SongInfo{}) {
		mg.rounds = append(mg.rounds, GuessRound{Song: info})
	}
}

func (mg *multiguessGame) ResetTimer() {
	mg.mu.Lock()
	defer mg.mu.Unlock()

	mg.lastTimestamp = time.Now()
}

var parentheticalsRegexp = regexp.MustCompile(`([[(].+[])])+`)

func normalize(title string) (normalized string) {
	return strings.ToLower(strings.TrimSpace(title))
}

func calcSimilarity(target, guess string) (score float64) {
	originalTarget := target
	originalGuess := guess

	// Normalize both strings for a fair comparison
	target = normalize(parentheticalsRegexp.ReplaceAllLiteralString(target, ""))
	guess = normalize(parentheticalsRegexp.ReplaceAllLiteralString(guess, ""))

	// Otherwise, check similarity
	dist := editDistance(target, guess)
	diff := float64(dist) / float64(utf8.RuneCountInString(target))
	score = max(0.0, 1.0 - diff)

	// Give back some bonus points if they matched the parentheticals
	targetParen := normalize(parentheticalsRegexp.FindString(originalTarget))
	if len(targetParen) > 0 {
		guessParen := normalize(parentheticalsRegexp.FindString(originalGuess))
		if targetParen == guessParen {
			// Your bonus points are proportional to half the length of the parenthetical
			// (e.g. perfectly guessing an equal-length title and parenthetical gives you
			// a score of 1.5)
			parenLen := utf8.RuneCountInString(targetParen) -
				strings.Count(targetParen, "(") -
				strings.Count(targetParen, ")") -
				strings.Count(targetParen, "[") -
				strings.Count(targetParen, "]")
			score += float64(parenLen) / float64(len(target)) / 2.0
		}
	}

	return score
}

func editDistance(targetStr string, guessStr string) (distance int) {
	target := []rune(targetStr)
	guess := []rune(guessStr)

	// Quick exit: Empty strings are easy
	if len(guess) == 0 {
		return len(target)
	} else if len(target) == 0 {
		return len(guess)
	}

	editDistances := make([][]int, len(target) + 1)
	for i := range editDistances {
		editDistances[i] = make([]int, len(guess) + 1)
	}

	// Distance from target[0:i] to "" is i insertions
	for i := range len(target) + 1 {
		editDistances[i][0] = i
	}

	// Distance from "" to guess[0:i] is i deletions
	for i := range len(guess) + 1 {
		editDistances[0][i] = i
	}

	// Distance from target[0:i] to guess[0:j] is:
	//  - ED(i-1, j-1) if target[i] == target[j]
	//  - Otherwise, minimum of:
	//    * Inserting target[i]: ED(i-1, j) + 1
	//    * Deleting guess[j]:   ED(i, j-1) + 1
	//    * Editing guess[j]:    ED(i-1, j-1) + 1
	for i := 1; i <= len(target); i++ {
		for j := 1; j <= len(guess); j++ {
			if target[i-1] == guess[j-1] {
				editDistances[i][j] = editDistances[i-1][j-1]
			} else {
				insertion := 1 + editDistances[i-1][j]
				deletion := 1 + editDistances[i][j-1]
				edition := 1 + editDistances[i-1][j-1]
				editDistances[i][j] = min(insertion, min(deletion, edition))
			}
		}
	}

	return editDistances[len(target)][len(guess)]
}
