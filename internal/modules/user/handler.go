package user

import (
	"net/http"

	"github.com/Alkush-Pipania/carter-go/internal/middleware"
	"github.com/Alkush-Pipania/carter-go/pkg/errors"
)

type UserHandler struct {
	userService Service
}

func NewUserHandler(userService Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		errors.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := UserProfileResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Image:     user.Image,
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt,
	}

	errors.RespondWithJSON(w, http.StatusOK, response)
}
