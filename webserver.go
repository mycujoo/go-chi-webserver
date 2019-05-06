package webserver

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/mycujoo/go-chi-webserver/middleware"
	h "github.com/mycujoo/go-chi-webserver/pkg/handler"
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
	router.Get("/", h.Ping)

	// NoRoute handler for catching all incorrect routes
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.Render(w, r, ErrNotFound())
	})

	return http.ListenAndServe(addr, router)
}
