package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"moviemicroservice.com/src/gen"
	grpcHandler "moviemicroservice.com/src/modules/ratings/internal/handler/grpc"
	"moviemicroservice.com/src/modules/ratings/internal/repository/memory"
	"moviemicroservice.com/src/modules/ratings/internal/service/ratings"
	"moviemicroservice.com/src/pkg/discovery"
	"moviemicroservice.com/src/pkg/discovery/consul"
)

const serviceName = "ratings"

func main() {
	var port int
	flag.IntVar(&port, "port", 8082, "ratings service port")
	flag.Parse()

	log.Printf("ratings service starting up on port %d...", port)

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
	service := ratings.New(repo)
	handler := grpcHandler.New(service)

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatalf("failed to listen on port. %v", err)
	}

	server := grpc.NewServer()
	gen.RegisterRatingServiceServer(server, handler)
	server.Serve(listener)
}
