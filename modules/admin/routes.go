package admin

import (
	"movie-platform/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.RequireRole("admin", "super_admin"))
		r.Get("/dashboard", h.Dashboard)
		r.Route("/movies", func(r chi.Router) {
			r.Post("/", h.CreateMovie)
			r.Put("/{id}", h.UpdateMovie)
			r.Delete("/{id}", h.DeleteMovie)
		})
		r.Route("/series", func(r chi.Router) {
			r.Post("/", h.CreateSeries)
			r.Put("/{id}", h.UpdateSeries)
			r.Delete("/{id}", h.DeleteSeries)
		})
	})
}
