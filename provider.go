package psa

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
)

// Provider implements all the functions we need.
type Provider interface {
	BuildAuthURL(state string) string
	GetCodeURL(r *http.Request) string
	GetToken(code string) (*oauth2.Token, error)
	GetIdentity(*oauth2.Token) string
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
		tok, err := p.GetToken(code)
		if err != nil {
			fmt.Printf("Error provider.go L33: %+v\n", err)
		}
		resp := p.GetIdentity(tok)
		w.Write([]byte(resp))
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
		tok, err := p.GetToken(code)
		if err != nil {
			panic(err) // For now
		}
		_ = p.GetIdentity(tok)
	})
}
