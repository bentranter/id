package id

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
)

// User contains the bare minimum info we need to identify
// someone from a provider.
type User struct {
	Email string
	ID    string
	Name  string
}

// Provider implements all the functions we need.
type Provider interface {
	BuildAuthURL(state string) string
	GetCodeURL(r *http.Request) string
	GetToken(code string) (*oauth2.Token, error)
	GetIdentity(*oauth2.Token) (*User, error)
}

// Authorize builds the auth url and redirects a user to
// it.
func Authorize(p Provider) http.Handler {
	url := p.BuildAuthURL("state")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})
}

// Callback handles the callback part of the flow.
func Callback(p Provider, redirectURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := p.GetCodeURL(r)
		tok, err := p.GetToken(code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user, err := p.GetIdentity(tok)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cookie, err := GenToken(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, cookie)

		http.Redirect(w, r, "http://localhost:3000/"+redirectURL, http.StatusFound)
	})
}

// HTTPRouterAuthorize is the same thing as the regular
// `Authorize`, but for Julien Schmidt's HttpRouter.
//
// This should be moved into it's own package.
func HTTPRouterAuthorize(p Provider) httprouter.Handle {
	url := p.BuildAuthURL("state")
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})
}

// HTTPRouterCallback is the thing as `Callback` but for
// Julien Schmidt's HttpRouter
func HTTPRouterCallback(p Provider, redirectURL string) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		code := p.GetCodeURL(r)
		tok, err := p.GetToken(code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		user, err := p.GetIdentity(tok)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		cookie, err := GenToken(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.SetCookie(w, cookie)

		http.Redirect(w, r, "http://localhost:3000/"+redirectURL, http.StatusFound)
	})
}
