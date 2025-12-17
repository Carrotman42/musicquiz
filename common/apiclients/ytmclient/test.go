//go:build ignore

package main

import (
	"chowski3/common/apiclients/ytmclient"
	"chowski3/common/auth"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/kkdai/youtube/v2"
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
	fmt.Printf("Title:      '%s'\n", song.VideoDetails.Title)
	fmt.Printf("Artist:     '%s'\n", song.VideoDetails.Author)
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
	fmt.Printf("Album:      '%s'\n", album)
	d, err := strconv.Atoi(song.Microformat.MicroformatDataRenderer.VideoDetails.DurationSeconds)
	if err != nil {
		log.Fatal("error parsing duration as int: ", err)
	}
	fmt.Printf("Duration:   %dm%ds\n", d / 60, d % 60)

	// Adapted from https://github.com/kkdai/youtube/blob/master/decipher.go#L15
	var streamOption *ytmclient.StreamingFormat
	for _, format := range song.StreamingData.Formats {
		if strings.Contains(format.MimeType, "audio/mp4") {
			streamOption = &format
			break
		}
	}
	for _, format := range song.StreamingData.AdaptiveFormats {
		if strings.Contains(format.MimeType, "audio/mp4") {
			streamOption = &format
			break
		}
	}
	for _, format := range song.StreamingData.Formats {
		if strings.Contains(format.MimeType, "audio/mp4") {
			streamOption = &format
			break
		}
	}
	if streamOption == nil {
		log.Fatalf("no mp4 audio stream found in %v,", song)
	}
	params, err := url.ParseQuery(streamOption.SignatureCipher)
	if err != nil {
		log.Fatal("error parsing stream signature url params: ", err)
	}
	uri, err := url.Parse(params.Get("url"))
	if err != nil {
		log.Fatal("error parsing streaming url: ", err)
	}
	// TODO:
	// 1. Fetch embed URL from https://youtube.com/embed/${VIDEO_ID}?hl=en
	// 2. Rip out playerPath=/s/player/$(PLAYER_ID}/player_ias.vflset/${LANG}/base.js from the body
	// 3. Fetch https://youtube.com${PLAYER_PATH}
	// 4. Parse the javascript-object-part of the body into javascript decryption instructions
	// 5. Run all of the operations on the "s" param and re-add it on the "sp" param
	// 6. Decrypt the "v" param using the "n" function in that javascript, queryencode the result, and use it to replace the `RawQuery` on the URI

	fmt.Printf("Stream:     %s\n", streamOption.SignatureCipher)
	fmt.Printf("Stream URI: %s\n", uri)


	// TODO: Testing youtubedr

	client := youtube.Client{}

	video, err := client.GetVideo("https://www.youtube.com/watch?v=" + songId)
	if err != nil {
		panic(err)
	}

	// vid := youtube.Video{
	// 	Title: song.VideoDetails.Title,
	// 	Author: song.VideoDetails.Author,
	// 	ID: song.VideoDetails.VideoId,
	// 	Duration: time.Duration(d) * time.Second,
	// 	Formats: []youtube.Format{
	// 		{
	// 			Cipher: streamOption.SignatureCipher,
	// 		},
	// 	},
	// }
	// fmt.Printf("manual vid: %v\n", vid)

	formats := video.Formats.Type("audio/mp4")
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	// // streamUrl, err := client.GetStreamURLContext(ctx, video, &formats[0])
	// streamUrl, err := client.GetStreamURL(&vid, &vid.Formats[0])
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("Stream URI (ydr): %s\n", streamUrl)
	// stream, err := ytm.StreamSong(ctx, streamUrl)
	// if err != nil {
	// 	panic(err)
	// }
	// defer stream.Close()

	file, err := os.Create("testaudio.mp4")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Printf("Downloading to %s...\n", file.Name())
	_, err = io.Copy(file, stream)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Downloaded to %s\n", file.Name())
}
