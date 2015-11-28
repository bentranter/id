package id

import (
	"fmt"
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
				// Errors from jwt-go
				case jwt.ErrInvalidKey:
					http.Error(w, err.Error(), http.StatusUnauthorized)
				case jwt.ErrHashUnavailable:
					http.Error(w, err.Error(), http.StatusNotFound)
				case jwt.ErrNoTokenInRequest:
					http.Error(w, err.Error(), http.StatusForbidden)

				// Other errors
				default:
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		}
	})
}

// Verify checks to make sure there is a cookie with a
// valid JWT.
func Verify(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("id")
	if err != nil {
		return err
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return "SECURE_STRING_HERE", nil
	})

	fmt.Printf("Token: %+v\n", token)

	return nil
}

// Verified is just a simple check to make sure that a
// user is authenticated.
func Verified(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("You're authenticated"))
	return nil
}
