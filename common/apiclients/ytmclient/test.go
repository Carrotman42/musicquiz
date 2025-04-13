//go:build ignore

package main

import (
	"chowski3/common/apiclients/ytmclient"
	"chowski3/common/auth"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

var (
	secrets = auth.MustMakeSecretStore("/tmp/put-this-in-a-smarter-spot-if-you-want", auth.NewAppID("bananapants"))

	ytScopes = []string{
		"https://www.googleapis.com/auth/youtube",
	}

	videoIdFlag = flag.String("video_id", "", "Video ID of the YouTube Music song to test-query for")
)

func main() {
	flag.Parse()

	ctx := context.Background()
	hc, err := auth.GoogleOAuthDance(ctx, secrets, ytScopes)
	if err != nil {
		log.Fatal("auth failure: ", err)
	}

	// Setup the initial magic SOCS=CAI cookie value
	cookies, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal("cookiejar failure: ", err)
	}
	ytmurl, err := url.Parse("https://music.youtube.com")
	if err != nil {
		log.Fatal("couldn't parse ytm URL?! ", err)
	}
	cookies.SetCookies(ytmurl, []*http.Cookie{{Name: "SOCS", Value: "CAI", Path: "/", Domain: "music.youtube.com"}})

	hc = &http.Client{Jar: cookies}
	ytm := ytmclient.New(hc)


	//v=h_r1CR6Q8z0&si=fWwqWbiquAjI99iW
	// fmt.Println(ytm.GetSong(ctx, "h_r1CR6Q8z0"))
	var songId string
	if *videoIdFlag != "" {
		songId = *videoIdFlag
	} else {
		songId = "ZRJdVTXkdGI"
	}
	song, err := ytm.GetSong(ctx, songId)
	if err != nil {
		log.Fatal("error getting song ZRJdVTXkdGI: ", err)
	}

	fmt.Printf("===== Song Info: '%s' =====\n", song.VideoDetails.VideoId)
	fmt.Printf("Title:    '%s'\n", song.VideoDetails.Title)
	fmt.Printf("Artist:   '%s'\n", song.VideoDetails.Author)
	var album string = ""
	for _, tag := range song.Microformat.MicroformatDataRenderer.Tags {
		// Need a `contains` check because for songs with two artists (e.g. 'Billy Strings & Brian Sutton'), there's a tag for each of them: ['Billy Strings', 'Brian Sutton', 'Live at the Legion', 'Randall Collins / Done Gone']
		if !strings.Contains(song.VideoDetails.Author, tag) && tag != song.VideoDetails.Title {
			if album == "" {
				album = tag
			} else {
				log.Printf("got more than one album title? '%s' and '%s'\n", album, tag)
			}
		}
	}
	if album == "" {
		log.Fatalln("couldn't extract album title; wrong number of tags")
	}
	fmt.Printf("Album:    '%s'\n", album)
	d, err := strconv.Atoi(song.Microformat.MicroformatDataRenderer.VideoDetails.DurationSeconds)
	if err != nil {
		log.Fatal("error parsing duration as int: ", err)
	}
	fmt.Printf("Duration: %dm%ds\n", d / 60, d % 60)
	fmt.Printf("Stream:   %s\n", song.StreamingData.ServerAbrStreamingUrl)
}
