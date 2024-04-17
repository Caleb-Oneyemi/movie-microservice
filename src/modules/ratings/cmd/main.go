package main

import (
	"log"
	"net/http"

	httpHandler "moviemicroservice.com/src/modules/ratings/internal/handler/http"
	"moviemicroservice.com/src/modules/ratings/internal/repository/memory"
	"moviemicroservice.com/src/modules/ratings/internal/service/ratings"
)

func main() {
	log.Println("ratings service starting up on port 8082...")

	repo := memory.New()
	service := ratings.New(repo)
	handler := httpHandler.New(service)

	http.Handle("/api/v1/ratings", http.HandlerFunc(handler.Handle))
	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
