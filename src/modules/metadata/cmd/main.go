package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	httpHandler "moviemicroservice.com/src/modules/metadata/internal/handler/http"
	"moviemicroservice.com/src/modules/metadata/internal/repository/memory"
	"moviemicroservice.com/src/modules/metadata/internal/service/metadata"
	"moviemicroservice.com/src/pkg/discovery"
	"moviemicroservice.com/src/pkg/discovery/consul"
)

const serviceName = "metadata"

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "metadata service port")
	flag.Parse()

	log.Printf("metadata service starting up on port %d...", port)

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

			time.Sleep(3 * time.Second)
		}
	}()

	//deregister once process terminates
	defer registry.Deregister(ctx, instanceID, serviceName)

	repo := memory.New()
	service := metadata.New(repo)
	handler := httpHandler.New(service)

	http.Handle("/api/v1/metadata", http.HandlerFunc(handler.Get))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}

}
