package main

import (
	"log"
	"net/http"

	"url-shortener/internal/handlers"
	"url-shortener/internal/repository"
	"url-shortener/internal/services"
)

func main() {
	// Wire dependencies
	repo := repository.NewMemoryRepo()
	service := services.NewURLService(repo)
	handler := &handlers.Handler{Service: service}

	// Setup router
	router := handlers.NewRouter(handler)

	log.Println("Server started at http://localhost:8000")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatal(err)
	}
}
