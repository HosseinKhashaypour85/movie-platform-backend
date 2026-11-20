package watch_history

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Route("/watch-history", func(r chi.Router) {
		r.Get("/episodes", h.EpisodeHistory)
		r.Put("/episodes/{episodeID}", h.SaveEpisode)
		r.Delete("/episodes/{episodeID}", h.DeleteEpisodeHistory)
		r.Get("/", h.List)
		r.Put("/{movieID}", h.Save)
		r.Get("/{movieID}", h.Get)
		r.Delete("/{movieID}", h.Delete)
	})
}
