package collection

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

type Service interface {
	GetCollectionsByUserID(ctx context.Context, userID string) ([]db.Collection, error)
	CreateCollection(ctx context.Context, userID string, name string) (*db.Collection, error)
	UpdateCollection(ctx context.Context, id string, name string) (*db.Collection, error)
	DeleteCollection(ctx context.Context, id string) error
}

type Handler struct {
	service   Service
	validator *validator.Validate
}

func NewHandler(service Service, validator *validator.Validate) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
	}
}

func (h *Handler) GetCollections(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	collections, err := h.service.GetCollectionsByUserID(r.Context(), userID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to retrieve collections")
		return
	}

	resp := ToCollectionResponses(collections)

	response.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req CreateCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	collection, err := h.service.CreateCollection(r.Context(), userID, req.Name)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to create collection")
		return
	}

	response.WriteJSON(w, http.StatusCreated, ToCollectionResponse(collection))
}

func (h *Handler) UpdateCollection(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	collection, err := h.service.UpdateCollection(r.Context(), id, req.Name)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to update collection")
		return
	}

	response.WriteJSON(w, http.StatusOK, ToCollectionResponse(collection))
}

func (h *Handler) DeleteCollection(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteCollection(r.Context(), id); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to delete collection")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
