package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"movie-platform/internal/config"
	"movie-platform/internal/database"
	"movie-platform/internal/httpx"
	"movie-platform/internal/middleware"
	"movie-platform/modules/admin"
	"movie-platform/modules/auth"
	"movie-platform/modules/catalog"
	"movie-platform/modules/engagement"
	"movie-platform/modules/users"
	"movie-platform/modules/watch_history"

	"github.com/go-chi/chi/v5"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.Load()
	if err != nil {
		logger.Error("configuration failed", "error", err)
		os.Exit(1)
	}

	db, err := database.Open(cfg)
	if err != nil {
		logger.Error("database failed", "error", err)
		os.Exit(1)
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("database pool failed", "error", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	if cfg.Environment != "production" {
		if err := db.AutoMigrate(&users.User{}, &auth.Session{}, &watch_history.WatchHistory{}, &catalog.Genre{}, &catalog.Movie{}, &catalog.Series{}, &catalog.Season{}, &catalog.Episode{}, &engagement.Favorite{}, &engagement.Rating{}, &admin.AuditLog{}); err != nil {
			logger.Error("migration failed", "error", err)
			os.Exit(1)
		}
	}

	userRepo := users.NewRepository(db)
	userHandler := users.NewHandler(db)
	authHandler := auth.NewHandler(auth.NewService(userRepo, auth.NewSessionRepository(db), cfg))
	historyHandler := watch_history.NewHandler(watch_history.NewService(watch_history.NewRepository(db)))
	catalogHandler := catalog.NewHandler(catalog.NewService(catalog.NewRepository(db)))
	engagementHandler := engagement.NewHandler(engagement.NewService(engagement.NewRepository(db)))
	adminHandler := admin.NewHandler(db)

	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.Recover(logger), middleware.SecurityHeaders, middleware.BodyLimit(1<<20), middleware.Timeout(15*time.Second), middleware.RateLimit(120, time.Minute), middleware.CORS(cfg.AllowedOrigins))
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		httpx.JSON(w, 200, map[string]string{"service": "movie-platform-api", "status": "running"})
	})
	router.Get("/health/live", func(w http.ResponseWriter, r *http.Request) { httpx.JSON(w, 200, map[string]string{"status": "alive"}) })
	router.Get("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(ctx); err != nil {
			httpx.Error(w, 503, "not_ready", "dependencies unavailable")
			return
		}
		httpx.JSON(w, 200, map[string]string{"status": "ready"})
	})

	router.Route("/api/v1", func(r chi.Router) {
		auth.RegisterRoutes(r, authHandler)
		catalog.RegisterRoutes(r, catalogHandler)
		engagement.RegisterPublicRoutes(r, engagementHandler)

		r.Group(func(protected chi.Router) {
			protected.Use(middleware.Auth(cfg.JWTSecret, cfg.JWTIssuer))
			admin.RegisterRoutes(protected, adminHandler)
		})

		r.Group(func(p chi.Router) {
			p.Use(middleware.Auth(cfg.JWTSecret, cfg.JWTIssuer))
			users.RegisterRoutes(p, userHandler)
			watch_history.RegisterRoutes(p, historyHandler)
			engagement.RegisterProtectedRoutes(p, engagementHandler)
		})
	})

	server := &http.Server{Addr: ":" + cfg.Port, Handler: router, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 1 << 20}
	go func() {
		logger.Info("server starting", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()
	stop, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	<-stop.Done()
	ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("shutdown failed", "error", err)
	}
}
