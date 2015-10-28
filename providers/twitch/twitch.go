package twitch

import (
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

const (
	// ScopeUserRead allows the client to access the user's
	// email address and id (for example)
	ScopeUserRead string = "user_read"
)

// New returns a new provider. Some providers have their
// endpoints as part of the Oauth2 package.
func New() *Provider {
	return &Provider{
		config: &oauth2.Config{
			ClientID:     os.Getenv("TWITCH_KEY"),
			ClientSecret: os.Getenv("TWICTH_SECRET"),
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://api.twitch.tv/kraken/oauth2/authorize",
				TokenURL: "https://api.twitch.tv/kraken/oauth2/token",
			},
			RedirectURL: "http://localhost:3000/auth/twitch/callback",
			Scopes:      []string{ScopeUserRead},
		},
		IdentityURL: "",
	}
}

// Provider holds all the info for our Oauth2 provider.
type Provider struct {
	config      *oauth2.Config
	IdentityURL string
}

// BuildAuthURL builds the authenticartion endpoint that we
// redirect our users to. State needs to be an unguessable
// string (probability of guessing < 2^128).
//
// It might be possible here to write some helper that
// generates a cryptographically secure state string, and
// uses that an a nonce. Since the provider instance is
// "alive" across our middleware, we could tack the nonce
// onto the middleware, and verify it once the identity
// provider returns it?
func (p *Provider) BuildAuthURL(state string) string {
	return p.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// GetCodeURL gets the short-lived access code from the
// callback URL that we can exchange for an access
// token.
func (p *Provider) GetCodeURL(r *http.Request) string {
	return ""
}

func (p *Provider) GetAccessToken() string {
	return ""
}

func (p *Provider) GetIdentity() string {
	return ""
}

// // Callback handles the rest.
// func Callback(w http.ResponseWriter, r *http.Request) {
// 	code := r.URL.Query().Get("code")
// 	config := psa.NewConfig(os.Getenv("TWITCH_KEY"), os.Getenv("TWITCH_SECRET"), "http://localhost:3000/auth/twitch/callback", authURL, tokenURL)
// 	url := config.BuildTokenURL(code)
// 	req, err := http.NewRequest("POST", tokenURL, url)
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		log.Fatalf("It's broken: %s\n", err)
// 	}
// 	log.Printf("Response: %+v\n", resp)
// 	defer resp.Body.Close()
// }
