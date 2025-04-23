package data

import (
	"fmt"
)

type Player struct {
	Name string

	// NOTE: *Player values are marshaled directly (via GuessRound et.
	// al.), so don't put any interesting in here without fixing that.
}

// Call [Clone] to copy.
type SongInfo struct {
	Title   string
	Artist  string
	Album   string
	VideoID string
}

func (si SongInfo) Clone() SongInfo {
	// Right now a shallow copy is sufficient, but I don't want to promise
	// that API forever - seems like we might have more structure later.
	return si
}

func (si SongInfo) String() string {
	if si.Artist == "" && si.Album == "" {
		// Sadly this is what we have all the time right now.
		return si.Title
	}
	if si.Artist == "" {
		return fmt.Sprintf("%s by %s", si.Title, si.Artist)
	}
	return fmt.Sprintf("%s by %s on %s", si.Title, si.Artist, si.Album)
}
