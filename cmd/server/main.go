package main

import (
	"log"
	"net/http"
	"os"

	"url-shortener/internal/handlers"
	"url-shortener/internal/repository"
	"url-shortener/internal/services"
	"url-shortener/pkg/database"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	useMongo := os.Getenv("USE_MONGO")

	// ✅ declare once here
	var service *services.URLService

	if useMongo == "true" {
		db := database.ConnectMongo()
		mongoRepo := repository.NewMongoRepo(db, os.Getenv("MONGO_COLLECTION"))
		service = services.NewURLService(mongoRepo) // ✅ just assign
		log.Println("🚀 Running with MongoDB repo")
	} else {
		repo := repository.NewMemoryRepo()
		service = services.NewURLService(repo) // ✅ just assign
		log.Println("🚀 Running with in-memory repo")
	}

	handler := &handlers.Handler{Service: service}

	// Setup router
	router := handlers.NewRouter(handler)

	log.Println("Server started at http://localhost:8000")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatal(err)
	}
}
