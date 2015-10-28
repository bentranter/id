package main

import (
	"log"
	"net/http"

	"github.com/bentranter/psa"
	"github.com/bentranter/psa/providers/twitch"

	"github.com/gorilla/mux"
	"github.com/gorilla/pat"

	"github.com/julienschmidt/httprouter"
)

func main() {
	// Initialize any provider
	provider := twitch.New()

	// Bare http
	http.HandleFunc("/", psa.Authorize(provider))
	http.HandleFunc("/auth/twitch/callback", psa.Callback(provider))

	// Default mux
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", psa.Authorize(provider))
	serveMux.HandleFunc("/auth/twitch/callback", psa.Callback(provider))

	// Gorilla's Pat
	p := pat.New()
	p.Get("/", psa.Authorize(provider))
	p.Get("/auth/twitch/callback", psa.Callback(provider))

	// Gorilla's Mux
	m := mux.NewRouter()
	m.HandleFunc("/", psa.Authorize(provider))
	m.HandleFunc("/auth/twitch/callback", psa.Callback(provider))

	// Julien Schmidt's httprouter
	r := httprouter.New()
	r.GET("/", psa.HTTPRouterAuthorize(provider))
	r.GET("/auth/twitch/callback", psa.HTTPRouterCallback(provider))

	log.Fatal(http.ListenAndServe(":3000", nil))
}
