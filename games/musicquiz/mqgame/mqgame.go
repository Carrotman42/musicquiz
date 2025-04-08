package mqgame

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"
	"sync"
)

func NewState(mp MusicPlayer) *State {
	return &State{
		musicPlayer: mp,
		players:     make(map[string]*Player),
		gameState:   new(multiguessGame),
	}
}

type State struct {
	musicPlayer MusicPlayer
	gameState   *multiguessGame

	mu        sync.Mutex
	players   map[string]*Player
	songState SongState
	// incremented any time a modification to the state of the system was
	// maybe made; supports hanging polls for auto-refreshing the browser.
	// It may be incremented spuriously, but that'll just cause harmless
	// refreshes to users waiting around.
	stateCounter int
	stateChans   []stateChan
}

type MusicPlayer interface {
	Play()
	Pause()
	//Seek(deltaSeconds int)
	NextSong(context.Context) (SongInfo, error)
}

// TODO: properly marshal an ongoing round, right now it will (probably) be
// considered finished when restored.
func (s *State) MarshalText() ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	m := marshaledState{
		s.gameState.Rounds(),
	}
	return json.Marshal(m)
}

type marshaledState struct {
	Rounds []GuessRound
}

// RestoreState returns a new *State from the result of [State.MarshalText].
// TODO: support ongoing rounds, right now they are canceled I think.
func RestoreState(mp MusicPlayer, marshaled []byte) (*State, error) {
	var m marshaledState
	if err := json.Unmarshal(marshaled, &m); err != nil {
		return nil, fmt.Errorf("RestoreState: bad marshaled text: %w", err)
	}
	nm := make(normmap)
	for i := range m.Rounds {
		m.Rounds[i].normalizePlayers(nm)
	}
	return &State{
		musicPlayer: mp,
		players:     make(map[string]*Player),
		gameState: &multiguessGame{
			rounds: m.Rounds,
		},
	}, nil
}

type Player struct {
	Name string

	// NOTE: *Player values are marshaled directly (via GuessRound et.
	// al.), so don't put any interesting in here without fixing that.
}

func (s *State) Player(name string) (*Player, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.players[name]
	return p, ok
}

func (s *State) NewPlayer(name string) (*Player, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.players[name]; ok {
		return nil, fmt.Errorf("player %q already in the game!", name)
	}
	p := &Player{
		Name: name,
	}
	s.players[name] = p
	s.lockedIncStateCounter()
	return p, nil
}

func (s *State) KickPlayer(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.players[name]; !ok {
		return fmt.Errorf("KickPlayer: player %q not found in game", name)
	}

	delete(s.players, name)
	s.lockedCheckRoundDone()
	return nil
}

func (s *State) Players() []*Player {
	s.mu.Lock()
	defer s.mu.Unlock()

	return slices.SortedFunc(maps.Values(s.players), func(a, b *Player) int {
		return cmp.Compare(a.Name, b.Name)
	})
}

func (s *State) StateCounter() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stateCounter
}

type stateChan struct {
	atState int
	toClose chan struct{}
}

func (s *State) lockedIncStateCounter() {
	s.stateCounter++
	cur := s.stateCounter
	for len(s.stateChans) > 0 {
		sc := s.stateChans[0]
		if sc.atState < cur {
			break
		}
		close(sc.toClose)
		s.stateChans = s.stateChans[1:]
	}
}

func (s *State) ChanForState(closeAtCounter int) <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stateCounter >= closeAtCounter {
		ch := make(chan struct{})
		close(ch)
		return ch
	}

	pos, ok := slices.BinarySearchFunc(s.stateChans, closeAtCounter, func(cur stateChan, target int) int {
		return cmp.Compare(cur.atState, target)
	})
	if ok {
		return s.stateChans[pos].toClose
	}
	ch := make(chan struct{})
	s.stateChans = slices.Insert(s.stateChans, pos, stateChan{closeAtCounter, ch})
	return ch
}

type SongState struct {
	Playing bool
}

func (s *State) SongState() SongState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.songState
}

// Maybe should expose a controller object or whatever instead of manually
// forwarding all these methods.

func (s *State) NumRounds() int {
	return s.gameState.NumRounds()
}

func (s *State) CurrentRound() GuessRound {
	return s.gameState.CurrentRound()
}

func (s *State) PriorRound() GuessRound {
	return s.gameState.PriorRound()
}

func (s *State) Scoreboard() Scoreboard {
	return s.gameState.Scoreboard()
}

// Rounds returns a deep copy of all the rounds.
func (s *State) Rounds() []GuessRound {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.gameState.Rounds()
}

func (s *State) PlayerGuess(p *Player, title string) (gotIt bool, _ error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// increment for good measure, even if there's an error below.
	s.lockedIncStateCounter()

	gotIt, err := s.gameState.PlayerGuess(p, title)
	if err == nil {
		s.lockedCheckRoundDone()
	}
	return gotIt, err
}

func (s *State) PlayerPass(p *Player) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// increment for good measure, even if there's an error below.
	s.lockedIncStateCounter()

	err := s.gameState.PlayerPass(p)
	if err == nil {
		s.lockedCheckRoundDone()
	}
	return err
}

func (s *State) lockedCheckRoundDone() {
	// Because a player may have been kicked, we have to actually check subsets.
	allPlayers := make(map[*Player]bool, len(s.players))
	for _, p := range s.players {
		allPlayers[p] = true
	}
	for _, p := range s.gameState.CurrentRound().FinishedPlayers() {
		delete(allPlayers, p)
	}
	if len(allPlayers) == 0 {
		s.gameState.SetCurrentSong(SongInfo{})
		s.lockedIncStateCounter()
	}
}

func (s *State) PlayerContest(p *Player, title string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// increment for good measure, even if there's an error below.
	s.lockedIncStateCounter()
	allPlayers := slices.Collect(maps.Values(s.players))
	return s.gameState.PlayerContest(p, title, allPlayers)
}

func (s *State) PlayerContestVote(other *Player, guessPlayerName, title string, shouldCount bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// increment for good measure, even if there's an error below.
	s.lockedIncStateCounter()
	return s.gameState.PlayerContestVote(other, guessPlayerName, title, shouldCount)
}

var ErrRoundAlreadyStarted = errors.New("can't start next round: prior round is still ongoing!")

func (s *State) BeginRound(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// increment for good measure, even if there's an error below.
	// But do it when returning, not now: doing it now will immediately
	// cause waiters to refresh, including the calling this method - which
	// causes this method to get cancelled before the player has a chance
	// to run.
	defer s.lockedIncStateCounter()

	if !s.gameState.CurrentRound().IsZero() {
		return ErrRoundAlreadyStarted
	}

	nextSong, err := s.musicPlayer.NextSong(ctx)
	if err != nil {
		return fmt.Errorf("BeginRound: %w", err)
	}
	s.gameState.SetCurrentSong(nextSong)
	return nil
}

var (
	ErrAlreadyPlaying = errors.New("already playing")
	ErrAlreadyPaused  = errors.New("already paused")
)

// SECRET BACKDOOR FUNCTION, should delete
func (s *State) Play() error {
	//return fmt.Errorf("nah")
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.songState.Playing {
		return ErrAlreadyPlaying
	}
	if s.musicPlayer != nil {
		s.musicPlayer.Play()
	}
	s.songState.Playing = true
	return nil
}

// SECRET BACKDOOR FUNCTION, should delete
func (s *State) Pause() error {
	//return fmt.Errorf("nah")
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.songState.Playing {
		return ErrAlreadyPaused
	}
	if s.musicPlayer != nil {
		s.musicPlayer.Pause()
	}
	s.songState.Playing = false
	return nil
}
