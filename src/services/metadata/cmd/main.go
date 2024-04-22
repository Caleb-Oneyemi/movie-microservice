package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
	"moviemicroservice.com/src/gen"
	"moviemicroservice.com/src/pkg/discovery"
	"moviemicroservice.com/src/pkg/discovery/consul"
	grpcHandler "moviemicroservice.com/src/services/metadata/internal/handler/grpc"
	"moviemicroservice.com/src/services/metadata/internal/repository/pg"
	"moviemicroservice.com/src/services/metadata/internal/service/metadata"
)

const serviceName = "metadata"

func main() {
	f, err := os.Open("../config/base.yaml")
	if err != nil {
		panic(err)
	}

	var config Config
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		panic(err)
	}

	port := config.Api.Port

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

	db_ctx := context.Background()

	db_url := os.Getenv("METADATA_DB_URL")
	repo, err := pg.New(db_ctx, db_url)
	if err != nil {
		panic(err)
	}

	println("metadata db connection successful")

	defer repo.CloseConnection(db_ctx)

	service := metadata.New(repo)
	handler := grpcHandler.New(service)

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatalf("failed to listen on port. %v", err)
	}

	server := grpc.NewServer()
	reflection.Register(server)

	gen.RegisterMetadataServiceServer(server, handler)
	server.Serve(listener)
}