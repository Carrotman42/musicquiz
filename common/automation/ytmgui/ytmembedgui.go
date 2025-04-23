// tymgui exposes a global interface to control the website of YouTube Music
// via keyboard controls. Thus it requires you to be logged in and ready to
// play music.
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

func InitEmbedGui(server *http.ServeMux, mp *EmbedMusicPlayer, addr string, browserCommand *string) {
	// Serve the HTML/JS for the host page
	server.HandleFunc("/host", func(wr http.ResponseWriter, req *http.Request) {
		hostTemplate.Execute(wr, "unused initial data")
	})

	// Create a middleman interface that shuttles actions between the YTM iframe and the rest of the app/game
	server.Handle("/host/comms", websocket.Handler(mp.connectHost))
	// server.HandleFunc("/host/wait/nextYtmAction", func(wr http.ResponseWriter, req *http.Request) {
	// 	action := <-mp.ch
	// 	_, err := wr.Write([]byte(action))
	// 	if err != nil {
	// 		http.Error(wr, "error writing action string", 500)
	// 	}
	// })

	// Launch the first host page
	log.Printf("Launching browser at 'http://%s/host'...", addr)
	cmd := exec.Command("/bin/sh", "-c", *browserCommand + " http://" + addr + "/host")
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
