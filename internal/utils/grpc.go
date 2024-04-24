package utils

import (
	"context"
	"math/rand"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"moviemicroservice.com/pkg/discovery"
)

// selects a random service instance and returns a connection to ir
func GetServiceConnection(ctx context.Context, serviceName string, registry discovery.Registry) (*grpc.ClientConn, error) {
	addrs, err := registry.GetServiceAddresses(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	index := rand.Intn(len(addrs))
	return grpc.Dial(addrs[index], grpc.WithTransportCredentials(insecure.NewCredentials()))
}
