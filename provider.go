package psa

import (
	"net/http"
)

// Provider implements all the functions we need.
type Provider interface {
	BuildAuthURL() string
	GetCode(r *http.Request) string
	GetAccessToken() string
	GetIdentity() string
}

// Authorize builds the auth url and redirects a user to
// it.
func Authorize(p Provider) http.Handler {
	url := p.BuildAuthURL()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})
}

// Callback handles the callback part of the flow.
func Callback(p Provider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := p.GetCode(r)
		http.Get(code)
	})
}
