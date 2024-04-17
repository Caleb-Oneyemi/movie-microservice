package main

import (
	"log"
	"net/http"

	metadataGw "moviemicroservice.com/src/modules/movies/internal/gateway/metadata/http"
	ratingsGw "moviemicroservice.com/src/modules/movies/internal/gateway/ratings/http"
	httpHandler "moviemicroservice.com/src/modules/movies/internal/handler/http"
	"moviemicroservice.com/src/modules/movies/internal/services/movies"
)

func main() {
	log.Println("movie gateway starting up on port on port 8083...")

	metadataGateway := metadataGw.New("localhost:8081")
	ratingGateway := ratingsGw.New("localhost:8082")

	srv := movies.New(ratingGateway, metadataGateway)

	h := httpHandler.New(srv)
	http.Handle("/movies", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
