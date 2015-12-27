package id

import (
	"crypto/rand"
	"errors"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	signingKey = generateRandomBytes()

	// ErrTokenInvalid means the token wasn't valid based
	// on the value of its signature.
	ErrTokenInvalid = errors.New("Invalid token")
	// ErrCannotParseToken means that the token could not
	//be parsed.
	ErrCannotParseToken = errors.New("Cannot parse token")
	// ErrInvalidSigningMethod means that the method used
	// to sign the token was not the expected method, or
	// is an invalid method.
	ErrInvalidSigningMethod = errors.New("Invalid signing method")
	// ErrNoSigningKey means that a signing key doesn't
	// exist.
	ErrNoSigningKey = errors.New("No signing key")
)

// AuthInit overrides the randomly generated signing key.
// Useful for users who want to use the same key across
// server restarts, so users don't lose their session.
func AuthInit(key []byte) {
	signingKey = key
}

// GenToken generates a new JSON web token from a user.
func GenToken(user *User) (*http.Cookie, error) {
	jwt := jwt.New(jwt.SigningMethodHS256)
	tokExp := time.Now().Add(time.Hour * 72).Unix()

	// Claims defined in the spec
	jwt.Claims["sub"] = user.ID
	jwt.Claims["exp"] = tokExp
	jwt.Claims["iat"] = time.Now().Unix()

	// These are optional/not in spec. They're used to
	// to determine who's signed in, and their role
	jwt.Claims["name"] = user.Name
	jwt.Claims["email"] = user.Email
	jwt.Claims["id"] = user.ID
	jwt.Claims["role"] = "user"

	// Sanity check to make sure signing key is longer
	// than zero
	if len(signingKey) == 0 {
		return nil, ErrNoSigningKey
	}

	tokStr, err := jwt.SignedString(signingKey)
	if err != nil {
		return nil, err
	}

	// Maybe use the access token expiry time in the raw
	// expires...
	return &http.Cookie{
		Name:       "id",
		Value:      tokStr,
		Path:       "/",
		RawExpires: string(tokExp),
		// Eventually, you'll need `secure` to be true
		HttpOnly: true,
	}, nil
}

// Verify checks to make sure there is a cookie with a
// valid JWT.
func Verify(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("id")
	if err != nil {
		return err
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if ok == false {
			return nil, ErrInvalidSigningMethod
		}
		return signingKey, nil
	})
	if err != nil {
		return err
	}

	if token.Valid == false {
		return ErrTokenInvalid
	}

	return nil
}

// ExpireCookie sets the expiry on the cookie. It will not send the request.
func ExpireCookie(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("user")
	cookie.Value = ""
	cookie.RawExpires = string(time.UnixDate)
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}

// Verified is just a simple check to make sure that a
// user is authenticated.
func Verified(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("You're authenticated"))
	return nil
}

func generateRandomBytes() []byte {
	// Use 32 bytes (256 bits) to satisfy the requirement
	// for the HMAC key length.
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		// If this errors, it means that something is wrong
		// the system's CSPRNG, which indicates a critical
		// operating system failure. Panic and crash here
		panic(err)
	}
	return b
}
