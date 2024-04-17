package main

import (
	"log"
	"net/http"

	httpHandler "moviemicroservice.com/src/modules/metadata/internal/handler/http"
	"moviemicroservice.com/src/modules/metadata/internal/repository/memory"
	"moviemicroservice.com/src/modules/metadata/internal/service/metadata"
)

func main() {
	log.Println("metadata service starting up on port 8081...")
	repo := memory.New()
	service := metadata.New(repo)
	handler := httpHandler.New(service)

	http.Handle("/api/v1/metadata", http.HandlerFunc(handler.Get))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
