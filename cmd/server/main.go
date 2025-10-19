package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"url-shortener/internal/handlers"
	"url-shortener/internal/middleware"
	"url-shortener/internal/repository"
	"url-shortener/internal/services"
	"url-shortener/pkg/database"
)

func main() {
	_ = godotenv.Load()

	useMongo := os.Getenv("USE_MONGO")

	var urlService *services.URLService
	if useMongo == "true" {
		db := database.ConnectMongo()
		repo := repository.NewMongoRepo(db, os.Getenv("MONGO_COLLECTION"))
		urlService = services.NewURLService(repo)
		log.Println("🚀 Running with MongoDB repo")
	} else {
		repo := repository.NewMemoryRepo()
		urlService = services.NewURLService(repo)
		log.Println("🚀 Running with in-memory repo")
	}

	handler := &handlers.Handler{Service: urlService}

	// build router (http.Handler)
	router := handlers.NewRouter(handler) // returns http.Handler (e.g., *http.ServeMux)

	// create rate limiter from ENV and wrap router
	rl := middleware.NewRateLimiterFromEnv()
	limitedHandler := rl.Middleware(router)

	addr := ":8000"
	server := &http.Server{
		Addr:         addr,
		Handler:      limitedHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Printf("Server started at http://localhost%s\n", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
