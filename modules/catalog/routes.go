package catalog

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Route("/catalog", func(r chi.Router) {
		r.Get("/genres", h.Genres)
		r.Get("/movies", h.Movies)
		r.Get("/movies/{id}", h.Movie)
		r.Get("/series", h.Series)
		r.Get("/series/{id}", h.SeriesItem)
		r.Get("/series/{id}/seasons", h.Seasons)
		r.Get("/seasons/{id}/episodes", h.Episodes)
	})
}
