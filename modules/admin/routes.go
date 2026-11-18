package admin

import (
 "movie-platform/internal/middleware"
 "github.com/go-chi/chi/v5"
)
func RegisterRoutes(r chi.Router,h *Handler){
 r.Route("/admin",func(r chi.Router){
  r.Use(middleware.RequireRole("admin","super_admin"))
  r.Get("/dashboard",h.Dashboard)
 })
}
