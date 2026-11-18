package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"movie-platform/internal/middleware"
	"movie-platform/modules/users"
	"movie-platform/modules/watch_history"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 1. بارگذاری متغیرهای محیطی از فایل .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or couldn't be loaded, using system env")
	}

	// 2. اتصال به دیتابیس با استفاده از اطلاعات .env
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// 3. اجرای مایگریشن جداول
	if err := db.AutoMigrate(
		&users.User{},
		&watch_history.WatchHistory{},
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// 4. تنظیم روتر و مسیرها
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("api is running successfully\n"))
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/api/v1/profile", func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
			if !ok {
				http.Error(w, "Invalid user context", http.StatusUnauthorized)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(map[string]any{
				"message": "Authenticated successfully",
				"user_id": userID,
			}); err != nil {
				log.Printf("encode profile response: %v", err)
			}
		})
	})

	// 5. تنظیم پورت و اجرای سرور
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("server is running on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
