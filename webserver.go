package webserver

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	h "github.com/mycujoo/go-chi-webserver/pkg/handler"
)

// SetupRouter function for creating router instance
func SetupRouter(env string) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
	if env != "production" {
		router.Use(middleware.Logger)
	}
	router.Use(middleware.DefaultCompress)
	router.Use(middleware.Recoverer)

	router.Get("/", h.Ping)
	router.NotFound(h.NoRoute)

	return router
}

// Listen and serve HTTP server
func Listen(addr string, router *chi.Mux) {
	http.ListenAndServe(addr, router)
}
