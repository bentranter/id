package main

import (
	"log"
	"net/http"

	"github.com/bentranter/id"
	"github.com/bentranter/id/providers/google"

	// "github.com/gorilla/mux"
	// "github.com/gorilla/pat"

	// "github.com/julienschmidt/httprouter"
)

func main() {
	// Initialize any provider
	provider := google.New()

	// Bare http
	http.Handle("/auth/gplus/authorize", id.Authorize(provider))
	http.Handle("/auth/gplus/callback", id.Callback(provider, "auth/restricted"))
	http.Handle("/auth/restricted", id.Middleware(id.Verify, id.Verified))

	// // Default mux
	// serveMux := http.NewServeMux()
	// serveMux.HandleFunc("/auth/twitch/authorize", id.Authorize(provider))
	// serveMux.HandleFunc("/auth/twitch/callback", id.Callback(provider))

	// // Gorilla's Pat
	// p := pat.New()
	// p.Get("/auth/twitch/authorize", id.Authorize(provider))
	// p.Get("/auth/twitch/callback", id.Callback(provider))

	// // Gorilla's Mux
	// m := mux.NewRouter()
	// m.HandleFunc("/auth/twitch/authorize", id.Authorize(provider))
	// m.HandleFunc("/auth/twitch/callback", id.Callback(provider))

	// // Julien Schmidt's httprouter
	// r := httprouter.New()
	// r.GET("/httprouter/auth/gplus/authorize", id.HTTPRouterAuthorize(provider))
	// r.GET("/httprouter/auth/gplus/callback", id.HTTPRouterCallback(provider))

	log.Fatal(http.ListenAndServe(":3000", nil))
}
