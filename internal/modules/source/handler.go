package source

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/Alkush-Pipania/carter-go/pkg/response"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	service Service
}

type Service interface {
	GetSourcesByCollectionID(ctx context.Context, collectionID string) ([]db.Source, error)
	CreateSource(ctx context.Context, userID string, req CreateSourceRequest) (db.Source, error)
}

func NewHandler(svc Service) *handler {
	return &handler{
		service: svc,
	}
}

func (h *handler) GetSourcesByCollectionID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	collectionID := chi.URLParam(r, "id")
	if collectionID == "" {
		response.WriteError(w, http.StatusBadRequest, "Collection ID is required")
		return
	}
	sources, err := h.service.GetSourcesByCollectionID(ctx, collectionID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to retrieve sources")
		return
	}

	response.WriteJSON(w, http.StatusOK, ToSourceResponses(sources))
}

func (h *handler) CreateSource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := r.Header.Get("user_id")
	if userID == "" {
		response.WriteError(w, http.StatusUnauthorized, "User ID is required")
		return
	}

	var req CreateSourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.CollectionID == "" {
		response.WriteError(w, http.StatusBadRequest, "Collection ID is required")
		return
	}
	if req.Type == "" {
		response.WriteError(w, http.StatusBadRequest, "Source type is required")
		return
	}
	if req.Title == "" {
		response.WriteError(w, http.StatusBadRequest, "Title is required")
		return
	}

	// Validate source type
	validTypes := map[string]bool{"link": true, "pdf": true, "ppt": true, "doc": true, "note": true}
	if !validTypes[req.Type] {
		response.WriteError(w, http.StatusBadRequest, "Invalid source type. Must be one of: link, pdf, ppt, doc, note")
		return
	}

	// For link type, original_url is required
	if req.Type == "link" && req.OriginalUrl == "" {
		response.WriteError(w, http.StatusBadRequest, "Original URL is required for link type")
		return
	}

	// Create source via service
	source, err := h.service.CreateSource(ctx, userID, req)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to create source")
		return
	}

	response.WriteJSON(w, http.StatusCreated, ToSourceResponse(source))
}
