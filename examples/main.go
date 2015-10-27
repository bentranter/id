package main

import (
	"log"
	"net/http"

	"github.com/bentranter/psa"
	"github.com/bentranter/psa/providers/twitch"
)

func main() {
	provider := twitch.New()

	http.Handle("/", psa.Authorize(provider))
	http.Handle("/auth/twitch/callback", psa.Callback(provider))

	log.Fatal(http.ListenAndServe(":3000", nil))
}
