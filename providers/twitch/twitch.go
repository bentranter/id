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
	// ScopeUserRead provides read access to non-public user information, such
	// as their email address.
	ScopeUserRead string = "user_read"
	// ScopeUserBlocksEdit provides the ability to ignore or unignore on
	// behalf of a user.
	ScopeUserBlocksEdit string = "user_blocks_edit"
	// ScopeUserBlocksRead provides read access to a user's list of ignored
	// users.
	ScopeUserBlocksRead string = "user_blocks_read"
	// ScopeUserFollowsEdit provides access to manage a user's followed
	// channels.
	ScopeUserFollowsEdit string = "user_follows_edit"
	// ScopeChannelRead provides read access to non-public channel information,
	// including email address and stream key.
	ScopeChannelRead string = "channel_read"
	// ScopeChannelEditor provides write access to channel metadata (game,
	// status, etc).
	ScopeChannelEditor string = "channel_editor"
	// ScopeChannelCommercial provides access to trigger commercials on
	// channel.
	ScopeChannelCommercial string = "channel_commercial"
	// ScopeChannelStream provides the ability to reset a channel's stream key.
	ScopeChannelStream string = "channel_stream"
	// ScopeChannelSubscriptions provides read access to all subscribers to
	// your channel.
	ScopeChannelSubscriptions string = "channel_subscriptions"
	// ScopeUserSubscriptions provides read access to subscriptions of a user.
	ScopeUserSubscriptions string = "user_subscriptions"
	// ScopeChannelCheckSubscription provides read access to check if a user is
	// subscribed to your channel.
	ScopeChannelCheckSubscription string = "channel_check_subscription"
	// ScopeChatLogin provides the ability to log into chat and send messages.
	ScopeChatLogin string = "chat_login"
)

// New returns a new provider. Some providers have their
// endpoints as part of the Oauth2 package. Twitch is not
// one of them, so it must be entered manually.
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
// It's essentially a wrapper around the Oauth2 config
// to have an identity URL field.
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
// callback URL sp we can exchange it for an access token.
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

	// Twitch doesn't follow the Oauth2 spec correctly. The
	// spec states that the access token should be passed
	// as a `Bearer` token, but Twitch wants an `Oauth`
	// token, so we must include this header.
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

// readBody is just a convenience method for getting a
// `*psa.User` out of the JSON response from the Twitch
// API.
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
