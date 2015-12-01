package id

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

// Handler is the function definition for our middleware.
type Handler func(w http.ResponseWriter, r *http.Request) error

// Middleware executes all our middleware.
func Middleware(handlers ...Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range handlers {
			err := handler(w, r)
			if err != nil {
				switch err {

				case ErrTokenInvalid:
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				case ErrCannotParseToken:
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				case ErrInvalidSigningMethod:
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return

				// Errors from jwt-go
				case jwt.ErrInvalidKey:
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				case jwt.ErrHashUnavailable:
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				case jwt.ErrNoTokenInRequest:
					http.Error(w, err.Error(), http.StatusForbidden)
					return

				// Other errors
				default:
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
	})
}
