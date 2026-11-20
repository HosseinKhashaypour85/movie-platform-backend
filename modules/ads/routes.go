package ads

import (
	"github.com/go-chi/chi/v5"
	"movie-platform/internal/middleware"
)

func RegisterPublicRoutes(r chi.Router, h *Handler) {
	r.Get("/ads", h.Serve)
	r.Post("/ads/{id}/events", h.Track)
}
func RegisterAdminRoutes(r chi.Router, h *Handler) {
	r.Route("/ads", func(a chi.Router) {
		a.Use(middleware.RequireRole("admin", "super_admin"))
		a.Get("/", h.ListCampaigns)
		a.Post("/", h.CreateCampaign)
		a.Put("/{id}", h.UpdateCampaign)
	})
}
