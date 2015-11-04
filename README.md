You should be able to do this:

```go
package main

import "github.com/bentranter/thrill/providers/twitch"

func main() {
    provider := twitch.New(clientID, clientSecret, callbackURL)

    http.HandleFunc("/auth/twitch/authorize", provider.Authorize)
    http.HandleFunc("/auth/twitch/callback", provider.Callback)

    http.ListenAndServe(":3000", nil)
}
```

It works with Gorilla's Mux, Pat, Julien Schmidt's HttpRouter, and of course the standard `net/http` package.

Need to figure out middleware approach... check Matt and Julien's stuff.

* Figure out why it won't work with Twitch
* Solve LinkedIn opaque URL problem (think you can tweak http.Client)
* Add param for redirect after callback URL fires


**FIGURE OUT WHY THIS SOMETIMES PANICS???**

Maybe it's a race condition? The trace shows `psa.genToken(0x0, 0x0, 0x0)`, and it only happens sometimes... Can I lock it?

# WARNING

This is extremely insecure. All errors need to be handled surrounding the token, the verification of the token isn't correct, the HMAC might not be occurring correctly, a token is still generated when `psa.User.ID/Email/Name` is `nil`... basically it's bad.
