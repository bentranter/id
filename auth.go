package id

import (
	"errors"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	signingKey []byte

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
)

// AuthInit reads the public and private keys used to sign
// the JSON web tokens, and can optionally set the
// `signingKey`. The `signingKey` is only used for HMAC, so
// it currently isn't being used.
func AuthInit(key ...[]byte) {
	switch len(key) {
	case 1:
		signingKey = key[0]
	default:
		signingKey = []byte("DEFAULT_SIGNING_KEY")
	}
}

// GenToken generates a new JSON web token from a user.
func GenToken(user *User) (*http.Cookie, error) {
	jwt := jwt.New(jwt.SigningMethodHS256)
	tokExp := time.Now().Add(time.Hour * 72).Unix()

	// Claims defined in the spec
	jwt.Claims["iss"] = "YOUR_SITE_NAME_OR_URI"
	jwt.Claims["sub"] = user.ID
	jwt.Claims["aud"] = "YOUR_SITE_NAME_OR_URI"
	jwt.Claims["exp"] = tokExp
	jwt.Claims["iat"] = time.Now().Unix()
	jwt.Claims["jti"] = "state" // Figure out what to do about this... it's technically used to prevent replay attacks

	// These are optional/not in spec
	jwt.Claims["name"] = user.Name
	jwt.Claims["email"] = user.Email
	jwt.Claims["id"] = user.ID
	jwt.Claims["role"] = "user"

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
		RawExpires: string(tokExp), // might horribly mess up
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

// Verified is just a simple check to make sure that a
// user is authenticated.
func Verified(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("You're authenticated"))
	return nil
}
