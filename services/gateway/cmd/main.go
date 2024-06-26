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

	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
	"moviemicroservice.com/gen"
	"moviemicroservice.com/pkg/discovery"
	"moviemicroservice.com/pkg/discovery/consul"
	metadata "moviemicroservice.com/services/gateway/internal/api/metadata/grpc"
	ratings "moviemicroservice.com/services/gateway/internal/api/ratings/grpc"
	grpcHandler "moviemicroservice.com/services/gateway/internal/handler/grpc"
	"moviemicroservice.com/services/gateway/internal/services/movies"
)

const serviceName = "gateway"

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
	log.Printf("Starting the gateway service on port %d", port)

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

	const limit = 100
	const burst = 100
	limiter := newLimiter(limit, burst)

	server := grpc.NewServer(grpc.UnaryInterceptor(ratelimit.UnaryServerInterceptor(limiter)))
	reflection.Register(server)
	gen.RegisterMovieServiceServer(server, handler)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		s := <-sigChan
		cancel()
		log.Printf("Received signal %v, attempting graceful shutdown", s)
		server.GracefulStop()
		log.Println("Gracefully stopped the gateway service")
	}()

	if err := server.Serve(listener); err != nil {
		panic(err)
	}
	wg.Wait()

}

type limiter struct {
	l *rate.Limiter
}

func newLimiter(limit int, burst int) *limiter {
	return &limiter{rate.NewLimiter(rate.Limit(limit), burst)}
}

func (l *limiter) Limit() bool {
	return l.l.Allow()
}
