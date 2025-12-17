// tymgui exposes a global interface to control the website of YouTube Music
// via keyboard controls. Thus it requires you to be logged in and ready to
// play music.
package ytmgui

import (
	"chowski3/common/oslevelinput"
	"chowski3/games/musicquiz/mqgame/data"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func New(wr oslevelinput.Writer) *Player {
	return &Player{wr}
}

// TODO: state tracking to help prevent incorrect app action. However, we have
// to carefully consider how easy it is for us to get out of sync.
type Player struct {
	wr oslevelinput.Writer
}

func (p *Player) keypress(code oslevelinput.EventCode) {
	p.wr.Keypress(code)
}

func (p *Player) Play() error {
	p.keypress(oslevelinput.KEY_SPACE)
	return nil // TODO: Can't detect errors, but I don't care
}

func (p *Player) Pause() error {
	p.keypress(oslevelinput.KEY_SPACE)
	return nil // TODO: Can't detect errors, but I don't care
}

func (p *Player) ChangeSong(song data.SongInfo) error {
	return fmt.Errorf("Cannot call this ChangeSong method on a keyboard/mouse-based player")
}

// Doesn't pause first.
func (p *Player) Seek(deltaSeconds int) {
	// Forwards:
	for deltaSeconds >= 10 {
		p.keypress(oslevelinput.KEY_L)
		deltaSeconds -= 10
	}
	for deltaSeconds > 0 {
		release := p.wr.Hold(oslevelinput.KEY_LEFTSHIFT)
		p.keypress(oslevelinput.KEY_L)
		release()
		deltaSeconds--
	}
	if deltaSeconds == 0 {
		return
	}
	// Backwards!
	for deltaSeconds <= -10 {
		p.keypress(oslevelinput.KEY_H)
		deltaSeconds += 10
	}
	for deltaSeconds < 0 {
		release := p.wr.Hold(oslevelinput.KEY_LEFTSHIFT)
		p.keypress(oslevelinput.KEY_H)
		release()
		deltaSeconds++
	}
	if deltaSeconds != 0 {
		panic("BUG")
	}
}

func (p *Player) NextSong(ctx context.Context) (data.SongInfo, error) {
	lastSong, err := p.SongInfo()
	if err != nil && !errors.Is(err, ErrNoSongInTitle) {
		return data.SongInfo{}, fmt.Errorf("NextSong: looking for prior song: %w", err)
	}
	p.keypress(oslevelinput.KEY_J)
	if song, err := p.waitForSongToNotBe(ctx, lastSong); err != nil {
		return data.SongInfo{}, fmt.Errorf("NextSong: %w", err)
	} else {
		return song, nil
	}
}

var ErrNoSongInTitle = errors.New("there were no firefox windows with YouTube Music playing anything")

// Returns zero value's data.SongInfo on error.
func (p *Player) SongInfo() (ret data.SongInfo, retErr error) {
	//defer func() {
	//	fmt.Printf("data.SongInfo -> (%v, %v)\n", ret, retErr)
	//}()
	out, err := exec.Command("xdotool", "search", "--name", "--classname", "--class", "firefox").CombinedOutput()
	if err != nil {
		return data.SongInfo{}, fmt.Errorf("data.SongInfo: xdotool search: %v; output:\n%s", err, out)
	}
	lines := strings.Split(string(out), "\n")
	//log.Printf("%d firefox windows: %q", len(lines), lines)
	for _, line := range lines {
		if line == "" {
			continue
		}
		cmd := exec.Command("xdotool", "getwindowname", string(line))
		title, err := cmd.CombinedOutput()
		//fmt.Printf("%s -> (%q, %v)\n", cmd, title, err)
		if err != nil {
			if strings.Contains(string(title), "BadWindow (invalid Window parameter)") {
				//log.Printf("Skip bad window %v", line)
				continue
			}

			return data.SongInfo{}, fmt.Errorf("%q: %v; output=%q", cmd, err, title)
		}
		if name, ok := strings.CutSuffix(string(title), " - YouTube Music — Mozilla Firefox\n"); ok {
			if name == "" {
				// debugging... seems unlikely but who knows
				return data.SongInfo{}, fmt.Errorf("weird: empty song title: %q", string(title))
			}
			// TODO: check other lines for accidental duplicates?
			return data.SongInfo{Title: name}, nil
		}
		//log.Printf("Ignore non-song title: %q", title)
	}
	return data.SongInfo{}, ErrNoSongInTitle
}

// wait may be valid even if there is an error.
// TODO: a way to force-kill firefox, if needed?
func (p *Player) OpenPlaylist(ctx context.Context, url string) (wait func() error, _ error) {
	// Could do xdg-open. Probably should.
	firefoxCmd := exec.Command("firefox", url)
	if err := firefoxCmd.Start(); err != nil {
		return nil, fmt.Errorf("OpenPlaylist: failed firefox fork: %w", err)
	}

	for range time.Tick(time.Second) {
		if ctx.Err() != nil {
			return firefoxCmd.Wait, fmt.Errorf("OpenPlaylist: while waiting for firefox to start: %w", ctx.Err())
		}

		cmd := exec.CommandContext(ctx, "bash")
		cmd.Stdin = strings.NewReader("xdotool search --name --classname --class firefox  | xargs -l xdotool getwindowname | grep -- 'YouTube Music — Mozilla Firefox'")
		stdout, stderr := new(strings.Builder), new(strings.Builder)
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		if err := cmd.Run(); err != nil {
			log.Printf("OpenPlaylist: run bash hackiness: err=%v; stderr:\n%s", err, stderr)
			log.Printf("retry in a bit...")
			continue
		}
		if stdout.Len() != 0 {
			log.Printf("OpenPlaylist: stdout: %v\n\nstderr: %v", stdout, stderr)
			return firefoxCmd.Wait, nil
		}
	}
	panic("unreachable")
}

// Workaround for not being able to use the keyboard to hit play when a playlist is first opened.
func (p *Player) UseMouseToHitPlay(mouse *oslevelinput.MouseWriter) {
	mouse.SetCursorLocation(250, 325)
	mouse.Click(oslevelinput.BTN_LEFT)
}

// TODO: allow the playlist to be blank, and if so then just scan for the
// existing window. We should keep the window ID handy instead of repeatedly
// querying for all windows, since that's a waste and makes things
// slower/heaviweight.  But, ideally, we'd read the inside state of the browser
// so that we can get album and artist info too.
func (p *Player) InitialSetup(ctx context.Context, mouse *oslevelinput.MouseWriter, playlist string) error {
	log.Print("Opening playlist...")
	wait, err := p.OpenPlaylist(ctx, playlist)
	if err != nil {
		return err
	}
	go func() {
		err := wait()
		log.Fatalf("OpenPlaylist: wait returned error: %v", err)
	}()
	log.Print("Playlist opened!")
	// TODO: do a better job with this; right now we just guess that it takes ~forever.
	log.Printf("EXTRA WAIT TIME FOR SLOWDING")
	time.Sleep(30 * time.Second)

	log.Printf("CLICK WORKAROUND")
	p.UseMouseToHitPlay(mouse)

	if _, err := p.waitForSongToNotBe(ctx, data.SongInfo{}); err != nil {
		return fmt.Errorf("InitialSetup while waiting for first song title to show up: %w", err)
	}

	log.Println("Wait again because play music is slow...")
	time.Sleep(5 * time.Second)
	// pause (have to use keyboard because the status bar now covers the playlist play button...
	p.keypress(oslevelinput.KEY_SPACE)
	// shuffle
	p.keypress(oslevelinput.KEY_S)
	// good to go!
	return nil
}

func (p *Player) waitForSongToNotBe(ctx context.Context, not data.SongInfo) (ret data.SongInfo, retErr error) {
	//fmt.Printf("waitForSongToNotBe(%v)...\n", not)
	//defer func() {
	//	fmt.Printf("waitForSongToNotBe(%v) -> (%v, %v)\n", not, ret,retErr)
	//}()
	for ticker := time.Tick(100 * time.Millisecond); ; {
		select {
		case <-ctx.Done():
			return data.SongInfo{}, fmt.Errorf("waitForSongToNotBe: %w", ctx.Err())
		case <-ticker:
		}
		if info, err := p.SongInfo(); err == nil && info != not {
			return info, nil
		} else if err == nil || errors.Is(err, ErrNoSongInTitle) {
			continue
		} else {
			return data.SongInfo{}, fmt.Errorf("waitForSongToNotBe: %w", err)
		}
	}
}

func (p *Player) InitGui(server *http.ServeMux, addr string, browserCommand *string) {}
