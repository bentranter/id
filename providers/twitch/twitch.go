package twitch

import (
	"net/http"
	"os"
)

const (
	// ScopeUserRead allows the client to access the user's
	// email address and id (for example)
	ScopeUserRead string = "user_read"
)

// New returns a new provider.
func New() *Provider {
	return &Provider{
		ClientID:     os.Getenv("TWITCH_KEY"),
		ClientSecret: os.Getenv("TWITCH_SECRET"),
		CallbackURL:  "http://localhost:3000/auth/twitch/callback",
		AuthURL:      "https://api.twitch.tv/kraken/oauth2/authorize",
		TokenURL:     "https://api.twitch.tv/kraken/oauth2/token",
		IdentityURL:  "https://api.twitch.tv/kraken/user",
	}
}

// Provider does this.. might call this Twitch instead
type Provider struct {
	ClientID     string
	ClientSecret string
	CallbackURL  string
	AuthURL      string
	TokenURL     string
	IdentityURL  string
}

// BuildAuthURL builds the authenticartion endpoint that we
// redirect our users to.
func (p *Provider) BuildAuthURL() string {
	return p.AuthURL
}

// GetCode gets the short-lived access code from the
// callback URL that we can exchange for an access
// token.
func (p *Provider) GetCode(r *http.Request) string {
	code := r.URL.Query().Get("code")
	return code
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
