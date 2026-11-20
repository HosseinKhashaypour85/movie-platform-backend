package main

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"movie-platform/modules/media"
)

type gateway struct {
	storage *media.S3Presigner
	client  *http.Client
}

func (g *gateway) serveObject(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/")
	if key == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// PresignGet creates a short-lived origin URL. It is used only server-side;
	// the client never receives or gets redirected to the storage URL.
	originURL, err := g.storage.PresignGet(key, time.Now().UTC(), 60*time.Second)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, originURL, nil)
	if err != nil {
		http.Error(w, "bad gateway", http.StatusBadGateway)
		return
	}
	for _, name := range []string{"Range", "If-Range"} {
		if value := r.Header.Get(name); value != "" {
			req.Header.Set(name, value)
		}
	}

	resp, err := g.client.Do(req)
	if err != nil {
		http.Error(w, "media origin unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		if resp.StatusCode == http.StatusNotFound {
			http.NotFound(w, r)
		} else {
			http.Error(w, "media origin unavailable", http.StatusBadGateway)
		}
		return
	}

	for _, name := range []string{
		"Content-Type", "Content-Length", "Content-Range",
		"Accept-Ranges", "ETag", "Last-Modified",
	} {
		if value := resp.Header.Get(name); value != "" {
			w.Header().Set(name, value)
		}
	}
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(resp.StatusCode)
	if r.Method == http.MethodHead {
		return
	}
	if _, err := io.Copy(w, resp.Body); err != nil {
		// The response may already be partially written. Log only; never expose
		// origin URLs or storage credentials in the client response.
		slog.Warn("media stream interrupted", "error", err)
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	storage, err := media.NewS3PresignerFromEnv()
	if err != nil || storage == nil {
		logger.Error("private S3 storage configuration is required", "error", err)
		os.Exit(1)
	}
	signingKey := []byte(os.Getenv("PLAYBACK_SIGNING_KEY"))
	if len(signingKey) < 32 {
		logger.Error("PLAYBACK_SIGNING_KEY must contain at least 32 bytes")
		os.Exit(1)
	}

	client := &http.Client{
		Timeout: 0, // HLS objects stream; enforce timeouts at the server/proxy layer.
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return errors.New("media origin redirects are disabled")
		},
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   20,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 15 * time.Second,
			ExpectContinueTimeout: time.Second,
		},
	}
	g := &gateway{storage: storage, client: client}

	router := chi.NewRouter()
	router.Get("/health/live", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("alive"))
	})
	router.Get("/health/ready", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	})
	router.Handle("/*", media.HLSCookieMiddleware(media.EdgeAuthConfig{
		SigningKey: signingKey,
	})(http.HandlerFunc(g.serveObject)))

	port := strings.TrimSpace(os.Getenv("MEDIA_GATEWAY_PORT"))
	if port == "" {
		port = "8090"
	}
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       90 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	go func() {
		logger.Info("media gateway starting", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("media gateway failed", "error", err)
			os.Exit(1)
		}
	}()

	stop, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	<-stop.Done()
	ctx, done := context.WithTimeout(context.Background(), 15*time.Second)
	defer done()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("media gateway shutdown failed", "error", err)
	}
}
