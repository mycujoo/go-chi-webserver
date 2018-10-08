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

	if env == "production" {
		router.Use(middleware.RealIP)
		router.Use(middleware.RequestID)
		router.Use(middleware.Recoverer)
	} else {
		router.Use(middleware.Logger)
	}

	router.Use(middleware.DefaultCompress)
	router.Use(middleware.SetHeader("Content-type", "application/json; charset=utf-8"))
	router.Use(middleware.SetHeader("Vary", "Origin"))
	router.Use(middleware.SetHeader("X-Content-Type-Options", "nosniff"))
	router.Use(middleware.SetHeader("X-Frame-Options", "SAMEORIGIN"))
	router.Use(middleware.SetHeader("X-XSS-Protection", "1; mode=block"))

	router.Get("/", h.Ping)
	router.NotFound(h.NoRoute)

	return router
}

// Listen and serve HTTP server
func Listen(addr string, router *chi.Mux) {
	http.ListenAndServe(addr, router)
}
