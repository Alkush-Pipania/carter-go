package authentication

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Alkush-Pipania/carter-go/pkg/logger"
	"go.uber.org/zap"
)

type ServiceInter interface {
	Login(ctx context.Context, email string, password string) (*LoginResponse, error)
	Register(ctx context.Context, password string, email string) error
	Logout(ctx context.Context, sessionID string) error
}

type Handler struct {
	service ServiceInter
}

func NewHandler(service ServiceInter) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Handling login request", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Login handler: invalid request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		logger.Warn("Login handler: authentication failed", zap.String("email", req.Email), zap.Error(err))
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

	logger.Info("Login handler: login successful", zap.String("user_id", resp.UserID))
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Handling logout request", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	sessionID, err := readSessionCookie(r)
	if err != nil {
		logger.Warn("Logout handler: no valid session cookie", zap.Error(err))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.service.Logout(r.Context(), sessionID); err != nil {
		logger.Error("Logout handler: logout failed", zap.String("session_id", sessionID), zap.Error(err))
		http.Error(w, "logout failed", http.StatusInternalServerError)
		return
	}

	clearSessionCookie(w)
	logger.Info("Logout handler: logout successful", zap.String("session_id", sessionID))
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Handling register request", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Register handler: invalid request body", zap.Error(err))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.Register(r.Context(), req.Password, req.Email); err != nil {
		logger.Error("Register handler: registration failed", zap.String("email", req.Email), zap.Error(err))
		http.Error(w, "registration failed", http.StatusBadRequest)
		return
	}

	logger.Info("Register handler: registration successful", zap.String("email", req.Email))
	w.WriteHeader(http.StatusCreated)
}
