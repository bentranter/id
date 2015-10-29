package main

import (
	"log"
	"net/http"

	"github.com/bentranter/psa"
	"github.com/bentranter/psa/providers/google"

	"github.com/gorilla/mux"
	"github.com/gorilla/pat"

	//"github.com/julienschmidt/httprouter"
)

func main() {
	// Initialize any provider
	provider := google.New()

	// Bare http
	http.Handle("/auth/gplus/authorize", psa.Authorize(provider))
	http.Handle("/auth/gplus/callback", psa.Callback(provider))

	// Default mux
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/auth/twitch/authorize", psa.Authorize(provider))
	serveMux.HandleFunc("/auth/twitch/callback", psa.Callback(provider))

	// Gorilla's Pat
	p := pat.New()
	p.Get("/auth/twitch/authorize", psa.Authorize(provider))
	p.Get("/auth/twitch/callback", psa.Callback(provider))

	// Gorilla's Mux
	m := mux.NewRouter()
	m.HandleFunc("/auth/twitch/authorize", psa.Authorize(provider))
	m.HandleFunc("/auth/twitch/callback", psa.Callback(provider))

	// Julien Schmidt's httprouter
	// r := httprouter.New()
	// r.GET("/httprouter/auth/twitch/authorize", psa.HTTPRouterAuthorize(provider))
	// r.GET("/httprouter/auth/twitch/callback", psa.HTTPRouterCallback(provider))

	log.Fatal(http.ListenAndServe(":3000", nil))
}
