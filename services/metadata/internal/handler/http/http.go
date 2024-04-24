package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"moviemicroservice.com/services/metadata/internal/repository"
	"moviemicroservice.com/services/metadata/internal/service/metadata"
)

type Handler struct {
	service *metadata.Service
}

func New(service *metadata.Service) *Handler {
	return &Handler{service}
}

func (h *Handler) Get(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := req.Context()
	metadata, err := h.service.Get(ctx, id)

	if err != nil && errors.Is(err, repository.ErrNotFound) {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("Error getting metadata from repo: %v\n", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(res).Encode(metadata); err != nil {
		log.Printf("Response encode error: %v\n", err)
	}
}
