package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"moviemicroservice.com/src/modules/movies/internal/gateway"
	"moviemicroservice.com/src/modules/ratings/pkg/models"
)

type Gateway struct {
	address string
}

func New(address string) *Gateway {
	return &Gateway{address}
}

func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID models.RecordID, recordType models.RecordType) (float64, error) {
	req, err := http.NewRequest(http.MethodGet, g.address+"/api/v1/ratings", nil)
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
		return 0, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("non-2xx response: %v", resp)
	}

	var v float64
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, err
	}

	return v, nil
}

func (g *Gateway) PutRating(ctx context.Context, recordID models.RecordID, recordType models.RecordType, rating *models.Rating) error {
	req, err := http.NewRequest(http.MethodPut, g.address+"/api/v1/ratings", nil)
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
