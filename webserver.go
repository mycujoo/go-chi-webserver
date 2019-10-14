package webserver

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/mycujoo/go-chi-webserver/middleware"
)

// SetupRouter function for creating router instance
func SetupRouter(env string) *chi.Mux {
	router := chi.NewRouter()

	router.Use(chiMiddleware.RealIP)
	router.Use(chiMiddleware.RequestID)
	if env != "production" {
		router.Use(chiMiddleware.Logger)
	}
	router.Use(chiMiddleware.DefaultCompress)
	router.Use(chiMiddleware.Recoverer)

	// NoRoute handler for catching all incorrect routes
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.Render(w, r, ErrNotFound())
	})

	return router
}

func EnableMetrics(name string, r *chi.Mux) {
	m := middleware.NewPrometheus(name)
	r.Use(m)
	r.Handle("/metrics", prometheus.Handler())
}

// Listen and serve HTTP server.
// All default routes must be defined after middlewares.
func Listen(addr string, router *chi.Mux) error {
	return http.ListenAndServe(addr, router)
}
