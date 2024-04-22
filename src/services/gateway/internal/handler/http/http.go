package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	gateway "moviemicroservice.com/src/services/gateway/internal/services/movies"
)

type Handler struct {
	service *gateway.Service
}

func New(service *gateway.Service) *Handler {
	return &Handler{service}
}

func (h *Handler) GetMovieDetails(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	details, err := h.service.Get(req.Context(), id)

	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("Get error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(details); err != nil {
		log.Printf("Encode error: %v\n", err)
	}
}
