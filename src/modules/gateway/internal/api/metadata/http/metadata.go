package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	api "moviemicroservice.com/src/modules/gateway/internal/api"
	"moviemicroservice.com/src/modules/metadata/pkg/models"
	"moviemicroservice.com/src/pkg/discovery"
)

type Api struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Api {
	return &Api{registry}
}

func (g *Api) Get(ctx context.Context, id string) (*models.MetaData, error) {
	addrs, err := g.registry.GetServiceAddresses(ctx, "metadata")
	if err != nil {
		return nil, err
	}

	//address at random index between 0 and addrs length
	randomAddress := addrs[rand.Intn(len(addrs))]
	url := "http://" + randomAddress + "/api/v1/metadata"

	log.Printf("Calling metadata service. Request: GET " + url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
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
		return nil, api.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", resp)
	}

	var meta *models.MetaData
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return nil, err
	}

	return meta, nil
}
