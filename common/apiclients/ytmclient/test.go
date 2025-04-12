//go:build ignore

package main

import (
	"chowski3/common/apiclients/ytmclient"
	"chowski3/common/auth"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var (
	secrets = auth.MustMakeSecretStore("/tmp/put-this-in-a-smarter-spot-if-you-want", auth.NewAppID("bananapants"))

	ytScopes = []string{
		"https://www.googleapis.com/auth/youtube",
	}
)

func main() {
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
	fmt.Println(ytm.GetSong(ctx, "ZRJdVTXkdGI"))
}
