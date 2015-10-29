package linkedin

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
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
			ClientID:     os.Getenv("LINKEDIN_KEY"),
			ClientSecret: os.Getenv("LINKEDIN_SECRET"),
			Endpoint:     linkedin.Endpoint,
			RedirectURL:  "http://localhost:3000/auth/linkedin/callback",
			Scopes:       nil,
		},
		IdentityURL: "https//api.linkedin.com/v1/people/~:(id,first-name,last-name,headline,location:(name),picture-url,email-address)",
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
func (p *Provider) GetIdentity(tok *oauth2.Token) (string, error) {
	client := p.config.Client(oauth2.NoContext, tok)
	resp, err := client.Get(p.IdentityURL)
	fmt.Printf("Res: %+v\nErr: %s\n", resp, err)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return "", nil
}
