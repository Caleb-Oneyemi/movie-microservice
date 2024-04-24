package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"moviemicroservice.com/services/ratings/internal/repository"
	"moviemicroservice.com/services/ratings/internal/service/ratings"
	"moviemicroservice.com/services/ratings/pkg/models"
)

type Handler struct {
	service *ratings.Service
}

func New(service *ratings.Service) *Handler {
	return &Handler{service}
}

func (h *Handler) Handle(res http.ResponseWriter, req *http.Request) {
	recordID := models.RecordID(req.FormValue("id"))
	if recordID == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	recordType := models.RecordType(req.FormValue("type"))
	if recordType == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := req.Context()

	switch req.Method {
	case http.MethodGet:
		ratings, err := h.service.GetAggregatedRatings(ctx, recordType, recordID)
		if err != nil && errors.Is(err, repository.ErrNotFound) {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(res).Encode(ratings); err != nil {
			log.Printf("Response encode error: %v\n", err)
		}
	case http.MethodPut:
		userID := models.UserID(req.FormValue("userId"))
		if userID == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		value, err := strconv.ParseFloat(req.FormValue("value"), 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := h.service.Put(ctx, recordType, recordID, &models.Rating{UserID: string(userID), Value: models.RatingValue(value)}); err != nil {
			log.Printf("Repository put error: %v\n", err)
			res.WriteHeader(http.StatusInternalServerError)
		}

	default:
		res.WriteHeader(http.StatusBadRequest)
	}
}
