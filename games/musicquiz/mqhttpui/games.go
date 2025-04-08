package mqhttpui

import (
	"chowski3/games/musicquiz/mqgame"
	"fmt"
	"html/template"
	"io"
	"math/rand/v2"
	"net/http"
)

type MultiguessGame struct {
}

func (mg MultiguessGame) RenderPage(wr io.Writer, s *mqgame.State, p *mqgame.Player, req *http.Request) error {
	if req.URL.Query().Has("scoreboard") {
		return scoreboardTemplate.Execute(wr, scoreboardData{p, s})
	}
	cloned, err := gameTemplate.Clone()
	if err != nil {
		return fmt.Errorf("unexpected template.Clone failure: %w", err)
	}
	// TODO: is there a better way to curry this top-level value into
	// sub-templates like playersSummary?
	// TODO: is this clone "inefficient"? Probably don't care too much but
	// it's worth a benchmark for future reference.
	cloned = cloned.Funcs(template.FuncMap{
		"currentPlayer": func() *mqgame.Player { return p },
	})
	return cloned.Execute(wr, gameData{p, s})
}

var builtinTemplateFuncs = template.FuncMap{
	"currentPlayer": func() *mqgame.Player {
		panic("shouldn't be called - should be overridden before executing")
	},
	"plus1": func(x int) int { return x + 1 },
	"randConfettiChar": func() string {
		return confettiChars[rand.N(len(confettiChars))]
	},
}

var confettiChars = []string{
	"🥎",
	"🧶",
	"💎",
	"🔔",
	"🥁",
	"💿",
	"💡",
	"📈",
	"🪠",
	"✅",
	"🇺🇸",
	"🐵", "🐒", "🦍", "🦧", "🐶", "🐕", "🦮", "🐩", "🐺", "🦊", "🦝", "🐱", "🐈", "🦁", "🐯", "🐅", "🐆", "🐴", "🐎", "🦄", "🦓", "🦌", "🦬", "🐮", "🐂", "🐃", "🐄", "🐷", "🐖", "🐗", "🐽", "🐏", "🐑", "🐐", "🐪", "🐫", "🦙", "🦒", "🐘", "🦣", "🦏", "🦛", "🐭", "🐁", "🐀", "🐹", "🐰", "🐇", "🐿", "🦫", "🦔", "🦇", "🐻", "🐨", "🐼", "🦥", "🦦", "🦨", "🦘", "🦡", "🐾", "🦃", "🐔", "🐓", "🐣", "🐤", "🐥", "🐦", "🐧", "🕊", "🦅", "🦆", "🦢", "🦉", "🦤", "🪶", "🦩", "🦚", "🦜", "🐳", "🐋", "🐬", "🦭", "🐟", "🐠", "🐡", "🦈", "🐙", "🐚", "🐌", "🦋", "🐛", "🐜", "🐝", "🪲", "🐞", "🦗", "🪳", "🕷", "🕸", "🦂", "🦟", "🪰", "🪱", "🦠", "💐", "🌸", "💮", "🏵", "🌹", "🥀", "🌺", "🌻", "🌼", "🌷", "🍇", "🍈", "🍉", "🍊", "🍋", "🍌", "🍍", "🥭", "🍎", "🍏", "🍐", "🍑", "🍒", "🍓", "🫐", "🥝", "🍅", "🫒", "🥥", "🥑", "🍆", "🥔", "🥕", "🌽", "🌶", "🫑", "🥒", "🥬", "🥦", "🧄", "🧅", "🍄", "🥜", "🌰", "🍞", "🥐", "🥖", "🫓", "🥨", "🥯", "🥞", "🧇", "🧀", "🍖", "🍗", "🥩", "🥓", "🍔", "🍟", "🍕", "🌭", "🥪", "🌮", "🌯", "🫔", "🥙", "🧆", "🥚", "🍳", "🥘", "🍲", "🫕", "🥣", "🥗", "🍿", "🧈", "🧂", "🥫", "🦀", "🦞", "🦐", "🦑", "🦪", "🍦", "🍧", "🍨", "🍩", "🍪", "🎂", "🍰", "🧁", "🥧", "🍫", "🍬", "🍭", "🍮", "🍯", "🍼", "🥛", "☕", "🫖", "🍵", "🍶", "🍾", "🍷", "🍸", "🍹", "🍺", "🍻", "🥂", "🥃", "🥤", "🧋", "🧃", "🧉", "🧊", "🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘", "🌙", "🌚", "🌛", "🌜", "🌡", "☀", "🌝", "🌞", "🪐", "⭐", "🌟", "🌠", "🌌", "☁", "⛅", "⛈", "🌤", "🌥", "🌦", "🌧", "🌨", "🌩", "🌪", "🌫", "🌬", "🌀", "🌈", "🌂", "☂", "☔", "⛱", "⚡", "❄", "☃", "⛄", "☄", "🔥", "💧", "🌊", "🎃", "🎄", "🎆", "🎇", "🧨", "✨", "🎈", "🎉", "🎊", "🎋", "🎍", "🎎", "🎏", "🎐", "🎑", "🧧", "🎀", "🎁", "🎗", "🎟", "🎫", "🎼", "🎵", "🎶", "🎙", "🎚", "🎛", "🎤", "🎧", "📻",
}

type gameData struct {
	Player *mqgame.Player

	State *mqgame.State
}

var gameTemplate = template.Must(template.New("gameTemplate").Funcs(builtinTemplateFuncs).Parse(`<html>
{{ define "confetti" }}
	<script type="text/javascript" src="https://cdn.jsdelivr.net/npm/canvas-confetti@1.9.3/dist/confetti.browser.min.js"></script>
	<script type="text/javascript">
		window.addEventListener("load", function() {
			var scalar = 4;
			var pineapple = confetti.shapeFromText({ text: {{.}}, scalar });

			confetti({
				shapes: [pineapple],
				scalar
			});
		});
	</script>
{{ end }}
{{ define "playersSummary" }}
	{{ $round := . }}
	<ul>
		{{ range $round.AllPlayers }}
			{{ $player := . }}
			{{ $roundState := $round.PlayerState $player }}
			<li>
				{{ $player.Name}}
				<ol>
					{{ range $roundState.Guesses }}
						<li>
							{{ $desc := .Description $round.Song }} 
							{{/* set $contestedVotes and $color */}}
							{{ $contestedVotes := "" }}
							{{ $color := "" }}
							{{ if not $desc.IsContested }}
								{{ $contestedVotes = "" }}
								{{ if $desc.Correct }}
									{{ $color = "green" }}
								{{ else }}
									{{ $color = "red" }}
								{{ end }}
							{{ else }}
								{{ $contestedVotes = printf "contested: %v" $desc.ContestedString }}
								
								{{ if $desc.CorrectedByContest }}
									{{ $color = "lime" }}
								{{ else }}
									{{ $color = "red" }}
								{{ end }}
							{{ end }}
							(<font color="{{ $color }}">{{ $desc.ScorePercent }}%</font>)&nbsp;&nbsp;<code>{{ .Guess }}</code>
							{{ if eq $player (currentPlayer) }}
								{{ if $desc.IsContested }}
									<br>{{ $contestedVotes }}
								{{ else if not $desc.Correct }}
									<form method=POST action=/game/action/contest style="display: inline;">
										<input type=hidden name=guessValue value="{{ .Guess }}" />
										<input type=submit value="Contest" />
									</form>
								{{ end }}
							{{ else if $desc.IsContested }}
								{{/* Contested AND not current player */}}

								{{ $vote := $desc.FindContestedVote (currentPlayer) }}
								{{ if $vote.Voted }}
									<br>{{ $contestedVotes }}
									(you voted {{if $vote.Vote}}ya{{else}}nah{{end}})
								{{ else if eq $vote.Player (currentPlayer) }}
									{{/* don't show $contestedVotes until they vote */}}
									<br>
									<blink>Should this count?</blink>
									<form method=POST action=/game/action/contest-vote style="display: inline;">
										<input type=hidden name=guessPlayer value="{{ $player.Name }}" />
										<input type=hidden name=guessValue value="{{ .Guess }}" />
										<input type=hidden name=vote value="will be set later" />
										<button onclick="
											this.parentElement.querySelector('input[name=vote]').value='true';
											this.parentElement.submit();
										">ya</button>
										<button onclick="
											this.parentElement.querySelector('input[name=vote]').value='false';
											this.parentElement.submit();
										">nah</button>
									</form>
								{{ else }}
									{{/* else this player wasn't found in votes list, so they are new I guess */}}
								{{ end }}
							{{ end }}
						</li>
					{{ else }}
						{{if not $roundState.Passed}}
							<li>(no guesses)</li>
						{{ end }}
					{{ end }}
					{{ if $roundState.Passed }}
						<li>(pass)</li>
					{{ end }}
				</ol>
			</li>
		{{ end }}
	</ul>
{{ end }}

<meta name="viewport" content="width=device-width, initial-scale=1">
<script src="/static/xhr.js"></script>
<script>
function xhrCheckNextRound() {
	{{/* TODO: technically this is a race, since we could see a different state between now and when the template finishes executing; it's probably fine! */}}
	xhrDo("GET", "/game/wait/state", {"next_state": "{{plus1 .State.StateCounter}}"}).then(
		(xhr) => { window.location.reload(); },
		(xhr) => {
			if (xhr.chowski_aborted) {
				console.log("hanging refresh poll cancelled; status="+xhr.status+"; responseText="+xhr.responseText+";");
				return;
			}
			var desc = '';
			for (const [key, val] of Object.entries(xhr)) {
				if (desc != '') {
					desc += '\n';
				}
				desc += key+':'+val;
			}
			alert("failure to check next round, you'll have to manually refresh; err="+xhr+"\n"+desc);
		},
	);
}
</script>
<h3><a href=/game?scoreboard>SCOREBOARD</a></h3>

{{ $round := .State.CurrentRound }}
{{ if $round.IsZero }}
	<script>xhrCheckNextRound();</script>
	{{ $players := .State.Players }}
	{{ $priorRound := .State.PriorRound }}
	{{ if $priorRound.IsZero }}
		<h2>Waiting to start game</h2>
		{{ len $players }} logged in player{{if ne 1 (len $players)}}s{{end}}:
		<ul>
			{{range $players}}<li>{{.}}{{ if eq . (currentPlayer) }}&nbsp;(YOU!){{end}}</li>{{end}}
		<br>
		<form method=post action=/game/action/begin-round>
			<input type=submit value="Start game!">
		</form>
	{{ else }}
		{{ $s := $priorRound.PlayerState .Player }}
		{{ if $s.Correct }}
			{{ template "confetti" (randConfettiChar) }}
		{{ end }}
		Round not started yet...<br><br>
		Last round <small>({{ $priorRound.Song }})</small>:
		{{ template "playersSummary" $priorRound }}
		{{ with $priorRound.OutstandingContests }}
			<h3>Players who need to vote</h3>
			<ul>
				{{ range . }}<li>{{.}}</li>{{end}}
			</ul>
		{{ else }}
			<br>
			<form method=post action=/game/action/begin-round>
				<input type=submit value="Next round!">
			</form>
		{{ end }}
	{{ end }}
{{ else }}
	{{ $s := $round.PlayerState .Player }}
	{{ if $s.StillGuessing }}
		Guess the song, quick!
		<form method=post action=/game/action/guess>
			<label>Title: <input name=title></label>
			<input type=submit value=Guess>
		</form>
		<form method=post action=/game/action/pass>
			<input type=submit value=
				{{- if eq (len $s.Guesses) 0 -}}
					"Skip (don't know this song)"
				{{- else -}}
					"Give up"
				{{- end -}}
			>
		</form>
		<ol>
			{{ range $round.PlayerSummary .Player }}
				<li>{{.}}</li>
			{{ end }}
		</ol>
	{{ else }}
		<script>xhrCheckNextRound();</script>
		{{ if $s.Correct }}
			{{ template "confetti" (randConfettiChar) }}
			<font color="dark green">Nice job, you got it right!</font>
		{{ else }}
			<font color="maroon">You'll get next one.</font>
		{{ end }}
		<br><br>
		Other players are currently desperately guessing for {{ $round.Song }}:
		{{ template "playersSummary" $round }}
	{{ end }}
{{ end }}

</div>
</html>`))

type scoreboardData struct {
	Player *mqgame.Player
	*mqgame.State
}

var scoreboardTemplate = template.Must(template.New("scoreboardTemplate").Parse(`<html>
<meta name="viewport" content="width=device-width, initial-scale=1">

{{ $scoreboard := .Scoreboard }}
<a href=/game>BACK</a><hr>

{{/*
{{ with $scoreboard.Ranking }}
	<h3>Ranking</h3>
	<ol>{{ range . }}
		<li>{{ .Player.Name }}: {{ .Total }}</li>
	{{ end }}</ol>
{{ else }}
	No one has made any guesses yet!
{{ end }}
*/}}

{{ $currentRound := .CurrentRound }}
{{ $allPlayers := .Players }}
<table border=1 cellspacing=0>
	<tr>
		<th>Song</th>
		{{ range $allPlayers }}<th>{{.}}</th>{{end}}
	</tr>
	{{ range $scoreboard.Rounds }}<tr>
		<td>
			{{ if eq .Round.Song $currentRound.Song }}
				(ongoing round)
			{{ else }}
				{{ .Round.Song }}
			{{ end }}
		</td>
		{{ $byPlayer := .ScoresByPlayer }}
		{{ range $allPlayers }}<td>
			{{ $info := index $byPlayer . }}
			{{ printf "%0.2g" $info.Score }}{{ if gt $info.Delay 0 }}&nbsp;<small>({{ $info.DelayString }})</small>{{ end }}
		</td>{{ end }}
	</tr>{{ end }}
	<tr>
		<th align=right>Total:</th>
		{{ range $allPlayers }}
			<td>{{ printf "%0.2g" ($scoreboard.PlayerTotal .) }}</td>
		{{ end }}
	</tr>
</table>
</html>`))

//<button style="max-width: 400px; width: 75%; aspect-ratio: 1/1;" onclick=
