package user

import "github.com/go-chi/chi/v5"

func Routes(h *UserHandler) chi.Router {
	r := chi.NewRouter()

	r.Get("/me", h.GetProfile)

	return r
}
