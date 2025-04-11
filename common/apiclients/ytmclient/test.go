//go:build ignore

package main

import (
	"chowski3/common/apiclients/ytmclient"
	"chowski3/common/auth"
	"context"
	"fmt"
	"log"
	"net/http"
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

	hc = http.DefaultClient
	ytm := ytmclient.New(hc)

	//v=h_r1CR6Q8z0&si=fWwqWbiquAjI99iW
	fmt.Println(ytm.GetSong(ctx, "h_r1CR6Q8z0"))
}
