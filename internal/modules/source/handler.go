package source

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Alkush-Pipania/carter-go/internal/middleware"
	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/Alkush-Pipania/carter-go/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service   Service
	validator *validator.Validate
}

type Service interface {
	GetSourcesByCollectionID(ctx context.Context, collectionID string) ([]db.Source, error)
	CreateSource(ctx context.Context, userID string, req CreateSourceRequest) (db.Source, error)
	GetSourceByID(ctx context.Context, sourceID string) (db.Source, error)
	DeleteSource(ctx context.Context, sourceID string) error
}

func NewHandler(svc Service, validator *validator.Validate) *Handler {
	return &Handler{
		service:   svc,
		validator: validator,
	}
}

func (h *Handler) GetSourcesByCollectionID(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) CreateSource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req CreateSourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Type-specific validation
	switch req.Type {
	case "link":
		if req.OriginalUrl == "" {
			response.WriteError(w, http.StatusBadRequest, "Original URL is required for link type")
			return
		}
	case "note":
		if req.Content == "" {
			response.WriteError(w, http.StatusBadRequest, "Content is required for note type")
			return
		}
	}

	// Create source via service (handles type routing internally)
	source, err := h.service.CreateSource(ctx, userID, req)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, ToSourceResponse(&source))
}

func (h *Handler) GetSourceByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sourceID := chi.URLParam(r, "id")
	if sourceID == "" {
		response.WriteError(w, http.StatusBadRequest, "Source ID is required")
		return
	}

	source, err := h.service.GetSourceByID(ctx, sourceID)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "Source not found")
		return
	}

	response.WriteJSON(w, http.StatusOK, ToSourceResponse(&source))
}

func (h *Handler) DeleteSource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sourceID := chi.URLParam(r, "id")
	if sourceID == "" {
		response.WriteError(w, http.StatusBadRequest, "Source ID is required")
		return
	}

	if err := h.service.DeleteSource(ctx, sourceID); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to delete source")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
