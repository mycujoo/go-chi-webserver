package handler

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/mycujoo/go-chi-webserver/pkg/error"
)

// Ping handler for returning health status
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK"))
}

// NoRoute handler for catching all incorrect routes
func NoRoute(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, error.ErrNotFound())
}
