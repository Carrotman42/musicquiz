// ytmembedgui exposes a global interface to control the website of YouTube Music
// via an embedded iframe player.
//
// ===== WARNING =====
// Unfortunately, this doesn't work on most videos -- many monetized videos or music
// videos are blocked from being embedded, or 'blocked in your country', such that
// while this embedded player may *technically* work, it doesn't work in practice
// to allow any music quiz that you'd actually want to play (you could guess at white
// noise, I suppose).
// ===== WARNING =====

package ytmgui

import (
	"chowski3/games/musicquiz/mqgame/data"
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync/atomic"
	"time"

	"golang.org/x/net/websocket"
)

// Error codes for websocketing
const ok = "OK"
const wait = "WAIT"

type EmbedMusicPlayer struct {
	hostConnected atomic.Bool
	commands      chan string
	errors        chan string
}

func NewEmbedPlayer() *EmbedMusicPlayer {
	return &EmbedMusicPlayer{
		commands: make(chan string),
		errors:   make(chan string),
	}
}

func (mp *EmbedMusicPlayer) InitGui(server *http.ServeMux, addr string, browserCommand *string) {
	// Serve the HTML/JS for the host page
	server.HandleFunc("/host", func(wr http.ResponseWriter, req *http.Request) {
		hostTemplate.Execute(wr, "unused initial data")
	})

	// Create a middleman interface that shuttles actions between the YTM iframe and the rest of the app/game
	server.Handle("/host/comms", websocket.Handler(mp.connectHost))

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

func (mp *EmbedMusicPlayer) checkError() error {
	err := <-mp.errors
	if err != ok {
		return fmt.Errorf(err)
	}
	return nil
}

func (mp *EmbedMusicPlayer) Play() error {
	log.Printf("Sending 'play' to websocket")
	mp.commands <- "play"
	log.Printf("Sent 'play' to websocket, waiting for return")
	err := mp.checkError()
	log.Printf("Got return %v from websocket", err)
	return err

}
func (mp *EmbedMusicPlayer) Pause() error {
	mp.commands <- "pause"
	return mp.checkError()
}
func (mp *EmbedMusicPlayer) NextSong(context.Context) (data.SongInfo, error) {
	return data.SongInfo{}, fmt.Errorf("Cannot call NextSong method on an embed player")
}
func (mp *EmbedMusicPlayer) ChangeSong(song data.SongInfo) error {
	mp.commands <- "next:" + song.VideoID
	return mp.checkError()
}

func (mp *EmbedMusicPlayer) connectHost(conn *websocket.Conn) {
	// Wait until we're the only host
	for !mp.hostConnected.CompareAndSwap(false, true) {
		conn.Write([]byte(wait))
		time.Sleep(5 * time.Second)
	}

	// We're the only host! Pass messages back and forth until the connection is
	// closed.
	for {
		cmd := <-mp.commands
		_, err := conn.Write([]byte(cmd))
		if err != nil {
			log.Printf("failed to read websocket: %v\n", err)
			mp.errors <- fmt.Sprintf("%v", err)
			break
		}

		data := make([]byte, 1000)
		n, err := conn.Read(data)
		if err != nil {
			log.Printf("failed to read websocket: %v\n", err)
			mp.errors <- fmt.Sprintf("%v", err)
			break
		}
		resp := string(data[:n])
		if resp != ok {
			log.Printf("javascript error: %s\n", resp)
		}
		mp.errors <- resp
	}

	// Websocket is closed, give up our spot as host (errors have already been sent back)
	mp.hostConnected.Store(false)
}

//go:embed ytmembedgui.html
var hostTemplateStr string
var hostTemplate = template.Must(template.New("hostTemplate").Parse(hostTemplateStr))
