package twitch

import (
	"net/http"
	"os"

	"github.com/bentranter/psa"
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

type Provider struct {
	ClientID     string
	ClientSecret string
	CallbackURL  string
	AuthURL      string
	TokenURL     string
	IdentityURL  string
}

func (p *Provider) Authorize(w http.ResponseWriter, r *http.Request) {
	psa.BuildAuthURL(p)
}

// Authorize does the whole process. It authenticates with
// the auth provider, and sets a cookie for our middleware
// to deal with.
// func Authorize(w http.ResponseWriter, r *http.Request) {
// 	// Redirect them to the sign in URL
// 	config := psa.NewConfig(os.Getenv("TWITCH_KEY"), os.Getenv("TWITCH_SECRET"), "http://localhost:3000/auth/twitch/callback", authURL, tokenURL)
// 	url := config.BuildAuthURL()
// 	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
// }

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
