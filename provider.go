package psa

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Provider implements all the functions we need.
type Provider interface {
	BuildAuthURL(state string) string
	GetCodeURL(r *http.Request) string
	GetAccessToken() string
	GetIdentity() string
}

// Authorize builds the auth url and redirects a user to
// it.
func Authorize(p Provider) http.HandlerFunc {
	url := p.BuildAuthURL("state")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})
}

// Callback handles the callback part of the flow.
func Callback(p Provider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := p.GetCodeURL(r)
		http.Get(code)
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
func HTTPRouterCallback(p Provider) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		code := p.GetCodeURL(r)
		http.Get(code)
	})
}
