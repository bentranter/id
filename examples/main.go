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
	http.Handle("/auth/gplus/authorize", id.Authorize(provider))
	http.Handle("/auth/gplus/callback", id.Callback(provider, "auth/restricted"))
	http.Handle("/auth/restricted", id.Middleware(id.Verify, id.Verified))

	// Default mux
	serveMux := http.NewServeMux()
	serveMux.Handle("/auth/gplus/authorize", id.Authorize(provider))
	serveMux.Handle("/auth/gplus/callback", id.Callback(provider, "auth/restricted"))
	serveMux.Handle("/auth/restricted", id.Middleware(id.Verify, id.Verified))

	// Gorilla's Pat. Requires type assertion.
	p := pat.New()
	p.Get("/auth/gplus/authorize", id.Authorize(provider).(http.HandlerFunc))
	p.Get("/auth/gplus/callback", id.Callback(provider, "auth/restricted").(http.HandlerFunc))

	// Gorilla's Mux
	m := mux.NewRouter()
	m.Handle("/auth/gplus/authorize", id.Authorize(provider))
	m.Handle("/auth/gplus/callback", id.Callback(provider, "auth/restricted"))

	// Julien Schmidt's httprouter
	r := httprouter.New()
	r.GET("/httprouter/auth/gplus/authorize", id.HTTPRouterAuthorize(provider))
	r.GET("/httprouter/auth/gplus/callback", id.HTTPRouterCallback(provider, "auth/restricted"))

	log.Printf("Serving HTTP on port 3000")
	log.Fatal(http.ListenAndServe(":3000", serveMux))
}
