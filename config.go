package psa

import (
	"log"
	"net/url"
)

// Config is our Oauth config
type Config struct {
	ClientID     string
	ClientSecret string
	CallbackURL  string
	AuthURL      string
	TokenURL     string
	Scopes       []string
}

// BuildAuthURL creates an AuthURL from the config
func (c *Config) BuildAuthURL() string {
	u, err := url.Parse(c.AuthURL)
	if err != nil {
		log.Fatalf("Couldn't parse the AuthURL provided: %s\n", err)
	}

	q := u.Query()
	q.Set("client_id", c.ClientID)
	q.Set("redirect_uri", c.CallbackURL)
	q.Set("scope", "user_read")
	q.Set("state", "state")
	q.Set("response_type", "code")

	u.RawQuery = q.Encode()
	return u.String()
}

// BuildTokenURL does what it says it does.
func (c *Config) BuildTokenURL(code string) string {
	u, err := url.Parse(c.TokenURL)
	if err != nil {
		log.Fatalf("Couldn't parse the TokenURL provided: %s\n", err)
	}

	q := u.Query()
	q.Set("client_id", c.ClientID)
	q.Set("client_secret", c.ClientSecret)
	q.Set("grant_type", code)
	q.Set("scope", "user_read")
	q.Set("state", "state")

	u.RawQuery = q.Encode()
	return u.String()
}

// NewConfig creates a new config for a provider.
func NewConfig(ClientID, ClientSecret, CallbackURL, AuthURL, TokenURL string, Scopes ...string) *Config {
	return &Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		CallbackURL:  CallbackURL,
		AuthURL:      AuthURL,
		TokenURL:     TokenURL,
		Scopes:       Scopes,
	}
}
