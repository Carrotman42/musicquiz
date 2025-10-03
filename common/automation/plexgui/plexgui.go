package plexgui

import (
	"chowski3/common/apiclients/plexclient"
	"chowski3/games/musicquiz/mqgame/data"
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"log"
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
	plex   *plexclient.PlexConnection
	tracks []plexclient.Track

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

	return &PlexPlayer{
		plex:          plex,
		tracks:        tracks,
		hostConnected: atomic.Bool{},
		commands:      make(chan string),
		errors:        make(chan string),
	}, nil
}

func InitPlexGui(server *http.ServeMux, player *PlexPlayer, addr string, browserCommand *string) {
	// Serve the HTML/JS for the host page
	server.HandleFunc("/host", func(wr http.ResponseWriter, req *http.Request) {
		hostTemplate.Execute(wr, "unused initial data")
	})

	// Create a middleman interface that shuttles actions between the YTM iframe and the rest of the app/game
	server.Handle("/host/comms", websocket.Handler(player.connectHost))

	// Launch the first host page
	log.Printf("Launching browser at 'http://%s/host'...", addr)
	cmd := exec.Command("/bin/sh", "-c", *browserCommand+" http://"+addr+"/host")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Start(); err != nil {
		log.Fatalf("error opening host page in a browser: %v\n", err)
	}
	log.Printf("Launched: %d", cmd.Process.Pid)
}

func (player *PlexPlayer) connectHost(conn *websocket.Conn) {
	// Wait until we're the only host
	for !player.hostConnected.CompareAndSwap(false, true) {
		conn.Write([]byte(wait))
		time.Sleep(5 * time.Second)
	}

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
	player.hostConnected.Store(false)
}

// TODO DO NOT SUBMIT: implement commands
// TODO DO NOT SUBMIT: implement http page so that we can actually do command-y things

func (mp *PlexPlayer) checkError() error {
	err := <-mp.errors
	if err != ok {
		return fmt.Errorf(err)
	}
	return nil
}

func (mp *PlexPlayer) Play() error {
	log.Printf("Sending 'play' to websocket")
	mp.commands <- "play"
	log.Printf("Sent 'play' to websocket, waiting for return")
	err := mp.checkError()
	log.Printf("Got return %v from websocket", err)
	return err

}
func (mp *PlexPlayer) Pause() error {
	mp.commands <- "pause"
	return mp.checkError()
}
func (mp *PlexPlayer) NextSong(context.Context) (data.SongInfo, error) {
	return data.SongInfo{}, fmt.Errorf("Cannot call NextSong method on an embed player")
}
func (mp *PlexPlayer) ChangeSong(song data.SongInfo) error {
	mp.commands <- "next:" + song.VideoID
	return mp.checkError()
}

//go:embed plexgui.html
var hostTemplateStr string
var hostTemplate = template.Must(template.New("hostTemplate").Parse(hostTemplateStr))
