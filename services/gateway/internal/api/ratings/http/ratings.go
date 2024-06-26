package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"moviemicroservice.com/pkg/discovery"
	api "moviemicroservice.com/services/gateway/internal/api"
	"moviemicroservice.com/services/ratings/pkg/models"
)

type Api struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Api {
	return &Api{registry}
}

func (g *Api) GetAggregatedRating(ctx context.Context, recordID models.RecordID, recordType models.RecordType) (float64, error) {
	addrs, err := g.registry.GetServiceAddresses(ctx, "ratings")
	if err != nil {
		return 0, err
	}

	//address at random index between 0 and addrs length
	randomAddress := addrs[rand.Intn(len(addrs))]
	url := "http://" + randomAddress + "/api/v1/ratings"

	log.Printf("Calling ratings service. Request: GET " + url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	req = req.WithContext(ctx)

	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", fmt.Sprintf("%v", recordType))
	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return 0, api.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("non-2xx response: %v", resp)
	}

	var v float64
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, err
	}

	return v, nil
}

func (g *Api) PutRating(ctx context.Context, recordID models.RecordID, recordType models.RecordType, rating *models.Rating) error {
	addrs, err := g.registry.GetServiceAddresses(ctx, "ratings")
	if err != nil {
		return err
	}

	//address at random index between 0 and addrs length
	randomAddress := addrs[rand.Intn(len(addrs))]
	url := "http://" + randomAddress + "/api/v1/ratings"

	log.Printf("Calling ratings service. Request: PUT " + url)

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}

	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", fmt.Sprintf("%v", recordType))
	values.Add("userId", string(rating.UserID))
	values.Add("value", fmt.Sprintf("%v", rating.Value))
	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("non-2xx response: %v", resp)
	}

	return nil
}
