package main

import (
	"log"
	"net/http"

	"github.com/bentranter/psa"
	"github.com/bentranter/psa/providers/facebook"

	// "github.com/gorilla/mux"
	// "github.com/gorilla/pat"

	// "github.com/julienschmidt/httprouter"
)

func main() {
	// Initialize any provider
	provider := facebook.New()

	// Bare http
	http.Handle("/auth/facebook/authorize", psa.Authorize(provider))
	http.Handle("/auth/facebook/callback", psa.Callback(provider))
	http.Handle("/auth/restricted", psa.Middleware(psa.Verify, psa.Verified))

	// // Default mux
	// serveMux := http.NewServeMux()
	// serveMux.HandleFunc("/auth/twitch/authorize", psa.Authorize(provider))
	// serveMux.HandleFunc("/auth/twitch/callback", psa.Callback(provider))

	// // Gorilla's Pat
	// p := pat.New()
	// p.Get("/auth/twitch/authorize", psa.Authorize(provider))
	// p.Get("/auth/twitch/callback", psa.Callback(provider))

	// // Gorilla's Mux
	// m := mux.NewRouter()
	// m.HandleFunc("/auth/twitch/authorize", psa.Authorize(provider))
	// m.HandleFunc("/auth/twitch/callback", psa.Callback(provider))

	// // Julien Schmidt's httprouter
	// r := httprouter.New()
	// r.GET("/httprouter/auth/gplus/authorize", psa.HTTPRouterAuthorize(provider))
	// r.GET("/httprouter/auth/gplus/callback", psa.HTTPRouterCallback(provider))

	log.Fatal(http.ListenAndServe(":3000", nil))
}
