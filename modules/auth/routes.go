package auth

import (
	"github.com/go-chi/chi/v5"
	"movie-platform/internal/middleware"
)

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Post("/auth/register", h.Register)
	r.Post("/auth/login", h.Login)
	r.Post("/auth/refresh", h.Refresh)
	r.Group(func(p chi.Router) {
		p.Use(middleware.Auth(h.Service.Config.JWTSecret, h.Service.Config.JWTIssuer))
		p.Post("/auth/logout", h.Logout)
		p.Post("/auth/logout-all", h.LogoutAll)
	})
}
