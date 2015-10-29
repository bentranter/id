package twitch

import (
	"fmt"
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
func (p *Provider) GetIdentity(tok *oauth2.Token) string {
	client := p.config.Client(oauth2.NoContext, tok)

	req, err := http.NewRequest("GET", p.IdentityURL, nil)
	req.Header.Add("Accept", "application/vnd.twitchtv.v3+json")
	fmt.Printf("Token?: %+v\n", tok)
	req.Header.Add("Authorization", "OAuth "+tok.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return "Twitch: 200 Ok.\n\nIt worked!"
	case 400:
		return fmt.Sprintf("Twitch: 400 Bad Request\n\nFor some reason, fetching a token failed. It's likely you when over a rate limit.\n\n%+v\n", resp)
	case 401:
		return fmt.Sprintf("Twitch: 401 Unauthorized.\n\nPlease double check that the IdentityURL is valid, that the following headers are set:\n\nAccept: application/vnd.twitchtv.v3+json\nAuthorization: OAuth %s\n\nIf you continue to have trouble, file an issue on GitHub.\n", tok.AccessToken)
	case 404:
		return "Twitch: 404 Not Found.\n\nThe requested resource wasn't available. There might be a problem with this user's Twitch account.\n"
	default:
		return fmt.Sprintf("Twitch: %s %s\n", resp.StatusCode, resp.Status)
	}
}
