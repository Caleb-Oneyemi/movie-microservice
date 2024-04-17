package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	metadata "moviemicroservice.com/src/modules/gateway/internal/api/metadata/http"
	ratings "moviemicroservice.com/src/modules/gateway/internal/api/ratings/http"
	httpHandler "moviemicroservice.com/src/modules/gateway/internal/handler/http"
	"moviemicroservice.com/src/modules/gateway/internal/services/movies"
	"moviemicroservice.com/src/pkg/discovery"
	"moviemicroservice.com/src/pkg/discovery/consul"
)

const serviceName = "gateway"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()

	log.Printf("Starting the movie service on port %d", port)

	//start consul on localhost:8500
	registry, err := consul.New("localhost:8500")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}

	//continuously ping consul every 3 seconds in goroutine
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	//deregister once process terminates
	defer registry.Deregister(ctx, instanceID, serviceName)

	metadataApi := metadata.New(registry)
	ratingsApi := ratings.New(registry)

	srv := movies.New(ratingsApi, metadataApi)
	h := httpHandler.New(srv)

	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}

}
