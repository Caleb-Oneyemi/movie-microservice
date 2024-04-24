package grpc

import (
	"context"

	"moviemicroservice.com/gen"
	grpcUtils "moviemicroservice.com/internal/utils"
	"moviemicroservice.com/pkg/discovery"
	"moviemicroservice.com/services/ratings/pkg/models"
)

type Api struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Api {
	return &Api{registry: registry}
}

func (g *Api) GetAggregatedRatings(ctx context.Context, recordID models.RecordID, recordType models.RecordType) (float64, error) {
	conn, err := grpcUtils.GetServiceConnection(ctx, "ratings", g.registry)
	if err != nil {
		return 0, err
	}

	defer conn.Close()

	client := gen.NewRatingServiceClient(conn)
	resp, err := client.GetAggregatedRatings(ctx, &gen.GetAggregatedRatingsRequest{RecordId: string(recordID), RecordType: string(recordType)})

	if err != nil {
		return 0, err
	}

	return resp.Value, nil
}
