package grpc

import (
	"context"

	"moviemicroservice.com/gen"
	grpcUtils "moviemicroservice.com/internal/utils"
	"moviemicroservice.com/pkg/discovery"
	"moviemicroservice.com/services/metadata/pkg/models"
)

type Api struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Api {
	return &Api{registry: registry}
}

func (g *Api) Get(ctx context.Context, id string) (*models.MetaData, error) {
	conn, err := grpcUtils.GetServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	client := gen.NewMetadataServiceClient(conn)
	resp, err := client.GetMetadata(ctx, &gen.GetMetadataRequest{MovieId: id})

	if err != nil {
		return nil, err
	}

	return models.MetadataFromProto(resp.Data), nil
}
