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
