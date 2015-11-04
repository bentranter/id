package twitch

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/bentranter/psa"
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
			ClientSecret: os.Getenv("TWITCH_SECRET"),
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://api.twitch.tv/kraken/oauth2/authorize",
				TokenURL: "https://api.twitch.tv/kraken/oauth2/token",
			},
			RedirectURL: "http://localhost:3000/auth/twitch/callback",
			Scopes:      []string{ScopeUserRead},
		},
		IdentityURL: "https://api.twitch.tv/kraken/user",
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
	return r.URL.Query().Get("code")
}

// GetToken gets the access and refresh tokens from the
// provider. Should probably be called `GetTokens`.
func (p *Provider) GetToken(code string) (*oauth2.Token, error) {
	tok, err := p.config.Exchange(oauth2.NoContext, code)
	return tok, err
}

// GetIdentity get's the client's identity from the
// provider. I'm not sure what every provider returns...
// I should look at what they do is Passport.js
//
// For all the providers that love to do weird stuff,
func (p *Provider) GetIdentity(tok *oauth2.Token) (*psa.User, error) {
	client := p.config.Client(oauth2.NoContext, tok)

	req, err := http.NewRequest("GET", p.IdentityURL, nil)
	req.Header.Add("Accept", "application/vnd.twitchtv.v3+json")
	req.Header.Add("Authorization", "OAuth"+tok.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	user := readBody(resp.Body)
	return user, nil
}

func readBody(r io.Reader) *psa.User {
	user := struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}{}

	err := json.NewDecoder(r).Decode(&user)
	if err != nil {
		panic(err)
	}

	return &psa.User{
		Name:  user.Name,
		Email: user.Email,
		ID:    user.ID,
	}
}
