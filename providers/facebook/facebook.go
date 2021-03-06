package facebook

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/bentranter/id"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

const (
// ScopeUserRead allows the client to access the user's
// email address and id (for example)
// ScopeUserRead string = "user_read"
)

// New returns a new provider. Some providers have their
// endpoints as part of the Oauth2 package.
func New() *Provider {
	return &Provider{
		config: &oauth2.Config{
			ClientID:     os.Getenv("FACEBOOK_KEY"),
			ClientSecret: os.Getenv("FACEBOOK_SECRET"),
			Endpoint:     facebook.Endpoint,
			RedirectURL:  "http://localhost:3000/auth/facebook/callback",
			Scopes:       []string{"email"},
		},
		IdentityURL: "https://graph.facebook.com/me?fields=email,first_name,last_name,link,bio,id,name,picture,location",
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

// GetIdentity gets the client's identity from the
// provider.
func (p *Provider) GetIdentity(tok *oauth2.Token) (*id.User, error) {
	client := p.config.Client(oauth2.NoContext, tok)

	resp, err := client.Get(p.IdentityURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	user := readBody(resp.Body)
	return user, nil
}

func readBody(r io.Reader) *id.User {
	user := struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}{}

	err := json.NewDecoder(r).Decode(&user)
	if err != nil {
		panic(err)
	}

	return &id.User{
		Name:  user.Name,
		Email: user.Email,
		ID:    user.ID,
	}
}
