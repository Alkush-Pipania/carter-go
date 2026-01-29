package upload

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Alkush-Pipania/carter-go/internal/middleware"
	"github.com/Alkush-Pipania/carter-go/pkg/response"
	"github.com/Alkush-Pipania/carter-go/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	RequestUploadURL(ctx context.Context, userID string, req PresignUploadRequest) (*PresignUploadResponse, error)
	ConfirmUpload(ctx context.Context, userID string, sourceID string) error
}

type Handler struct {
	service   Service
	validator *validator.Validate
}

func NewHandler(svc Service, validator *validator.Validate) *Handler {
	return &Handler{
		service:   svc,
		validator: validator,
	}
}

// RequestUploadURL handles requests for presigned S3 upload URLs
func (h *Handler) RequestUploadURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok || userID == "" {
		response.WriteError(w, http.StatusUnauthorized, "User ID is required")
		return
	}

	var req PresignUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate struct using go-playground/validator
	if err := h.validator.Struct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate content type
	if !utils.IsAllowedContentType(req.ContentType) {
		response.WriteError(w, http.StatusBadRequest, "Unsupported content type. Allowed: PDF, PPT, DOC")
		return
	}

	result, err := h.service.RequestUploadURL(ctx, userID, req)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to generate upload URL")
		return
	}

	response.WriteJSON(w, http.StatusOK, result)
}

// ConfirmUpload handles upload confirmation after client uploads to S3
func (h *Handler) ConfirmUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok || userID == "" {
		response.WriteError(w, http.StatusUnauthorized, "User ID is required")
		return
	}

	sourceID := chi.URLParam(r, "id")
	if sourceID == "" {
		response.WriteError(w, http.StatusBadRequest, "Source ID is required")
		return
	}

	err := h.service.ConfirmUpload(ctx, userID, sourceID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to confirm upload")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"status": "confirmed"})
}
