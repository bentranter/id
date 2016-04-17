![id (Golang)](https://github.com/bentranter/id/raw/master/assets/logo.png "id")

[![GoDoc](https://godoc.org/github.com/bentranter/id?status.svg)](https://godoc.org/github.com/bentranter/id)

Sessionless, passwords was authentication. 


### How

JSON web tokens + OAuth identity providers.

All you need to do is this:

```go
package main

import (
    "net/http"

    "github.com/bentranter/id"
    "github.com/bentranter/id/providers/facebook"
)

func main() {
    provider := facebook.New("<your-client-id>", "<your-client-secret>", "<your-client-callback-url>")

    http.Handle("/auth/facebook/authorize", id.Authorize(provider))
    http.Handle("/auth/facebook/callback", id.Callback(provider, "<your-redirect-url>"))

    http.ListenAndServe(":3000", nil)
}
```

It works with Gorilla's Mux, Pat, Julien Schmidt's HttpRouter, and of course the standard `net/http` package.

### Looking for a more "normal" password flow?

Checkout [Password](https://github.com/bentranter/password)
