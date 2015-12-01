package main

import (
	"log"
	"net/http"

	"github.com/bentranter/id"
	"github.com/bentranter/id/providers/google"

	"github.com/gorilla/mux"
	"github.com/gorilla/pat"

	"github.com/julienschmidt/httprouter"
)

func main() {
	// Initialize any provider
	provider := google.New()

	// Bare http
	http.Handle("/auth/google/authorize", id.Authorize(provider))
	http.Handle("/auth/google/callback", id.Callback(provider, "auth/restricted"))
	http.Handle("/auth/restricted", id.Middleware(id.Verify, id.Verified))

	// Default mux
	serveMux := http.NewServeMux()
	serveMux.Handle("/auth/google/authorize", id.Authorize(provider))
	serveMux.Handle("/auth/google/callback", id.Callback(provider, "auth/restricted"))
	serveMux.Handle("/auth/restricted", id.Middleware(id.Verify, id.Verified))

	// Gorilla's Pat. Requires type assertion.
	p := pat.New()
	p.Get("/auth/google/authorize", id.Authorize(provider).(http.HandlerFunc))
	p.Get("/auth/google/callback", id.Callback(provider, "auth/restricted").(http.HandlerFunc))

	// Gorilla's Mux
	m := mux.NewRouter()
	m.Handle("/auth/google/authorize", id.Authorize(provider))
	m.Handle("/auth/google/callback", id.Callback(provider, "auth/restricted"))

	// Julien Schmidt's httprouter
	r := httprouter.New()
	r.GET("/httprouter/auth/google/authorize", id.HTTPRouterAuthorize(provider))
	r.GET("/httprouter/auth/google/callback", id.HTTPRouterCallback(provider, "auth/restricted"))

	log.Printf("Serving HTTP on port 3000")
	log.Fatal(http.ListenAndServe(":3000", serveMux))
}
