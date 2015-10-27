package main

import (
	"log"
	"net/http"

	"github.com/bentranter/psa/providers/twitch"
)

func main() {
	http.HandleFunc("/", twitch.Authorize)
	http.HandleFunc("/auth/twitch/callback", twitch.Callback)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
