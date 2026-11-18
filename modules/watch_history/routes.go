package watch_history

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Route("/watch-history", func(r chi.Router) {
		r.Get("/", h.List)
		r.Put("/{movieID}", h.Save)
		r.Get("/{movieID}", h.Get)
		r.Delete("/{movieID}", h.Delete)
	})
}
