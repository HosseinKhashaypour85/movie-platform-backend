package admin

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"movie-platform/internal/middleware"
)

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.RequireRole("admin", "super_admin"))
		r.Get("/dashboard", h.Dashboard)
		r.Get("/audit-logs", h.ListAuditLogs)
		r.Use(AuditMiddleware(h.db, slog.Default()))
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
		r.Route("/seasons", func(r chi.Router) {
			r.Post("/", h.CreateSeason)
			r.Put("/{id}", h.UpdateSeason)
			r.Delete("/{id}", h.DeleteSeason)
		})
		r.Route("/episodes", func(r chi.Router) {
			r.Post("/", h.CreateEpisode)
			r.Put("/{id}", h.UpdateEpisode)
			r.Delete("/{id}", h.DeleteEpisode)
		})
	})
}
