package main

import (
	"log"
	"net/http"

	metadata "moviemicroservice.com/src/modules/gateway/internal/api/metadata/http"
	ratings "moviemicroservice.com/src/modules/gateway/internal/api/ratings/http"
	httpHandler "moviemicroservice.com/src/modules/gateway/internal/handler/http"
	"moviemicroservice.com/src/modules/gateway/internal/services/movies"
)

func main() {
	log.Println("movie gateway starting up on port on port 8083...")

	metadataApi := metadata.New("localhost:8081")
	ratingsApi := ratings.New("localhost:8082")

	srv := movies.New(ratingsApi, metadataApi)

	h := httpHandler.New(srv)
	http.Handle("/movies", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
