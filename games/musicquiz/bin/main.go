// TODO: Prevent BeginRound at API level, not just UI level.
// Keep track of youtube music window, and requery more quickly.
//
// Different confetti pieces.
// better fuzzy matching.
//
//	ignore "feat" parentheticals.
package main

import (
	"chowski3/common/automation/ytmgui"
	"chowski3/common/oslevelinput"
	"chowski3/games/musicquiz/mqgame"
	"chowski3/games/musicquiz/mqhttpui"
	"context"
	"errors"
	"flag"
	"log"
	"maps"
	"os"
	"slices"
	"time"
)

var (
	playlistURL   = flag.String("playlist-url", "https://music.youtube.com/playlist?list=PLN7UIWiGb1K5h8_iuyoKhmD-N9z0Jo3lG", "full URL for playlist to play")
	skipInit      = flag.Bool("skip-init", false, "Assume that YouTube Music is in the foreground and ready to go; causes other various flags to be ignored")
	persistFile   = flag.String("persist-file", "", "Restore from this file (if it exists), as well as store state to this file regularly")
	persistPeriod = flag.Duration("persist-period", 30*time.Second, "Period between persisting state to --persist-file")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	wrs, err := oslevelinput.OpenAllWrite()
	if err != nil {
		log.Fatal(err)
	}
	var keyboard oslevelinput.Writer
	if len(wrs) == 1 {
		for _, keyboard = range wrs {
		}
	} else {
		if len(flag.Args()) != 1 {
			log.Fatalf("WHICH TO PICK? %q", slices.Sorted(maps.Keys(wrs)))
		}
		var ok bool
		if keyboard, ok = wrs[flag.Args()[0]]; !ok {
			log.Fatalf("%q not found; options: %v", flag.Args()[0], slices.Sorted(maps.Keys(wrs)))
		}
	}
	mouse, err := oslevelinput.OpenMouseWriter("/dev/input/event0")
	if err != nil {
		log.Fatal("opening mouse: ", err)
	}

	ytm := ytmgui.New(keyboard)

	if *skipInit {
		log.Printf("Skipping initialization; please have ytm open in its own window, in the foreground, with the playlist ready to go; don't forget shuffle!")
	} else {
		setupCtx, cf := context.WithTimeout(ctx, 2*time.Minute)
		defer cf()
		if err := ytm.InitialSetup(setupCtx, mouse, *playlistURL); err != nil {
			log.Fatal(err)
		}
	}

	var state *mqgame.State
	if *persistFile == "" {
		state = mqgame.NewState(ytm)
		goto ready
	}
	if data, err := os.ReadFile(*persistFile); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatalf("Failed reading %v: %v", *persistFile, err)
		}
		log.Printf("Creating fresh state (file %v had error %v)", *persistFile, err)
		state = mqgame.NewState(ytm)
	} else if state, err = mqgame.RestoreState(ytm, data); err != nil {
		log.Fatalf("Failed to restore from %v: %v", *persistFile, err)
	} else {
		log.Printf("Restored state from %v", *persistFile)
	}
	go func() {
		for range time.Tick(*persistPeriod) {
			txt, err := state.MarshalText()
			if err != nil {
				log.Printf("persist state: failed to marshal: %v", err)
				continue
			}
			if err := os.WriteFile(*persistFile, txt, 0666); err != nil {
				log.Printf("persist state: failed to WriteFile: %v", err)
				continue
			}
		}
	}()

ready:
	mqhttpui.Run(mqhttpui.MultiguessGame{}, state, "192.168.86.33:1123", "192.168.86.33")
}
