package authentication

import (
	"context"
	"encoding/json"
	"net/http"
)

type Service interface {
	Login(ctx context.Context, email string, password string) (*LoginResponse, error)
	Register(ctx context.Context, password string, email string) error
	Logout(ctx context.Context, sessionID string) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    resp.SessionID,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID := readSessionCookie(r)
	if sessionID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.service.Logout(r.Context(), sessionID); err != nil {
		http.Error(w, "logout failed", http.StatusInternalServerError)
		return
	}

	clearSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.Register(r.Context(), req.Password, req.Email); err != nil {
		http.Error(w, "registration failed", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
