package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/cors" // <-- 1. IMPORT "cors"

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

	// --- 2. ADD CORS MIDDLEWARE ---
	// IMPORTANT: Change "https://your-frontend.vercel.app" to your *actual* Vercel URL
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"https://your-frontend.vercel.app", "http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	// Wrap your rate-limited handler with the CORS handler
	finalHandler := c.Handler(limitedHandler)
	// --- END CORS SECTION ---

	// --- 3. GET PORT FROM HEROKU ENV ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default for local testing
	}
	addr := ":" + port // Construct the address for ListenAndServe
	// --- END PORT SECTION ---

	server := &http.Server{
		Addr:         addr,         // <-- 4. USE THE DYNAMIC 'addr'
		Handler:      finalHandler, // <-- 5. USE THE FINAL 'finalHandler' (with CORS)
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// --- 6. UPDATE LOG MESSAGE ---
	log.Printf("Server started on port %s\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
