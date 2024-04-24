package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
	"moviemicroservice.com/gen"
	"moviemicroservice.com/pkg/discovery"
	"moviemicroservice.com/pkg/discovery/consul"
	grpcHandler "moviemicroservice.com/services/ratings/internal/handler/grpc"
	"moviemicroservice.com/services/ratings/internal/repository/pg"
	"moviemicroservice.com/services/ratings/internal/service/ratings"
)

const serviceName = "ratings"

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

	log.Printf("ratings service starting up on port %d...", port)

	//start consul on localhost:8500
	registry, err := consul.New("localhost:8500")
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
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

	db_ctx, db_cancel := context.WithCancel(context.Background())

	db_url := os.Getenv("RATINGS_DB_URL")
	repo, err := pg.New(db_ctx, db_url)
	if err != nil {
		panic(err)
	}

	println("ratings db connection successful")

	defer repo.CloseConnection(db_ctx)

	service := ratings.New(repo, nil)
	handler := grpcHandler.New(service)

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatalf("failed to listen on port. %v", err)
	}

	server := grpc.NewServer()
	reflection.Register(server)
	gen.RegisterRatingServiceServer(server, handler)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		s := <-sigChan
		cancel()
		db_cancel()
		log.Printf("Received signal %v, attempting graceful shutdown", s)
		server.GracefulStop()
		log.Println("Gracefully stopped the ratings service")
	}()

	if err := server.Serve(listener); err != nil {
		panic(err)
	}
	wg.Wait()
}
