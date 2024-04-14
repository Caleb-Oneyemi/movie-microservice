package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"moviemicroservice.com/src/modules/metadata/pkg/models"
	"moviemicroservice.com/src/modules/movies/internal/gateway"
)

type Gateway struct {
	address string
}

func New(address string) *Gateway {
	return &Gateway{address}
}

func (g *Gateway) Get(ctx context.Context, id string) (*models.MetaData, error) {
	req, err := http.NewRequest(http.MethodGet, g.address+"/api/v1/metadata", nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	queries := req.URL.Query()
	queries.Add("id", id)
	req.URL.RawQuery = queries.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", resp)
	}

	var meta *models.MetaData
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return nil, err
	}

	return meta, nil
}
