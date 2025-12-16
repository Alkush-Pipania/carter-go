package collection

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/go-chi/chi/v5"
)

type Service interface {
	GetCollectionsByUserID(ctx context.Context, userID string) ([]db.Collection, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetCollections(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	collections, err := h.service.GetCollectionsByUserID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve collections")
		return
	}

	response := ToCollectionResponses(collections)

	writeJSON(w, http.StatusOK, response)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}
