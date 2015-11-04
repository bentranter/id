package google

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/bentranter/psa"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	// ScopeProfile allows the client to access the user's
	// profile info.
	ScopeProfile string = "profile"

	//ScopeEmail allows the cient to access the user's
	// email address.
	ScopeEmail string = "email"
)

// New returns a new provider. Some providers have their
// endpoints as part of the Oauth2 package, and Google is
// one of them.
func New() *Provider {
	return &Provider{
		config: &oauth2.Config{
			ClientID:     os.Getenv("GPLUS_KEY"),
			ClientSecret: os.Getenv("GPLUS_SECRET"),
			Endpoint:     google.Endpoint,
			RedirectURL:  "http://localhost:3000/auth/gplus/callback",
			Scopes:       []string{ScopeEmail, ScopeProfile},
		},
		IdentityURL: "https://www.googleapis.com/oauth2/v2/userinfo",
	}
}

// Provider holds all the info for our Oauth2 provider.
type Provider struct {
	config      *oauth2.Config
	IdentityURL string
}

// BuildAuthURL builds the authentication endpoint that we
// redirect our users to.
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
// provider.
func (p *Provider) GetToken(code string) (*oauth2.Token, error) {
	tok, err := p.config.Exchange(oauth2.NoContext, code)
	return tok, err
}

// GetIdentity gets the client's identity from the
// provider.
func (p *Provider) GetIdentity(tok *oauth2.Token) (*psa.User, error) {
	client := p.config.Client(oauth2.NoContext, tok)

	resp, err := client.Get(p.IdentityURL)
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
