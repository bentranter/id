package id

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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
			fmt.Fprintf(w, "ERROR: %s\n", err)
		}
		user, err := p.GetIdentity(tok)
		if err != nil {
			fmt.Fprintf(w, "ERROR: %s\n", err)
		}

		cookie := genToken(user)
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
			fmt.Fprintf(w, "Error: %s\n", err)
		}
		user, err := p.GetIdentity(tok)
		cookie := genToken(user)
		http.SetCookie(w, cookie)

		http.Redirect(w, r, "http://localhost:3000/"+redirectURL, http.StatusFound)
	})
}

// I NEED TO READ THE SPEC:
//
// http://tools.ietf.org/html/rfc7519
func genToken(user *User) *http.Cookie {
	jwt := jwt.New(jwt.SigningMethodHS256)

	// Claims defined in the spec
	jwt.Claims["iss"] = "YOUR_SITE_NAME_OR_URI"
	jwt.Claims["sub"] = user.ID
	jwt.Claims["aud"] = "YOUR_SITE_NAME_OR_URI"
	jwt.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	jwt.Claims["iat"] = time.Now().Unix()
	jwt.Claims["jti"] = "state" // Figure out what to do about this... it's techinically used to prevent replay attacks

	// These are optional/not in spec
	jwt.Claims["name"] = user.Name
	jwt.Claims["email"] = user.Email
	jwt.Claims["id"] = user.ID
	jwt.Claims["role"] = "user"

	tokStr, err := jwt.SignedString([]byte("SECURE_KEY_HERE"))
	if err != nil {
		fmt.Printf("Error signing string: %s\n", err)
	}

	// Maybe use the access token expiry time in the raw
	// expires...
	return &http.Cookie{
		Name:       "id",
		Value:      tokStr,
		Path:       "/",
		RawExpires: "0",
		// Eventually, you'll need `secure` to be true
		HttpOnly: true,
	}
}
