package mqhttpui

import (
	"bytes"
	"chowski3/common/ksuite/kjs"
	"chowski3/games/musicquiz/mqgame"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type UI interface {
	RenderPage(io.Writer, *mqgame.State, *mqgame.Player, *http.Request) error
}

func Run(ui UI, gamestate *mqgame.State, addr string, domain string) {
	s := server{ui, gamestate}
	// Non-logged-in handlers:
	http.HandleFunc("/login", s.loginHandler(domain))
	http.HandleFunc("/static/xhr.js", kjs.XHRHandler)

	// Logged-in handlers:
	http.HandleFunc("/", s.loggedInHandler(s.index))
	http.HandleFunc("/game", s.loggedInHandler(s.game))
	http.HandleFunc("/game/action/{action}", s.loggedInHandler(s.gameAction))
	http.HandleFunc("/game/wait/state", s.loggedInHandler(s.waitState))
	http.HandleFunc("/secret-backdoor", s.loggedInHandler(s.secretBackdoor))

	log.Printf("Serving on %v with cookie domain %v", addr, domain)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (s server) loggedInHandler(cb func(http.ResponseWriter, *http.Request, *mqgame.Player) error) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		p, err := s.getPlayer(req)
		if err != nil {
			fmt.Println("getPlayer failed: ", err)
			if err2 := playerNotLoggedInTemplate.Execute(resp, playerNotLoggedInData{err}); err2 != nil {
				log.Printf("Error rendering not-logged-in template: %v", err2)
			}
			return
		}
		if err := cb(resp, req, p); err != nil {
			log.Printf("Error for %q in path %q: %v", p.Name, req.URL.Path, err)
		}
	}
}

type server struct {
	UI
	*mqgame.State
}

const (
	playerNameCookie = "player-name-2"
)

func (s server) getPlayer(req *http.Request) (*mqgame.Player, error) {
	c, err := req.Cookie(playerNameCookie)
	if err != nil {
		return nil, fmt.Errorf("player not logged in: %w", err)
	}
	p, ok := s.Player(c.Value)
	if ok {
		return p, nil
	}
	if true {
		p, err := s.NewPlayer(c.Value)
		if err != nil {
			return nil, fmt.Errorf("problem auto-logging in %q (from cookie): %v", c.Value, err)
		}
		return p, nil
	}

	return nil, fmt.Errorf("player not logged in: %q (from cookie) not found in game", c.Value)
}

type playerNotLoggedInData struct {
	Err error
}

var playerNotLoggedInTemplate = template.Must(template.New("playerNotLoggedInTemplate").Parse(`<html>
<meta name="viewport" content="width=device-width, initial-scale=1">
<h1>Log in plz</h1>
<form method=POST action=/login>
	<label>Name: <input name=name></label><br>
	<input type=submit value=Yay!>
</form>
</html>`))

func (s server) loginHandler(domain string) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.Error(wr, "must POST bro", 401)
			return
		}
		if err := req.ParseForm(); err != nil {
			http.Error(wr, "Failed to parse form", 400)
			return
		}
		name := req.Form.Get("name")
		if name == "" {
			http.Error(wr, "must 'name' bro", 401)
			return
		}
		_, err := s.NewPlayer(name)
		if err != nil {
			http.Error(wr, fmt.Sprintf("Error setting name to %q: %v", name, err), 401)
			return
		}
		http.SetCookie(wr, &http.Cookie{
			//SameSite: http.SameSiteStrictMode,
			Path:    "/",
			Expires: time.Now().Add(300 * 24 * time.Hour),
			MaxAge:  300 * 24 * 60 * 60,
			Domain:  domain,
			Name:    playerNameCookie,
			Value:   name,
		})
		//if name == "kev" {
		//	http.Redirect(wr, req, "/secret-backdoor", http.StatusFound)
		//	return
		//}

		http.Redirect(wr, req, "/game", http.StatusFound)
	}
}

func (s server) index(wr http.ResponseWriter, req *http.Request, p *mqgame.Player) error {
	http.Redirect(wr, req, "/game", http.StatusFound)
	return nil
}

var secretBackdoorTmpl = template.Must(template.New("secretBackdoorTmpl").Parse(`<html>
<meta name="viewport" content="width=device-width, initial-scale=1">

{{ define "action" }}
<form method=post><input type=hidden name=action value="{{.}}"><input type=submit value="{{.}}"></form>
{{ end }}

{{ template "action" "play" }}
{{ template "action" "pause" }}
<hr>

<ul>
	{{ range .Players }}
		<li>
			{{ .Name }}
			<form method=post style="display: inline;">
				<input type=hidden name=action value="kick">
				<input type=hidden name=player_name value="{{.Name}}">
				<input type=submit value="Kick">
			</form>
		</li>
	{{ end }}
</ul>

</html>
`))

func (s server) secretBackdoor(wr http.ResponseWriter, req *http.Request, p *mqgame.Player) error {
	if p.Name != "kev" {
		http.Redirect(wr, req, "/", http.StatusFound)
		return nil
	}
	if req.Method != "POST" {
		goto renderPage
	}
	if err := req.ParseForm(); err != nil {
		http.Redirect(wr, req, "/", http.StatusFound)
		return err
	}
	{
		var err error
		switch act := req.Form.Get("action"); act {
		default:
			err = fmt.Errorf("unknown action %q", act)
		case "play":
			err = s.Play()
		case "pause":
			err = s.Pause()
		case "kick":
			err = s.KickPlayer(req.Form.Get("player_name"))
		}
		if err != nil {
			fmt.Fprintf(wr, "Error: %v", err)
			return fmt.Errorf("failed secret backdoor action: %v", err)
		}
	}

renderPage:
	buf := new(bytes.Buffer)
	if err := secretBackdoorTmpl.Execute(buf, s.State); err != nil {
		http.Error(wr, "Template error: "+err.Error(), 500)
		return err
	}
	if _, err := buf.WriteTo(wr); err != nil {
		return fmt.Errorf("buffer write: %w", err)
	}
	return nil
}

func (s server) game(wr http.ResponseWriter, req *http.Request, p *mqgame.Player) error {
	buf := new(bytes.Buffer)
	if err := s.UI.RenderPage(buf, s.State, p, req); err != nil {
		http.Error(wr, "Template error: "+err.Error(), 500)
		return err
	}
	if _, err := buf.WriteTo(wr); err != nil {
		return fmt.Errorf("buffer write: %w", err)
	}
	return nil
}

func (s server) gameAction(wr http.ResponseWriter, req *http.Request, p *mqgame.Player) error {
	if req.Method != "POST" {
		http.Error(wr, "must POST bro", http.StatusMethodNotAllowed)
		return fmt.Errorf("want POST for game action, got %q", req.Method)
	}
	if err := req.ParseForm(); err != nil {
		http.Error(wr, "Failed to parse form", 400)
		return err
	}
	var err error
actionSwitch:
	switch act := req.PathValue("action"); act {
	default:
		http.Error(wr, "unknown/missing action", 400)
		return fmt.Errorf("unknown/missing action in request (player=%v): action=%q", p.Name, act)
	case "pass":
		err = s.PlayerPass(p)
	case "guess":
		_, err = s.PlayerGuess(p, req.Form.Get("title"))
	case "begin-round":
		err = s.BeginRound(req.Context())
		if errors.Is(err, mqgame.ErrRoundAlreadyStarted) {
			// Assume race condition with another player;
			// just allow the redirect to /game happen below.
			err = nil
		}
	case "contest":
		err = s.PlayerContest(p, req.Form.Get("guessValue"))
	case "contest-vote":
		var valid bool
		switch vote := req.Form.Get("vote"); vote {
		case "true":
			valid = true
		case "false":
			valid = false
		default:
			err = fmt.Errorf("unknown vote %q", vote)
			break actionSwitch
		}
		err = s.PlayerContestVote(p, req.Form.Get("guessPlayer"), req.Form.Get("guessValue"), valid)
	}
	if err != nil {
		http.Error(wr, err.Error(), 400)
		return err
	}
	http.Redirect(wr, req, "/game", http.StatusFound)
	return nil
}

func (s server) waitState(wr http.ResponseWriter, req *http.Request, p *mqgame.Player) error {
	nextStr := req.URL.Query().Get("next_state")
	fmt.Printf("%v is waiting for round %v (currently %d)\n", p.Name, nextStr, s.StateCounter())
	next, err := strconv.Atoi(nextStr)
	if err != nil {
		return fmt.Errorf("invalid next_round %q: %v", nextStr, err)
	}
	ctx := req.Context()
	select {
	case <-ctx.Done():
		return fmt.Errorf("out of time waiting for round %v to start: %w", next, ctx.Err())
	case <-s.ChanForState(next):
		// Ready for next round!
		return nil
	}
}
