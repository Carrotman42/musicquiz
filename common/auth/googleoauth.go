package auth

import (
	_ "embed"
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"
)

//go:embed clientsecret.txt
var clientSecret string

// Uses a hardcoded client ID/secret.
// Uses stdout/stdin to interact with the user.
// TODO: use the context for cancelling waiting for user response.
func GoogleOAuthDance(ctx context.Context, secrets SecretStore, scopes []string) (*http.Client, error) {
	conf := &oauth2.Config{
		ClientID:     "615383202793-khpjq1lt7aabgnq087cv5utu5qqchv1d.apps.googleusercontent.com",
		ClientSecret: clientSecret,
		RedirectURL:  fmt.Sprintf("http://%v/auth/oauth_redirect", "localhost"), // *domain),
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}
	if token, err := googleAccessTokenSecret.Read(secrets); err == nil {
		return conf.Client(ctx, token), nil
	}
	url := conf.AuthCodeURL("state")
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)
	sc := bufio.NewScanner(os.Stdin)
	if !sc.Scan() {
		return nil, fmt.Errorf("auth client: cancelled")
	}
	code := sc.Text()

	t, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}
	if err := googleAccessTokenSecret.Write(secrets, t); err != nil {
		log.Printf("WARNING: failed to persist oauth token; you'll have to log in again next time; error: %v", err)
	}
	fmt.Println("Logged in!")
	return conf.Client(ctx, t), nil
}

var googleAccessTokenSecret = NewSecret[*oauth2.Token]("auth.GoogleOAuthDance")

// GoogleADC supports auth via "application default credentials".
func GoogleADC(ctx context.Context, scopes []string) (*http.Client, error) {
	hc, ep, err := htransport.NewClient(ctx, option.WithScopes(scopes...))
	if err != nil {
		return nil, err
	}
	// TODO: plumb this endpoint correctly, in case we want to hook into here for testing.
	_ = ep
	// TODO: validate that this token is still valid.
	return hc, err
}
