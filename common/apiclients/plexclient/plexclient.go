package plexclient

import (
	"chowski3/common/auth"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"time"
)

// ======================================== public API ==============================

// Everything you need for a Plex connection
type PlexConnection struct {
	// Shared HTTP client
	client *http.Client

	// Token to the global plex.tv API
	globalToken string

	// Individual plex media server
	url string
	// Access token to the individual plex media server (TODO: when does it expire?)
	hostToken string
}

// A playlist of tracks
type Playlist struct {
	Name string
	NumTracks int

	key string
}

// An individual music track
type Track struct {
	Artist string
	Title  string
	Album  string

	key string
}

// Create and authenticate a new connection to the Plex server.
func Connect() (*PlexConnection, error) {
	// 0. Create a stable HTTP client that we can set some permanent(ish) headers on
	cookies, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Jar: cookies}

	// 1. Get an access token to the global Plex server
	globalToken, err := getAuthToken(client)
	if err != nil {
		fmt.Println("Error getting global auth token:", err)
		return nil, err
	}

	conn := PlexConnection{
		client:      client,
		globalToken: *globalToken,
	}

	// 2. Pick an individual plex media server
	servers, err := conn.listPlexMediaServers()
	if err != nil {
		fmt.Println("Error listing all servers:", err)
		return nil, err
	}
	fmt.Println("Please select a server:")
	for i := range len(servers) {
		fmt.Printf("  %d) %s (%s)\n", i, servers[i].Name, servers[i].PublicAddress)
	}
	var selection int
	_, err = fmt.Scan(&selection)
	if err != nil || selection < 0 || selection >= len(servers) {
		fmt.Println("Invalid selection: ", selection, err)
		return nil, errors.New("Invalid server selection")
	}

	// 3. From that server, pick an appropriate connection
	var c *PlexConn = nil
	// Prefer local & non-relay
	for i := range len(servers[selection].Connections) {
		server := servers[selection].Connections[i]
		if server.Local && !server.Relay {
			c = &server
			break
		}
	}
	// If none of those, take any non-relay
	if c == nil {
		for i := range len(servers[selection].Connections) {
			server := servers[selection].Connections[i]
			if !server.Relay {
				c = &server
				break
			}
		}
	}
	// If STILL none, take anything
	if c == nil {
		c = &servers[selection].Connections[0]
	}

	conn.url = c.URI
	conn.hostToken = servers[selection].AccessToken
	return &conn, nil
}

func (plex *PlexConnection) ListPlaylists() (*[]Playlist, error) {
	var resp PlexResponse
	err := pmsReq(plex, "GET", "/playlists", &resp)
	if err != nil {
		fmt.Println("Error listing playlists:", err)
		return nil, err
	}

	playlists := make([]Playlist, resp.MediaContainer.Size)
	for i := range resp.MediaContainer.Size {
		p := resp.MediaContainer.Metadata[i]
		playlists[i] = Playlist{Name: p.Title, NumTracks: p.LeafCount, key: p.Key}
	}
	return &playlists, nil
}

func (plex *PlexConnection) ListPlaylistContents(playlist Playlist) (*[]Track, error) {
	// List its contents
	var playlistContents PlexResponse
	err := pmsReq(plex, "GET", playlist.key, &playlistContents)
	if err != nil {
		fmt.Println("Error getting playlist contents:", err)
		return nil, err
	}
	tracks := make([]Track, playlistContents.MediaContainer.Size)
	for i := range playlistContents.MediaContainer.Size {
		t := playlistContents.MediaContainer.Metadata[i]
		dumped, _ := json.Marshal(t)
		if i < 3 {
			fmt.Println("Track:", string(dumped))
		}
		tracks[i] = Track{Artist: t.GrandparentTitle, Title: t.Title, Album: t.ParentTitle, key: t.Key}
	}
	return &tracks, nil
}

// TODO DO NOT SUBMIT: Just return the stream URL?
func (plex *PlexConnection) DownloadSong(song string, outfile string) {
	// Search for a song
	var output PlexResponse
	u := "/library/all?mediaQuery=" + url.PathEscape("type=track&title="+song)
	fmt.Println("Querying", u)
	err := pmsReq(plex, "GET", u, &output)
	if err != nil {
		fmt.Println("Error searching for a song:", err)
	}
	key := output.MediaContainer.Metadata[0].Key
	fmt.Printf("Searching for 'Guy I Used To Be': %d results -- '%s' (track #%d) (key=%s)\n", output.MediaContainer.Size, output.MediaContainer.Metadata[0].Title, output.MediaContainer.Metadata[0].Index, key)

	// Transcode it live? IT WORKS!!
	offset := 0
	u = fmt.Sprintf("/music/:/transcode/universal/start?offset=%d&path=%s", offset, url.PathEscape(key))
	fmt.Println("Transcoding with url:", u)
	bytes, err := pmsReqBytes(plex, "POST", u)
	if err != nil {
		fmt.Println("Error streaming song:", song, err)
		return
	}

	err = os.WriteFile(outfile, *bytes, 0644)
	if err != nil {
		fmt.Println("Error writing to "+outfile, err)
	}
	fmt.Println("Wrote mp3 bytes to " + outfile)
}

// Create a request for an individual API (e.g. "POST", "/downloadQueue") on an individual plex media server
func pmsReq[T any](plex *PlexConnection, method string, path string, output *T) error {
	_, err := ajax(plex.client, method, plex.url+path, http.Header{
		"Accept":                   {"application/json"},
		"X-Plex-Token":             {plex.hostToken},
		"X-Plex-Client-Identifier": {plexClientId()},
		"X-Plex-Platform":          {"Chrome"},
	}, http.NoBody, output)
	return err
}

// Create a request for an individual API (e.g. "POST", "/downloadQueue") on an individual plex media server that returns bytes
func pmsReqBytes(plex *PlexConnection, method string, path string) (*[]byte, error) {
	return ajax[any](plex.client, method, plex.url+path, http.Header{
		"Accept":                   {"application/json"},
		"X-Plex-Token":             {plex.hostToken},
		"X-Plex-Client-Identifier": {plexClientId()},
		"X-Plex-Platform":          {"Chrome"},
	}, http.NoBody, nil)
}

// ======================================== private APIs ==================================

// "The name of your app; ex 'My Cool App'"
func plexProductName() string {
	return "Chowski Music Quiz"
}

// "A random string or UUID is sufficient here"
func plexClientId() string {
	var hostname, err = os.Hostname()
	if err != nil {
		hostname = "unknowndevice"
	}
	return "ChowskiMusicQuiz-" + hostname
}

// PinResponse struct to unmarshal the JSON response
type Pin struct {
	ID        int       `json:"id"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// TokenResponse struct to unmarshal the JSON response
type TokenResponse struct {
	AuthToken string `json:"authToken"`
}

type PlexConn struct {
	Protocol string `json:"protocol"`
	Address  string `json:"address"`
	Port     int    `json:"port"`
	URI      string `json:"uri"`
	Local    bool   `json:"local"`
	Relay    bool   `json:"relay"`
	IPv6     bool   `json:"IPv6"`
}

// PlexDevice represents a single Plex Media Server or other device accessible to the user.
// The top-level JSON is an array of these structs.
type PlexDevice struct {
	Name                   string     `json:"name"`
	Product                string     `json:"product"`
	ProductVersion         string     `json:"productVersion"`
	Platform               string     `json:"platform"`
	PlatformVersion        string     `json:"platformVersion"`
	Device                 string     `json:"device"`
	ClientIdentifier       string     `json:"clientIdentifier"`
	Provides               string     `json:"provides"`
	OwnerID                *int64     `json:"ownerId"`
	SourceTitle            *string    `json:"sourceTitle"`
	PublicAddress          string     `json:"publicAddress"`
	AccessToken            string     `json:"accessToken"`
	SearchEnabled          bool       `json:"searchEnabled"`
	CreatedAt              time.Time  `json:"createdAt"`
	LastSeenAt             time.Time  `json:"lastSeenAt"`
	Owned                  bool       `json:"owned"`
	Home                   bool       `json:"home"`
	Synced                 bool       `json:"synced"`
	Relay                  bool       `json:"relay"`
	Presence               bool       `json:"presence"`
	HTTPSRequired          bool       `json:"httpsRequired"`
	PublicAddressMatches   bool       `json:"publicAddressMatches"`
	DNSRebindingProtection bool       `json:"dnsRebindingProtection"`
	NATLoopbackSupported   bool       `json:"natLoopbackSupported"`
	Connections            []PlexConn `json:"connections"`
}
type PlexDevices []PlexDevice

type PlexResponse struct {
	MediaContainer MediaContainer `json:"MediaContainer"`
}

type MediaContainer struct {
	Metadata        []Metadata `json:"Metadata"`
	AllowSync       bool       `json:"allowSync"`
	Identifier      string     `json:"identifier"`
	MediaTagPrefix  string     `json:"mediaTagPrefix"`
	MediaTagVersion int        `json:"mediaTagVersion"`
	Size            int        `json:"size"`
}

type Metadata struct {
	Genre                []TagItem `json:"Genre"`
	Image                []Image   `json:"Image"`
	AddedAt              int       `json:"addedAt"`
	Art                  string    `json:"art"`
	GrandparentArt       string    `json:"grandparentArt"`
	GrandparentGUID      string    `json:"grandparentGuid"`
	GrandparentKey       string    `json:"grandparentKey"`
	GrandparentRatingKey string    `json:"grandparentRatingKey"`
	GrandparentThumb     string    `json:"grandparentThumb"`
	GrandparentTitle     string    `json:"grandparentTitle"`
	GUID                 string    `json:"guid"`
	Index                int       `json:"index"`
	Key                  string    `json:"key"`
	LastViewedAt         int       `json:"lastViewedAt"`
	LeafCount int `json:"leafCount"`
	LibrarySectionID     int       `json:"librarySectionID"`
	LibrarySectionKey    string    `json:"librarySectionKey"`
	LibrarySectionTitle  string    `json:"librarySectionTitle"`
	ParentGUID           string    `json:"parentGuid"`
	ParentIndex          int       `json:"parentIndex"`
	ParentKey            string    `json:"parentKey"`
	ParentRatingKey      string    `json:"parentRatingKey"`
	ParentStudio         string    `json:"parentStudio"`
	ParentThumb          string    `json:"parentThumb"`
	ParentTitle          string    `json:"parentTitle"`
	ParentYear           int       `json:"parentYear"`
	RatingCount          int       `json:"ratingCount"`
	RatingKey            string    `json:"ratingKey"`
	Smart                bool      `json:"smart"`
	Summary              string    `json:"summary"`
	Thumb                string    `json:"thumb"`
	Title                string    `json:"title"`
	Type                 string    `json:"type"`
	UpdatedAt            int       `json:"updatedAt"`
	ViewCount            int       `json:"viewCount"`
}

type TagItem struct {
	Tag string `json:"tag"`
}

type Image struct {
	Alt  string `json:"alt"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

var (
	secrets         = auth.MustMakeSecretStore("/tmp/musicquiz/plex", auth.NewAppID("plexmusicquiz"))
	plexTokenSecret = auth.NewSecret[*string]("auth.PlexTokenSecret")
)

func getAuthToken(client *http.Client) (*string, error) {
	// Auth instructions: https://developer.plex.tv/pms/#section/API-Info/Authenticating-with-Plex
	//
	// I'm lazy so I'm doing the legacy auth token flow. JWTs require per-device signing keys and are hard.

	// 0. Check whether existing access token exists and is still valid
	if token, err := plexTokenSecret.Read(secrets); err == nil {
		req, err := http.NewRequest("GET", "https://plex.tv/api/v2/user", http.NoBody)
		if err != nil {
			fmt.Println("Error creating token validity check http request:", err)
			return nil, err
		}
		req.Header = http.Header{
			"Accept":                   {"application/json"},
			"X-Plex-Product":           {plexProductName()},
			"X-Plex-Client-Identifier": {plexClientId()},
			"X-Plex-Token":             {*token},
		}
		resp, err := client.Do(req)
		if err == nil {
			if resp.StatusCode == 200 {
				fmt.Println("Using saved token (still valid)")
				return token, nil
			} else {
				fmt.Printf("Token invalid (status %d), getting a new one\n", resp.StatusCode)
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				if err == nil {
					fmt.Println("Full response:", string(body))
				} else {
					fmt.Println("Err:", err)
				}
			}
		} else {
			fmt.Println("Couldn't check token validity, assuming invalid and getting a new one:", err)
		}
	}

	// 1. Request a PIN, if we don't already have one
	pin := Pin{}
	_, err := ajax(client, "POST", "https://plex.tv/api/v2/pins?strong=true", http.Header{
		"Accept":                   {"application/json"},
		"X-Plex-Product":           {plexProductName()},
		"X-Plex-Client-Identifier": {plexClientId()},
	}, http.NoBody, &pin)
	if err != nil {
		fmt.Println("Error getting PIN:", err)
		return nil, err
	}

	// 2. Build an auth URL
	authUrl := "https://app.plex.tv/auth#" +
		"?clientID=" + url.PathEscape(plexClientId()) +
		"&code=" + pin.Code +
		// TODO: register plex auth response URL at the top-level? otherwise, poll
		// "&forwardUrl=" + url.PathEscape("https://localhost/TODO") +
		"&context%5Bdevice%5D%5Bproduct%5D=" + url.PathEscape(plexProductName())
	fmt.Println("====================================================================================")
	fmt.Println("====================================================================================")
	fmt.Println("=== Please authorize access to Plex at: ", authUrl)
	fmt.Println("====================================================================================")
	fmt.Println("====================================================================================")
	fmt.Println("")

	// 3. Poll for the token
	for range 120 {
		fmt.Print(".")

		tokenResponse := TokenResponse{}
		_, err = ajax(client, "GET", "https://plex.tv/api/v2/pins/"+strconv.Itoa(pin.ID), http.Header{
			"Accept":                   {"application/json"},
			"X-Plex-Client-Identifier": {plexClientId()},
		}, http.NoBody, &tokenResponse)

		if tokenResponse.AuthToken != "" {
			fmt.Println("Found plex.tv auth token:", tokenResponse.AuthToken)
			err = plexTokenSecret.Write(secrets, &tokenResponse.AuthToken)
			if err != nil {
				fmt.Println("Failed to persist plex token, you may need to reauth next time:", err)
			}
			return &tokenResponse.AuthToken, nil
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println(".")

	return nil, errors.New("Timed out getting Plex auth token")
}

func (conn *PlexConnection) listPlexMediaServers() (PlexDevices, error) {

	// 4. Get individual Plex Media Server tokens?
	devices := PlexDevices{}
	_, err := ajax(conn.client, "GET", "https://clients.plex.tv/api/v2/resources?includeHttps=1&includeRelay=1&includeIPv6=1", http.Header{
		"Accept":                   {"application/json"},
		"X-Plex-Client-Identifier": {plexClientId()},
		"X-Plex-Token":             {conn.globalToken},
	}, http.NoBody, &devices)
	if err != nil {
		fmt.Println("Error getting all plex devices:", err)
		return nil, err
	}
	/*
			e.g.:
		   [
		   	{
		   		"name":"arensonjr-desktop",
		   		"product":"Plex Media Server",
		   		"productVersion":"1.43.0.10162-b67a664b6",
		   		"platform":"Linux",
		   		"platformVersion":"6.16.8-arch3-1 (#1 SMP PREEMPT_DYNAMIC Mon, 22 Sep 2025 22:08:35 +0000)",
		   		"device":"Gigabyte Technology Co., Ltd. B560 DS3H AC-Y1",
		   		"clientIdentifier":"8a0ff8b0a8b3cde5a2fd610236d5a006dbf17ad9",
		   		"provides":"server",
		   		"ownerId":null,
		   		"sourceTitle":null,
		   		"publicAddress":"67.190.126.127",
		   		"accessToken":"zBsc3v-K4av5-W1C9kVX",
		   		"searchEnabled":true,
		   		"createdAt":"2024-05-11T21:28:33Z",
		   		"lastSeenAt":"2025-09-28T22:55:36Z",
		   		"owned":true,
		   		"home":false,
		   		"synced":false,
		   		"relay":true,
		   		"presence":true,
		   		"httpsRequired":false,
		   		"publicAddressMatches":false,
		   		"dnsRebindingProtection":false,
		   		"natLoopbackSupported":true,
		   		"connections":[
		   			{"protocol":"https","address":"172.17.0.1","port":32400,"uri":"https://172-17-0-1.4a92409553f2449f9142e1d5ad56aa9f.plex.direct:32400","local":true,"relay":false,"IPv6":false},
		   			{"protocol":"https","address":"192.168.1.137","port":32400,"uri":"https://192-168-1-137.4a92409553f2449f9142e1d5ad56aa9f.plex.direct:32400","local":true,"relay":false,"IPv6":false},
		   			{"protocol":"https","address":"2601:280:5f00:eb50::ee78","port":32400,"uri":"https://2601-0280-5f00-eb50-0000-0000-0000-ee78.4a92409553f2449f9142e1d5ad56aa9f.plex.direct:32400","local":false,"relay":false,"IPv6":true},
		   			{"protocol":"https","address":"2601:280:5f00:eb50:bf54:e005:b3dc:dd12","port":32400,"uri":"https://2601-0280-5f00-eb50-bf54-e005-b3dc-dd12.4a92409553f2449f9142e1d5ad56aa9f.plex.direct:32400","local":false,"relay":false,"IPv6":true},
		   			{"protocol":"https","address":"67.190.126.127","port":32400,"uri":"https://67-190-126-127.4a92409553f2449f9142e1d5ad56aa9f.plex.direct:32400","local":false,"relay":false,"IPv6":false},
		   			{"protocol":"https","address":"69.164.192.240","port":8443,"uri":"https://69-164-192-240.4a92409553f2449f9142e1d5ad56aa9f.plex.direct:8443","local":false,"relay":true,"IPv6":false}
		   		]},
		   	{
		   		"name":"plexmediaserver","product":"Plex Media Server","productVersion":"1.42.1.10060-4e8b05daf","platform":"Linux","platformVersion":"24.04.3 LTS (Noble Numbat)","device":"Intel Corporation NUC7i3BNB","clientIdentifier":"794ff96a93bc98a87f4a608d8e56080714191ed4","provides":"server","ownerId":14436974,"sourceTitle":"doolbneerg","publicAddress":"50.35.70.33","accessToken":"xWRF_La-BaGT5cahnfak","searchEnabled":true,"createdAt":"2024-04-07T20:03:00Z","lastSeenAt":"2025-09-28T21:54:48Z","owned":false,"home":false,"synced":false,"relay":true,"presence":true,"httpsRequired":true,"publicAddressMatches":false,"dnsRebindingProtection":false,"natLoopbackSupported":true,"connections":[{"protocol":"https","address":"172.21.0.1","port":32400,"uri":"https://172-21-0-1.8a6b1392e7694d7faf9ba82498b2cc86.plex.direct:32400","local":true,"relay":false,"IPv6":false},{"protocol":"https","address":"172.19.0.1","port":32400,"uri":"https://172-19-0-1.8a6b1392e7694d7faf9ba82498b2cc86.plex.direct:32400","local":true,"relay":false,"IPv6":false},{"protocol":"https","address":"172.17.0.1","port":32400,"uri":"https://172-17-0-1.8a6b1392e7694d7faf9ba82498b2cc86.plex.direct:32400","local":true,"relay":false,"IPv6":false},{"protocol":"https","address":"172.18.0.1","port":32400,"uri":"https://172-18-0-1.8a6b1392e7694d7faf9ba82498b2cc86.plex.direct:32400","local":true,"relay":false,"IPv6":false},{"protocol":"https","address":"192.168.0.115","port":32400,"uri":"https://192-168-0-115.8a6b1392e7694d7faf9ba82498b2cc86.plex.direct:32400","local":true,"relay":false,"IPv6":false},{"protocol":"https","address":"50.35.70.33","port":17456,"uri":"https://50-35-70-33.8a6b1392e7694d7faf9ba82498b2cc86.plex.direct:17456","local":false,"relay":false,"IPv6":false},{"protocol":"https","address":"69.164.192.240","port":8443,"uri":"https://69-164-192-240.8a6b1392e7694d7faf9ba82498b2cc86.plex.direct:8443","local":false,"relay":true,"IPv6":false}]}]
	*/

	return devices, nil
}

func ajax[T any](client *http.Client, method string, url string, headers http.Header, body io.Reader, output *T) (*[]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return nil, err
	}
	req.Header = headers

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error executing HTTP request:", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		fmt.Println("Non-200 response from HTTP request:", resp.Status)
		return nil, fmt.Errorf("Non-200 status code: %s (from URL %s)", resp.Status, url)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	// TODO: For raw api debugging
	// log.Println("Request: ", method, url)
	// log.Println("Response:", string(respBody))

	if output != nil {
		err = json.Unmarshal(respBody, output)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			fmt.Println("Body is:", string(respBody))
			return nil, err
		}
		return nil, nil
	} else {
		return &respBody, nil
	}
}
