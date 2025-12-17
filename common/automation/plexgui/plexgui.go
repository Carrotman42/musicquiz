package plexgui

import (
	"chowski3/common/apiclients/plexclient"
	"chowski3/games/musicquiz/mqgame/data"
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/net/websocket"
)

type PlexPlayer struct {
	// Initialized, authenticated Plex API
	plex      *plexclient.PlexConnection
	tracks    []plexclient.Track
	nextTrack int

	// Websocket state
	hostConnected atomic.Bool
	commands      chan string
	errors        chan string
}

// Error codes for websocketing
const ok = "OK"
const wait = "WAIT"

func NewPlexPlayer(playlist string) (*PlexPlayer, error) {
	// Connect to plex
	plex, err := plexclient.Connect()
	if err != nil {
		fmt.Println("Error connecting to Plex server:", err)
		return nil, err
	}

	// Find the tracks to run the quiz on
	var tracks []plexclient.Track
	playlists, err := plex.ListPlaylists()
	if err != nil {
		fmt.Println("Error listing Plex playlists:", err)
		return nil, err
	}
	for i := range len(*playlists) {
		if strings.Contains((*playlists)[i].Name, playlist) {
			t, err := plex.ListPlaylistContents((*playlists)[i])
			if err != nil {
				fmt.Println("Error listing playlist contents:", err)
				return nil, err
			}
			tracks = *t
			break
		}
	}
	rand.Shuffle(len(tracks), func(i, j int) {
		tracks[i], tracks[j] = tracks[j], tracks[i]
	})

	return &PlexPlayer{
		plex:          plex,
		tracks:        tracks,
		nextTrack:     0,
		hostConnected: atomic.Bool{},
		commands:      make(chan string),
		errors:        make(chan string),
	}, nil
}

func (player *PlexPlayer) InitGui(server *http.ServeMux, addr string, browserCommand *string) {
	// Serve the HTML/JS for the host page
	server.HandleFunc("/host", func(wr http.ResponseWriter, req *http.Request) {
		hostTemplate.Execute(wr, "unused initial data")
	})

	// Create a middleman interface that shuttles actions between the YTM iframe and the rest of the app/game
	server.Handle("/host/comms", websocket.Handler(player.connectHost))

	// Launch the first host page
	url := "http://" + addr + "/host"
	log.Printf("Launching browser at '%s'...", url)
	cmd := exec.Command("/bin/sh", "-c", *browserCommand+" "+url)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Start(); err != nil {
		log.Fatalf("error opening host page in a browser: %v\n", err)
	}
	log.Printf("Launched: %d", cmd.Process.Pid)
}

// TODO: This is buggy, refreshing the host page breaks the game
func (player *PlexPlayer) connectHost(conn *websocket.Conn) {
	// Wait until we're the only host
	for !player.hostConnected.CompareAndSwap(false, true) {
		conn.Write([]byte(wait))
		time.Sleep(5 * time.Second)
	}

	// Give up our spot as host if/when websocket closes for any reason
	defer player.hostConnected.Store(false)

	// We're the only host! Pass messages back and forth until the connection is
	// closed.
	for {
		cmd := <-player.commands
		_, err := conn.Write([]byte(cmd))
		if err != nil {
			log.Printf("failed to read websocket: %v\n", err)
			player.errors <- fmt.Sprintf("%v", err)
			break
		}

		data := make([]byte, 1000)
		n, err := conn.Read(data)
		if err != nil {
			log.Printf("failed to read websocket: %v\n", err)
			player.errors <- fmt.Sprintf("%v", err)
			break
		}
		resp := string(data[:n])
		if resp != ok {
			log.Printf("javascript error: %s\n", resp)
		}
		player.errors <- resp
	}

	// Websocket is closed, give up our spot as host (errors have already been sent back)
	defer player.hostConnected.Store(false)
}

func (plex *PlexPlayer) checkError() error {
	err := <-plex.errors
	if err != ok {
		return fmt.Errorf(err)
	}
	return nil
}

func (plex *PlexPlayer) Play() error {
	log.Printf("Sending 'play' to websocket")
	plex.commands <- "play"
	log.Printf("Sent 'play' to websocket, waiting for return")
	err := plex.checkError()
	log.Printf("Got return %v from websocket", err)
	return err

}

func (plex *PlexPlayer) Pause() error {
	plex.commands <- "pause"
	return plex.checkError()
}

func (plex *PlexPlayer) NextSong(context.Context) (data.SongInfo, error) {
	track := plex.tracks[plex.nextTrack]
	// (loop back to the beginning if we hit the end of the track list, just to make sure it never crashes?)
	plex.nextTrack = (plex.nextTrack + 1) % len(plex.tracks)

	cmd := "next:" + plex.plex.GetStreamUrl(&track)
	log.Println("Sending 'next song' command:", cmd)
	plex.commands <- cmd

	return data.SongInfo{
		Title:  track.Title,
		Album:  track.Album,
		Artist: track.Artist,
	}, plex.checkError()
}

func (plex *PlexPlayer) ChangeSong(song data.SongInfo) error {
	return fmt.Errorf("Cannot call NextSong method on a Plex player")
}

//go:embed plexgui.html
var hostTemplateStr string
var hostTemplate = template.Must(template.New("hostTemplate").Parse(hostTemplateStr))
