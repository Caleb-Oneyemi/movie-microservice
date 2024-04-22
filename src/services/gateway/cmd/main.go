package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"moviemicroservice.com/src/gen"
	"moviemicroservice.com/src/pkg/discovery"
	"moviemicroservice.com/src/pkg/discovery/consul"
	metadata "moviemicroservice.com/src/services/gateway/internal/api/metadata/grpc"
	ratings "moviemicroservice.com/src/services/gateway/internal/api/ratings/grpc"
	grpcHandler "moviemicroservice.com/src/services/gateway/internal/handler/grpc"
	"moviemicroservice.com/src/services/gateway/internal/services/movies"
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

	service := movies.New(ratingsApi, metadataApi)
	handler := grpcHandler.New(service)

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	reflection.Register(server)

	gen.RegisterMovieServiceServer(server, handler)
	if err := server.Serve(listener); err != nil {
		panic(err)
	}

}
